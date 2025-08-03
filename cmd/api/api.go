package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type application struct {
	config config
}

type config struct {
	addr string
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
