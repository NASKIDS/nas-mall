package service

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/redis/go-redis/v9"

	redisClient "github.com/naskids/nas-mall/app/payment/biz/dal/redis"
	"github.com/naskids/nas-mall/app/payment/infra/rpc"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
	payment "github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
)

type CancelChargeService struct {
	ctx context.Context
} // NewCancelChargeService new CancelChargeService
func NewCancelChargeService(ctx context.Context) *CancelChargeService {
	return &CancelChargeService{ctx: ctx}
}

// Run 取消订单
func (s *CancelChargeService) Run(req *payment.CancelChargeReq) (resp *payment.CancelChargeResp, err error) {
	// 参数验证
	if req.OrderId == "" {
		return nil, kerrors.NewBizStatusError(400, "订单ID不能为空")
	}
	if req.UserId == 0 {
		return nil, kerrors.NewBizStatusError(400, "用户ID不能为空")
	}

	// 构建Redis锁的key
	lockKey := fmt.Sprintf("payment:cancel:lock:%s", req.OrderId)
	// 尝试获取锁，过期时间10秒
	success, err := redisClient.RedisClient.SetNX(s.ctx, lockKey, "1", 10*time.Second).Result()
	if err != nil {
		return nil, kerrors.NewBizStatusError(500, "获取锁失败: "+err.Error())
	}
	if !success {
		return nil, kerrors.NewBizStatusError(429, "操作太频繁，请稍后再试")
	}

	// 确保锁会被释放
	defer redisClient.RedisClient.Del(s.ctx, lockKey)

	//查询订单
	orderResp, err := rpc.OrderClient.GetOrderStatus(s.ctx, &order.GetOrderStatusReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
	})
	if err != nil {
		if err == redis.Nil || err.Error() == "record not found" {
			return nil, kerrors.NewBizStatusError(404, "订单不存在")
		}
		return nil, kerrors.NewBizStatusError(500, "查询订单失败: "+err.Error())
	}

	// 检查订单状态
	if orderResp.Status == "paid" {
		// 已支付的订单不能取消
		return &payment.CancelChargeResp{Success: false}, nil
	}

	// 更新订单状态为已取消
	_, err = rpc.OrderClient.MarkOrderCanceled(s.ctx, &order.MarkOrderCanceledReq{
		UserId:  req.UserId,
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, kerrors.NewBizStatusError(500, "取消订单失败: "+err.Error())
	}

	return &payment.CancelChargeResp{Success: true}, nil
}
