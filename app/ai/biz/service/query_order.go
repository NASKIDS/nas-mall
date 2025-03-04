package service

import (
	"context"

	ai "github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

type QueryOrderService struct {
	ctx context.Context
} // NewQueryOrderService new QueryOrderService
func NewQueryOrderService(ctx context.Context) *QueryOrderService {
	return &QueryOrderService{ctx: ctx}
}

// Run create note info
func (s *QueryOrderService) Run(req *ai.QueryOrderRequest) (resp *order.ListOrderResp, err error) {
	// Finish your business logic.

	return
}
