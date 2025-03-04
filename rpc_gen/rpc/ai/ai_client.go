package ai

import (
	"context"

	ai "github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"

	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/ai/aimodelservice"
)

type RPCClient interface {
	KitexClient() aimodelservice.Client
	Service() string
	QueryOrder(ctx context.Context, Req *ai.QueryOrderRequest, callOptions ...callopt.Option) (r *order.ListOrderResp, err error)
	SimulateAutoOrder(ctx context.Context, Req *ai.SimulateAutoOrderRequest, callOptions ...callopt.Option) (stream aimodelservice.AiModelService_SimulateAutoOrderClient, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := aimodelservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient aimodelservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() aimodelservice.Client {
	return c.kitexClient
}

func (c *clientImpl) QueryOrder(ctx context.Context, Req *ai.QueryOrderRequest, callOptions ...callopt.Option) (r *order.ListOrderResp, err error) {
	return c.kitexClient.QueryOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) SimulateAutoOrder(ctx context.Context, Req *ai.SimulateAutoOrderRequest, callOptions ...callopt.Option) (stream aimodelservice.AiModelService_SimulateAutoOrderClient, err error) {
	return c.kitexClient.SimulateAutoOrder(ctx, Req, callOptions...)
}
