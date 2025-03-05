package order

import (
	"context"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"

	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() orderservice.Client
	Service() string
	PlaceOrder(ctx context.Context, Req *order.PlaceOrderReq, callOptions ...callopt.Option) (r *order.PlaceOrderResp, err error)
	ListOrder(ctx context.Context, Req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error)
	MarkOrderPaid(ctx context.Context, Req *order.MarkOrderPaidReq, callOptions ...callopt.Option) (r *order.MarkOrderPaidResp, err error)
	MarkOrderCanceled(ctx context.Context, Req *order.MarkOrderCanceledReq, callOptions ...callopt.Option) (r *order.MarkOrderCanceledResp, err error)
	DeleteOrder(ctx context.Context, Req *order.DeleteOrderReq, callOptions ...callopt.Option) (r *order.DeleteOrderResp, err error)
	GetOrderByID(ctx context.Context, Req *order.GetOrderReq, callOptions ...callopt.Option) (r *order.GetOrderResp, err error)
	GetOrderStatus(ctx context.Context, Req *order.GetOrderStatusReq, callOptions ...callopt.Option) (r *order.GetOrderStatusResp, err error)
	ScheduledOrderCancel(ctx context.Context, Req *order.ScheduledOrderCancelReq, callOptions ...callopt.Option) (r *order.ScheduledOrderCancelResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := orderservice.NewClient(dstService, opts...)
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
	kitexClient orderservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() orderservice.Client {
	return c.kitexClient
}

func (c *clientImpl) PlaceOrder(ctx context.Context, Req *order.PlaceOrderReq, callOptions ...callopt.Option) (r *order.PlaceOrderResp, err error) {
	return c.kitexClient.PlaceOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) ListOrder(ctx context.Context, Req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error) {
	return c.kitexClient.ListOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) MarkOrderPaid(ctx context.Context, Req *order.MarkOrderPaidReq, callOptions ...callopt.Option) (r *order.MarkOrderPaidResp, err error) {
	return c.kitexClient.MarkOrderPaid(ctx, Req, callOptions...)
}

func (c *clientImpl) MarkOrderCanceled(ctx context.Context, Req *order.MarkOrderCanceledReq, callOptions ...callopt.Option) (r *order.MarkOrderCanceledResp, err error) {
	return c.kitexClient.MarkOrderCanceled(ctx, Req, callOptions...)
}

func (c *clientImpl) DeleteOrder(ctx context.Context, Req *order.DeleteOrderReq, callOptions ...callopt.Option) (r *order.DeleteOrderResp, err error) {
	return c.kitexClient.DeleteOrder(ctx, Req, callOptions...)
}

func (c *clientImpl) GetOrderByID(ctx context.Context, Req *order.GetOrderReq, callOptions ...callopt.Option) (r *order.GetOrderResp, err error) {
	return c.kitexClient.GetOrderByID(ctx, Req, callOptions...)
}

func (c *clientImpl) GetOrderStatus(ctx context.Context, Req *order.GetOrderStatusReq, callOptions ...callopt.Option) (r *order.GetOrderStatusResp, err error) {
	return c.kitexClient.GetOrderStatus(ctx, Req, callOptions...)
}

func (c *clientImpl) ScheduledOrderCancel(ctx context.Context, Req *order.ScheduledOrderCancelReq, callOptions ...callopt.Option) (r *order.ScheduledOrderCancelResp, err error) {
	return c.kitexClient.ScheduledOrderCancel(ctx, Req, callOptions...)
}
