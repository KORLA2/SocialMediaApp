package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) HealthCheck(ctx *gin.Context) {

	ctx.JSON(200, gin.H{"Healthy": "end POint OK"})

}
