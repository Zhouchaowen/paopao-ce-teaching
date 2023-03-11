package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"paopao-ce-teaching/internal/middleware"
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

	// 用户注册
	r.POST("/auth/register", api.Register)

	// 获取用户基本信息
	r.GET("/user/profile", api.GetUserProfile)

	// 鉴权路由组
	authApi := r.Group("/").Use(middleware.JWT())
	{
		// 获取当前用户信息
		authApi.GET("/user/info", api.GetUserInfo)

		// 修改密码
		authApi.POST("/user/password", api.ChangeUserPassword)

		// 修改昵称
		authApi.POST("/user/nickname", api.ChangeNickname)
	}

	{
		// 发布动态
		authApi.POST("/post", api.CreatePost)

		// 删除动态
		authApi.DELETE("/post", api.DeletePost)

		// 获取动态详情
		r.GET("/post", api.GetPost)

		// 获取用户动态列表
		authApi.GET("/user/posts", api.GetUserPosts)
	}

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
