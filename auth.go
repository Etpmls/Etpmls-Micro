package em

import (
	"context"
	em_library "github.com/Etpmls/Etpmls-Micro/library"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type auth struct {

}

// Verify that the token is valid
// 验证token是否有效
func (this *auth) VerifyToken(token string) (bool, error) {
	// Parse token
	_, err := JwtToken.ParseToken(token, em_library.Config.App.Key)
	if err != nil {
		LogInfo.Output(MessageWithLineNum(err.Error()))
		return false, err
	}

	return true, nil
}

// Parse the token
// 解析token
func (this *auth) ParseToken(tokenString string) (interface{}, error) {
	return JwtToken.ParseToken(tokenString, em_library.Config.App.Key)
}

// Create a general token
// 创建通用token
func (this *auth) CreateGeneralToken(userId int, username string) (string, error) {
	return JwtToken.CreateToken(&jwt.StandardClaims{
		Id: strconv.Itoa(userId),                                                          // 用户ID
		ExpiresAt: time.Now().Add(time.Second * em_library.Config.App.TokenExpirationTime).Unix(), // 过期时间 - 12个小时
		Issuer:    username,                                                                    // 发行者
	}, em_library.Config.App.Key)
}

func (this *auth) GetIdByToken(token string) (int, error) {
	return JwtToken.GetIdByToken(token, em_library.Config.App.Key)
}

func (this *auth) GetIssuerByToken(token string) (string, error) {
	return JwtToken.GetIssuerByToken(token, em_library.Config.App.Key)
}

// Get token from ctx
// 从ctx获取令牌
func (this *auth) GetTokenFromCtx(ctx context.Context) (string, error) {
	var request request
	return request.GetValueFromCtx(ctx, "token")
}

// Get token from header
// 从header获取令牌
func (this *auth) Rpc_GetTokenFromHeader(ctx context.Context) (string, error) {
	var request request
	return request.Rpc_GetValueFromHeader(ctx, "token")
}

/*

// Verify that the token has access permissions
// 验证token是否具有访问权限
func (this *auth)VerifyPermissions(token string, fullMethodName string) (error) {
	// Get Claims
	// 获取Claims
	tmp, err := em_library.JwtToken.ParseToken(token)
	tk, ok := tmp.(*jwt.Token)
	if !ok || err != nil {
		LogInfo.Output(em_utils.MessageWithLineNum("Get Claims failed!" + err.Error()))
		return err
	}

	// Determine whether the role has the corresponding permissions
	// 判断所属角色是否有相应的权限
	if claims,ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		if userId, ok := claims["jti"].(string); ok {
			id, err := strconv.Atoi(userId)
			if err == nil {
				b := NewClient().AuthCheck(EA.AuthServiceName, fullMethodName, uint(id))
				if b {
					return nil
				}
			} else {
				LogInfo.Output(em_utils.MessageWithLineNum(err.Error()))
			}
		}
	}

	return errors.New("PermissionDenied")
}*/