package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"testing"
)

// 由于我们生成token采用了读取配置文件的方式
// 这里的测试无法正常实现
func TestGenerateToken(t *testing.T) {
	claims := MyClaims{
		UserId:   1,
		Username: "sni",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 0,
			IssuedAt:  0,
		},
	}
	accessToken, err := GenerateAccessToken(claims)
	refreshToken, err := GenerateRefreshToken(claims)
	if err != nil {
		t.Errorf("token生成失败:%v", err)
	}

	fmt.Println(map[string]interface{}{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}

func Test1(t *testing.T) {

}

func TestParseToken(t *testing.T) {
	claims := MyClaims{
		UserId:   1,
		Username: "sni",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 0,
			IssuedAt:  0,
		},
	}
	token, _ := GenerateAccessToken(claims)
	tokenInfo, err := ParseToken(token)
	if err != nil {
		// 解析token错误
	}

	id := tokenInfo.UserId
	username := tokenInfo.Username

	fmt.Printf("id:%v,\tusername:%v", id, username)
}

func TestRefreshToken(t *testing.T) {
	claims := MyClaims{
		UserId:   1,
		Username: "sni",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 0,
			IssuedAt:  0,
		},
	}
	token, _ := GenerateRefreshToken(claims)

	refreshToken, err := RefreshToken(token)
	if err != nil {

	}

	fmt.Println("refreshToken:", refreshToken)
}
