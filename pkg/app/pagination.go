package app

import (
	"github.com/gin-gonic/gin"
	"paopao-ce-teaching/pkg/convert"
)

func GetPage(c *gin.Context) int {
	page := convert.StrTo(c.Query("page")).MustInt()
	if page <= 0 {
		return 1
	}

	return page
}

func GetPageSize(c *gin.Context) int {
	pageSize := convert.StrTo(c.Query("page_size")).MustInt()
	if pageSize <= 0 {
		return 10
	}
	if pageSize > 10 {
		return 20
	}

	return pageSize
}

func GetPageOffset(c *gin.Context) (offset, limit int) {
	page := convert.StrTo(c.Query("page")).MustInt()
	if page <= 0 {
		page = 1
	}

	limit = convert.StrTo(c.Query("page_size")).MustInt()
	if limit <= 0 {
		limit = 10
	} else if limit > 10 {
		limit = 20
	}
	offset = (page - 1) * limit
	return
}
