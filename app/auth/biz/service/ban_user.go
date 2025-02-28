package service

import (
	"context"

	"github.com/naskids/nas-mall/app/auth/biz/model"
	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

type BanUserService struct {
	ctx       context.Context
	userStore model.AuthUserStore
} // NewBanUserService new BanUserService
func NewBanUserService(ctx context.Context) *BanUserService {
	return &BanUserService{ctx: ctx, userStore: model.DefaultAuthUserStore()}
}

// Run create note info
func (s *BanUserService) Run(req *auth.BanUserReq) (resp *auth.BanUserResp, err error) {
	// Finish your business logic.

	return
}
