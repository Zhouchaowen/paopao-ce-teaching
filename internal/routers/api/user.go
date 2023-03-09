package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"net/http"
	"paopao-ce-teaching/internal/services"
	"paopao-ce-teaching/pkg/app"
	"paopao-ce-teaching/pkg/errors"
	jwt "paopao-ce-teaching/pkg/jwt"
	"unicode/utf8"
)

// Login 用户登录
func Login(c *gin.Context) {
	param := services.AuthRequest{}

	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		log.Errorf("Login.ShouldBind errs: %v", errs)
		app.ToErrorResponse(c, errors.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	user, err := services.DoLogin(&param)
	if err != nil {
		log.Errorf("service.DoLogin err: %v", err)
		app.ToErrorResponse(c, err.(*errors.Error))
		return
	}

	token, err := jwt.GenerateToken(user)
	if err != nil {
		log.Errorf("app.GenerateToken err: %v", err)
		app.ToErrorResponse(c, errors.UnauthorizedAuthFailed)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Register 用户注册
func Register(c *gin.Context) {
	param := services.RegisterRequest{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		log.Errorf("Register.ShouldBind errs: %v", errs)
		app.ToErrorResponse(c, errors.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	// 用户名检查
	err := services.ValidUsername(param.Username)
	if err != nil {
		log.Errorf("service.ValidUsername errs: %v", err)
		app.ToErrorResponse(c, err.(*errors.Error))
		return
	}

	// 密码检查
	err = services.CheckPassword(param.Password)
	if err != nil {
		log.Errorf("service.CheckPassword errs: %v", err)
		app.ToErrorResponse(c, err.(*errors.Error))
		return
	}

	user, err := services.Register(
		param.Username,
		param.Password,
	)

	if err != nil {
		log.Errorf("service.Register errs: %v", err)
		app.ToErrorResponse(c, errors.UserRegisterFailed)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}

// GetUserProfile 获取用户基本信息
func GetUserProfile(c *gin.Context) {
	username := c.Query("username")

	user, err := services.GetUserByUsername(username)
	if err != nil {
		log.Errorf("service.GetUserByUsername err: %v\n", err)
		app.ToErrorResponse(c, errors.NoExistUsername)
		return
	}

	app.ToResponse(c, gin.H{
		"id":       user.ID,
		"nickname": user.Nickname,
		"username": user.Username,
		"status":   user.Status,
		"avatar":   user.Avatar,
		"is_admin": user.IsAdmin,
	})
}

// GetUserInfo 获取用户基本信息
func GetUserInfo(c *gin.Context) {
	username, exists := c.Get("USERNAME")
	if !exists {
		app.ToErrorResponse(c, errors.UnauthorizedAuthNotExist)
		return
	}

	user, err := services.GetUserByUsername(username.(string))
	if err != nil {
		app.ToErrorResponse(c, errors.UnauthorizedAuthNotExist)
		return
	}

	phone := ""
	if user.Phone != "" && len(user.Phone) == 11 {
		phone = user.Phone[0:3] + "****" + user.Phone[7:]
	}

	app.ToResponse(c, gin.H{
		"id":       user.ID,
		"nickname": user.Nickname,
		"username": user.Username,
		"status":   user.Status,
		"avatar":   user.Avatar,
		"balance":  user.Balance,
		"phone":    phone,
		"is_admin": user.IsAdmin,
	})
}

// ChangeUserPassword 修改密码
func ChangeUserPassword(c *gin.Context) {
	param := services.ChangePasswordReq{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		log.Errorf("app.BindAndValid errs: %v", errs)
		app.ToErrorResponse(c, errors.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	username, exists := c.Get("USERNAME")
	if !exists {
		app.ToErrorResponse(c, errors.UnauthorizedAuthNotExist)
		return
	}

	user, err := services.GetUserByUsername(username.(string))
	if err != nil {
		app.ToErrorResponse(c, errors.UnauthorizedAuthNotExist)
		return
	}

	// 密码检查
	err = services.CheckPassword(param.Password)
	if err != nil {
		log.Errorf("service.Register err: %v", err)
		app.ToErrorResponse(c, err.(*errors.Error))
		return
	}

	// 旧密码校验
	if !services.ValidPassword(user.Password, param.OldPassword, user.Salt) {
		app.ToErrorResponse(c, errors.ErrorOldPassword)
		return
	}

	// 更新入库
	user.Password, user.Salt = services.EncryptPasswordAndSalt(param.Password)
	services.UpdateUserInfo(user)

	app.ToResponse(c, nil)
}

// ChangeNickname 修改昵称
func ChangeNickname(c *gin.Context) {
	param := services.ChangeNicknameReq{}
	valid, errs := app.BindAndValid(c, &param)
	if !valid {
		log.Errorf("app.BindAndValid errs: %v", errs)
		app.ToErrorResponse(c, errors.InvalidParams.WithDetails(errs.Errors()...))
		return
	}

	if utf8.RuneCountInString(param.Nickname) < 2 || utf8.RuneCountInString(param.Nickname) > 12 {
		app.ToErrorResponse(c, errors.NicknameLengthLimit)
		return
	}

	username, exists := c.Get("USERNAME")
	if !exists {
		app.ToErrorResponse(c, errors.UnauthorizedAuthNotExist)
		return
	}

	user, err := services.GetUserByUsername(username.(string))
	if err != nil {
		app.ToErrorResponse(c, errors.UnauthorizedAuthNotExist)
		return
	}

	// 执行绑定
	user.Nickname = param.Nickname
	services.UpdateUserInfo(user)

	app.ToResponse(c, nil)
}
