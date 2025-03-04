package auth

import (
	"context"

	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth/authservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() authservice.Client
	Service() string
	DeliverToken(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryTokenResp, err error)
	RefreshToken(ctx context.Context, Req *auth.RefreshTokenReq, callOptions ...callopt.Option) (r *auth.RefreshTokenResp, err error)
	VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyTokenResp, err error)
	BanUser(ctx context.Context, Req *auth.BanUserReq, callOptions ...callopt.Option) (r *auth.BanUserResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := authservice.NewClient(dstService, opts...)
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
	kitexClient authservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() authservice.Client {
	return c.kitexClient
}

func (c *clientImpl) DeliverToken(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryTokenResp, err error) {
	return c.kitexClient.DeliverToken(ctx, Req, callOptions...)
}

func (c *clientImpl) RefreshToken(ctx context.Context, Req *auth.RefreshTokenReq, callOptions ...callopt.Option) (r *auth.RefreshTokenResp, err error) {
	return c.kitexClient.RefreshToken(ctx, Req, callOptions...)
}

func (c *clientImpl) VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyTokenResp, err error) {
	return c.kitexClient.VerifyTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) BanUser(ctx context.Context, Req *auth.BanUserReq, callOptions ...callopt.Option) (r *auth.BanUserResp, err error) {
	return c.kitexClient.BanUser(ctx, Req, callOptions...)
}
