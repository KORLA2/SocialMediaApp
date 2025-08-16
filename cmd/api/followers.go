package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type FollowPayload struct {
	UserID int `json:"userid"`
}

func (a *application) FollowUserHandler(c *gin.Context) {

	ctx := c.Request.Context()
	followerUser, err := GetUserFromContext(c)

	if err != nil {
		a.InternalServerError(c, "Follower User Context Not FOund", err)
		return
	}
	var payload FollowPayload

	if err := c.BindJSON(&payload); err != nil {
		a.BadRequest(c, "Cannot Bind Followers json", err)
		return
	}
	if err := a.store.Followers.Create(ctx, followerUser.ID, payload.UserID); err != nil {
		a.InternalServerError(c, "Cannot Create a Follow Request", err)
		return
	}
	a.Success(c, "Followed User", followerUser.ID, http.StatusOK)

}
func (a *application) UnfollowUserHandler(c *gin.Context) {

	ctx := c.Request.Context()
	followerUser, err := GetUserFromContext(c)

	if err != nil {
		a.InternalServerError(c, "Follower User Context Not Found", err)
		return
	}
	var payload FollowPayload

	if err := c.BindJSON(&payload); err != nil {
		a.BadRequest(c, "Cannot Bind Followers json", err)
		return
	}
	if err := a.store.Followers.Delete(ctx, followerUser.ID, payload.UserID); err != nil {
		a.InternalServerError(c, "Cannot Unfollow", err)
		return
	}
	a.Success(c, "UnFollowed User", followerUser.ID, http.StatusOK)

}
