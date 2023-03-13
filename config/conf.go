package config

import (
	"github.com/joeshaw/envdecode"
)

type (
	Config struct {
		Environment  string `env:"ENV"`
		InstanceName string `env:"INSTANCE_NAME,default=kata-peya-api"`
		PprofEnabled bool   `env:"PPROF_ENABLED,default=false"`
		Server       struct {
			Port     int    `env:"SERVER_PORT"`
			BasePath string `env:"BASE_PATH,default=api/kata-peya/v1"`
		}
		Otel struct {
			ExporterEndpoint string `env:"OTEL_EXPORTER_JAEGER_ENDPOINT"`
		}
		DB struct {
			Host     string `env:"DB_HOST,default=localhost"`
			Port     int    `env:"DB_PORT,default=3306"`
			User     string `env:"DB_USER,default=root"`
			Password string `env:"DB_PASSWORD,default=mysql"`
			Database string `env:"DB_DATABASE,default=pets"`
		}
	}
)

var cfg Config

func init() {
	if err := envdecode.Decode(&cfg); err != nil {
		panic(err)
	}
}

func Get() Config {
	return cfg
}
