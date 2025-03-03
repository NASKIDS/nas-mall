package model

import (
	"context"

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

// TODO impl
// 要存的东西，access token， 角色，refresh key ，refresh key 的版本, refresh key 黑名单（缓存）， access Key 白名单

func GetUser(ctx context.Context, db *gorm.DB, cache *redis.Client, userId uint64) (user AuthUser, err error) {
	return AuthUser{
		UserID:         3,
		Role:           "visitor",
		RefreshVersion: 0,
	}, nil
}

func UpdateRefreshVersion(ctx context.Context, db *gorm.DB, cache *redis.Client, userId uint64, version uint64) error {
	// TODO implement me
	panic("implement me")
}

func GetRefreshVersion(ctx context.Context, db *gorm.DB, cache *redis.Client, userId uint64) (uint64, error) {
	// TODO implement me
	panic("implement me")
}
