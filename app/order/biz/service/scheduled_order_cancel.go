package service

import (
	"context"
	"fmt"

	"github.com/naskids/nas-mall/app/order/infra/mq"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

type ScheduledOrderCancelService struct {
	ctx context.Context
} // NewScheduledOrderCancelService new ScheduledOrderCancelService
func NewScheduledOrderCancelService(ctx context.Context) *ScheduledOrderCancelService {
	return &ScheduledOrderCancelService{ctx: ctx}
}

// Run create note info
func (s *ScheduledOrderCancelService) Run(req *order.ScheduledOrderCancelReq) (resp *order.ScheduledOrderCancelResp, err error) {
	// Finish your business logic.
	fmt.Println("开始执行定时取消订单")
	if req.UserId == 0 {
		err = fmt.Errorf("user_id不能为空")
		return nil, err
	}
	if req.OrderId == "" {
		err = fmt.Errorf("order_id不能为空")
		return nil, err
	}
	if req.ScheduledTime <= 0 {
		err = fmt.Errorf("scheduled_time不能小于0")
		return nil, err
	}
	err = mq.PublishOrderCancelMessage(req.OrderId, req.UserId, req.ScheduledTime)
	if err != nil {
		err = fmt.Errorf("发布订单取消消息失败: %v", err)
		return nil, err
	}
	fmt.Println("发布订单取消消息成功")

	return &order.ScheduledOrderCancelResp{}, nil
}
