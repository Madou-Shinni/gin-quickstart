package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/Madou-Shinni/gin-quickstart/pkg/response"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/cache"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/md5"
	"github.com/Madou-Shinni/go-logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const cacheExpire = time.Second * 100

func Cache(ca cache.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		pre := c.Request.URL.Path
		key := pre + md5.Md5To16(c.Request.RequestURI)

		respWrite := &responseCacheWriter{ResponseWriter: c.Writer}
		c.Writer = respWrite

		if c.Request.Method == "GET" {
			// 获取缓存
			res, err := ca.Get(key)
			if err != nil || res == "" {
				// 缓存空
				c.Next()
				// 获取接口返回值，设置缓存
				resp := response.Response{}
				json.Unmarshal([]byte(respWrite.body.String()), &resp)
				if resp.Code != http.StatusOK {
					// 不设置缓存
					return
				}
				err = ca.Set(key, resp.Data, cacheExpire)
				if err != nil {
					logger.Error("cache err", zap.Error(err))
				}
				return
			}
			// 获取缓存返回
			var data interface{}
			json.Unmarshal([]byte(res.(string)), &data)
			response.Success(c, data)
			c.Abort()
			return
		} else {
			// 删除缓存
			c.Next()
			err := ca.DelByPrefix(pre)
			if err != nil {
				logger.Error("cache err", zap.Error(err))
			}
			return
		}
	}
}

// responseCacheWriter
type responseCacheWriter struct {
	gin.ResponseWriter

	body bytes.Buffer
}

func (w *responseCacheWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *responseCacheWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
