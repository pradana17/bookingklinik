package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Paginator struct {
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	Page       int   `json:"page"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int64 `json:"total_pages"`
}

func Pagination(c *gin.Context) (*Paginator, error) {
	limit := 10 //default limit
	page := 1   //default page

	var err error

	limitStr := c.DefaultQuery("limit", fmt.Sprintf("%d", limit))
	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	pageStr := c.DefaultQuery("page", fmt.Sprintf("%d", page))
	page, err = strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	if offset < 0 {
		offset = 0
	}

	return &Paginator{Limit: limit, Offset: offset, Page: page}, nil
}
