package utils

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Paginator struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Page   int `json:"page"`
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

	return &Paginator{Limit: limit, Offset: offset, Page: page}, nil
}
