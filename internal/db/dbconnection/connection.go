// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package dbconnection

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"githum.com/anupam111/image-store/internal/config"
)

//go:generate mockgen -source ../../config/config.go -package dbconnection -destination connection_mock.go

type Pool struct {
	DB *sqlx.DB
}

type Connector struct{}

func setConnector(dbConfig *config.DBConfig) {
	if dbConfig.Connector == nil {
		dbConfig.Connector = &Connector{}
	}
}

func New(dbConfig *config.DBConfig) *Pool {
	setConnector(dbConfig)
	dbConnection := fmt.Sprintf(
		"host=%s port=%d user=%s password='%s' dbname=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Name,
	)

	connection := dbConfig.Connector.Connect(dbConfig.DriverName, dbConnection)

	return &Pool{
		DB: connection,
	}
}

func (d *Connector) Connect(driverName, dbConnection string) *sqlx.DB {
	return sqlx.MustConnect(driverName, dbConnection)
}
