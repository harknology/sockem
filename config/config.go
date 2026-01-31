package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var (
	LOG_FORMAT    string
	SECRET_KEY    string
	CLIENT_KEY    string
	ALLOWED_HOSTS []string
	BUFFER_SIZE   int
	PORT          uint64
	HOST          string
)

func init() {
	var err error

	SECRET_KEY = os.Getenv("SOCKEM_SECRET_KEY")
	LOG_FORMAT = os.Getenv("SOCKEM_LOG_FORMAT")
	CLIENT_KEY = os.Getenv("SOCKEM_CLIENT_KEY")
	ALLOWED_HOSTS = strings.Split(os.Getenv("SOCKEM_ALLOWED_HOSTS"), " ")

	bufSize, found := os.LookupEnv("SOCKEM_BUFFER_SIZE")
	if found {
		BUFFER_SIZE, err = strconv.Atoi(bufSize)
		if err != nil {
			slog.Error("convert", "variable", "SOCKEM_BUFFER_SIZE", "error", err)
		}
	} else {
		BUFFER_SIZE = 1024
	}

	port, found := os.LookupEnv("SOCKEM_PORT")
	if found {
		PORT, err = strconv.ParseUint(port, 10, 0)
		if err != nil {
			slog.Error("convert", "variable", "SOCKEM_BUFFER_SIZE", "error", err)
		}
	} else {
		PORT = 3005
	}

	host, found := os.LookupEnv("SOCKEM_HOST")
	if found {
		HOST = host
	} else {
		HOST = "0.0.0.0"
	}
}
