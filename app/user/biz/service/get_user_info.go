package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"

	"github.com/naskids/nas-mall/app/user/biz/model"
	user "github.com/naskids/nas-mall/rpc_gen/kitex_gen/user"

	"github.com/naskids/nas-mall/app/user/biz/dal/mysql"
)

type GetUserInfoService struct {
	ctx context.Context
} // NewGetUserInfoService new GetUserInfoService
func NewGetUserInfoService(ctx context.Context) *GetUserInfoService {
	return &GetUserInfoService{ctx: ctx}
}

// Run create note info
func (s *GetUserInfoService) Run(req *user.GetUserInfoReq) (resp *user.GetUserInfoResp, err error) {
	// Finish your business logic.
	klog.Infof("GetUserInfoReq:%+v", req)
	// 根据 Id 查询用户
	userRow, err := model.GetFullMessageById(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		return
	}
	// 2. 验证密码
	if req.Password != "" {
		err = bcrypt.CompareHashAndPassword(
			[]byte(userRow.PasswordHashed),
			[]byte(req.Password),
		)
		// 密码验证失败
		if err != nil {
			return
		}
	}
	return &user.GetUserInfoResp{
		UserId:    req.UserId,
		Email:     userRow.Email,
		CreatedAt: uint64(userRow.CreatedAt.Unix()),
		UpdatedAt: uint64(userRow.UpdatedAt.Unix()),
		DeletedAt: uint64(userRow.DeletedAt.Unix()),
	}, nil
}
