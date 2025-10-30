package config

import (
	"github.com/kelseyhightower/envconfig"
)

type MainConfig struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"tjTestProject"`
	ServicePort int    `envconfig:"HTTP_PORT" default:"8000"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`

	PostgresHost     string `envconfig:"PGSQL_HOST" default:"localhost"`
	PostgresPort     string `envconfig:"PGSQL_PORT" default:"5432"`
	DBType           string `envconfig:"DB_TYPE" default:"postgres"`
	PostgresUsername string `envconfig:"PGSQL_USERNAME" default:"admin"`
	PostgresPassword string `envconfig:"PGSQL_PASSWORD" default:"admin"`
	DBName           string `envconfig:"PGSQL_DBNAME" default:"transjkt"`
	LogMode          bool   `envconfig:"LOG_MODE" default:"true"`
	MaxIdleConns     int    `envconfig:"MAX_IDLE_CONNS" default:"10"`
	MaxOpenConns     int    `envconfig:"MAX_OPEN_CONNS" default:"10"`
	ConnMaxLifetime  int    `envconfig:"CONN_MAX_LIFETIME" default:"60"`

	MQTTUrl     string `envconfig:"MQTT_URL" default:"tcp://localhost:1883"`
	RabbitMQUrl string `envconfig:"RABBIT_URL" default:"amqp://admin:admin123@localhost:5672/"`

	SignatureKey string `envconfig:"JWT_SECRET_KEY" default:"secret"`
}

func Get() *MainConfig {
	var c MainConfig
	envconfig.MustProcess("", &c)
	return &c
}
