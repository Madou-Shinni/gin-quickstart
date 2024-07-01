package tools

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strconv"
)

const (
	UserIdKey = "userId"
	RoleIdKey = "roleId"
	ExpKey    = "exp" // 过期时间key
)

var (
	ErrorUserInfo = errors.New("用户异常，请重新登录")
)

// GenToken 生成token map中key=exp过期时间
func GenToken(mapClaims jwt.MapClaims, signed string) (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err = claims.SignedString([]byte(signed))
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserIdFromJwt 解析token
func GetUserIdFromJwt(tokenStr string, signed string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signed), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return 0, jwt.ErrTokenInvalidClaims
	} else {
		userIdStr := claims[UserIdKey].(string)
		userId, err := strconv.ParseUint(userIdStr, 10, 64)
		if err != nil {
			return 0, ErrorUserInfo
		}
		return uint(userId), nil
	}
}

// GetRoleIdFromJwt 解析token
func GetRoleIdFromJwt(tokenStr string, signed string) (uint, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signed), nil
	})
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return 0, jwt.ErrTokenInvalidClaims
	} else {
		roleIdStr := claims[RoleIdKey].(string)
		roleId, err := strconv.ParseUint(roleIdStr, 10, 64)
		if err != nil {
			return 0, ErrorUserInfo
		}
		return uint(roleId), nil
	}
}

// GetClaimsFromJwt 解析token
func GetClaimsFromJwt(tokenStr string, signed string) (jwt.MapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(signed), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	} else {
		return claims, nil
	}
}
