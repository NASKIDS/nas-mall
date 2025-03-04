package service

import (
	"context"

	ai "github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai"
)

type SimulateAutoOrderService struct {
	ctx context.Context
}

// NewSimulateAutoOrderService new SimulateAutoOrderService
func NewSimulateAutoOrderService(ctx context.Context) *SimulateAutoOrderService {
	return &SimulateAutoOrderService{ctx: ctx}
}

func (s *SimulateAutoOrderService) Run(req *ai.SimulateAutoOrderRequest, stream ai.AiModelService_SimulateAutoOrderServer) (err error) {
	return
}
