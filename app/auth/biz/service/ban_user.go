package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"

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
	bannedCount := int32(0) // 初始化封禁计数器

	for _, id := range req.UserIds {
		// 1. 加入用户黑名单
		if err := redis.RedisClient.SAdd(s.ctx, "auth:user_blacklist", id).Err(); err != nil {
			klog.Errorf("加入黑名单失败 user_id=%d: %v", id, err)
			continue // 跳过无法加入黑名单的用户
		}

		// 2. 删除该用户所有 access token 白名单
		userTokensKey := fmt.Sprintf("auth:user_tokens:%d", id)
		tokens, err := redis.RedisClient.SMembers(s.ctx, userTokensKey).Result()
		if err != nil {
			klog.Errorf("获取用户 token 列表失败 user_id=%d: %v", id, err)
		} else if len(tokens) > 0 {
			// 批量删除所有 access token
			if err := redis.RedisClient.Del(s.ctx, tokens...).Err(); err != nil {
				klog.Errorf("删除 access token 失败 user_id=%d: %v", id, err)
				continue
			}
			// 删除用户 token 集合
			if err := redis.RedisClient.Del(s.ctx, userTokensKey).Err(); err != nil {
				klog.Errorf("删除用户 token 集合失败 user_id=%d: %v", id, err)
				continue
			}
		}

		// 3. 删除用户认证数据
		if err := mysql.DB.Where("user_id = ?", id).Delete(&model.AuthUser{}).Error; err != nil {
			klog.Errorf("删除数据库认证信息失败 user_id=%d: %v", id, err)
			continue
		}
		if err := redis.RedisClient.Del(s.ctx, fmt.Sprintf("auth:user:%d", id)).Err(); err != nil {
			klog.Errorf("删除缓存认证信息失败 user_id=%d: %v", id, err)
			continue
		}
		bannedCount++ // 只有成功加入黑名单才计数
	}

	return &auth.BanUserResp{
		BannedCount: bannedCount, // 返回实际成功封禁的数量
	}, nil
}
