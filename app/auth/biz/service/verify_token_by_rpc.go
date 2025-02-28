package service

import (
	"context"
	"fmt"

	"github.com/naskids/nas-mall/app/auth/biz/model"
	"github.com/naskids/nas-mall/common/token"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type VerifyTokenByRPCService struct {
	ctx        context.Context
	tokenMaker *token.Maker
	userStore  model.AuthUserStore
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx, tokenMaker: token.DefaultMaker(), userStore: model.DefaultAuthUserStore()}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyTokenResp, err error) {
	claims, err := s.tokenMaker.ParseAccessToken(req.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("remote failed to verify token: [%w]", err)
	}
	// 检查缓存白名单有该token
	//if not in white list {
	//	return &auth.VerifyTokenResp{IsVaild: false}, nil
	//}

	userId := claims["uid"].(uint64)
	role := claims["rol"].(string)

	// 检查用户黑名单没有
	//if in user blacklist {
	//	return &auth.VerifyTokenResp{IsVaild: false}, nil
	//}

	// 从存储中获取 user 最新信息
	user, err := s.userStore.GetUser(userId)
	if err != nil {
		return nil, fmt.Errorf("get user by id failed: [%w]", err)
	}
	// 角色信息一致
	return &auth.VerifyTokenResp{
		IsValid: role == user.Role,
	}, nil
}
