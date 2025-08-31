package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (a *application) FollowUserHandler(c *gin.Context) {

	ctx := c.Request.Context()
	followerUser, _ := strconv.Atoi(c.Param("userID"))
	User, err := GetUserFromContext(c)

	if err != nil {
		a.InternalServerError(c, "Follower User Context Not FOund", err)
		return
	}

	if User.ID == followerUser {
		a.BadRequest(c, "Cannot Do this", fmt.Errorf("You cannot follow yoursef"))
		return
	}

_, err = a.store.Users.GetUserByID(ctx, followerUser)
	if err != nil {
		a.Unauthorized(c, "token is invalid", err)
		return
	}

	if err := a.store.Followers.Create(ctx, followerUser, User.ID); err != nil {
		a.InternalServerError(c, "Cannot Create a Follow Request", err)
		return
	}
	a.Success(c, "Followed User", followerUser, http.StatusOK)

}
func (a *application) UnfollowUserHandler(c *gin.Context) {

	ctx := c.Request.Context()
	followerUser, _ := strconv.Atoi(c.Param("userID"))
	User, err := GetUserFromContext(c)


	if err != nil {
		a.InternalServerError(c, "Follower User Context Not Found", err)
		return
	}
	_, err = a.store.Users.GetUserByID(ctx, followerUser)
	if err != nil {
		a.Unauthorized(c, "token is invalid", err)
		return
	}


	if err := a.store.Followers.Delete(ctx, followerUser, User.ID); err != nil {
		a.InternalServerError(c, "Cannot Unfollow", err)
		return
	}
	a.Success(c, "UnFollowed User", followerUser, http.StatusOK)

}
