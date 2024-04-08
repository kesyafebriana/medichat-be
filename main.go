package main

import (
	"context"
	"log"
	"medichat-be/config"
	"medichat-be/cryptoutil"
	"medichat-be/database"
	"medichat-be/handler"
	"medichat-be/logger"
	"medichat-be/middleware"
	"medichat-be/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loaiding env config: %v", err)
	}

	l := logrus.New()
	log := logger.FromLogrus(l)

	db, err := database.ConnectPostgresDB(conf.DatabaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database %v", err)
	}
	defer db.Close()

	jwtProvider := cryptoutil.NewJWTProviderHS256(
		conf.JWTIssuer,
		conf.JWTSecret,
		conf.JWTLifespan,
	)

	pingHandler := handler.NewPingHandler()

	requestIDMid := middleware.RequestIDHandler()
	loggerMid := middleware.Logger(log)
	corsHandler := middleware.CorsHandler()
	errorHandler := middleware.ErrorHandler()
	authenticator := middleware.Authenticator(jwtProvider)

	router := server.SetupServer(server.SetupServerOpts{
		PingHandler: pingHandler,

		RequestID:     requestIDMid,
		Authenticator: authenticator,
		CorsHandler:   corsHandler,
		Logger:        loggerMid,
		ErrorHandler:  errorHandler,
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error running server: %v", err)
		}
	}()

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
