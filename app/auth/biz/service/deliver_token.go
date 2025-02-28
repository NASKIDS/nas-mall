package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/utils"

	"github.com/naskids/nas-mall/app/auth/biz/model"
	"github.com/naskids/nas-mall/common/token"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type DeliverTokenService struct {
	ctx        context.Context
	tokenMaker *token.Maker
	userStore  model.AuthUserStore
} // NewDeliverTokenService new DeliverTokenService
func NewDeliverTokenService(ctx context.Context) *DeliverTokenService {
	return &DeliverTokenService{ctx: ctx, tokenMaker: token.DefaultMaker(), userStore: model.DefaultAuthUserStore()}
}

// Run create note info
func (s *DeliverTokenService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryTokenResp, err error) {
	// 1. 验证用户身份
	user, err := s.userStore.GetUser(req.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 生成令牌
	var accessToken string
	var refreshToken string
	accessToken, err = s.tokenMaker.GenerateAccessToken(utils.H{"uid": user.ID, "rol": user.Role})
	if err != nil {
		return nil, fmt.Errorf("access token gen err: [%w]", err)
	}
	refreshToken, err = s.tokenMaker.GenerateRefreshToken(utils.H{"uid": user.ID, "rol": user.Role, "ver": user.RefreshVersion})
	if err != nil {
		return nil, fmt.Errorf("refresh token gen err: [%w]", err)
	}

	// TODO 3. 持久化 token 到 redis
	return &auth.DeliveryTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
