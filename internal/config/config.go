package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server   *ServerConfig
	Postgres *PostgresConfig
	Minio    *MinioConfig
	Mailer   *MailerConfig
	Auth     *AuthConfig
}

type ServerConfig struct {
	Port         int
	IdleTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type PostgresConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type MinioConfig struct {
	Endpoint     string
	User         string
	Password     string
	AvatarBucket string
}

type MailerConfig struct {
	SenderEmail    string
	SenderPassword string
	RelayHost      string
	RelayPort      int
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	JWTKey          []byte
}

func GetConfig() Config {
	return Config{
		Server: &ServerConfig{
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 1*time.Minute),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
			Port:         getIntEnv("PORT", 8080),
		},
		Postgres: &PostgresConfig{
			DSN:          getEnv("DB_DSN", "postgres://postgres:postgres@localhost:5432/database?sslmode=disable"),
			MaxOpenConns: getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnv("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  getDurationEnv("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
		Minio: &MinioConfig{
			Endpoint:     getEnv("MINIO_ENDPOINT", "localhost:9000"),
			User:         getEnv("MINIO_USER", "minio"),
			Password:     getEnv("MINIO_PASS", "minio123"),
			AvatarBucket: getEnv("AVATAR_BUCKET", "avatars"),
		},
		Mailer: &MailerConfig{
			SenderEmail:    getEnv("MAILER_SENDER_EMAIL", "Y2b9l@example.com"),
			SenderPassword: getEnv("MAILER_SENDER_PASSWORD", "gymshark"),
			RelayHost:      getEnv("MAILER_RELAY_HOST", "smtp.gmail.com"),
			RelayPort:      getIntEnv("MAILER_RELAY_PORT", 587),
		},
		Auth: &AuthConfig{
			AccessTokenTTL:  getDurationEnv("ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: getDurationEnv("REFRESH_TOKEN_TTL", 24*time.Hour),
			JWTKey:          []byte(getEnv("JWT_KEY", "secret")),
		},
	}
}

func getIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func getEnv(key, defaultValue string) string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	return valueStr
}
