package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *application) GetUserFeedHandler(c *gin.Context) {

	ctx := c.Request.Context()

	feed, err := a.store.Posts.Feed(ctx, 1)

	if err != nil {
		a.InternalServerError(c, "Cannot Fetch User Feed", err)
		return
	}

	a.Success(c, "users Feeed", feed, http.StatusOK)

}
