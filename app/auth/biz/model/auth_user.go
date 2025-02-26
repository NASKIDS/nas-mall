package model

type AuthUser interface {
	GetUser(id uint64) (user User, err error)
	UpdateRefreshVersion(id uint64, i any) error
	GetRefreshVersion(id uint64) (uint64, error)
}

type User struct {
	ID             uint64
	RefreshVersion uint64
}

// TODO impl
