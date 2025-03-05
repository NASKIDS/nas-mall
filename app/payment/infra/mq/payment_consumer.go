package mq

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/naskids/nas-mall/app/payment/biz/service"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
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
	cancelChargeResp, err := service.NewCancelChargeService(context.Background()).Run(&payment.CancelChargeReq{
		OrderId: cancelMsg.OrderID,
		UserId:  cancelMsg.UserID,
	})
	if err != nil || !cancelChargeResp.Success {
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
