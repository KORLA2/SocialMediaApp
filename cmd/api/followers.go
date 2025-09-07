package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/KORLA2/SocialMedia/models"
	"github.com/gin-gonic/gin"
)

// FollowRequest    godoc
//
//	@Summary		Follows a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.User.ID
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user/{id}/follow [put]
func (a *application) FollowUserHandler(c *gin.Context) {

	ctx := c.Request.Context()
	User, err := GetUserFromContext(c)
	if err != nil {
		a.InternalServerError(c, "User Context Not Found", err)
	}
	followUser, err := GetFollowUserFromContext(c)
	if err != nil {
		a.InternalServerError(c, "Follow User Context Not Found", err)
	}

	if User.ID == followUser.ID {
		a.BadRequest(c, "Cannot Do this", fmt.Errorf("you cannot follow yoursef"))
		return
	}

	if err := a.store.Followers.Create(ctx, followUser.ID, User.ID); err != nil {
		a.InternalServerError(c, "Cannot Create a Follow Request", err)
		return
	}
	a.Success(c, "Followed User", followUser.ID, http.StatusOK)

}

// UnFollowRequest    godoc
//
//	@Summary		UnFollows a user
//	@Description	UnFollows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.User.ID
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user/{id}/unfollow [put]
func (a *application) UnfollowUserHandler(c *gin.Context) {

	ctx := c.Request.Context()
	User, err := GetUserFromContext(c)
	if err != nil {
		a.InternalServerError(c, "User Context Not Found", err)
	}
	followUser, err := GetFollowUserFromContext(c)
	if err != nil {
		a.InternalServerError(c, "Follow User Context Not Found", err)
	}

	if err := a.store.Followers.Delete(ctx, followUser.ID, User.ID); err != nil {
		a.InternalServerError(c, "Cannot Unfollow", err)
		return
	}
	a.Success(c, "UnFollowed User", followUser.ID, http.StatusOK)

}

func (a *application) FollowUserContextMiddleware(c *gin.Context) {
	ctx := c.Request.Context()
	followUser, _ := strconv.Atoi(c.Param("userID"))

	user, err := a.store.Users.GetUserByID(ctx, followUser)

	if err != nil {
		a.Unauthorized(c, "token is invalid", err)
		return
	}
	c.Set("FollowUser", user)
}

func GetFollowUserFromContext(c *gin.Context) (*models.User, error) {
	FollowUserInterface, exists := c.Get("FollowUser")
	if !exists {

		return nil, fmt.Errorf("follow user context not set")

	}

	return FollowUserInterface.(*models.User), nil
}
