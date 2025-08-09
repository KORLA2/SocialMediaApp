package main

import (
	"net/http"
	"time"

	"github.com/KORLA2/SocialMedia/internal/store"
	"github.com/gin-gonic/gin"
)

type application struct {
	config config
	store  *store.Storage
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.GET("/v1/health", app.HealthCheck)

	return router
}

func (app *application) Run(mux http.Handler) error {

	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	return server.ListenAndServe()

}
