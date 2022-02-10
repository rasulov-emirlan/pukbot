package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

const (
	port        = "PORT"
	databaseURL = "DATABASE_URL"
	botToken    = "BOT_TOKEN"

	googleFsClientID     = "GOOGLE_CLIENT_ID"
	googleFsProjectID    = "GOOGLE_PROJECT_ID"
	googleFsAuthURI      = "GOOGLE_AUTH_URI"
	googleFsTokenURI     = "GOOGLE_TOKEN_URI"
	googleFsAuthProvider = "GOOGLE_AUTH_PROVIDER"
	googleFsClientSecret = "GOOGLE_CLIENT_SECRET"
	googleFsRedirectURIs = "GOOGLE_REDIRECT_URI"

	googleFsAccessToken  = "GOOGLE_ACCESS_TOKEN"
	googleFsRefreshToken = "GOOGLE_REFRESH_TOKEN"
	googleFsTokenType    = "GOOGLE_TOKEN_TYPE"
	googleFsExpiry       = "GOOGLE_EXPIRY"
)

type (
	GoogleFS struct {
		Credentials struct {
			ClientID                string
			ProjectID               string
			AuthURI                 string
			TokenURI                string
			AuthProviderx509CertURL string
			ClientSecret            string
			RedirectURIs            []string
		}

		Token struct {
			AccessToken  string    `json:"access_token"`
			RefreshToken string    `json:"refresh_token"`
			TokenType    string    `json:"token_type"`
			Expiry       time.Time `json:"expiry"`
		}
	}

	Config struct {
		ServerPort  string
		DatabaseURL string
		BotToken    string

		GoogleFS GoogleFS
	}
)

func NewConfig(fromFile bool, filenames ...string) (Config, error) {
	var cfg Config

	cfg.ServerPort = os.Getenv(port)
	cfg.DatabaseURL = os.Getenv(databaseURL)
	cfg.BotToken = os.Getenv(botToken)

	cfg.GoogleFS.Credentials.ClientID = os.Getenv(googleFsClientID)
	cfg.GoogleFS.Credentials.ProjectID = os.Getenv(googleFsProjectID)
	cfg.GoogleFS.Credentials.AuthURI = os.Getenv(googleFsAuthURI)
	cfg.GoogleFS.Credentials.TokenURI = os.Getenv(googleFsTokenURI)
	cfg.GoogleFS.Credentials.AuthProviderx509CertURL = os.Getenv(googleFsAuthProvider)
	cfg.GoogleFS.Credentials.ClientSecret = os.Getenv(googleFsClientSecret)
	cfg.GoogleFS.Credentials.RedirectURIs = strings.Split(os.Getenv(googleFsRedirectURIs), " ")

	if cfg.GoogleFS.Credentials.ClientID == "" ||
		cfg.GoogleFS.Credentials.ProjectID == "" ||
		cfg.GoogleFS.Credentials.AuthURI == "" ||
		cfg.GoogleFS.Credentials.TokenURI == "" ||
		cfg.GoogleFS.Credentials.AuthProviderx509CertURL == "" ||
		cfg.GoogleFS.Credentials.ClientSecret == "" ||
		cfg.GoogleFS.Credentials.RedirectURIs == nil {
		return Config{}, errors.New("no google credentials info")
	}

	cfg.GoogleFS.Token.AccessToken = os.Getenv(googleFsAccessToken)
	cfg.GoogleFS.Token.RefreshToken = os.Getenv(googleFsRefreshToken)
	cfg.GoogleFS.Token.TokenType = os.Getenv(googleFsTokenType)

	if cfg.GoogleFS.Token.AccessToken == "" || cfg.GoogleFS.Token.RefreshToken == "" || cfg.GoogleFS.Token.TokenType == "" {
		return Config{}, errors.New("no token info")
	}
	t, err := time.ParseInLocation("2006-01-02T15:04:05.999999999Z07:00", os.Getenv(googleFsExpiry), time.Now().Location())
	cfg.GoogleFS.Token.Expiry = t

	return cfg, err
}
