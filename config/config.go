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

	JWTIssuer string

	AdminAccessSecret  string
	AdminRefreshSecret string

	UserAccessSecret  string
	UserRefreshSecret string

	DoctorAccessSecret  string
	DoctorRefreshSecret string

	PharmacyManagerAccessSecret  string
	PharmacyManagerRefreshSecret string

	AccessTokenLifespan        time.Duration
	RefreshTokenLifespan       time.Duration
	ResetPasswordTokenLifespan time.Duration

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

	ret.AdminAccessSecret = os.Getenv("ADMIN_ACCESS_SECRET")
	ret.AdminRefreshSecret = os.Getenv("ADMIN_REFRESH_SECRET")

	ret.UserAccessSecret = os.Getenv("USER_ACCESS_SECRET")
	ret.UserRefreshSecret = os.Getenv("USER_REFRESH_SECRET")

	ret.DoctorAccessSecret = os.Getenv("DOCTOR_ACCESS_SECRET")
	ret.DoctorRefreshSecret = os.Getenv("DOCTOR_REFRESH_SECRET")

	ret.PharmacyManagerAccessSecret = os.Getenv("PHARMACY_MANAGER_ACCESS_SECRET")
	ret.PharmacyManagerRefreshSecret = os.Getenv("PHARMACY_MANAGER_REFRESH_SECRET")

	s = os.Getenv("ACCESS_TOKEN_LIFESPAN")
	i, err := strconv.Atoi(s)
	if err != nil {
		return Config{}, err
	}
	ret.AccessTokenLifespan = time.Duration(i) * time.Minute

	s = os.Getenv("REFRESH_TOKEN_LIFESPAN")
	i, err = strconv.Atoi(s)
	if err != nil {
		return Config{}, err
	}
	ret.RefreshTokenLifespan = time.Duration(i) * time.Minute

	s = os.Getenv("RESET_PASSWORD_TOKEN_LIFESPAN")
	i, err = strconv.Atoi(s)
	if err != nil {
		return Config{}, err
	}
	ret.ResetPasswordTokenLifespan = time.Duration(i) * time.Minute

	ret.GoogleAPIClientID = os.Getenv("GOOGLE_API_CLIENT_ID")
	ret.GoogleAPIClientSecret = os.Getenv("GOOGLE_API_CLIENT_SECRET")
	ret.GoogleAPIRedirectURL = os.Getenv("GOOGLE_API_REDIRECT_URL")

	return ret, nil
}
