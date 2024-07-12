package middleware

import (
	"bytes"
	"fmt"
	"github.com/Madou-Shinni/go-logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		form := ""
		bodyStr := ""

		if strings.Contains(c.ContentType(), "form") {
			// 解析表单数据里面的数据
			if err := c.Request.ParseForm(); err != nil {
				logger.Error("c.Request.ParseForm()", zap.Error(err))
			}
			// 读取表单数据
			form = c.Request.PostForm.Encode()
		}

		// 判断请求类型是否是json
		if strings.Contains(c.ContentType(), "application/json") {
			defer c.Request.Body.Close()
			body, _ := ioutil.ReadAll(c.Request.Body)
			//注意：重新赋值必须这样否则无法从context重在获取数据
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			bodyStr += string(body)
		}

		c.Next()

		cost := time.Since(start).Milliseconds()
		// 接口耗时 小于1000显示单位毫秒 大于1000显示单位秒
		costTime := map[bool]string{true: fmt.Sprintf("%vms", cost), false: fmt.Sprintf("%vs", float64(cost)/float64(1000))}[cost < 1000]

		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Any("form", form),
			zap.Any("json-body", bodyStr),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("cost", costTime),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
