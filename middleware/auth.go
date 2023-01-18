package middleware

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/jwt"
	"github.com/gin-gonic/gin"
)

// jwt认证
// 解析请求头中的token
// 解析成功则通过，失败返回错误信息
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("token")
		if token == "" {
			response.Error(c, constant.CODE_NO_PERMISSIONS, constant.CODE_NO_PERMISSIONS.Msg())
			c.Abort()
			return
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			response.Error(c, constant.CODE_NO_PERMISSIONS, constant.CODE_NO_PERMISSIONS.Msg())
			c.Abort()
			return
		}

		// 将token的信息保存到上下文中
		c.Set(constant.TokenKey, claims)

		c.Next()
	}
}
