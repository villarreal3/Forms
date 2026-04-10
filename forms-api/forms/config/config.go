package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	ServerPort      string
	EmailsDir       string // Directorio para almacenar correos
	ImagesDir       string // Directorio para almacenar imágenes
	LoginServiceURL string // URL del servicio de login (legacy, mantenido para compatibilidad)
	UsersAPIURL     string // URL del servicio de usuarios (users-api)
	JWTSecret       string // Secreto JWT compartido - DEBE coincidir con users-api
}

func Load() (*Config, error) {
	// Intentar cargar archivo .env si existe (no es crítico si no existe)
	_ = godotenv.Load()

	cfg := &Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "3306"),
		DBUser:          getEnv("DB_USER", "root"),
		DBPassword:      getEnv("DB_PASSWORD", "daniel3"),
		DBName:          getEnv("DB_NAME", "app_db"),
		ServerPort:      getEnv("SERVER_PORT", "8081"),
		EmailsDir:       getEnv("EMAILS_DIR", "src/emails"),
		ImagesDir:       getEnv("IMAGES_DIR", "src/img"),
		LoginServiceURL: getEnv("LOGIN_SERVICE_URL", "http://localhost:8080"),
		UsersAPIURL:     getEnv("USERS_API_URL", "http://users-api:3000"),
		JWTSecret:       getEnv("JWT_SECRET", "your_jwt_secret_key_change_this_in_production"),
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=Local",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
