package service

import (
	"context"
	"errors"
	"time"

	"github.com/naskids/nas-mall/app/auth/biz/model"
	"github.com/naskids/nas-mall/app/auth/utils/token"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type RefreshTokenService struct {
	ctx             context.Context
	tokenMaker      token.Maker
	userStore       model.AuthUser
	tokenDuration   time.Duration
	refreshDuration time.Duration
} // NewRefreshTokenService new RefreshTokenService
func NewRefreshTokenService(ctx context.Context) *RefreshTokenService {
	return &RefreshTokenService{ctx: ctx}
}

// Run create note info
func (s *RefreshTokenService) Run(req *auth.RefreshTokenReq) (resp *auth.RefreshTokenResp, err error) {
	// Finish your business logic.
	// 1. 解析刷新令牌
	userID, tokenVersion, err := s.tokenMaker.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. 验证版本号
	currentVersion, err := s.userStore.GetRefreshVersion(userID)
	if err != nil || currentVersion != tokenVersion {
		return nil, errors.New("stale refresh token")
	}

	// 3. 生成新令牌并更新版本
	newAccess, _ := s.tokenMaker.GenerateAccessToken(userID, s.tokenDuration)
	newRefresh, _ := s.tokenMaker.GenerateRefreshToken(userID, currentVersion+1, s.refreshDuration)

	if err := s.userStore.UpdateRefreshVersion(userID, currentVersion+1); err != nil {
		return nil, errors.New("failed to update token")
	}

	return &auth.RefreshTokenResp{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}
