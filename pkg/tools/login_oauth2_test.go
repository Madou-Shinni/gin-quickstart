package tools

import "testing"

func TestLoginGithub(t *testing.T) {
	req := LoginGithubReq{
		ClientId:     "client_id",
		ClientSecret: "client_secret",
		Code:         "code",
	}
	resp, err := LoginGithub(req)
	if err != nil {
		t.Errorf("LoginGithub() error = %v", err)
		return
	}
	t.Logf("resp: %v", resp)
}
