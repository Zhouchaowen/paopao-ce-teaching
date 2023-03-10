package app

import (
	"net/http"
	"os"
	"paopao-ce-teaching/pkg/errors"

	"github.com/gin-gonic/gin"
)

func ToResponse(ctx *gin.Context, data interface{}) {
	hostname, _ := os.Hostname()
	if data == nil {
		data = gin.H{
			"code":      0,
			"msg":       "success",
			"tracehost": hostname,
		}
	} else {
		data = gin.H{
			"code":      0,
			"msg":       "success",
			"data":      data,
			"tracehost": hostname,
		}
	}
	ctx.JSON(http.StatusOK, data)
}

type Pager struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalRows int64 `json:"total_rows"`
}

func ToResponseList(ctx *gin.Context, list interface{}, totalRows int64) {
	ToResponse(ctx, gin.H{
		"list": list,
		"pager": Pager{
			Page:      GetPage(ctx),
			PageSize:  GetPageSize(ctx),
			TotalRows: totalRows,
		},
	})
}

func ToErrorResponse(ctx *gin.Context, err *errors.Error) {
	response := gin.H{"code": err.Code(), "msg": err.Msg()}
	details := err.Details()
	if len(details) > 0 {
		response["details"] = details
	}

	ctx.JSON(err.StatusCode(), response)
}
