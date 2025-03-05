package service

import (
	"context"

	"github.com/cloudwego/kitex/pkg/klog"

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
	klog.Info(req.UserMessage)
	defer func() {
		_ = stream.Close()
	}()
	for i := 0; i < 3; i++ {
		err := stream.Send(&ai.SimulateAutoOrderResponse{AssistantMessage: "hi!"})
		if err != nil {
			return err
		}
	}
	return
}
