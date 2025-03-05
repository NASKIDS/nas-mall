// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"

	"github.com/naskids/nas-mall/app/checkout/biz/dal/redis"
	"github.com/naskids/nas-mall/app/checkout/infra/mq"
	"github.com/naskids/nas-mall/app/checkout/infra/rpc"
	checkout "github.com/naskids/nas-mall/rpc_gen/kitex_gen/checkout"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/email"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/order"
	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

const (
	// 锁前缀和过期时间
	LockPrefix     = "lock:order:cancel:"
	LockExpiration = 10 * time.Second
)

/*
	Run

1. get cart
2. calculate cart
3. create order
4. empty cart
5. pay
6. change order result
7. finish
*/
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	// Finish your business logic.
	// Idempotent
	// get cart
	// cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
	// if err != nil {
	// 	klog.Error(err)
	// 	err = fmt.Errorf("GetCart.err:%v", err)
	// 	return
	// }
	// if cartResult == nil || cartResult.Cart == nil || len(cartResult.Cart.Items) == 0 {
	// 	err = errors.New("cart is empty")
	// 	return
	// }
	var (
		oi    []*order.OrderItem
		total float32
	)
	// for _, cartItem := range cartResult.Cart.Items {
	// 	productResp, resultErr := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{Id: cartItem.ProductId})
	// 	if resultErr != nil {
	// 		klog.Error(resultErr)
	// 		err = resultErr
	// 		return
	// 	}
	// 	if productResp.Product == nil {
	// 		continue
	// 	}
	// 	p := productResp.Product
	// 	cost := p.Price * float32(cartItem.Quantity)
	// 	total += cost
	// 	oi = append(oi, &order.OrderItem{
	// 		Item: &cart.CartItem{ProductId: cartItem.ProductId, Quantity: cartItem.Quantity},
	// 		Cost: cost,
	// 	})
	// }
	// create order
	orderReq := &order.PlaceOrderReq{
		UserId:       req.UserId,
		UserCurrency: "USD",
		OrderItems:   oi,
		Email:        req.Email,
	}
	if req.Address != nil {
		addr := req.Address
		zipCodeInt, _ := strconv.Atoi(addr.ZipCode)
		orderReq.Address = &order.Address{
			StreetAddress: addr.StreetAddress,
			City:          addr.City,
			Country:       addr.Country,
			State:         addr.State,
			ZipCode:       int32(zipCodeInt),
		}
	}

	orderResult, err := rpc.OrderClient.PlaceOrder(s.ctx, orderReq)
	if err != nil {
		err = fmt.Errorf("PlaceOrder.err:%v", err)
		return
	}
	klog.Info("orderResult", orderResult)

	var orderId string
	if orderResult != nil && orderResult.Order != nil {
		orderId = orderResult.Order.OrderId
	}
	// publish order finish message
	mq.PublishOrderCancelMessage(orderId, req.UserId, 120)
	// charge
	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
		},
	}
	mq.PublishPaymentCancelMessage(orderId, req.UserId, 120)
	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
	fmt.Printf("paymentResult:%v, %v", paymentResult, err)
	if err != nil {
		err = fmt.Errorf("charge.err:%v", err)
		return
	}
	// 加锁
	lockKey := LockPrefix + orderId
	locked, err := redis.RedisClient.SetNX(s.ctx, lockKey, time.Now().String(), LockExpiration).Result()
	if err != nil {
		return nil, err
	}

	if !locked {
		log.Printf("支付[%s]正在被其他实例处理", orderId)
		return nil, errors.New("payment is canceling") // 其他实例正在处理，订单已取消
	}

	// 确保锁释放
	defer redis.RedisClient.Del(s.ctx, lockKey)
	rpc.PaymentClient.CreatePaymentLog(s.ctx, &payment.CreatePaymentLogReq{
		UserId:        req.UserId,
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
		Amount:        total,
	})

	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@example.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You just created an order in CloudWeGo shop",
		Content:     "You just created an order in CloudWeGo shop",
	})
	msg := &nats.Msg{Subject: "email", Data: data, Header: make(nats.Header)}

	// otel inject
	otel.GetTextMapPropagator().Inject(s.ctx, propagation.HeaderCarrier(msg.Header))

	_ = mq.Nc.PublishMsg(msg)

	klog.Info(paymentResult)
	// change order state
	klog.Info(orderResult)

	_, err = rpc.OrderClient.MarkOrderPaid(s.ctx, &order.MarkOrderPaidReq{UserId: req.UserId, OrderId: orderId})
	if err != nil {
		klog.Error(err)
		return
	}
	// empty cart
	// emptyResult, err := rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
	// if err != nil {
	// 	err = fmt.Errorf("EmptyCart.err:%v", err)
	// 	return
	// }
	// klog.Info(emptyResult)
	resp = &checkout.CheckoutResp{
		OrderId:       orderId,
		TransactionId: paymentResult.TransactionId,
	}
	return
}
