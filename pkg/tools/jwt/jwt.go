package jwt

import (
	"errors"
	"github.com/Madou-Shinni/gin-quickstart/internal/conf"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const Security = "#eqw1"

// MyClaims 自定义结构体，并内嵌jwt.StandardClaims
// jwt.StandardClaims只包含官方字段
// 需要自定义需要的字段
type MyClaims struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// 生成access token
// param: claims需要保存在token中的信息
// return: token信息 and error
func GenerateAccessToken(claims MyClaims) (string, error) {
	// 当前时间
	now := time.Now()
	// 过期时间
	accessExpire := now.Unix() + claims.ExpiresAt
	//   签发人
	issuer := claims.Issuer
	// 密钥
	secret := Security

	claims = MyClaims{
		UserId:   claims.UserId,
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: accessExpire, // 过期时间
			Issuer:    issuer,       // 签名
		},
	}
	// 根据签名生成token，NewWithClaims(加密方式,claims) ==》 头部，载荷，签证
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	return accessToken, err
}

// 生成refresh token
// param: claims需要保存在token中的信息
// return: token信息 and error
func GenerateRefreshToken(claims MyClaims) (string, error) {
	// 当前时间
	now := time.Now()
	// 过期时间
	refreshExpire := now.Unix() + claims.ExpiresAt
	//   签发人
	issuer := claims.Issuer
	// 密钥
	secret := Security

	claims = MyClaims{
		UserId:   claims.UserId,
		Username: claims.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshExpire, // 过期时间
			Issuer:    issuer,        // 签名
		},
	}
	// 根据签名生成token，NewWithClaims(加密方式,claims) ==》 头部，载荷，签证
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	return refreshToken, err
}

// 解析token
// param: t 生成的token字符串
// return: 自定义的jwt信息结构体指针*MyClaims 和 error
func ParseToken(t string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(t, MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return conf.Conf.Secret, nil
	})
	if err != nil {
		// token解析失败
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("That's not even a token")
			} else if v.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, errors.New("Token is expired")
			} else if v.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("Token not active yet")
			} else {
				return nil, errors.New("Couldn't handle this token:")
			}
		}
	}
	if token != nil {
		// 校验token
		if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New("Couldn't handle this token:")

	} else {
		return nil, errors.New("Couldn't handle this token:")

	}
}

// 刷新token
// param: t 已经生成的token
// return: token and error
func RefreshToken(t string) (string, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}
	token, err := jwt.ParseWithClaims(t, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return conf.Conf.Secret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		jwt.TimeFunc = time.Now
		// 延长一小时
		claims.StandardClaims.ExpiresAt = time.Now().Add(1 * time.Hour).Unix()
		return GenerateAccessToken(*claims)
	}
	return "", errors.New("Couldn't handle this token:")
}
