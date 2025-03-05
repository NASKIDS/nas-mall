package service

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/kitex/pkg/kerrors"

	"github.com/naskids/nas-mall/app/payment/biz/dal/redis"
	"github.com/naskids/nas-mall/app/payment/infra/rpc"
	orderClient "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
	payment "github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
)

type CancelChargeService struct {
	ctx context.Context
} // NewCancelChargeService new CancelChargeService
func NewCancelChargeService(ctx context.Context) *CancelChargeService {
	return &CancelChargeService{ctx: ctx}
}

const (
	PaymentCancelLockPrefix = "lock:order:cancel:"
	LockExpiration          = 10 * time.Second
)

// Run 取消订单
func (s *CancelChargeService) Run(req *payment.CancelChargeReq) (resp *payment.CancelChargeResp, err error) {
	// 参数验证
	if req.OrderId == "" {
		return nil, kerrors.NewBizStatusError(400, "订单ID不能为空")
	}
	if req.UserId == 0 {
		return nil, kerrors.NewBizStatusError(400, "用户ID不能为空")
	}

	lockKey := PaymentCancelLockPrefix + req.OrderId
	locked, err := redis.RedisClient.SetNX(s.ctx, lockKey, time.Now().String(), LockExpiration).Result()
	if err != nil {
		return nil, err
	}

	if !locked {
		log.Printf("支付[%s]正在被其他实例处理", req.OrderId)
		return &payment.CancelChargeResp{Success: true}, nil // 其他实例正在处理，视为成功
	}

	// 确保锁释放
	defer redis.RedisClient.Del(s.ctx, lockKey)

	// 查询订单
	order, err := rpc.OrderClient.GetOrderStatus(s.ctx, &orderClient.GetOrderStatusReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, err
	}

	// 检查订单状态，只有未支付的订单才能取消
	if order.Status != "placed" {
		log.Printf("订单[%s]当前状态为[%s]，无需取消", req.OrderId, order.Status)
		return &payment.CancelChargeResp{Success: true}, nil // 非未支付状态，无需取消，视为成功
	}

	// 执行订单取消操作
	_, err = rpc.OrderClient.MarkOrderCanceled(s.ctx, &orderClient.MarkOrderCanceledReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
	})
	if err != nil {
		return &payment.CancelChargeResp{Success: false}, err
	}

	log.Printf("订单[%s]已成功取消", req.OrderId)

	return &payment.CancelChargeResp{Success: true}, nil
}
