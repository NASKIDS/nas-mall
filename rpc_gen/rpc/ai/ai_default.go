package ai

import (
	"context"

	ai "github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai/aimodelservice"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
)

func QueryOrder(ctx context.Context, req *ai.QueryOrderRequest, callOptions ...callopt.Option) (resp *order.ListOrderResp, err error) {
	resp, err = defaultClient.QueryOrder(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "QueryOrder call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SimulateAutoOrder(ctx context.Context, Req *ai.SimulateAutoOrderRequest, callOptions ...callopt.Option) (stream aimodelservice.AiModelService_SimulateAutoOrderClient, err error) {
	stream, err = defaultClient.SimulateAutoOrder(ctx, Req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SimulateAutoOrder call failed,err =%+v", err)
		return nil, err
	}
	return stream, nil
}
