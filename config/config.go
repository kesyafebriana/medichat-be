package config

import (
	"encoding/base64"
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
	ServerAddr         string
	WebDomain          string
	FEDomain           string
	DatabaseURL        string
	FEVerificationURL  string
	FEResetPasswordURL string

	AuthEmailUsername string
	AuthEmailPassword string
	EmailSender       string

	SessionKey []byte

	JWTIssuer string

	AdminAccessSecret           string
	UserAccessSecret            string
	DoctorAccessSecret          string
	PharmacyManagerAccessSecret string
	RefreshSecret               string

	AccessTokenLifespan        time.Duration
	RefreshTokenLifespan       time.Duration
	ResetPasswordTokenLifespan time.Duration
	VerifyEmailTokenLifespan   time.Duration

	GoogleAPIClientID     string
	GoogleAPIClientSecret string
	GoogleAPIRedirectURL  string

	IsRelease bool
}

func InitConfig() error {
	return godotenv.Load()
}

func LoadConfig() (Config, error) {
	ret := Config{}

	ret.ServerAddr = os.Getenv("SERVER_ADDR")
	ret.WebDomain = os.Getenv("WEB_DOMAIN")
	ret.FEDomain = os.Getenv("FE_DOMAIN")
	ret.DatabaseURL = os.Getenv("DATABASE_URL")
	ret.FEVerificationURL = os.Getenv("FE_VERIFICATION_URL")
	ret.FEResetPasswordURL = os.Getenv("FE_RESET_PASSWORD_URL")

	ret.AuthEmailUsername = os.Getenv("AUTH_EMAIL_USERNAME")
	ret.AuthEmailPassword = os.Getenv("AUTH_EMAIL_PASSWORD")
	ret.EmailSender = os.Getenv("EMAIL_SENDER")

	s := os.Getenv("SESSION_KEY")
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return Config{}, err
	}
	ret.SessionKey = b

	ret.JWTIssuer = os.Getenv("JWT_ISSUER")

	ret.AdminAccessSecret = os.Getenv("ADMIN_ACCESS_SECRET")
	ret.UserAccessSecret = os.Getenv("USER_ACCESS_SECRET")
	ret.DoctorAccessSecret = os.Getenv("DOCTOR_ACCESS_SECRET")
	ret.PharmacyManagerAccessSecret = os.Getenv("PHARMACY_MANAGER_ACCESS_SECRET")
	ret.RefreshSecret = os.Getenv("REFRESH_SECRET")

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

	s = os.Getenv("VERIFY_EMAIL_TOKEN_LIFESPAN")
	i, err = strconv.Atoi(s)
	if err != nil {
		return Config{}, err
	}
	ret.VerifyEmailTokenLifespan = time.Duration(i) * time.Minute

	ret.GoogleAPIClientID = os.Getenv("GOOGLE_API_CLIENT_ID")
	ret.GoogleAPIClientSecret = os.Getenv("GOOGLE_API_CLIENT_SECRET")
	ret.GoogleAPIRedirectURL = os.Getenv("GOOGLE_API_REDIRECT_URL")

	ret.IsRelease = os.Getenv("MEDICHAT_RELEASE") != ""

	return ret, nil
}
