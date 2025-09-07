package main

import (
	"log"
	"net/http"

	"github.com/KORLA2/SocialMedia/internal/store"
	"github.com/gin-gonic/gin"
)

// GetUserFeed godoc
//
//	@Summary		Fetches a user feed
//	@Description	Fetches a user feed based on the users they follow
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Number of items to return"
//	@Param			offset	query		int		false	"Offset for pagination"
//	@Param			sort	query		string	false	"Sort order (asc/desc)"
//	@Param			search	query		string	false	"Search term"
//	@Success		200		{array}		models.UserFeed
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user/feed [get]
func (a *application) GetUserFeedHandler(c *gin.Context) {

	pq := store.PaginatedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
		Search: "",
	}
	pq, err := pq.Parse(c)
	if err != nil {
		a.InternalServerError(c, "Cannot Parse Feed Query", err)
		return
	}

	if err := validate.Struct(pq); err != nil {
		a.BadRequest(c, "Cannot Validate Page Query Struct", err)
		return
	}

	log.Print(pq)
	ctx := c.Request.Context()

	feed, err := a.store.Posts.Feed(ctx, 1, pq)

	if err != nil {
		a.InternalServerError(c, "Cannot Fetch User Feed", err)
		return
	}

	a.Success(c, "users Feeed", feed, http.StatusOK)

}
