package em

import (
	"context"
)

type auth struct {

}

// Verify that the token is valid
// 验证token是否有效
func (this *auth)VerifyToken(token string) (bool, error) {
	// Parse token
	_, err := JwtToken.ParseToken(token)
	if err != nil {
		LogInfo.Output(MessageWithLineNum(err.Error()))
		return false, err
	}

	return true, nil
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