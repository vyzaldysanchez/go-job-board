package main

import (
	"encoding/json"
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

func DefaultPostgressConfig() PostgressConfig {
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

type Config struct {
	Port     int             `json:"port"`
	Env      string          `json:"env"`
	Pepper   string          `json:"pepper"`
	HMACKey  string          `json:"hmacKey"`
	Database PostgressConfig `json:"database"`
	Mailgun  MailgunConfig   `json:"mailgun"`
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

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func LoadConfig(configRequired bool) Config {
	if !configRequired {
		fmt.Println("config.json not required, using default config for development")
		return DefaultConfig()
	}
	f, err := os.Open("config.json")
	if err != nil {
		if configRequired {
			panic(err)
		}
		fmt.Println("config.json not found, using default config for development")
		return DefaultConfig()
	}
	var c Config
	decoder := json.NewDecoder(f)

	err = decoder.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Succesfull Loaded config.json")
	return c
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
