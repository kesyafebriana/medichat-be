package server

import (
	"medichat-be/apperror"
	"medichat-be/handler"

	"github.com/gin-gonic/gin"
)

type SetupServerOpts struct {
	PingHandler *handler.PingHandler
	ChatHandler *handler.ChatHandler
	
	RequestID     gin.HandlerFunc
	Authenticator gin.HandlerFunc
	CorsHandler   gin.HandlerFunc
	Logger        gin.HandlerFunc
	ErrorHandler  gin.HandlerFunc
}

func SetupServer(opts SetupServerOpts) *gin.Engine {
	router := gin.New()
	router.ContextWithFallback = true
	router.Use(
		opts.RequestID,
		opts.Logger,
		gin.Recovery(),
		opts.CorsHandler,
		opts.ErrorHandler,
	)

	apiV1Group := router.Group("/api/v1")

	apiV1Group.GET(
		"/ping",
		opts.PingHandler.Ping,
	)

	chatGroup := apiV1Group.Group("/chat")

	chatGroup.POST("/send", opts.ChatHandler.Chat)

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
