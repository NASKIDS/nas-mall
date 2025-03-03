package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/utils"

	"github.com/naskids/nas-mall/app/auth/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/auth/biz/dal/redis"
	"github.com/naskids/nas-mall/app/auth/biz/model"
	"github.com/naskids/nas-mall/common/token"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type DeliverTokenService struct {
	ctx context.Context
} // NewDeliverTokenService new DeliverTokenService
func NewDeliverTokenService(ctx context.Context) *DeliverTokenService {
	return &DeliverTokenService{ctx: ctx}
}

// Run create note info
func (s *DeliverTokenService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryTokenResp, err error) {
	// 1. 验证用户身份
	user, err := model.GetUser(s.ctx, mysql.DB, redis.RedisClient, req.UserId)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 2. 生成令牌
	var accessToken string
	var refreshToken string
	accessToken, err = token.Maker.GenerateAccessToken(utils.H{"uid": user.ID, "rol": user.Role})
	if err != nil {
		return nil, fmt.Errorf("access token gen err: [%w]", err)
	}
	refreshToken, err = token.Maker.GenerateRefreshToken(utils.H{"uid": user.ID, "rol": user.Role, "ver": user.RefreshVersion})
	if err != nil {
		return nil, fmt.Errorf("refresh token gen err: [%w]", err)
	}

	// TODO 3. 持久化 token 到 redis
	return &auth.DeliveryTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
