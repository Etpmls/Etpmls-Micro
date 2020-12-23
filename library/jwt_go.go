package em_library
//https://www.jianshu.com/p/0c60f665d5d7
//https://godoc.org/github.com/dgrijalva/jwt-go#NewWithClaims

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strconv"
)

type jwtGo struct {
}

func NewJwtGo() *jwtGo {
	return &jwtGo{}
}

// Create Token
// 创建令牌
func (j *jwtGo)CreateToken(c interface{}, secret string) (t string, err error) {
	claims, ok := c.(*jwt.StandardClaims)
	if !ok {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		return t, err
	}
	return ss, err
}


// Parse Token
// 解析令牌
func (j *jwtGo)ParseToken(tokenString string, secret string) (t interface{}, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if token == nil || err != nil {
		return t, err
	}

	if token.Valid {
		return token, err
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return t, errors.New("令牌格式错误！")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return t, errors.New("令牌已过期或尚未激活！")
		} else {
			return t, errors.New("不能处理该令牌！")
		}
	} else {
		return t, errors.New("不能处理该令牌！")
	}
}


// Get Username
// 获取用户名
func (j *jwtGo)GetIssuerByToken(tokenString string, secret string) (issuer string, err error) {
	tmp, err := j.ParseToken(tokenString, secret)
	if err != nil {
		return "", err
	}

	tk, ok := tmp.(*jwt.Token)
	if !ok {
		return "", err
	}

	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		issuer, ok := claims["iss"].(string)
		if !ok {
			return "", errors.New("令牌ID解析错误！")
		}

		return issuer, nil
	}

	return "", errors.New("当前token无效！")
}


// Get User ID
// 获取用户ID
func (j *jwtGo)GetIdByToken(tokenString string, secret string) (userId int, err error) {
	tmp, err := j.ParseToken(tokenString, secret)
	if err != nil {
		return 0, err
	}

	tk, ok := tmp.(*jwt.Token)
	if !ok {
		return 0, err
	}

	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		id, ok := claims["jti"].(string)
		if !ok {
			return 0, errors.New("令牌ID解析错误！")
		}

		userId, err := strconv.Atoi(id)
		if err != nil {
			return 0, errors.New("令牌中的ID不是数字！")
		}

		return userId, nil
	}

	return 0, errors.New("当前token无效！")
}



