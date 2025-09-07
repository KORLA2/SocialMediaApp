package main

import (
	"fmt"
	"net/http"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
)

// GetUser           godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user/{id} [get]
func (a *application) GetUserHandler(c *gin.Context) {

	followUser, err := GetFollowUserFromContext(c)

	if err != nil {
		a.InternalServerError(c, "Follow User Context Not set", err)
		return
	}

	a.Success(c, "Fecthed User Successfully ", *followUser, http.StatusOK)

}

func GetUserFromContext(ctx *gin.Context) (*models.User, error) {

	usersInterface, exists := ctx.Get("User")

	if !exists {

		return nil, fmt.Errorf("user context not set")

	}

	return usersInterface.(*models.User), nil

}
