package model

type AuthUser interface {
	GetUser(id uint64) (user User, err error)
	UpdateRefreshVersion(id uint64, i any) error
	GetRefreshVersion(id uint64) (uint64, error)
}

type User struct {
	ID             uint64
	Role           string
	RefreshVersion uint64
}

// TODO impl
// 要存的东西，access token 里面的验证码，角色，refresh key 里面的验证码，refresh key 的版本
