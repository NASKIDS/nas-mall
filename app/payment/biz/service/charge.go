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
	"math/rand"
	"strconv"
	"time"

	creditcard "github.com/durango/go-credit-card"
	"github.com/google/uuid"

	"github.com/naskids/nas-mall/rpc_gen/kitex_gen/payment"
)

type ChargeService struct {
	ctx context.Context
} // NewChargeService new ChargeService
func NewChargeService(ctx context.Context) *ChargeService {
	return &ChargeService{ctx: ctx}
}

// Run create note info
func (s *ChargeService) Run(req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {

	card := creditcard.Card{
		Number: req.CreditCard.CreditCardNumber,
		Cvv:    strconv.Itoa(int(req.CreditCard.CreditCardCvv)),
		Month:  strconv.Itoa(int(req.CreditCard.CreditCardExpirationMonth)),
		Year:   strconv.Itoa(int(req.CreditCard.CreditCardExpirationYear)),
	}

	err = card.Validate(true)
	// fmt.Printf("err:%v", err)

	// fmt.Printf("kerrors.NewBizStatusError(400, err.Error()):%v", kerrors.NewBizStatusError(400, err.Error()))
	if err != nil {
		return nil, err
	}
	payErr := randomPay()
	// fmt.Printf("payErr:%v", payErr)
	if payErr != nil {
		return nil, payErr
	}

	translationId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("translationId:%v", translationId)

	// err = model.CreatePaymentLog(mysql.DB, s.ctx, &model.PaymentLog{
	// 	UserId:        req.UserId,
	// 	OrderId:       req.OrderId,
	// 	TransactionId: translationId.String(),
	// 	Amount:        req.Amount,
	// 	PayAt:         time.Now(),
	// })
	// if err != nil {
	// 	return nil, err
	// }
	return &payment.ChargeResp{TransactionId: translationId.String()}, nil
}
func randomPay() error {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(100)
	if random < 50 {
		return errors.New("支付失败")
	}
	return nil
}
