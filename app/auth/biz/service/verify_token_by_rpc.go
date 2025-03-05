package service

import (
	"context"
	"fmt"

	"github.com/naskids/nas-mall/app/auth/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/auth/biz/dal/redis"
	"github.com/naskids/nas-mall/app/auth/biz/model"
	"github.com/naskids/nas-mall/common/token"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type VerifyTokenByRPCService struct {
	ctx context.Context
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyTokenResp, err error) {
	claims, err := token.Maker.ParseAccessToken(req.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("remote failed to verify token: [%w]", err)
	}
	// 检查白名单
	exists, _ := redis.RedisClient.Exists(s.ctx, fmt.Sprintf("auth:access:%s", req.AccessToken)).Result()
	if exists == 0 {
		return &auth.VerifyTokenResp{IsValid: false}, nil
	}

	userId := claims["uid"].(float64)
	role := claims["rol"].(string)

	// 检查用户黑名单
	isBanned, _ := redis.RedisClient.SIsMember(s.ctx, "auth:user_blacklist", userId).Result()
	if isBanned {
		return &auth.VerifyTokenResp{IsValid: false}, nil
	}

	// 从存储中获取 user 最新信息
	user, err := model.GetUser(s.ctx, mysql.DB, redis.RedisClient, uint64(userId))
	if err != nil {
		return nil, fmt.Errorf("get user by id failed: [%w]", err)
	}
	// 角色信息一致
	return &auth.VerifyTokenResp{
		IsValid: role == user.Role,
	}, nil
}
