package service

import (
	"context"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/naskids/nas-mall/app/user/biz/model"
	user "github.com/naskids/nas-mall/rpc_gen/kitex_gen/user"

	"github.com/naskids/nas-mall/app/user/biz/dal/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserService struct {
	ctx context.Context
} // NewUpdateUserService new UpdateUserService
func NewUpdateUserService(ctx context.Context) *UpdateUserService {
	return &UpdateUserService{ctx: ctx}
}

// Run create note info
func (s *UpdateUserService) Run(req *user.UpdateUserReq) (resp *user.UpdateUserResp, err error) {
	// Finish your business logic.
	klog.Infof("UpdateUserReq:%+v", req)
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
	updates := make(map[string]interface{})
	// 更新密码字段
	if req.NewPassword != nil {
		newPassword := *req.NewPassword // 解析
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			klog.Errorf("Password hashing error: %v", err)
			return nil, err
		}
		updates["password_hashed"] = string(hashedPassword)
		userRow.PasswordHashed = string(hashedPassword)
	}

	// 更新邮箱字段
	if req.NewEmail != nil {
		userRow.Email = *req.NewEmail
		updates["email"] = *req.NewEmail
	}

	// 执行数据库更新
	if err := model.UpdateUser(mysql.DB, s.ctx, req.UserId, updates); err != nil {
		klog.Errorf("Update user error: %v", err)
		return nil, err
	}
	return &user.UpdateUserResp{
		UserId:   userRow.ID,
		NewEmail: userRow.Email,
	}, nil
}
