package service

import (
	"context"
	"fmt"

	"github.com/naskids/nas-mall/app/auth/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/auth/biz/dal/redis"
	"github.com/naskids/nas-mall/app/auth/biz/model"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type BanUserService struct {
	ctx context.Context
} // NewBanUserService new BanUserService
func NewBanUserService(ctx context.Context) *BanUserService {
	return &BanUserService{ctx: ctx}
}

// Run create note info
func (s *BanUserService) Run(req *auth.BanUserReq) (resp *auth.BanUserResp, err error) {
	for _, id := range req.UserIds {
		// 加入黑名单
		redis.RedisClient.SAdd(s.ctx, "auth:user_blacklist", id)

		// 删除所有关联的 access token
		userTokensKey := fmt.Sprintf("auth:user_tokens:%d", id)
		tokens, _ := redis.RedisClient.SMembers(s.ctx, userTokensKey).Result()
		if len(tokens) > 0 {
			redis.RedisClient.Del(s.ctx, tokens...)     // 删除白名单中的 token
			redis.RedisClient.Del(s.ctx, userTokensKey) // 删除用户 token 集合
		}

		// 删除用户认证数据
		mysql.DB.Where("user_id = ?", id).Delete(&model.AuthUser{})
		redis.RedisClient.Del(s.ctx, fmt.Sprintf("auth:user:%d", id))
	}
	return &auth.BanUserResp{}, nil
}
