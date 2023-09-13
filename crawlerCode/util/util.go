package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Pagination is page util
func Pagination(ctx *gin.Context) (pageStr string, num int, err error) {
	limit := ctx.DefaultQuery("page_size", "8")
	pageNumber := ctx.DefaultQuery("page_number", "1")
	limitInt, err := strconv.Atoi(limit)

	if err != nil || limitInt < 0 {
		return "", 0, err

	}
	pageNumberInt, err := strconv.Atoi(pageNumber)
	if err != nil || pageNumberInt < 0 {
		return "", 0, err
	}
	if pageNumberInt != 0 {
		pageNumberInt--
	}
	offsetInt := limitInt * pageNumberInt
	pageStr = fmt.Sprintf(" limit %d offset %d", limitInt+1, offsetInt)
	return pageStr, limitInt, nil
}
