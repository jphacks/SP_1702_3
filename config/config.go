package config

import "github.com/kelseyhightower/envconfig"

type Configuration struct {
	DBUser   string `envconfig:"DB_USER" default:"otomo"`
	DBHost   string `envconfig:"DB_HOST" default:"localhost"`
	DBPort   string `envconfig:"DB_PORT" default:"5432"`
	DBName   string `envconfig:"DB_NAME" default:"otomo"`
	DBPasswd string `envconfig:"DB_PASSWD" default:"secret"`
}

var (
	config Configuration
)

const (
	prefix = "TBS"
)

func init() {
	envconfig.MustProcess(prefix, &config)
}

func reload() {
	envconfig.Process(prefix, &config)
}

func DBUser() string {
	return config.DBUser
}

func DBHost() string {
	return config.DBHost
}

func DBPort() string {
	return config.DBPort
}

func DBName() string {
	return config.DBName
}

func DBPasswd() string {
	return config.DBPasswd
}
