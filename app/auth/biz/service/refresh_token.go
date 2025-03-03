package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/utils"

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

	userID, tokenVersion, role := tk["uid"].(uint64), tk["ver"].(uint64), tk["rol"].(string)
	// 2. 验证用户
	user, err := model.GetUser(s.ctx, mysql.DB, redis.RedisClient, userID)
	if err != nil {
		return nil, err
	}
	if user.ID != userID || user.Role != role {
		return nil, fmt.Errorf("incorrect user info: [%w]", err)
	}
	var currentVersion uint64
	currentVersion, err = model.GetRefreshVersion(s.ctx, mysql.DB, redis.RedisClient, userID)
	if err != nil || currentVersion != tokenVersion {
		return nil, fmt.Errorf("stale refresh token: [%w]", err)
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
	err = model.UpdateRefreshVersion(s.ctx, mysql.DB, redis.RedisClient, userID, currentVersion+1)
	if err != nil {
		return nil, fmt.Errorf("failed to update token: [%w]", err)
	}

	// TODO  4. 持久化token 到 redis
	return &auth.RefreshTokenResp{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
