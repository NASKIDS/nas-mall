// Code generated by Kitex v0.9.1. DO NOT EDIT.

package orderservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	order "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	PlaceOrder(ctx context.Context, Req *order.PlaceOrderReq, callOptions ...callopt.Option) (r *order.PlaceOrderResp, err error)
	ListOrder(ctx context.Context, Req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error)
	MarkOrderPaid(ctx context.Context, Req *order.MarkOrderPaidReq, callOptions ...callopt.Option) (r *order.MarkOrderPaidResp, err error)
	MarkOrderCanceled(ctx context.Context, Req *order.MarkOrderCanceledReq, callOptions ...callopt.Option) (r *order.MarkOrderCanceledResp, err error)
	DeleteOrder(ctx context.Context, Req *order.DeleteOrderReq, callOptions ...callopt.Option) (r *order.DeleteOrderResp, err error)
	GetOrderByID(ctx context.Context, Req *order.GetOrderReq, callOptions ...callopt.Option) (r *order.GetOrderResp, err error)
	GetOrderStatus(ctx context.Context, Req *order.GetOrderStatusReq, callOptions ...callopt.Option) (r *order.GetOrderStatusResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kOrderServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kOrderServiceClient struct {
	*kClient
}

func (p *kOrderServiceClient) PlaceOrder(ctx context.Context, Req *order.PlaceOrderReq, callOptions ...callopt.Option) (r *order.PlaceOrderResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PlaceOrder(ctx, Req)
}

func (p *kOrderServiceClient) ListOrder(ctx context.Context, Req *order.ListOrderReq, callOptions ...callopt.Option) (r *order.ListOrderResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListOrder(ctx, Req)
}

func (p *kOrderServiceClient) MarkOrderPaid(ctx context.Context, Req *order.MarkOrderPaidReq, callOptions ...callopt.Option) (r *order.MarkOrderPaidResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MarkOrderPaid(ctx, Req)
}

func (p *kOrderServiceClient) MarkOrderCanceled(ctx context.Context, Req *order.MarkOrderCanceledReq, callOptions ...callopt.Option) (r *order.MarkOrderCanceledResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MarkOrderCanceled(ctx, Req)
}

func (p *kOrderServiceClient) DeleteOrder(ctx context.Context, Req *order.DeleteOrderReq, callOptions ...callopt.Option) (r *order.DeleteOrderResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeleteOrder(ctx, Req)
}

func (p *kOrderServiceClient) GetOrderByID(ctx context.Context, Req *order.GetOrderReq, callOptions ...callopt.Option) (r *order.GetOrderResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetOrderByID(ctx, Req)
}

func (p *kOrderServiceClient) GetOrderStatus(ctx context.Context, Req *order.GetOrderStatusReq, callOptions ...callopt.Option) (r *order.GetOrderStatusResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetOrderStatus(ctx, Req)
}
