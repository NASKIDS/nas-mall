package token

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/utils"
)

// VerifyAccessToken 验证并取出 access token 中的 claims，如果验证不通过返回 err
func VerifyAccessToken(tk string) (claims utils.H, err error) {
	tokenMaker := InitTokenMaker()
	claims, err = tokenMaker.ParseAccessToken(tk)
	if err != nil {
		return nil, fmt.Errorf("invaild token: [%w]", err)
	}
	return claims, nil
}
