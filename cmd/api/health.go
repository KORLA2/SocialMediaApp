package main

import (
	"github.com/gin-gonic/gin"
)

func (app *application) HealthCheck(ctx *gin.Context) {

	ctx.String(200, "Healthy end POint OK")

}
