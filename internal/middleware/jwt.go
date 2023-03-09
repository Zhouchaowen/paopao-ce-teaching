package middleware

import (
	"paopao-ce-teaching/pkg/app"
	"paopao-ce-teaching/pkg/errors"
	appJwt "paopao-ce-teaching/pkg/jwt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			token string
			ecode = errors.Success
		)
		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("Authorization")

			// 验证前端传过来的token格式，不为空，开头为Bearer
			if token == "" || !strings.HasPrefix(token, "Bearer ") {
				app.ToErrorResponse(c, errors.UnauthorizedTokenError)
				c.Abort()
				return
			}

			// 验证通过，提取有效部分（除去Bearer)
			token = token[7:]
		}
		if token == "" {
			ecode = errors.InvalidParams
		} else {
			claims, err := appJwt.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					ecode = errors.UnauthorizedTokenTimeout
				default:
					ecode = errors.UnauthorizedTokenError
				}
			} else {
				c.Set("UID", claims.UID)
				c.Set("USERNAME", claims.Username)
			}
		}

		if ecode != errors.Success {
			app.ToErrorResponse(c, ecode)
			c.Abort()
			return
		}

		c.Next()
	}
}
