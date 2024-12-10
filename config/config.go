package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type ConfigFile struct {
	Server struct {
		Host           string
		Port           int
		GinReleaseMode string
	}
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
		Prefix   string
	}
	Valkey struct {
		Host    string
		Port    int
		Channel string
	}
	App struct {
		JSecret             string
		ExtraLog            string
		SaveXml             string
		TkTime              string
		SwaggerBasePath     string
		TokenPrefBackoffice string
		TokenPref3rdParty   string
		TokenCheck          string
	}
	AdminUser struct {
		Username string
		Password string
	}
}

var Configvar ConfigFile

// Load the .env file and initialize the config
func (c *ConfigFile) Load() error {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	// Server configuration
	c.Server.Host = c.getEnv("SERVER_HOST", "127.0.0.1")
	c.Server.Port, err = strconv.Atoi(c.getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return fmt.Errorf("invalid server port: %v", err)
	}
	c.Server.GinReleaseMode = c.getEnv("GIN_RELEASE_MODE", "false")

	// Database configuration
	c.Database.Host = c.getEnv("DB_HOST", "127.0.0.1")
	c.Database.Port, err = strconv.Atoi(c.getEnv("DB_PORT", "5432"))
	if err != nil {
		return fmt.Errorf("invalid database port: %v", err)
	}
	c.Database.User = c.getEnv("DB_USER", "root")
	c.Database.Password = c.getEnv("DB_PASSWORD", "password")
	c.Database.Name = c.getEnv("DB_NAME", "fyc")
	c.Database.SSLMode = c.getEnv("SSLMode", "disable")
	c.Database.Prefix = c.getEnv("Prefix", "fyc")

	// Application configuration
	c.App.JSecret = c.getEnv("JWT_Secret", "0")
	c.App.ExtraLog = c.getEnv("ExtraLog", "false")
	c.App.SaveXml = c.getEnv("SaveXml", "true")
	c.App.TkTime = c.getEnv("ExpireTokenTime", "1")
	c.App.TokenPrefBackoffice = c.getEnv("TokenPrefBackoffice", "false")
	c.App.TokenPref3rdParty = c.getEnv("TokenPref3rdParty", "false")
	c.App.TokenCheck = c.getEnv("TokenCheck", "false")

	// Backoffice Admin user data
	c.AdminUser.Username = c.getEnv("USERNAME", "admin")
	c.AdminUser.Password = c.getEnv("PASSWORD", "admin")

	// Valkey Config
	c.Valkey.Host = c.getEnv("VALKEY_HOST", "127.0.0.1")
	c.Valkey.Port, err = strconv.Atoi(c.getEnv("VALKEY_PORT", "6379"))
	if err != nil {
		return fmt.Errorf("invalid Valkey port: %v", err)
	}
	c.Valkey.Channel = c.getEnv("VALKEY_CHANNEL", "fyc")

	// Swagger BasePath configuration
	BasePath := c.getEnv("SwaggerBasePath", "/")
	if len(BasePath) > 1 {
		// basePath := BasePath

		// Ensure it starts with /
		if !strings.HasPrefix(BasePath, "/") {
			BasePath = "/" + BasePath
		}
		c.App.SwaggerBasePath = strings.TrimSuffix(BasePath, "/")
	}

	return nil
}

// Helper function to get the environment variable with a default value
func (c *ConfigFile) getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
