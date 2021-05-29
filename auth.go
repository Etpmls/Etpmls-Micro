package em

import (
	"context"
	"github.com/Etpmls/Etpmls-Micro/v3/define"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type auth struct {

}

// Verify that the token is valid
// 验证token是否有效
func (this *auth) VerifyToken(token string) (bool, error) {
	k, err := Kv.ReadKey(em_define.KvAppKey)
	if err != nil {
		LogInfo.Path(err)
		return false ,err
	}

	// Parse token
	_, err = JwtToken.ParseToken(token, k)
	if err != nil {
		LogInfo.New(MessageWithLineNum(err.Error()))
		return false, err
	}

	return true, nil
}

// Parse the token
// 解析token
func (this *auth) ParseToken(tokenString string) (interface{}, error) {
	k, err := Kv.ReadKey(em_define.KvAppKey)
	if err != nil {
		LogInfo.Path(err)
		return nil ,err
	}

	return JwtToken.ParseToken(tokenString, k)
}

// Create a general token
// 创建通用token
func (this *auth) CreateGeneralToken(userId int, username string) (string, error) {
	m, err := Kv.List(em_define.KvApp)
	if err != nil {
		LogInfo.Path(err)
		return "" ,err
	}

	d, err := time.ParseDuration(m[em_define.KvAppTokenExpirationTime])
	if err != nil {
		LogInfo.Path("[Default: " +em_define.DefaultAppTokenExpirationTimeText+ "]" +em_define.KvAppTokenExpirationTime+ " is not configured or format is incorrect.", err)
		d = em_define.DefaultAppTokenExpirationTime
	}

	return JwtToken.CreateToken(&jwt.StandardClaims{
		Id: strconv.Itoa(userId),                                                          // 用户ID
		ExpiresAt: time.Now().Add(d).Unix(), // 过期时间 - 12个小时
		Issuer:    username,                                                                    // 发行者
	}, m[em_define.KvAppKey])
}

func (this *auth) GetIdByToken(token string) (int, error) {
	k, err := Kv.ReadKey(em_define.KvAppKey)
	if err != nil {
		LogInfo.Path(err)
		return 0 ,err
	}

	return JwtToken.GetIdByToken(token, k)
}

func (this *auth) GetIssuerByToken(token string) (string, error) {
	k, err := Kv.ReadKey(em_define.KvAppKey)
	if err != nil {
		LogInfo.Path(err)
		return "" ,err
	}

	return JwtToken.GetIssuerByToken(token, k)
}

// Get token from ctx
// 从ctx获取令牌
func (this *auth) GetTokenFromCtx(ctx context.Context) (string, error) {
	var request request
	return request.GetValueFromCtx(ctx, "token")
}

// Get token from header
// 从header获取令牌
func (this *auth) GetTokenFromHeader(ctx context.Context) (string, error) {
	var request request
	return request.GetValueFromHeader(ctx, "token")
}
