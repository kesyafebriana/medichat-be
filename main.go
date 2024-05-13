package main

import (
	"context"
	"fmt"
	"log"
	"medichat-be/apperror"
	"medichat-be/config"
	"medichat-be/constants"
	"medichat-be/cryptoutil"
	"medichat-be/database"
	"medichat-be/handler"
	"medichat-be/logger"
	"medichat-be/middleware"
	"medichat-be/repository/postgres"
	"medichat-be/server"
	"medichat-be/service"
	"medichat-be/util"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	firebase "firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func main() {
	config.InitConfig()
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading env config: %v", err)
	}

	l := logrus.New()

	if conf.IsRelease {
		infofile, err := os.OpenFile("/var/log/medichat.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer infofile.Close()
		l.SetOutput(infofile)
	}

	log := logger.FromLogrus(l)

	db, err := database.ConnectPostgresDB(conf.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	}
	defer db.Close()

	util.InitValidators()

	apperror.SetIncludeStackTrace(!conf.IsRelease)

	if conf.IsRelease {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	adminAccessProvider := cryptoutil.NewJWTProviderHS256(
		conf.JWTIssuer,
		conf.AdminAccessSecret,
		conf.AccessTokenLifespan,
	)

	userAccessProvider := cryptoutil.NewJWTProviderHS256(
		conf.JWTIssuer,
		conf.UserAccessSecret,
		conf.AccessTokenLifespan,
	)

	doctorAccessProvider := cryptoutil.NewJWTProviderHS256(
		conf.JWTIssuer,
		conf.DoctorAccessSecret,
		conf.AccessTokenLifespan,
	)

	pharmacyManagerAccessProvider := cryptoutil.NewJWTProviderHS256(
		conf.JWTIssuer,
		conf.PharmacyManagerAccessSecret,
		conf.AccessTokenLifespan,
	)

	anyAccessProvider := cryptoutil.NewJWTProviderAny([]cryptoutil.JWTProvider{
		adminAccessProvider,
		userAccessProvider,
		doctorAccessProvider,
		pharmacyManagerAccessProvider,
	})

	managerOrAdminAccessProvider := cryptoutil.NewJWTProviderAny([]cryptoutil.JWTProvider{
		adminAccessProvider,
		pharmacyManagerAccessProvider,
	})

	userOrAdminAccessProvider := cryptoutil.NewJWTProviderAny([]cryptoutil.JWTProvider{
		adminAccessProvider,
		userAccessProvider,
	})

	userOrManagerAccessProvider := cryptoutil.NewJWTProviderAny([]cryptoutil.JWTProvider{
		userAccessProvider,
		pharmacyManagerAccessProvider,
	})

	userOrManagerOrAdminAccessProvider := cryptoutil.NewJWTProviderAny([]cryptoutil.JWTProvider{
		adminAccessProvider,
		userAccessProvider,
		pharmacyManagerAccessProvider,
	})

	refreshProvider := cryptoutil.NewJWTProviderHS256(
		conf.JWTIssuer,
		conf.RefreshSecret,
		conf.RefreshTokenLifespan,
	)
	ctx := context.Background()
	sa := option.WithCredentialsFile("./serviceAccount.json")
	firebaseConfig := &firebase.Config{ProjectID: "rapunzel-medichat"}

	app, err := firebase.NewApp(ctx, firebaseConfig, sa)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Error connecting to firebase %v", err)
	}
	defer client.Close()

	cld, _ := util.NewCloudinarylProvider()

	passwordHasher := cryptoutil.NewPasswordHasherBcrypt(constants.HashCost)

	resetPasswordTokenProvider := cryptoutil.NewRandomTokenProvider(
		constants.ResetPasswordTokenByteLength,
	)
	verifyEmailTokenProvider := cryptoutil.NewRandomTokenProvider(
		constants.VerifyEmailTokenByteLength,
	)

	googleAuthProvider := cryptoutil.NewGoogleAuthProvider(cryptoutil.GoogleAuthProviderOpts{
		RedirectURL:  conf.GoogleAPIRedirectURL,
		ClientID:     conf.GoogleAPIClientID,
		ClientSecret: conf.GoogleAPIClientSecret,
	})
	googleAuthStateProvider := cryptoutil.NewRandomTokenProvider(
		constants.GoogleAuthStateByteLength,
	)

	appEmail, err := util.NewAppEmail(util.AppEmailOpts{
		FEVerivicationURL:  conf.FEVerificationURL,
		FEResetPasswordURL: conf.FEResetPasswordURL,
	})
	if err != nil {
		log.Fatalf("Error creating app email: %v", err)
	}
	emailProvider := util.NewGmailProvider(util.EmailProviderOpts{
		Username:    conf.AuthEmailUsername,
		Password:    conf.AuthEmailPassword,
		EmailSender: fmt.Sprintf(conf.EmailSender, conf.AuthEmailUsername),
	})

	dataRepository := postgres.NewDataRepository(db)

	chatService := service.NewChatService(service.ChatServiceOpts{
		DataRepository: dataRepository,
		Client:         client,
		Cloud:          cld,
	})

	accountService := service.NewAccountService(service.AccountServiceOpts{
		DataRepository:                dataRepository,
		PasswordHasher:                passwordHasher,
		AdminAccessProvider:           adminAccessProvider,
		UserAccessProvider:            userAccessProvider,
		DoctorAccessProvider:          doctorAccessProvider,
		PharmacyManagerAccessProvider: pharmacyManagerAccessProvider,
		RefreshProvider:               refreshProvider,
		RPTProvider:                   resetPasswordTokenProvider,
		RPTLifespan:                   conf.ResetPasswordTokenLifespan,
		VETProvider:                   verifyEmailTokenProvider,
		VETLifespan:                   conf.VerifyEmailTokenLifespan,
		AppEmail:                      appEmail,
		EmailProvider:                 emailProvider,
	})

	categoryService := service.NewCategoryService(service.CategoryServiceOpts{
		DataRepository: dataRepository,
		Cloud:          cld,
	})

	googleAuthService := service.NewOAuth2Service(service.OAuth2ServiceOpts{
		OAuth2Provider: googleAuthProvider,
	})
	googleService := service.NewGoogleService(service.GoogleServiceOpts{
		DataRepository: dataRepository,
		OAuth2Service:  googleAuthService,
		AccountService: accountService,
	})

	userService := service.NewUserService(service.UserServiceOpts{
		DataRepository: dataRepository,
		CloudProvider:  cld,
	})
	doctorService := service.NewDoctorService(service.DoctorServiceOpts{
		DataRepository: dataRepository,
		CloudProvider:  cld,
	})

	specializationService := service.NewSpecializationService(service.SpecializationServiceOpts{
		DataRepository: dataRepository,
	})

	productService := service.NewProductService(service.ProductServiceOpts{
		DataRepository: dataRepository,
		Cloud:          cld,
	})

	pharmacyService := service.NewPharmacyService(service.PharmacyServiceOpts{
		DataRepository: dataRepository,
	})

	pharmacyManagerService := service.NewPharmacyManagerService(service.PharmacyManagerServiceOpts{
		DataRepository: dataRepository,
		CloudProvider:  cld,
	})

	stockService := service.NewStockService(service.StockServiceOpts{
		DataRepository: dataRepository,
	})

	paymentService := service.NewPaymentService(service.PaymentServiceOpts{
		DataRepository: dataRepository,
		CloudProvider:  cld,
	})

	orderService := service.NewOrderService(service.OrderServiceOpts{
		DataRepository: dataRepository,
		CloudProvider:  cld,
	})

	accountHandler := handler.NewAccountHandler(handler.AccountHandlerOpts{
		AccountSrv: accountService,
		Domain:     conf.WebDomain,
	})
	categoryHandler := handler.NewCategoryHandler(handler.CategoryHandlerOpts{
		CategorySrv: categoryService,
		Domain:      conf.WebDomain,
	})
	pingHandler := handler.NewPingHandler()
	googleAuthHandler := handler.NewOAuth2Handler(handler.OAuth2HandlerOpts{
		OAuth2Service:       googleAuthService,
		RandomTokenProvider: googleAuthStateProvider,
	})
	googleHandler := handler.NewGoogleHandler(handler.GoogleHandlerOpts{
		GoogleSrv: googleService,
		Domain:    conf.WebDomain,
	})

	userHandler := handler.NewUserHandler(handler.UserHandlerOpts{
		UserSrv: userService,
	})
	doctorHandler := handler.NewDoctorHandler(handler.DoctorHandlerOpts{
		DoctorSrv: doctorService,
	})

	specializationHandler := handler.NewSpecializationHandler(handler.SpecializationHandlerOpts{
		SpecializationSrv: specializationService,
	})

	productHandler := handler.NewProductHandler(handler.ProductHandlerOpts{
		ProductSrv: productService,
	})

	pharmacyHandler := handler.NewPharmacyHandler(handler.PharmacyHandlerOpts{
		PharmacySrv: pharmacyService,
	})
	chatHandler := handler.NewChatHandler(chatService)

	pharmacyManagerHandler := handler.NewPharmacyManagerHandler(handler.PharmacyManagerHandlerOpts{
		PharmacyManagerSrv: pharmacyManagerService,
	})

	stockHandler := handler.NewStockHandler(handler.StockHandlerOpts{
		StockSrv: stockService,
	})

	paymentHandler := handler.NewPaymentHandler(handler.PaymentHandlerOpts{
		PaymentSrv: paymentService,
	})

	orderHandler := handler.NewOrderHandler(handler.OrderHandlerOpts{
		OrderSrv: orderService,
	})

	requestIDMid := middleware.RequestIDHandler()
	loggerMid := middleware.Logger(log)
	corsHandler := middleware.CorsHandler(conf.FEDomain)
	errorHandler := middleware.ErrorHandler()

	authenticator := middleware.Authenticator(anyAccessProvider)
	adminAuthenticator := middleware.Authenticator(adminAccessProvider)
	userAuthenticator := middleware.Authenticator(userAccessProvider)
	doctorAuthenticator := middleware.Authenticator(doctorAccessProvider)
	pharmacyManagerAuthenticator := middleware.Authenticator(pharmacyManagerAccessProvider)

	managerOrAdminAuthenticator := middleware.Authenticator(managerOrAdminAccessProvider)

	UserOrAdminAuthenticator := middleware.Authenticator(userOrAdminAccessProvider)

	userOrManagerAuthenticator := middleware.Authenticator(userOrManagerAccessProvider)
	userOrManagerOrAdminAuthenticator := middleware.Authenticator(userOrManagerOrAdminAccessProvider)

	router := server.SetupServer(server.SetupServerOpts{
		AccountHandler:         accountHandler,
		ChatHandler:            chatHandler,
		PingHandler:            pingHandler,
		GoogleAuthHandler:      googleAuthHandler,
		GoogleHandler:          googleHandler,
		UserHandler:            userHandler,
		DoctorHandler:          doctorHandler,
		SpecializationHandler:  specializationHandler,
		CategoryHandler:        categoryHandler,
		ProductHandler:         productHandler,
		PharmacyHandler:        pharmacyHandler,
		PharmacyManagerHandler: pharmacyManagerHandler,
		StockHandler:           stockHandler,
		PaymentHandler:         paymentHandler,
		OrderHandler:           orderHandler,

		SessionKey: conf.SessionKey,

		RequestID:                    requestIDMid,
		Authenticator:                authenticator,
		AdminAuthenticator:           adminAuthenticator,
		UserAuthenticator:            userAuthenticator,
		DoctorAuthenticator:          doctorAuthenticator,
		PharmacyManagerAuthenticator: pharmacyManagerAuthenticator,
		CorsHandler:                  corsHandler,
		Logger:                       loggerMid,
		ErrorHandler:                 errorHandler,

		ManagerOrAdminAuthenticator: managerOrAdminAuthenticator,

		UserOrAdminAuthenticator: UserOrAdminAuthenticator,

		UserOrManagerAuthenticator: userOrManagerAuthenticator,

		UserOrManagerOrAdminAuthenticator: userOrManagerOrAdminAuthenticator,
	})

	srv := &http.Server{
		Addr:    conf.ServerAddr,
		Handler: router,
	}

	log.Info("Starting Server...")

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error running server: %v", err)
		}
	}()

	log.Infof("Server started at %s", conf.ServerAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	shFinish := make(chan struct{})
	go func() {
		if err := srv.Shutdown(shCtx); err != nil {
			log.Fatalf("Error shutting down server: %v", err)
		}
		shFinish <- struct{}{}
	}()

	select {
	case <-shCtx.Done():
		log.Info("Shutdown timeout.")
	case <-shFinish:
		log.Info("Shutdown finished.")
	}
}
