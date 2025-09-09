package main

import (
	"github.com/gin-gonic/gin"
)

// Health Check godoc
//
//	@Summary		Health Check
//	@Description	Checks the health status of the API
//	@Tags			Health Check
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	"Healthy"
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/health [put]

func (app *application) HealthCheck(ctx *gin.Context) {

	ctx.JSON(200, gin.H{"Healthy": "end POint OK"})

}
