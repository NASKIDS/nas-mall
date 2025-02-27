package service

import (
	"context"
	"errors"
	"time"

	"github.com/naskids/nas-mall/app/auth/biz/model"
	"github.com/naskids/nas-mall/app/auth/utils/token"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type DeliverTokenService struct {
	ctx        context.Context
	tokenMaker token.Maker
	userStore  model.AuthUser
} // NewDeliverTokenService new DeliverTokenService
func NewDeliverTokenService(ctx context.Context) *DeliverTokenService {
	return &DeliverTokenService{ctx: ctx}
}

// Run create note info
func (s *DeliverTokenService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	// 1. 验证用户身份
	user, err := s.userStore.GetUser(req.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 生成令牌
	accessToken, _ := s.tokenMaker.GenerateAccessToken(user.ID, s.tokenDuration)
	refreshToken, _ := s.tokenMaker.GenerateRefreshToken(user.ID, user.RefreshVersion, s.refreshDuration)

	return &auth.DeliveryResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
