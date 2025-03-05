package payment

import (
	"context"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
	payment "github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
)

func Charge(ctx context.Context, req *payment.ChargeReq, callOptions ...callopt.Option) (resp *payment.ChargeResp, err error) {
	resp, err = defaultClient.Charge(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Charge call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func CancelCharge(ctx context.Context, req *payment.CancelChargeReq, callOptions ...callopt.Option) (resp *payment.CancelChargeResp, err error) {
	resp, err = defaultClient.CancelCharge(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CancelCharge call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func CreatePaymentLog(ctx context.Context, req *payment.CreatePaymentLogReq, callOptions ...callopt.Option) (resp *payment.CreatePaymentLogResp, err error) {
	resp, err = defaultClient.CreatePaymentLog(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "CreatePaymentLog call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
