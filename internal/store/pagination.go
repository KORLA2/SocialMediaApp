package store

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginatedQuery struct {
	Limit  int    `json:"limit" validate:"gte=1,lte=20"`
	Offset int    `json:"offset" validate:"gte=0"`
	Sort   string `json:"sort" valiadte:"oneof=asc desc"`
}

func (pq PaginatedQuery) Parse(c *gin.Context) (PaginatedQuery, error) {

	limitString := c.Query("limit")
	offsetString := c.Query("offset")
	sort := c.Query("sort")
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		return pq, err

	}
	offset, err := strconv.Atoi(offsetString)
	if err != nil {
		return pq, err

	}

	pq.Limit = limit
	pq.Offset = offset
	pq.Sort = sort

			return pq, nil;


}
