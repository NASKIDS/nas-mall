package model

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/naskids/nas-mall/common"
)

type AuthUser struct {
	common.Model
	UserID         uint64
	Role           string
	RefreshVersion uint64
}

func (u AuthUser) TableName() string {
	return "auth_user"
}

// GetUser 实现：带缓存的用户查询
func GetUser(ctx context.Context, db *gorm.DB, cache *redis.Client, userId uint64) (user AuthUser, err error) {
	cacheKey := fmt.Sprintf("auth:user:%d", userId)

	// 尝试从 Redis 缓存获取
	result, err := cache.HGetAll(ctx, cacheKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return AuthUser{}, fmt.Errorf("cache error: %w", err)
	}
	if len(result) > 0 {
		user.UserID = userId
		user.Role = result["role"]
		user.RefreshVersion, _ = strconv.ParseUint(result["refresh_version"], 10, 64)
		return user, nil
	}

	// 缓存未命中，查询数据库
	err = db.Where("user_id = ?", userId).First(&user).Error
	if err != nil {
		return AuthUser{}, fmt.Errorf("db error: %w", err)
	}

	// 更新 Redis 缓存
	cache.HSet(ctx, cacheKey, map[string]interface{}{
		"user_id":         user.UserID,
		"role":            user.Role,
		"refresh_version": user.RefreshVersion,
	})
	cache.Expire(ctx, cacheKey, time.Hour)

	return user, nil
}

// UpdateRefreshVersion 实现：版本号更新（DB + 缓存）
func UpdateRefreshVersion(ctx context.Context, db *gorm.DB, cache *redis.Client, userId uint64, version uint64) error {
	// 更新数据库
	err := db.Model(&AuthUser{}).
		Where("user_id = ?", userId).
		Update("refresh_version", version).Error
	if err != nil {
		return fmt.Errorf("db update failed: %w", err)
	}

	// 更新缓存
	cacheKey := fmt.Sprintf("auth:user:%d", userId)
	if err := cache.HSet(ctx, cacheKey, "refresh_version", version).Err(); err != nil {
		cache.Del(ctx, cacheKey) // 若缓存更新失败则清除缓存保证一致性
		return fmt.Errorf("cache update failed: %w", err)
	}

	return nil
}

// GetRefreshVersion 实现：优先从缓存获取版本
func GetRefreshVersion(ctx context.Context, db *gorm.DB, cache *redis.Client, userId uint64) (uint64, error) {
	cacheKey := fmt.Sprintf("auth:user:%d", userId)

	// 尝试从缓存获取
	versionStr, err := cache.HGet(ctx, cacheKey, "refresh_version").Result()
	if err == nil {
		return strconv.ParseUint(versionStr, 10, 64)
	}

	// 缓存未命中，查数据库
	var user AuthUser
	err = db.Select("refresh_version").
		Where("user_id = ?", userId).
		First(&user).Error
	if err != nil {
		return 0, fmt.Errorf("db error: %w", err)
	}

	// 回填缓存
	cache.HSet(ctx, cacheKey, "refresh_version", user.RefreshVersion)
	return user.RefreshVersion, nil
}
