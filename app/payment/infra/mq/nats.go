package mq

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	Nc        *nats.Conn
	JetStream nats.JetStreamContext
	err       error
)

// 初始化NATS连接和JetStream
func Init() {
	// 修改 NATS 连接地址为本地地址
	Nc, err = nats.Connect("nats://127.0.0.1:4222", nats.RetryOnFailedConnect(true))
	if err != nil {
		log.Printf("连接NATS失败: %v", err)
		panic(err)
	}
	log.Println("成功连接到NATS服务器")

	// 创建JetStream上下文
	JetStream, err = Nc.JetStream()
	if err != nil {
		log.Printf("创建JetStream上下文失败: %v", err)
		panic(err)
	}

	// 创建Stream（如果不存在）
	PaymentStreamConfig := &nats.StreamConfig{
		Name:      "PAYMENTS",
		Subjects:  []string{"payment.*"},
		MaxAge:    24 * time.Hour,
		Storage:   nats.FileStorage,
		Retention: nats.InterestPolicy,
		Replicas:  1,
	}

	// 设置较长的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = JetStream.AddStream(PaymentStreamConfig, nats.Context(ctx))
	if err != nil {
		if err == nats.ErrStreamNameAlreadyInUse {
			log.Println("Stream已存在，继续执行")
		} else {
			log.Printf("创建Stream失败: %v", err)
			panic(err)
		}
	}

	log.Println("成功初始化NATS JetStream")
}

const (
	// 主题
	PaymentCancelSubject = "payment.cancel" // 支付取消
)

// 支付取消消息
type PaymentCancelMessage struct {
	OrderID    string `json:"order_id"`
	UserID     uint64 `json:"user_id"`
	CreateTime int64  `json:"create_time"` // 消息创建时间戳
	ExpireTime int64  `json:"expire_time"` // 过期时间戳
}
