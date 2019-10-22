package main

import (
	"fmt"
	"os"
	"strconv"
)

type PostgressConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

func (c PostgressConfig) Dialect() string {
	return "postgres"
}

func (c PostgressConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Dbname)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Dbname)
}

func DefaultPostgressConfig() DatabaseConfig {
	port, err := strconv.Atoi(getEnvVar("PGPORT"))
	if err != nil {
		fmt.Println("Could not parse DB port")
	}
	return PostgressConfig{
		Host:     getEnvVar("PGHOST"),
		Port:     port,
		User:     getEnvVar("PGUSER"),
		Password: getEnvVar("PGPASSWORD"),
		Dbname:   getEnvVar("PGDATABASE"),
	}
}

type DatabaseConfig interface {
	Dialect() string
	ConnectionInfo() string
}

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmacKey"`
	Database DatabaseConfig `json:"-"`
	Mailgun  MailgunConfig  `json:"mailgun"`
}

func DefaultConfig() Config {
	return Config{
		Port:     5000,
		Env:      "dev",
		Pepper:   "mUGD8rTdJe",
		HMACKey:  "the-secret-key",
		Database: DefaultPostgressConfig(),
	}
}

func ProdConfig() Config {
	c := Config{
		Env:     "prod",
		Pepper:  getEnvVar("PASSWORD_PEPPER"),
		HMACKey: getEnvVar("HMAC_KEY"),
	}
	Port, err := strconv.Atoi(getEnvVar("PORT"))
	databaseUrl := getEnvVar("DATABASE_URL")
	if err == nil {
		c.Port = Port
	}
	if databaseUrl != "" {
		c.Database = HerokuPGDatabase{databaseUrl: databaseUrl}
	}
	return c
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func LoadConfig(isProd bool) Config {
	if !isProd {
		fmt.Println("config.json not required, using default config for development")
		return DefaultConfig()
	}

	return ProdConfig()
}

type HerokuPGDatabase struct {
	databaseUrl string
}

func (c HerokuPGDatabase) Dialect() string {
	return "postgres"
}

func (c HerokuPGDatabase) ConnectionInfo() string {
	return c.databaseUrl
}

type MailgunConfig struct {
	APIKey       string `json:"api_key"`
	PublicAPIKEY string `json:"public_api_key_key"`
	Domain       string `json:"domain"`
}

type OAuthConfig struct {
	ID       string `json:"id"`
	Secret   string `json:"secret"`
	AuthURL  string `json:"auth_url"`
	TokenURL string `json:"token_url"`
}

func getEnvVar(key string) string {
	val := os.Getenv(key)
	if val == "" {
		fmt.Fprintf(os.Stderr, "Could not find env %s\n", key)
	}
	return val
}
