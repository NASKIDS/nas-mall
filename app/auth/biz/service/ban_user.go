package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"

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
	ids := req.UserIds
	for _, id := range ids {
		// 加入用户黑名单
		klog.Info(id)
		// 删除token白名单

		// 删除用户 auth 信息
	}
	return
}
