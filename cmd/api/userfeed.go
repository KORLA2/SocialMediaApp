package main

import (
	"net/http"

	"github.com/KORLA2/SocialMedia/internal/store"
	"github.com/gin-gonic/gin"
)

func (a *application) GetUserFeedHandler(c *gin.Context) {

	pq := store.PaginatedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}
pq,err:=pq.Parse(c);
	if err!=nil{
		a.InternalServerError(c,"Cannot Parse Feed Query",err)
	return;
	}
	
	if err:=validate.Struct(pq);err!=nil{
		a.BadRequest(c,"Cannot Validate Page Query Struct",err)
	return;
	}
	
	ctx := c.Request.Context()

	feed, err := a.store.Posts.Feed(ctx, 1, pq)

	if err != nil {
		a.InternalServerError(c, "Cannot Fetch User Feed", err)
		return
	}

	a.Success(c, "users Feeed", feed, http.StatusOK)

}
