package config

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt"
)

type ApiConfig struct {
	Url string
}

type DbConfig struct {
	DataSourceName string
}

type TokenConfig struct {
	ApplicationName  string
	JwtSigningMethod *jwt.SigningMethodHMAC
	JwtSignatureKey  []byte
}

type FileConfig struct {
	FilePath string
}

type Config struct {
	ApiConfig
	DbConfig
	TokenConfig
	FileConfig
}

func (c *Config) readConfig() {
	api := os.Getenv("API_URL")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		dbHost,
		dbUser,
		dbPass,
		dbName,
		dbPort)

	filePath := os.Getenv("FILE_PATH")

	c.ApiConfig = ApiConfig{Url: api}
	c.DbConfig = DbConfig{DataSourceName: dsn}

	jwtApp := os.Getenv("JWT_APP")
	jwtSign := os.Getenv("JWT_SIGN")

	c.TokenConfig = TokenConfig{ApplicationName: jwtApp, JwtSigningMethod: jwt.SigningMethodHS256, JwtSignatureKey: []byte(jwtSign)}
	c.FileConfig = FileConfig{FilePath: filePath}
}

func NewConfig() Config {
	cfg := Config{}
	cfg.readConfig()
	return cfg
}
