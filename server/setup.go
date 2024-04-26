package server

import (
	"medichat-be/apperror"
	"medichat-be/handler"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SetupServerOpts struct {
	AccountHandler        *handler.AccountHandler
	PingHandler           *handler.PingHandler
	GoogleAuthHandler     *handler.OAuth2Handler
	GoogleHandler         *handler.GoogleHandler
	UserHandler           *handler.UserHandler
	DoctorHandler         *handler.DoctorHandler
	SpecializationHandler *handler.SpecializationHandler

	SessionKey []byte

	RequestID gin.HandlerFunc

	Authenticator       gin.HandlerFunc
	UserAuthenticator   gin.HandlerFunc
	DoctorAuthenticator gin.HandlerFunc

	CorsHandler  gin.HandlerFunc
	Logger       gin.HandlerFunc
	ErrorHandler gin.HandlerFunc
}

func SetupServer(opts SetupServerOpts) *gin.Engine {
	router := gin.New()
	router.ContextWithFallback = true

	sessionStore := cookie.NewStore(opts.SessionKey)

	router.Use(
		opts.RequestID,
		opts.Logger,
		gin.Recovery(),
		opts.CorsHandler,
		sessions.Sessions("session", sessionStore),
		opts.ErrorHandler,
	)

	apiV1Group := router.Group("/api/v1")

	apiV1Group.GET(
		"/ping",
		opts.PingHandler.Ping,
	)

	authGroup := apiV1Group.Group("/auth")
	authGroup.POST(
		"/register",
		opts.AccountHandler.Register,
	)
	authGroup.POST(
		"/login",
		opts.AccountHandler.Login,
	)
	authGroup.POST(
		"/forget-password",
		opts.AccountHandler.ForgetPassword,
	)
	authGroup.POST(
		"/reset-password",
		opts.AccountHandler.ResetPassword,
	)
	authGroup.GET(
		"/check-reset-token",
		opts.AccountHandler.CheckResetPasswordToken,
	)
	authGroup.POST(
		"/verify-token",
		opts.AccountHandler.GetVerifyEmailToken,
	)
	authGroup.POST(
		"/verify-email",
		opts.AccountHandler.VerifyEmail,
	)
	authGroup.GET(
		"/check-verify-token",
		opts.AccountHandler.CheckVerifyEmailToken,
	)
	authGroup.POST(
		"/refresh",
		opts.AccountHandler.RefreshTokens,
	)
	authGroup.GET(
		"/profile",
		opts.Authenticator,
		opts.AccountHandler.GetProfile,
	)

	googleGroup := apiV1Group.Group("/google")
	googleGroup.GET(
		"/auth",
		opts.GoogleAuthHandler.GetAuthURL,
	)
	googleGroup.GET(
		"/callback",
		opts.GoogleHandler.OAuth2Callback,
	)

	userGroup := apiV1Group.Group(
		"/users",
	)

	userProfileGroup := userGroup.Group(
		"/profile",
		opts.UserAuthenticator,
	)
	userProfileGroup.GET(
		".",
		opts.UserHandler.GetProfile,
	)
	userProfileGroup.POST(
		".",
		opts.UserHandler.CreateProfile,
	)
	userProfileGroup.PUT(
		".",
		opts.UserHandler.UpdateProfile,
	)
	userProfileGroup.POST(
		"/locations",
		opts.UserHandler.AddLocation,
	)
	userProfileGroup.PUT(
		"/locations",
		opts.UserHandler.UpdateLocation,
	)
	userProfileGroup.DELETE(
		"/locations/:id",
		opts.UserHandler.DeleteLocation,
	)

	doctorGroup := apiV1Group.Group(
		"/doctors",
	)
	doctorGroup.GET(
		".",
		opts.DoctorHandler.ListDoctors,
	)

	doctorProfileGroup := doctorGroup.Group(
		"/profile",
		opts.DoctorAuthenticator,
	)
	doctorProfileGroup.GET(
		".",
		opts.DoctorHandler.GetProfile,
	)
	doctorProfileGroup.POST(
		".",
		opts.DoctorHandler.CreateProfile,
	)
	doctorProfileGroup.PUT(
		".",
		opts.DoctorHandler.UpdateProfile,
	)
	doctorProfileGroup.POST(
		"/active-status",
		opts.DoctorHandler.SetActiveStatus,
	)

	specializationGroup := apiV1Group.Group(
		"/specializations",
	)
	specializationGroup.GET(
		".",
		opts.SpecializationHandler.GetAll,
	)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.Error(apperror.NewAppError(
			apperror.CodeNotFound,
			"route not found",
			nil,
		))
		ctx.Abort()
	})

	return router
}
