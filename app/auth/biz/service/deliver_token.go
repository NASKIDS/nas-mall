package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"

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

	var accessToken string
	var refreshToken string
	// 2. 检查是否存在有效 refresh_token
	refreshKey := fmt.Sprintf("auth:refresh:%d", user.UserID)
	existingRefresh, err := redis.RedisClient.Get(s.ctx, refreshKey).Result()
	if err == nil {
		// 验证旧 refresh_token 是否有效
		if claims, err := token.Maker.ParseRefreshToken(existingRefresh); err == nil {
			// 检查版本号是否匹配且未过期
			if version, ok := claims["ver"].(float64); ok && uint64(version) == user.RefreshVersion {
				// 复用旧 refresh_token
				refreshToken = existingRefresh
			}
		}
	}

	accessToken, err = token.Maker.GenerateAccessToken(utils.H{"uid": user.UserID, "rol": user.Role})
	// 3. 需要生成新令牌的情况
	if refreshToken == "" {
		if err != nil {
			return nil, fmt.Errorf("access token gen err: [%w]", err)
		}

		refreshToken, err = token.Maker.GenerateRefreshToken(utils.H{
			"uid": user.UserID,
			"rol": user.Role,
			"ver": user.RefreshVersion + 1,
		})
		if err != nil {
			return nil, fmt.Errorf("refresh token gen err: [%w]", err)
		}

		// 存储新 refresh_token 并设置过期时间（与令牌有效期一致）
		if err := redis.RedisClient.Set(s.ctx, refreshKey, refreshToken, token.Maker.RefreshKeyDuration).Err(); err != nil {
			klog.Errorf("存储 refresh_token 失败: %v", err)
		}
		err = model.UpdateRefreshVersion(s.ctx, mysql.DB, redis.RedisClient, user.UserID, user.RefreshVersion+1)
		if err != nil {
			return nil, fmt.Errorf("failed to update token: [%w]", err)
		}
	}

	// 4. 持久化 access_token
	accessKey := fmt.Sprintf("auth:access:%s", accessToken)
	if err := redis.RedisClient.Set(s.ctx, accessKey, user.UserID, token.Maker.AccessKeyDuration).Err(); err != nil {
		return nil, fmt.Errorf("存储 access_token 失败: %w", err)
	}
	userTokensKey := fmt.Sprintf("auth:user_tokens:%d", user.UserID)
	redis.RedisClient.SAdd(s.ctx, userTokensKey, accessKey)
	return &auth.DeliveryTokenResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
