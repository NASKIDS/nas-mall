package auth

import (
	"context"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"

	auth "github.com/naskids/nas-mall/rpc_gen/kitex_gen/auth"
)

func DeliverToken(ctx context.Context, req *auth.DeliverTokenReq, callOptions ...callopt.Option) (resp *auth.DeliveryResp, err error) {
	resp, err = defaultClient.DeliverToken(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeliverToken call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func RefreshToken(ctx context.Context, req *auth.RefreshTokenReq, callOptions ...callopt.Option) (resp *auth.RefreshTokenResp, err error) {
	resp, err = defaultClient.RefreshToken(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "RefreshToken call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq, callOptions ...callopt.Option) (resp *auth.VerifyResp, err error) {
	resp, err = defaultClient.VerifyTokenByRPC(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "VerifyTokenByRPC call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
