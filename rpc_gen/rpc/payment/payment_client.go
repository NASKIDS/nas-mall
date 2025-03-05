package payment

import (
	"context"

	payment "github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"

	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment/paymentservice"
)

type RPCClient interface {
	KitexClient() paymentservice.Client
	Service() string
	Charge(ctx context.Context, Req *payment.ChargeReq, callOptions ...callopt.Option) (r *payment.ChargeResp, err error)
	CancelCharge(ctx context.Context, Req *payment.CancelChargeReq, callOptions ...callopt.Option) (r *payment.CancelChargeResp, err error)
	CreatePaymentLog(ctx context.Context, Req *payment.CreatePaymentLogReq, callOptions ...callopt.Option) (r *payment.CreatePaymentLogResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := paymentservice.NewClient(dstService, opts...)
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
	kitexClient paymentservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() paymentservice.Client {
	return c.kitexClient
}

func (c *clientImpl) Charge(ctx context.Context, Req *payment.ChargeReq, callOptions ...callopt.Option) (r *payment.ChargeResp, err error) {
	r, err = c.kitexClient.Charge(ctx, Req, callOptions...)
	return
}

func (c *clientImpl) CancelCharge(ctx context.Context, Req *payment.CancelChargeReq, callOptions ...callopt.Option) (r *payment.CancelChargeResp, err error) {
	return c.kitexClient.CancelCharge(ctx, Req, callOptions...)
}

func (c *clientImpl) CreatePaymentLog(ctx context.Context, Req *payment.CreatePaymentLogReq, callOptions ...callopt.Option) (r *payment.CreatePaymentLogResp, err error) {
	return c.kitexClient.CreatePaymentLog(ctx, Req, callOptions...)
}
