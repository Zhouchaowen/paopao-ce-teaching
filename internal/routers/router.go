package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"paopao-ce-teaching/internal/routers/api"
)

func NewRouter() *gin.Engine {
	e := gin.Default()
	e.HandleMethodNotAllowed = true

	r := e.Group("/v1")

	// 获取version
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": "v0.0.1",
		})
	})

	// 用户登录
	r.POST("/auth/login", api.Login)

	// 默认404
	e.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "Not Found",
		})
	})

	// 默认405
	e.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"code": 405,
			"msg":  "Method Not Allowed",
		})
	})

	return e
}
