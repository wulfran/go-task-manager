package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var fatalf = log.Fatalf

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		fatalf("error while loading env file: %s", err)
	}

	cfg := Config{
		DB: DBConfig{
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
		},
		JWT: JWTConfig{
			Secret: os.Getenv("JWT_SECRET"),
		},
	}

	validate(cfg)

	return cfg
}

func validate(c Config) {
	err := c.DB.Validate()
	if err != nil {
		fatalf("DBConfig validation error: %s", err)
	}

	err = c.JWT.Validate()
	if err != nil {
		fatalf("JWTConfig validation error: %s", err)
	}
}
