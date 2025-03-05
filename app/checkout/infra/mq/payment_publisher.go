package mq

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

// PublishPaymentCancelMessage 发布支付取消消息
func PublishPaymentCancelMessage(orderID string, userID uint64, expireSeconds int64) error {
	// 计算时间
	now := time.Now().Unix()
	expireTime := now + expireSeconds

	// 创建订单取消消息
	cancelMsg := PaymentCancelMessage{
		OrderID:    orderID,
		UserID:     userID,
		CreateTime: now,
		ExpireTime: expireTime,
	}

	// 序列化消息
	msgBytes, err := json.Marshal(cancelMsg)
	if err != nil {
		return fmt.Errorf("序列化支付取消消息失败: %w", err)
	}

	// 创建消息对象并设置延迟时间
	msg := &nats.Msg{
		Subject: PaymentCancelSubject,
		Data:    msgBytes,
		Header:  nats.Header{},
	}
	// 设置延迟时间（单位: 秒）
	delayDuration := fmt.Sprintf("%ds", expireSeconds)
	msg.Header.Set("Nats-Msg-Delay", delayDuration)

	// 发布延迟消息到 NATS JetStream
	_, err = JetStream.PublishMsg(msg)
	if err != nil {
		return fmt.Errorf("发布支付取消消息失败: %w", err)
	}

	log.Printf("已发布支付[%s]超时取消消息，创建时间: %d, 过期时间: %d, 延迟: %s",
		orderID, now, expireTime, delayDuration)
	return nil
}
