package service

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	accessTTL  time.Duration
	refreshTTL time.Duration
	key        []byte
}

func ParseENV(log *slog.Logger) Config {
	cfg := Config{}
	cfg.accessTTL = 120 * time.Minute
	cfg.refreshTTL = 30 * time.Hour * 24
	cfg.key = []byte("secret")

	tmp, exists := os.LookupEnv("ACCESS_TTL")
	if exists {
		accessTTL, err := strconv.Atoi(tmp)
		if err != nil {
			log.Error("Cannot parse ACCESS_TTL to int")
		} else {
			cfg.accessTTL = time.Duration(accessTTL) * time.Minute
		}
	} else {
		log.Error("Cannot find environment variable ACCESS_TTL")
	}

	tmp, exists = os.LookupEnv("REFRESH_TTL")
	if exists {
		refreshTTL, err := strconv.Atoi(tmp)
		if err != nil {
			log.Error("Cannot parse REFRESH_TTL to int")
		} else {
			cfg.accessTTL = time.Duration(refreshTTL) * time.Minute
		}
	} else {
		log.Error("Cannot find environment variable REFRESH_TTL")
	}

	tmp, exists = os.LookupEnv("SECRET_KEY")
	if exists {
		cfg.key = []byte(tmp)
	} else {
		log.Error("Cannot find environment variable SECRET_KEY")
	}

	return cfg
}
