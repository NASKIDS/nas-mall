package auth

import (
	"context"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"

	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth/authservice"
)

type RPCClient interface {
	KitexClient() authservice.Client
	Service() string
	DeliverToken(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error)
	RefreshToken(ctx context.Context, Req *auth.RefreshTokenReq, callOptions ...callopt.Option) (r *auth.RefreshTokenResp, err error)
	VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error)
	RefreshIfExpired(ctx context.Context, accessToken, refreshToken string, callOptions ...callopt.Option) (r *auth.RefreshTokenResp, err error)
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

func (c *clientImpl) RefreshIfExpired(ctx context.Context, accessToken, refreshToken string, callOptions ...callopt.Option) (r *auth.RefreshTokenResp, err error) {
	// TODO : 导入公钥
	publicKey := lo.Must(paseto.NewV4AsymmetricPublicKeyFromHex("1eb9dbbbbc047c03fd70604e0071f0987e16b28b757225c11f00415d0e20b1a2"))

	if err != nil {
		return nil, err
	}
	require.JSONEq(t,
		"{\"data\":\"this is a signed message\",\"exp\":\"2022-01-01T00:00:00+00:00\"}",
		string(token.ClaimsJSON()),
	)
	require.Equal(t,
		"{\"kid\":\"zVhMiPBP9fRf2snEcT7gFTioeA9COcNy9DfgL1W60haN\"}",
		string(token.Footer()),
	)
	require.NoError(t, err)

}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() authservice.Client {
	return c.kitexClient
}

func (c *clientImpl) DeliverToken(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error) {
	return c.kitexClient.DeliverToken(ctx, Req, callOptions...)
}

func (c *clientImpl) RefreshToken(ctx context.Context, Req *auth.RefreshTokenReq, callOptions ...callopt.Option) (r *auth.RefreshTokenResp, err error) {
	return c.kitexClient.RefreshToken(ctx, Req, callOptions...)
}

func (c *clientImpl) VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error) {
	return c.kitexClient.VerifyTokenByRPC(ctx, Req, callOptions...)
}
