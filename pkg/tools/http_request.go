package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

// NewRequest 请求包装
func NewRequest(method HttpMethod, timeout time.Duration, url string, data map[string]interface{}, headers map[string]string) (body []byte, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	if method == "GET" {
		var query string
		for k, v := range data {
			query = fmt.Sprint(query, "&", k, "=", v)
		}
		query = removePrefix(query)
		url = fmt.Sprint(url, "?", query)
		data = nil
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, string(method), url, bytes.NewBuffer(marshal))
	if err != nil {
		return body, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return body, err
	}

	return body, err
}

// NewRequestV2 请求包装并且返回状态码
func NewRequestV2(method HttpMethod, timeout time.Duration, url string, data map[string]interface{}, headers map[string]string) (body []byte, code int, err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	if method == "GET" {
		var query string
		for k, v := range data {
			query = fmt.Sprint(query, "&", k, "=", v)
		}
		query = removePrefix(query)
		url = fmt.Sprint(url, "?", query)
		data = nil
	}

	// 如果是 POST 请求，并且 Content-Type 是 multipart/form-data
	if method == "POST" && headers["Content-Type"] == "multipart/form-data" {
		// 使用 multipart.Writer 创建 multipart 请求体
		bodyBuffer := &bytes.Buffer{}
		writer := multipart.NewWriter(bodyBuffer)

		// 将传入的字段数据添加到 multipart 请求体中
		for k, v := range data {
			switch v := v.(type) {
			case string:
				// 如果是字符串类型，使用 CreateFormField 添加
				err := writer.WriteField(k, v)
				if err != nil {
					return nil, 0, err
				}
			case []byte:
				// 如果是文件（例如[]byte类型的文件内容），使用 CreateFormFile 添加
				part, err := writer.CreateFormFile(k, "file") // file是文件的字段名，你可以根据需要改
				if err != nil {
					return nil, 0, err
				}
				_, err = part.Write(v) // 写入文件内容
				if err != nil {
					return nil, 0, err
				}
			default:
				err := writer.WriteField(k, fmt.Sprint(v))
				if err != nil {
					return nil, 0, err
				}
			}
		}

		// 关闭 writer，这将自动生成 multipart 的边界
		err := writer.Close()
		if err != nil {
			return nil, 0, err
		}

		// 构建请求
		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx, string(method), url, bodyBuffer)
		if err != nil {
			return nil, 0, err
		}

		// 设置 Content-Type 为 multipart/form-data，并传递 boundary
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// 处理请求头
		if headers != nil {
			for k, v := range headers {
				req.Header.Add(k, v)
			}
		}

		// 发送请求
		resp, err := client.Do(req)
		if err != nil {
			return nil, 0, err
		}

		// 读取响应体
		body, err = io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return body, resp.StatusCode, err
		}

		return body, resp.StatusCode, nil
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		return nil, 0, err
	}

	client := http.Client{}
	req, err := http.NewRequestWithContext(ctx, string(method), url, bytes.NewBuffer(marshal))
	if err != nil {
		return body, 0, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	body, err = io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return body, resp.StatusCode, err
	}

	return body, resp.StatusCode, err
}

func removePrefix(input string) string {
	// 如果字符串为空或只有一个字符，则返回空字符串
	if len(input) <= 1 {
		return ""
	}

	// 使用切片获取除第一个字符外的所有字符
	result := input[1:]
	return result
}
