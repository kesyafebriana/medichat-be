package server

import (
	"medichat-be/apperror"
	"medichat-be/handler"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SetupServerOpts struct {
	AccountHandler    *handler.AccountHandler
	PingHandler       *handler.PingHandler
	ChatHandler       *handler.ChatHandler
	GoogleAuthHandler *handler.OAuth2Handler
	GoogleHandler     *handler.GoogleHandler
	CategoryHandler   *handler.CategoryHandler

	SessionKey []byte

	RequestID gin.HandlerFunc

	Authenticator      gin.HandlerFunc
	AdminAuthenticator gin.HandlerFunc

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

	chatGroup := apiV1Group.Group("/chat")

	chatGroup.POST("/send", opts.ChatHandler.Chat)
	chatGroup.PATCH("/close", opts.ChatHandler.CloseRoom)
	chatGroup.POST("/create", opts.ChatHandler.CreateRoom)

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

	router.NoRoute(func(ctx *gin.Context) {
		ctx.Error(apperror.NewAppError(
			apperror.CodeNotFound,
			"route not found",
			nil,
		))
		ctx.Abort()
	})

	categoryGroup := apiV1Group.Group("/categories")
	categoryGroup.GET("/", opts.Authenticator, opts.CategoryHandler.GetCategories)
	categoryGroup.GET("/:slug", opts.Authenticator, opts.CategoryHandler.GetCategoryBySlug)
	categoryGroup.GET("/hierarchy", opts.Authenticator, opts.CategoryHandler.GetCategoriesHierarchy)
	categoryGroup.POST("/", opts.AdminAuthenticator, opts.CategoryHandler.CreateCategoryLevelOne)
	categoryGroup.POST("/:slug", opts.AdminAuthenticator, opts.CategoryHandler.CreateCategoryLevelTwo)
	categoryGroup.PATCH("/:slug", opts.AdminAuthenticator, opts.CategoryHandler.UpdateCategory)
	categoryGroup.DELETE("/:slug", opts.AdminAuthenticator, opts.CategoryHandler.DeleteCategory)

	return router
}
