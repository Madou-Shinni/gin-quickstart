package common

import (
	"errors"
	"github.com/Madou-Shinni/gin-quickstart/constants"
	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

// GetUserIdFromCtx 从上下文中获取用户id
func GetUserIdFromCtx(c *gin.Context) (uint, error) {
	u, ok := c.Get(constants.CtxUserIdKey)
	if !ok {
		return 0, ErrorUserNotLogin
	}

	return u.(uint), nil
}
