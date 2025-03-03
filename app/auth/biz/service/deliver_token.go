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
	accessToken, err = token.Maker.GenerateAccessToken(utils.H{"uid": user.UserID, "rol": user.Role})
	if err != nil {
		return nil, fmt.Errorf("access token gen err: [%w]", err)
	}
	refreshToken, err = token.Maker.GenerateRefreshToken(utils.H{"uid": user.UserID, "rol": user.Role, "ver": user.RefreshVersion})
	if err != nil {
		return nil, fmt.Errorf("refresh token gen err: [%w]", err)
	}

	// 持久化 access token 到 Redis 白名单
	accessKey := fmt.Sprintf("auth:access:%s", accessToken)
	if err := redis.RedisClient.Set(s.ctx, accessKey, user.UserID, token.Maker.AccessKeyDuration).Err(); err != nil {
		return nil, fmt.Errorf("存储 access_token 失败: %w", err)
	}

	// 维护用户 token 集合（用于批量管理）
	userTokensKey := fmt.Sprintf("auth:user_tokens:%d", user.UserID)
	redis.RedisClient.SAdd(s.ctx, userTokensKey, accessKey)
	return &auth.DeliveryTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
