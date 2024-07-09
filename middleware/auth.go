package middleware

import (
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools"
	"github.com/gin-gonic/gin"
)

// jwt认证
// 解析请求头中的token
// 解析成功则通过，失败返回错误信息
func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			response.Error(c, constant.CODE_NO_PERMISSIONS, constant.CODE_NO_PERMISSIONS.Msg())
			c.Abort()
			return
		}
		// 解析token
		claims, err := tools.GetClaimsFromJwt(token, conf.Conf.JwtConfig.Secret)
		if err != nil {
			response.Error(c, constant.CODE_NO_PERMISSIONS, err.Error())
			c.Abort()
			return
		}

		userId := claims[tools.UserIdKey]
		roleId := claims[tools.RoleIdKey]

		// 将解析的userId保存到上下文中
		c.Set(constants.CtxUserIdKey, userId)
		c.Set(constants.CtxRoleIdkEY, roleId)
		c.Next()
	}
}
