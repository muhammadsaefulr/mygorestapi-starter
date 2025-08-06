package config

import (
	"github.com/muhammadsaefulr/NimeStreamAPI/internal/shared/utils"
	"github.com/spf13/viper"
)

var (
	IsProd               bool
	ClientFeHost         string
	AppHost              string
	AppPort              int
	AppUrl               string
	DBHost               string
	DBUser               string
	DBPassword           string
	DBName               string
	DBPort               int
	RedisHost            string
	RedisPassword        string
	TMDbApiKey           string
	GoogleAppCredentials string
	JWTSecret            string
	JWTAccessExp         int
	JWTRefreshExp        int
	JWTResetPasswordExp  int
	JWTVerifyEmailExp    int
	SMTPHost             string
	SMTPPort             int
	SMTPUsername         string
	SMTPPassword         string
	EmailFrom            string
	GoogleClientID       string
	GoogleClientSecret   string
	RedirectURL          string
)

func init() {
	loadConfig()

	// server configuration
	IsProd = viper.GetString("APP_ENV") == "prod"
	AppHost = viper.GetString("APP_HOST")
	AppPort = viper.GetInt("APP_PORT")
	AppUrl = viper.GetString("APP_URL")
	ClientFeHost = viper.GetString("CLIENT_FE_HOST")

	// database configuration
	DBHost = viper.GetString("DB_HOST")
	DBUser = viper.GetString("DB_USER")
	DBPassword = viper.GetString("DB_PASSWORD")
	DBName = viper.GetString("DB_NAME")
	DBPort = viper.GetInt("DB_PORT")

	// Redis configuration
	RedisHost = viper.GetString("REDIS_HOST")
	RedisPassword = viper.GetString("REDIS_PASSWORD")

	GoogleAppCredentials = viper.GetString("GOOGLE_APPLICATION_CREDENTIALS")

	// apikey configuration
	TMDbApiKey = viper.GetString("TMDB_APIKEY")

	// jwt configuration
	JWTSecret = viper.GetString("JWT_SECRET")
	JWTAccessExp = viper.GetInt("JWT_ACCESS_EXP_MINUTES")
	JWTRefreshExp = viper.GetInt("JWT_REFRESH_EXP_DAYS")
	JWTResetPasswordExp = viper.GetInt("JWT_RESET_PASSWORD_EXP_MINUTES")
	JWTVerifyEmailExp = viper.GetInt("JWT_VERIFY_EMAIL_EXP_MINUTES")

	// SMTP configuration
	SMTPHost = viper.GetString("SMTP_HOST")
	SMTPPort = viper.GetInt("SMTP_PORT")
	SMTPUsername = viper.GetString("SMTP_USERNAME")
	SMTPPassword = viper.GetString("SMTP_PASSWORD")
	EmailFrom = viper.GetString("EMAIL_FROM")

	// oauth2 configuration
	GoogleClientID = viper.GetString("GOOGLE_CLIENT_ID")
	GoogleClientSecret = viper.GetString("GOOGLE_CLIENT_SECRET")
	RedirectURL = viper.GetString("REDIRECT_URL")
}

func loadConfig() {
	viper.SetConfigFile(".env")

	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	utils.Log.Info("Config loaded from environment variables")
}
