package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lackone/gin-scaffold/internal/contract"
	"github.com/lackone/gin-scaffold/internal/core"
	"github.com/lackone/gin-scaffold/pkg/app"
	"github.com/lackone/gin-scaffold/pkg/errcode"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		errCode := errcode.Success
		token := ""

		if s, exist := c.GetQuery("token"); exist {
			token = s
		} else {
			token = c.GetHeader("token")
		}

		if token == "" {
			errCode = errcode.InvalidParams
		} else {
			engine := c.MustGet("engine").(*core.Engine)
			container := engine.GetContainer()
			configService := container.MustMake(contract.ConfigKey).(contract.Config)

			_, err := app.ParseJWT(token, configService.GetString("jwt.secret"))
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					errCode = errcode.UnauthorizedTokenTimeout
				default:
					errCode = errcode.UnauthorizedTokenError
				}
			}
		}

		if errCode != errcode.Success {
			response := app.NewResponse(c)
			response.ToError(errCode)
			c.Abort()
			return
		}

		c.Next()
	}
}
