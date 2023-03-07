package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"net/http"
	"paopao-ce-teaching/internal/services"
	app "paopao-ce-teaching/pkg/jwt"
)

// Login 用户登录
func Login(c *gin.Context) {
	param := services.AuthRequest{}

	err := c.ShouldBind(&param)
	if err != nil {
		log.Errorf("Login.ShouldBind errs: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "10001",
			"msg":     "入参错误",
			"details": err.Error(),
		})
		return
	}

	user, err := services.DoLogin(c, &param)
	if err != nil {
		log.Errorf("service.DoLogin err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "10000",
			"msg":     "内部错误",
			"details": err.Error(),
		})
		return
	}

	token, err := app.GenerateToken(user)
	if err != nil {
		log.Errorf("app.GenerateToken err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    "10002",
			"msg":     "鉴权失败，Token 生成失败",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
