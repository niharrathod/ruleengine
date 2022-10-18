package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/niharrathod/ruleengine/app/config"
	"github.com/niharrathod/ruleengine/app/ext/datastore"
	"github.com/niharrathod/ruleengine/app/handler"
	"github.com/niharrathod/ruleengine/app/log"
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

// initialization activities
// Incase of init activity failure, log the error and exit (os.Exist(1))
func (app *appServer) init() {
	config.Initialize()
	log.Initialize()
	datastore.Initialize()
}

func (app *appServer) Run() {

	if config.IsProduction() {
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

	app.shutdown()
}

/*
To tear down the app. Order of tear down activities is important

 1. http listener - to stop incoming traffic
 2. close datastore connection
    # Add more activities here
    log sync should be last activity
*/
func (app *appServer) shutdown() {
	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Logger.Info("Application is shutting down")

	// stop http listener
	if err := app.httpserver.Shutdown(shutdownContext); err != nil {
		log.Logger.Error("Server Shutdown failed:", zap.String("error", err.Error()))
	}

	// close datastore connection
	datastore.Close(shutdownContext)

	// sync logs
	err := log.Logger.Sync()
	if err != nil {
		// logger sync failed, so back to basics :(
		fmt.Println(err.Error())
	}
}
