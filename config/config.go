package config

import (
	"encoding/base64"
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var (
	ErrMissingKey = errors.New("missing key")
)

type Config struct {
	ServerAddr  string
	DatabaseURL string

	SessionKey []byte

	JWTIssuer   string
	JWTSecret   string
	JWTLifespan time.Duration

	GoogleAPIClientID     string
	GoogleAPIClientSecret string
	GoogleAPIRedirectURL  string
}

func InitConfig() error {
	return godotenv.Load()
}

func LoadConfig() (Config, error) {
	ret := Config{}

	ret.ServerAddr = os.Getenv("SERVER_ADDR")

	ret.DatabaseURL = os.Getenv("DATABASE_URL")

	s := os.Getenv("SESSION_KEY")
	log.Println(s)
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return Config{}, err
	}
	ret.SessionKey = b

	ret.JWTIssuer = os.Getenv("JWT_ISSUER")

	ret.JWTSecret = os.Getenv("JWT_SECRET")

	s = os.Getenv("JWT_LIFESPAN")
	i, err := strconv.Atoi(s)
	if err != nil {
		return Config{}, err
	}
	ret.JWTLifespan = time.Duration(i) * time.Minute

	ret.GoogleAPIClientID = os.Getenv("GOOGLE_API_CLIENT_ID")
	ret.GoogleAPIClientSecret = os.Getenv("GOOGLE_API_CLIENT_SECRET")
	ret.GoogleAPIRedirectURL = os.Getenv("GOOGLE_API_REDIRECT_URL")

	return ret, nil
}
