package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort            string
	AppEnv             string
	DBPath             string
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	CORSOrigins        []string
	LogLevel           string
	RateLimitEnabled   bool
	SwaggerEnabled     bool
	TLSEnabled         bool
	TLSCertFile        string
	TLSKeyFile         string
}

var App *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		accessExpiry = 15 * time.Minute
	}
	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		refreshExpiry = 7 * 24 * time.Hour
	}

	// Parse CORS origins from comma-separated env var
	corsRaw := getEnv("CORS_ORIGINS", "")
	var corsOrigins []string
	for _, o := range strings.Split(corsRaw, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			corsOrigins = append(corsOrigins, trimmed)
		}
	}
	if len(corsOrigins) == 0 {
		log.Println("WARNING: CORS_ORIGINS not set. Defaulting to allowing all origins (not recommended for production)")
		corsOrigins = []string{"*"}
	}

	tlsEnabled := getEnvBool("TLS_ENABLED", false)
	tlsCertFile := getEnv("TLS_CERT_FILE", "")
	tlsKeyFile := getEnv("TLS_KEY_FILE", "")

	if tlsEnabled && (tlsCertFile == "" || tlsKeyFile == "") {
		log.Fatal("TLS_ENABLED requires TLS_CERT_FILE and TLS_KEY_FILE to be set")
	}

	jwtSecret := getEnv("JWT_SECRET", "")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required. Set it in environment variables or .env file")
	}
	if len(jwtSecret) < 32 {
		log.Fatal("JWT_SECRET must be at least 32 characters long for secure signing")
	}

	App = &Config{
		AppPort:            getEnv("APP_PORT", "8080"),
		AppEnv:             getEnv("APP_ENV", "development"),
		DBPath:             getEnv("DB_PATH", "./data.db"),
		JWTSecret:          jwtSecret,
		AccessTokenExpiry:  accessExpiry,
		RefreshTokenExpiry: refreshExpiry,
		CORSOrigins:        corsOrigins,
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		RateLimitEnabled:   getEnvBool("RATE_LIMIT_ENABLED", true),
		SwaggerEnabled:     getEnvBool("SWAGGER_ENABLED", false),
		TLSEnabled:         tlsEnabled,
		TLSCertFile:        tlsCertFile,
		TLSKeyFile:         tlsKeyFile,
	}

	// Configure global slog
	var logLevel slog.Level
	switch strings.ToLower(App.LogLevel) {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
