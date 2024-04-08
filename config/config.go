package config

import (
	"errors"
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
	JWTIssuer   string
	JWTSecret   string
	JWTLifespan time.Duration
}

func InitConfig() error {
	return godotenv.Load(".env")
}

func LoadConfig() (Config, error) {
	env, err := godotenv.Read(".env")
	if err != nil {
		return Config{}, err
	}

	ret := Config{}

	s, ok := env["SERVER_ADDR"]
	if !ok {
		return Config{}, ErrMissingKey
	}
	ret.ServerAddr = s

	s, ok = env["DATABASE_URL"]
	if !ok {
		return Config{}, ErrMissingKey
	}
	ret.DatabaseURL = s

	s, ok = env["JWT_ISSUER"]
	if !ok {
		return Config{}, ErrMissingKey
	}
	ret.JWTIssuer = s

	s, ok = env["JWT_SECRET"]
	if !ok {
		return Config{}, ErrMissingKey
	}
	ret.JWTSecret = s

	s, ok = env["JWT_LIFESPAN"]
	if !ok {
		return Config{}, ErrMissingKey
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return Config{}, err
	}
	ret.JWTLifespan = time.Duration(i) * time.Minute

	return ret, nil
}
