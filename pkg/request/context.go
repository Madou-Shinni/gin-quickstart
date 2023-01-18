package request

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/jwt"
	"github.com/gin-gonic/gin"
)

// 获取当前用户信息
// param: c gin上下文
// return: 自定义的jwt信息指针
func GetCurrentUser(c *gin.Context) *jwt.MyClaims {
	if token, exists := c.Get(constant.TokenKey); !exists {
		return nil
	} else {
		claims := token.(*jwt.MyClaims)
		return claims
	}
}
