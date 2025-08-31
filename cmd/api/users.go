package main

import (
	"fmt"
	"net/http"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
)

func (a *application) GetUserHandler(c *gin.Context) {

	user, err := GetUserFromContext(c)

	if err != nil {
		a.InternalServerError(c, "User Context Not set", err)
		return
	}

	a.Success(c, "Fecthed User Successfully ", *user, http.StatusOK)

}

func GetUserFromContext(ctx *gin.Context) (*models.User, error) {

	usersInterface, exists := ctx.Get("User")

	if !exists {

		return nil, fmt.Errorf("user context not set")

	}

	return usersInterface.(*models.User), nil

}
