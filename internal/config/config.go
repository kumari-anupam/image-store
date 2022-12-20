package config

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
)

// DBConnector ...
type DBConnector interface {
	Connect(driverName, dbConnection string) *sqlx.DB
}

// ImageStoreServiceConfig Configuration specific to Image store Service.
type ImageStoreServiceConfig struct {
	DBConfig      DBConfig
	ServiceConfig ServiceConfig
}

//ServiceConfig ...
type ServiceConfig struct {
	LogLevel     string `envconfig:"LOG_LEVEL"`
	Port         int    `envconfig:"PORT" default:"27006"`
	GinAccessLog bool   `envconfig:"GIN_ACCESS_LOG" default:"true"`
}

// DBConfig represents database configurations.
type DBConfig struct {
	Connector  DBConnector
	Password   string `envconfig:"DB_PASSWORD"`
	Host       string `envconfig:"DB_HOST" default:"localhost"`
	Name       string `envconfig:"DB_NAME"`
	Username   string `envconfig:"DB_USERNAME"`
	DriverName string `envconfig:"DB_DRIVER" default:"postgres"`
	Port       int    `envconfig:"DB_PORT" default:"5432"`
}

// GeImageStoreConfig Provides image-store service related all configurations.
func GeImageStoreConfig() (*ImageStoreServiceConfig, error) {
	var serviceConfig ServiceConfig
	if err := envconfig.Process("", &serviceConfig); err != nil {
		return nil, fmt.Errorf("error while reading db config, %w", err)
	}

	var dbConfig DBConfig
	if err := envconfig.Process("", &dbConfig); err != nil {
		return nil, fmt.Errorf("error while reading db config, %w", err)
	}

	return &ImageStoreServiceConfig{
		ServiceConfig: serviceConfig,
		DBConfig:      dbConfig,
	}, nil
}
