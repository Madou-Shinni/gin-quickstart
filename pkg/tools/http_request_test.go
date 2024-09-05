package tools

import (
	"testing"
	"time"
)

const Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVVUlEIjoiYTY2NDg5MGMtOWMyMS00ZDkzLWE5M2MtMGU5YzgzYmM0NzBkIiwiSUQiOjE5LCJVc2VybmFtZSI6Inh1d2VpIiwiTmlja05hbWUiOiLlvpDngpwiLCJBdXRob3JpdHlJZCI6ODg4LCJCdWZmZXJUaW1lIjo4NjQwMDAwMDAwMDAwMCwiZXhwIjoxNzA1MTExNTA5LCJpc3MiOiJxbVBsdXMiLCJuYmYiOjE3MDQzNDg4MDV9.ccS2irYOSqZMAKr0-DODCSEY4sFfq7aqwoNbVDWSoXA"

func TestNewRequest(t *testing.T) {
	headers := map[string]string{
		"X-Token": Token,
	}
	data := map[string]interface{}{
		"page":     2,
		"pageSize": 10,
	}
	resp, err := NewRequest("GET", time.Second, "http://localhost:9999", data, headers)
	if err != nil {
		return
	}

	t.Logf("resp: %s", string(resp))
}
