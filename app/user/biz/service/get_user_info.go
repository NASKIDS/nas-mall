package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/naskids/nas-mall/app/user/biz/model"
	user "github.com/naskids/nas-mall/rpc_gen/kitex_gen/user"
	"golang.org/x/crypto/bcrypt"

	"github.com/naskids/nas-mall/app/user/biz/dal/mysql"
)

type GetUserInfoService struct {
	ctx context.Context
} // NewGetUserInfoService new GetUserInfoService
func NewGetUserInfoService(ctx context.Context) *GetUserInfoService {
	return &GetUserInfoService{ctx: ctx}
}

/*
// 获取用户身份信息
message GetUserInfoReq {
    uint64 user_id = 1;
    string password = 2;  // 或使用 password 字段验证
}

message GetUserInfoResp {
    uint64 user_id = 1;
    string email = 2;
    uint64 created_at = 3;
    uint64 updated_at = 4;
    uint64 deleted_at = 5;
}
*/
// Run create note info
func (s *GetUserInfoService) Run(req *user.GetUserInfoReq) (resp *user.GetUserInfoResp, err error) {
	// Finish your business logic.
	klog.Infof("GetUserInfoReq:%+v", req)
	// 根据 Id 查询用户
	userRow, err := model.GetById(mysql.DB, s.ctx, req.UserId)
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
	return
}
