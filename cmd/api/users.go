package main

import (
	"fmt"
	"net/http"
	"strconv"

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

func (a *application) UsersContextMiddleWare(c *gin.Context) {

	ctx := c.Request.Context()

	userIDstring := c.Param("userID")
	userID, _ := strconv.Atoi(userIDstring)
	user, err := a.store.Users.GetUserByID(ctx, userID)

	if err != nil {
		a.InternalServerError(c, "Cannot Get User", err)
		return
	}
	c.Set("User", user)

}

func GetUserFromContext(ctx *gin.Context) (*models.User, error) {

	usersInterface, exists := ctx.Get("User")

	if !exists {

		return nil, fmt.Errorf("user context not set")

	}

	return usersInterface.(*models.User), nil

}
