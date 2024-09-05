package tools

import (
	"encoding/json"
	"time"
)

const (
	LoginGithubUrl = "https://github.com/login/oauth/access_token" // 登录github
	GithubUserUrl  = "https://api.github.com/user"                 // github用户信息
)

type (
	// LoginGithubReq 登录github
	LoginGithubReq struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		Code         string `json:"code"`
	}

	// LoginGithubResp 登录github返回
	LoginGithubResp struct {
		AccessToken string `json:"access_token"`
	}

	// GithubUserResp github用户信息
	GithubUserResp struct {
		ID     uint   `json:"id"`         // github用户id
		Login  string `json:"login"`      // github用户名
		Name   string `json:"name"`       // github用户昵称
		Email  string `json:"email"`      // github用户邮箱
		Avatar string `json:"avatar_url"` // github用户头像
		Bio    string `json:"bio"`        // github用户简介
	}
)

// LoginGithub 登录github
func LoginGithub(req LoginGithubReq) (result LoginGithubResp, err error) {
	headers := map[string]string{
		"Accept": "application/json",
	}
	data := map[string]interface{}{
		"client_id":     req.ClientId,
		"client_secret": req.ClientSecret,
		"code":          req.Code,
	}
	resp, err := NewRequest(GET, time.Second, LoginGithubUrl, data, headers)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return
	}

	return
}

func (r LoginGithubResp) GetUserInfo() (result *GithubUserResp, err error) {
	header := map[string]string{
		"Authorization": "token " + r.AccessToken,
	}

	resp, err := NewRequest(GET, time.Second, GithubUserUrl, nil, header)
	if err != nil {
		return
	}

	err = json.Unmarshal(resp, &result)
	if err != nil {
		return
	}

	return
}
