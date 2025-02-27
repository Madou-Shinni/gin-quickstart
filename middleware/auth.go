package middleware

import (
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools"
	"github.com/gin-gonic/gin"
	"net/http"
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

// UnameAndPwdAuth 用户名密码认证
func UnameAndPwdAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		amc := conf.Conf.AsynqConfig.Monitor
		username := c.Request.Header.Get("username")
		pwd := c.Request.Header.Get("pwd")

		// 如果没有认证或者认证失败，返回 401 错误
		if username != amc.Username || pwd != amc.Password {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 认证成功，继续处理请求
		c.Next()
	}
}

func ipWhitelistMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的客户端 IP 地址
		clientIP := c.ClientIP()
		// 定义允许访问的 IP 地址列表
		allowedIPs := []string{"192.168.1.100", "192.168.1.101"}

		// 检查客户端 IP 是否在允许的列表中
		isAllowed := false
		for _, ip := range allowedIPs {
			if clientIP == ip {
				isAllowed = true
				break
			}
		}

		// 如果 IP 地址不在白名单中，返回 403 Forbidden 错误
		if !isAllowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// 允许访问，继续处理请求
		c.Next()
	}
}
