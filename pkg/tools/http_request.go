package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

// NewRequest 请求包装
func NewRequest(method HttpMethod, url string, data map[string]interface{}, headers map[string]string) (body []byte, err error) {

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
	req, err := http.NewRequest(string(method), url, bytes.NewBuffer(marshal))
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

func removePrefix(input string) string {
	// 如果字符串为空或只有一个字符，则返回空字符串
	if len(input) <= 1 {
		return ""
	}

	// 使用切片获取除第一个字符外的所有字符
	result := input[1:]
	return result
}
