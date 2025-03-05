package service

import (
	"context"
	"time"

	"github.com/naskids/nas-mall/app/payment/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/payment/biz/model"
	payment "github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
)

type CreatePaymentLogService struct {
	ctx context.Context
} // NewCreatePaymentLogService new CreatePaymentLogService
func NewCreatePaymentLogService(ctx context.Context) *CreatePaymentLogService {
	return &CreatePaymentLogService{ctx: ctx}
}

// Run create note info
func (s *CreatePaymentLogService) Run(req *payment.CreatePaymentLogReq) (resp *payment.CreatePaymentLogResp, err error) {
	// Finish your business logic.
	err = model.CreatePaymentLog(mysql.DB, s.ctx, &model.PaymentLog{
		UserId:        req.UserId,
		OrderId:       req.OrderId,
		TransactionId: req.TransactionId,
		Amount:        req.Amount,
		PayAt:         time.Now(),
	})
	if err != nil {
		return nil, err
	}

	return &payment.CreatePaymentLogResp{}, nil
}
