package model

import "sync"

type AuthUserStore interface {
	GetUser(id uint64) (user User, err error)
	UpdateRefreshVersion(id uint64, i any) error
	GetRefreshVersion(id uint64) (uint64, error)
}

type User struct {
	ID             uint64
	Role           string
	RefreshVersion uint64
}

var (
	once             sync.Once
	defaultUserStore AuthUserStore
)

// TODO impl
// 要存的东西，access token， 角色，refresh key ，refresh key 的版本, refresh key 黑名单（缓存）， access Key 白名单

func DefaultAuthUserStore() AuthUserStore {
	once.Do(func() {
		defaultUserStore = new(AuthUserStoreImpl)
	})
	return defaultUserStore
}

var _ AuthUserStore = &AuthUserStoreImpl{}

type AuthUserStoreImpl struct {
}

func (a *AuthUserStoreImpl) GetUser(id uint64) (user User, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthUserStoreImpl) UpdateRefreshVersion(id uint64, i any) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthUserStoreImpl) GetRefreshVersion(id uint64) (uint64, error) {
	//TODO implement me
	panic("implement me")
}
