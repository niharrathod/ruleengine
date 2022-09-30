package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/niharrathod/ruleengine/app/handler"
	"github.com/niharrathod/ruleengine/app/internal/config"
	"github.com/niharrathod/ruleengine/app/internal/log"
	"go.uber.org/zap"

	ginzap "github.com/gin-contrib/zap"
)

type appServer struct {
	httpserver *http.Server
}

func New() *appServer {
	appServer := appServer{}
	appServer.init()
	return &appServer
}

func (app *appServer) init() {
	config.Initialize()
	log.Initialize()
	// add initiation activities here
}

func (app *appServer) Run() {

	if config.EnvironmentMode == config.ProductionMode {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(log.Logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(log.Logger, true))

	rest := router.Group("health")
	rest.GET("/check/", handler.HealthCheck())

	// add rest api and handler mapping here

	app.httpserver = &http.Server{
		Addr:    config.Server.Http.BindIp + ":" + strconv.Itoa(config.Server.Http.BindPort),
		Handler: router,
	}

	go app.startServer()
	app.prepAndWaitForShutDown()
}

func (app *appServer) startServer() {
	log.Logger.Info("Application is starting on " + config.Server.Http.BindIp + ":" + strconv.Itoa(config.Server.Http.BindPort))
	err := app.httpserver.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Logger.Panic("http server listen failed : " + err.Error())
	}
}

func (app *appServer) prepAndWaitForShutDown() {
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-quitSignal

	log.Logger.Info("Application is shutting down")

	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app.shutdown(shutdownContext)
}

func (app *appServer) shutdown(ctx context.Context) {

	// add context aware shutdown activities here

	if err := app.httpserver.Shutdown(ctx); err != nil {
		log.Logger.Error("Server Shutdown failed:", zap.String("error", err.Error()))
	}
	log.Logger.Sync()
}
