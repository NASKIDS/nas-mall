package main

import (
	"context"

	"github.com/naskids/nas-mall/app/ai/biz/service"
	ai "github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

// AiModelServiceImpl implements the last service interface defined in the IDL.
type AiModelServiceImpl struct{}

// QueryOrder implements the AiModelServiceImpl interface.
func (s *AiModelServiceImpl) QueryOrder(ctx context.Context, req *ai.QueryOrderRequest) (resp *order.ListOrderResp, err error) {
	resp, err = service.NewQueryOrderService(ctx).Run(req)

	return resp, err
}

func (s *AiModelServiceImpl) SimulateAutoOrder(req *ai.SimulateAutoOrderRequest, stream ai.AiModelService_SimulateAutoOrderServer) (err error) {
	ctx := context.Background()
	err = service.NewSimulateAutoOrderService(ctx).Run(req, stream)
	return
}
