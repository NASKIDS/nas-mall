package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/naskids/nas-mall/app/payment/biz/dal/redis"
	"github.com/naskids/nas-mall/app/payment/infra/rpc"
	orderClient "github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
)

const (
	// 锁前缀和过期时间
	PaymentCancelLockPrefix = "lock:order:cancel:"
	LockExpiration          = 10 * time.Second

	// 消费者配置
	PaymentCancelConsumerName = "payment-cancel-consumer"
	PaymentCancelDurable      = "payment-cancel-durable"
)

// InitPaymentCancelConsumer 初始化支付取消消费者
func InitPaymentCancelConsumer() error {
	// 创建JetStream消费者
	_, err := JetStream.QueueSubscribe(
		PaymentCancelSubject,
		"payment-cancel-group",
		handlePaymentCancelMessage,
		nats.DeliverNew(),                  // 只处理新消息
		nats.ManualAck(),                   // 手动确认消息
		nats.AckWait(30*time.Second),       // 确认超时时间
		nats.MaxAckPending(100),            // 最大未确认消息数
		nats.MaxDeliver(-1),                // 最大重试次数
		nats.Durable(PaymentCancelDurable), // 持久化订阅
	)

	if err != nil {
		return err
	}

	log.Println("支付取消消费者已成功启动")
	return nil
}

// 处理订单取消消息
func handlePaymentCancelMessage(msg *nats.Msg) {
	var cancelMsg PaymentCancelMessage
	if err := json.Unmarshal(msg.Data, &cancelMsg); err != nil {
		log.Printf("解析支付取消消息失败: %v", err)
		// 消息格式错误，拒绝消息
		msg.Nak()
		return
	}

	// 检查消息是否已过期（到了取消时间）
	now := time.Now().Unix()
	if now < cancelMsg.ExpireTime {
		log.Printf("支付[%s]未到取消时间，当前: %d, 过期时间: %d，延迟处理",
			cancelMsg.OrderID, now, cancelMsg.ExpireTime)
		// 消息未过期，稍后再处理
		// 使用延迟Nak来在一段时间后重新投递消息
		msg.NakWithDelay(time.Second * 60)
		return
	}

	// 处理订单取消
	if err := processPaymentCancel(cancelMsg); err != nil {
		log.Printf("处理支付[%s]取消失败: %v", cancelMsg.OrderID, err)
		// 处理失败，稍后重试
		msg.NakWithDelay(time.Second * 5)
		return
	}

	// 处理成功，确认消息
	if err := msg.Ack(); err != nil {
		log.Printf("确认消息失败: %v", err)
	}
}

// 处理订单取消逻辑
func processPaymentCancel(cancelMsg PaymentCancelMessage) error {
	ctx := context.Background()

	// 获取分布式锁，防止重复处理
	lockKey := PaymentCancelLockPrefix + cancelMsg.OrderID
	locked, err := redis.RedisClient.SetNX(ctx, lockKey, time.Now().String(), LockExpiration).Result()
	if err != nil {
		return err
	}

	if !locked {
		log.Printf("支付[%s]正在被其他实例处理", cancelMsg.OrderID)
		return nil // 其他实例正在处理，视为成功
	}

	// 确保锁释放
	defer redis.RedisClient.Del(ctx, lockKey)

	// 查询订单
	order, err := rpc.OrderClient.GetOrderStatus(ctx, &orderClient.GetOrderStatusReq{
		UserId:  cancelMsg.UserID,
		OrderId: cancelMsg.OrderID,
	})
	if err != nil {
		return err
	}

	// 检查订单状态，只有未支付的订单才能取消
	if order.Status != "placed" {
		log.Printf("订单[%s]当前状态为[%s]，无需取消", cancelMsg.OrderID, order.Status)
		return nil // 非未支付状态，无需取消，视为成功
	}

	// 执行订单取消操作
	_, err = rpc.OrderClient.MarkOrderCanceled(ctx, &orderClient.MarkOrderCanceledReq{
		UserId:  cancelMsg.UserID,
		OrderId: cancelMsg.OrderID,
	})
	if err != nil {
		return err
	}

	log.Printf("订单[%s]已成功取消，创建时间: %d, 过期时间: %d",
		cancelMsg.OrderID, cancelMsg.CreateTime, cancelMsg.ExpireTime)
	return nil
}
