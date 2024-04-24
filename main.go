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

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	config.InitConfig()
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading env config: %v", err)
	}

	l := logrus.New()
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

	refreshProvider := cryptoutil.NewJWTProviderHS256(
		conf.JWTIssuer,
		conf.RefreshSecret,
		conf.RefreshTokenLifespan,
	)

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
	})

	googleAuthService := service.NewOAuth2Service(service.OAuth2ServiceOpts{
		OAuth2Provider: googleAuthProvider,
	})
	googleService := service.NewGoogleService(service.GoogleServiceOpts{
		DataRepository: dataRepository,
		OAuth2Service:  googleAuthService,
		AccountService: accountService,
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

	requestIDMid := middleware.RequestIDHandler()
	loggerMid := middleware.Logger(log)
	corsHandler := middleware.CorsHandler(conf.FEDomain)
	errorHandler := middleware.ErrorHandler()

	authenticator := middleware.Authenticator(anyAccessProvider)

	router := server.SetupServer(server.SetupServerOpts{
		AccountHandler:    accountHandler,
		PingHandler:       pingHandler,
		GoogleAuthHandler: googleAuthHandler,
		GoogleHandler:     googleHandler,
		CategoryHandler:   categoryHandler,

		SessionKey: conf.SessionKey,

		RequestID:     requestIDMid,
		Authenticator: authenticator,
		CorsHandler:   corsHandler,
		Logger:        loggerMid,
		ErrorHandler:  errorHandler,
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
