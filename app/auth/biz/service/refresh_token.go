package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"

	"github.com/naskids/nas-mall/app/auth/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/auth/biz/dal/redis"
	"github.com/naskids/nas-mall/app/auth/biz/model"
	"github.com/naskids/nas-mall/common/token"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type RefreshTokenService struct {
	ctx context.Context
} // NewRefreshTokenService new RefreshTokenService
func NewRefreshTokenService(ctx context.Context) *RefreshTokenService {
	return &RefreshTokenService{ctx: ctx}
}

// Run create note info
func (s *RefreshTokenService) Run(req *auth.RefreshTokenReq) (resp *auth.RefreshTokenResp, err error) {
	// 1. 解析刷新令牌
	var tk utils.H
	tk, err = token.Maker.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: [%w]", err)
	}

	userID, tokenVersion, role := uint64(tk["uid"].(float64)), uint64(tk["ver"].(float64)), tk["rol"].(string)
	// 2. 验证用户
	user, err := model.GetUser(s.ctx, mysql.DB, redis.RedisClient, userID)
	if err != nil {
		return nil, err
	}
	if user.UserID != userID || user.Role != role {
		return nil, fmt.Errorf("incorrect user info: [%w]", err)
	}
	var currentVersion uint64
	currentVersion, err = model.GetRefreshVersion(s.ctx, mysql.DB, redis.RedisClient, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token version in db: [%w]", err)
	}
	if currentVersion > tokenVersion {
		return nil, fmt.Errorf("stale refresh token version: [%d]", tokenVersion)
	}
	// 3. 生成新令牌并更新版本
	newAccess, err := token.Maker.GenerateAccessToken(utils.H{"uid": userID, "rol": role})
	if err != nil {
		return nil, fmt.Errorf("access token gen err: [%w]", err)
	}
	newRefresh, err := token.Maker.GenerateRefreshToken(utils.H{"uid": userID, "rol": role, "ver": currentVersion + 1})
	if err != nil {
		return nil, fmt.Errorf("refresh token gen err: [%w]", err)
	}
	// 存储新 refresh_token 并设置过期时间（与令牌有效期一致）
	if err := redis.RedisClient.Set(
		s.ctx, fmt.Sprintf("auth:refresh:%d", user.UserID), newRefresh, token.Maker.RefreshKeyDuration,
	).Err(); err != nil {
		klog.Errorf("存储 refresh_token 失败: %v", err)
	}
	err = model.UpdateRefreshVersion(s.ctx, mysql.DB, redis.RedisClient, userID, currentVersion+1)
	if err != nil {
		return nil, fmt.Errorf("failed to update token: [%w]", err)
	}

	// 持久化新 access token
	newAccessKey := fmt.Sprintf("auth:access:%s", newAccess)
	if err := redis.RedisClient.Set(s.ctx, newAccessKey, userID, token.Maker.AccessKeyDuration).Err(); err != nil {
		return nil, fmt.Errorf("failed to store new  access_token : %w", err)
	}

	// 更新用户 token 集合
	userTokensKey := fmt.Sprintf("auth:user_tokens:%d", userID)
	redis.RedisClient.SAdd(s.ctx, userTokensKey, newAccessKey)
	return &auth.RefreshTokenResp{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
