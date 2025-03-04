package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	user "github.com/naskids/nas-mall/rpc_gen/kitex_gen/user"
)

type LogoutService struct {
	ctx context.Context
} // NewLogoutService new LogoutService
func NewLogoutService(ctx context.Context) *LogoutService {
	return &LogoutService{ctx: ctx}
}

// Run create note info
func (s *LogoutService) Run(req *user.LogoutReq) (resp *user.LogoutResp, err error) {
	// Finish your business logic.
	klog.Infof("LogoutReq:%+v", req)
	return &user.LogoutResp{}, nil
}
