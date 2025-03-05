package service

import (
	"context"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/naskids/nas-mall/app/order/biz/dal/mysql"
	"github.com/naskids/nas-mall/app/order/biz/model"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

type GetOrderStatusService struct {
	ctx context.Context
} // NewGetOrderStatusService new GetOrderStatusService
func NewGetOrderStatusService(ctx context.Context) *GetOrderStatusService {
	return &GetOrderStatusService{ctx: ctx}
}

// Run create note info
func (s *GetOrderStatusService) Run(req *order.GetOrderStatusReq) (resp *order.GetOrderStatusResp, err error) {
	// Finish your business logic.
	// 参数验证
	if req.UserId == 0 || req.OrderId == "" {
		err = fmt.Errorf("user_id或order_id不能为空")
		return nil, err
	}

	// 从数据库查询订单
	orderData, err := model.GetOrder(mysql.DB, s.ctx, req.UserId, req.OrderId)
	if err != nil {
		klog.Errorf("model.GetOrder error: %v", err)
		return nil, err
	}

	// 封装订单信息

	resp = &order.GetOrderStatusResp{
		UserId:  req.UserId,
		OrderId: req.OrderId,
		Status:  string(orderData.OrderState),
	}

	return resp, nil

}
