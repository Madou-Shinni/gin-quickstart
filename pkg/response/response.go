package response

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    constant.RspCode `json:"code,omitempty"`
	Message string           `json:"message,omitempty"`
	Data    interface{}      `json:"data,omitempty"`
}

type PageResponse struct {
	Total int64       `json:"total,omitempty"`
	List  interface{} `json:"list,omitempty"`
}

func Success(c *gin.Context, data ...interface{}) {
	var r Response
	r.Code = http.StatusOK
	r.Message = "success"

	if len(data) > 0 {
		r.Data = data[0]
	}

	c.JSON(http.StatusOK, r)
}

func Error(c *gin.Context, errorCode constant.RspCode, msg ...string) {
	var r Response

	r.Code = errorCode

	if len(msg) > 0 {
		r.Message = msg[0]
	}

	c.JSON(http.StatusOK, r)
}
