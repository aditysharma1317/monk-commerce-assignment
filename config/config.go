package config

import (
	"encoding/json"
	"io"
	"os"
)

var conf *Config

type Config struct {
	Env         string `json:"env"`
	AppName     string `json:"app_name"`
	Port        string `json:"port"`
	DatabaseURL string `json:"database_url"`
	LogLevel    string `json:"log_level"`
}

func ParseJSON(r io.Reader, v any) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}

func Get() *Config {
	return conf
}

func Set(c *Config) {
	conf = c

	if os.Getenv("DATABASE_URL") != "" {
		conf.DatabaseURL = os.Getenv("DATABASE_URL")
	}
}
