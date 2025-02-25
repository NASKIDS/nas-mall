package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/naskids/nas-mall/app/user/biz/model"
	user "github.com/naskids/nas-mall/rpc_gen/kitex_gen/user"
	"golang.org/x/crypto/bcrypt"

	"github.com/naskids/nas-mall/app/user/biz/dal/mysql"
)

type DeleteUserService struct {
	ctx context.Context
} // NewDeleteUserService new DeleteUserService
func NewDeleteUserService(ctx context.Context) *DeleteUserService {
	return &DeleteUserService{ctx: ctx}
}

/*
	// 删除用户
	message DeleteUserReq {
		uint64 user_id = 1;
		// 输入密码确认删除
		string password = 2;
	}
*/
// Run create note info
func (s *DeleteUserService) Run(req *user.DeleteUserReq) (resp *user.DeleteUserResp, err error) {
	// Finish your business logic.
	klog.Infof("DeleteUserReq:%+v", req)
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
	if err = model.DeleteUser(mysql.DB, s.ctx, req.UserId); err != nil {
		return
	}
	return &user.DeleteUserResp{}, nil
}
