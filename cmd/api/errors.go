package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *application) InternalServerError(c *gin.Context, errormessage string, err error) {
	c.AbortWithStatusJSON(500, gin.H{errormessage: err.Error()})

}
func (a *application) Success(c *gin.Context, Successmessage string, success any, Status int) {
	c.JSON(Status, gin.H{Successmessage: success})

}
func (a *application) BadRequest(c *gin.Context, errormessage string, err error) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{errormessage: err.Error()})

}
