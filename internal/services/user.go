package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"paopao-ce-teaching/internal/conf"
	"paopao-ce-teaching/internal/core/user"
	"paopao-ce-teaching/pkg/util"
	"strings"
)

type AuthRequest struct {
	Username string `json:"username" form:"username" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}

// DoLogin 用户认证
func DoLogin(ctx *gin.Context, param *AuthRequest) (*user.User, error) {
	u, err := user.GetUserByUsername(conf.DB, param.Username)
	if err != nil {
		return nil, fmt.Errorf("账户不存在")
	}

	if u.ID > 0 {
		// 对比密码是否正确
		if ValidPassword(u.Password, param.Password, u.Salt) {

			if u.Status == user.UserStatusClosed {
				return nil, fmt.Errorf("该账户已被封停")
			}

			return u, nil
		}

		return nil, fmt.Errorf("账户密码错误")
	}

	return nil, fmt.Errorf("账户不存在")
}

// ValidPassword 检查密码是否一致
func ValidPassword(dbPassword, password, salt string) bool {
	return strings.Compare(dbPassword, util.EncodeMD5(util.EncodeMD5(password)+salt)) == 0
}
