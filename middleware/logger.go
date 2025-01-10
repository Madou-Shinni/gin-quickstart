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
	"runtime"
	"strings"
	"time"
)

const (
	defaultSkipFrames = 4
	defaultMaxPC      = 4
	maxBodyLength     = 1024 // 最大打印长度
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
			// 如果请求体长度大于设定的最大限制，截取前 MaxBodyLength 个字符
			if len(body) > maxBodyLength {
				bodyStr = string(body[:maxBodyLength])
			} else {
				bodyStr = string(body)
			}
			//注意：重新赋值必须这样否则无法从context重在获取数据
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
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
						zap.String("stack", getRelevantStack(defaultSkipFrames)),
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

// 获取有效的堆栈信息 忽略前skipFrames层
func getRelevantStack(skipFrames int) string {
	stackBuf := make([]uintptr, defaultMaxPC) // 限制最多10层调用深度
	length := runtime.Callers(skipFrames, stackBuf[:])

	stack := strings.Builder{}
	for i := 0; i < length; i++ {
		pc := stackBuf[i]
		fn := runtime.FuncForPC(pc - 1) // 获取当前帧的函数信息
		if fn == nil {
			continue
		}
		file, line := fn.FileLine(pc - 1)
		// 过滤掉一些无关的系统文件或日志打印相关的代码
		if !strings.Contains(file, "runtime/") && !strings.Contains(file, "zap") {
			stack.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, fn.Name()))
		}
	}
	return stack.String()
}
