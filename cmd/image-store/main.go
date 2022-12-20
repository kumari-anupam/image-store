package main

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"githum.com/anupam111/image-store/internal/config"
	"githum.com/anupam111/image-store/internal/server"
	"os"
	"strconv"
)

func main() {
	// It loads env data.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error occured while reading env data : %v", err)
	}

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("error occured while converting string to int: %v", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("error occured while converting string to int: %v", err)
	}

	serverConfig := config.ServiceConfig{
		LogLevel: os.Getenv("LOG_LEVEL"),
		Port:     port,
	}

	dbConfig := config.DBConfig{
		Password:   os.Getenv("DB_PASSWORD"),
		Host:       os.Getenv("DB_HOST"),
		Name:       os.Getenv("DB_NAME"),
		Username:   os.Getenv("DB_USERNAME"),
		DriverName: os.Getenv("DB_DRIVER"),
		Port:       dbPort,
	}
	server := server.NewAppServer()

	imageStoreServiceConfig := &config.ImageStoreServiceConfig{
		DBConfig:      dbConfig,
		ServiceConfig: serverConfig,
	}

	fmt.Printf("%+v\n", imageStoreServiceConfig)
	server.ConfigureAndStart(imageStoreServiceConfig)
}
