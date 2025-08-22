package store

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginatedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" valiadte:"oneof=asc desc"`
	Search string `validate:"max=100"`
	Score  int
}

func (pq PaginatedQuery) Parse(c *gin.Context) (PaginatedQuery, error) {

	limitString := c.Query("limit")
	offsetString := c.Query("offset")
	sort := c.Query("sort")
	search := c.Query("search")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return pq, err

	}
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		return pq, err

	}

	if search == "" {
		log.Print("Serach Value is ", search)

	}
	pq.Limit = limit
	pq.Offset = offset
	pq.Sort = sort
	pq.Search = search

	return pq, nil

}
