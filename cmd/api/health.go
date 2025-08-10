package main

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func (app *application) HealthCheck(ctx *gin.Context) {

	// ctx.JSON(200, "Healthy end POint OK")


	ctx.Error(errors.New("new error"))
}
