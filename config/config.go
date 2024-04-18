package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	ErrMissingKey = errors.New("missing key")
)

type Config struct {
	ServerAddr        string
	DatabaseURL       string
	JWTIssuer         string
	JWTSecret         string
	JWTLifespan       time.Duration
	AuthEmail         string
	AuthEmailPassword string
	EmailSender       string
}

func InitConfig() error {
	return godotenv.Load()
}

func LoadConfig() (Config, error) {
	ret := Config{}

	ret.ServerAddr = os.Getenv("SERVER_ADDR")

	ret.DatabaseURL = os.Getenv("DATABASE_URL")

	ret.JWTIssuer = os.Getenv("JWT_ISSUER")

	ret.JWTSecret = os.Getenv("JWT_SECRET")

	ret.AuthEmail = os.Getenv("AUTH_EMAIL")

	ret.AuthEmailPassword = os.Getenv("AUTH_EMAIL_PASSWORD")

	ret.EmailSender = os.Getenv("EMAIL_SENDER")

	s := os.Getenv("JWT_LIFESPAN")
	i, err := strconv.Atoi(s)
	if err != nil {
		return Config{}, err
	}
	ret.JWTLifespan = time.Duration(i) * time.Minute

	return ret, nil
}
