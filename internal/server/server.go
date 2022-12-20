package server

import (
	"context"
	"errors"
	"fmt"
	"githum.com/anupam111/image-store/internal/apihandler"
	"githum.com/anupam111/image-store/internal/config"
	"githum.com/anupam111/image-store/internal/controller"
	"githum.com/anupam111/image-store/internal/db/dbconnection"
	"githum.com/anupam111/image-store/internal/db/dbhandler"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// AppServer ...
type AppServer struct {
	log    *log.Logger
	router *gin.Engine
	server *http.Server
}

// NewAppServer implements AppServer.
func NewAppServer() *AppServer {
	return &AppServer{}
}

// ConfigureAndStart Configures AppServerBase.
func (app *AppServer) ConfigureAndStart(config *config.ImageStoreServiceConfig) {
	app.log = configureLogger(config.ServiceConfig.LogLevel)
	app.router = gin.New()
	gin.EnableJsonDecoderDisallowUnknownFields()
	app.router.Use(gin.Recovery())
	app.router.HandleMethodNotAllowed = true

	base := app.router.Group("")
	if config.ServiceConfig.GinAccessLog {
		base.Use(gin.Logger())
	}

	app.setupRouter(config.DBConfig)
	app.Start(config.ServiceConfig)
}

func (app *AppServer) setupRouter(dbConfig config.DBConfig) {
	logger := log.New()
	v1router := app.router.Group("/v1")
	dbConnection := dbconnection.New(&dbConfig)
	dbHandler := dbhandler.NewDBHandler(logger, dbConnection)
	controller := controller.NewImageController(logger, dbHandler)
	handler := apihandler.NewAPIHandler(logger, controller)

	v1router.POST("/album", handler.CreateImageAlbum)
	v1router.POST("/album/images", handler.CreateImage)
	v1router.DELETE("/album/:albumName", handler.DeleteImageAlbum)
	v1router.DELETE("/album/images/:imageName", handler.DeleteImage)
	v1router.GET("/album/images/:imageName", handler.GetImageByID)
	v1router.GET("/album/images", handler.GetAlbumImages)
}

// Start starts the Server for real.
func (app *AppServer) Start(conf config.ServiceConfig) {
	log.Info("Starting image-store server...")
	app.startGinServer(conf)
	log.Info("image-store server started successfully ...")

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("Shutting down image-store server...")
	app.StopServer()
}

func (app *AppServer) startGinServer(conf config.ServiceConfig) {
	app.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: app.router,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := app.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("listen: %s\n", err)
			os.Exit(1)
		}
	}()
}

func (app *AppServer) StopServer() {
	// The context is used to inform the server it has 20 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		log.Errorf("Server Shutdown: %v", err)
	}

	log.Info("Server stopped successfully")
}

func configureLogger(logLevel string) *log.Logger {
	log.Infof("Setting global log level to '%s'", logLevel)
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Error("error while parsing the log level")
	}

	return &log.Logger{
		Level: level,
	}
}
