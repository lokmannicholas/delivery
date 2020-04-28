package config

import (
	"os"
)

var conf *Config

type Config struct {
	Port      string       `json:"port"`
	MapApiKey string       `json:"map-api-key"`
	Mysql     *MysqlConfig `json:"mysql"`
}

type MysqlConfig struct {
	Addr     string `json:"addr"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       string `json:"db"`
}

func Get() *Config {

	if conf == nil {
		loadConfig()
	}
	return conf
}

const MYSQL_ADDR = "MYSQL_ADDR"
const MYSQL_USER = "MYSQL_USER"
const MYSQL_PASSWORD = "MYSQL_PASSWORD"
const MYSQL_DB = "MYSQL_DB"
const PORT = "PORT"
const MAP_API_KEY = "MAP_API_KEY"

func loadConfig() {
	port := "8080"
	if len(os.Getenv(PORT)) > 0 {
		port = os.Getenv(PORT)
	}
	conf = &Config{
		Port:      port,
		MapApiKey: os.Getenv(MAP_API_KEY),
		Mysql: &MysqlConfig{
			Addr:     os.Getenv(MYSQL_ADDR),
			User:     os.Getenv(MYSQL_USER),
			Password: os.Getenv(MYSQL_PASSWORD),
			DB:       os.Getenv(MYSQL_DB),
		},
	}
}
