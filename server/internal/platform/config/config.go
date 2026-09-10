package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
}

func Load(envFile string) (Config, error) {
	if envFile == "" {
		return Config{}, errors.New(".env file path is required")

	}

	values, err := godotenv.Read(envFile)
	if err != nil {
		return Config{}, errors.New("could not read .env file")
	}

	required := []string{
		"PORT",
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_SSLMODE",
	}

	for _, key := range required {
		if values[key] == "" {
			return Config{}, fmt.Errorf("%s is required in .env", key)
		}
	}

	if err := validatePort("PORT", values["PORT"]); err != nil {
		return Config{}, err
	}

	if err := validatePort("POSTGRES_PORT", values["POSTGRES_PORT"]); err != nil {
		return Config{}, err
	}

	// Build the URL safely, including special characters in credentials.
	dbURL := url.URL{
		Scheme: "postgres",
		User: url.UserPassword(
			values["POSTGRES_USER"],
			values["POSTGRES_PASSWORD"],
		),
		Host: net.JoinHostPort(
			values["POSTGRES_HOST"],
			values["POSTGRES_PORT"],
		),
		Path: "/" + values["POSTGRES_DB"],
	}

	query := url.Values{}
	query.Set("sslmode", values["POSTGRES_SSLMODE"])
	dbURL.RawQuery = query.Encode()

	return Config{
		Port:        values["PORT"],
		DatabaseURL: dbURL.String(),
	}, nil
}

func validatePort(name, value string) error {
	port, err := strconv.Atoi(value)
	if err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("%s must be between 1 and 65535", name)
	}
	return nil
}
