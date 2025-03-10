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

package model

import (
	"context"

	"gorm.io/gorm"

	"github.com/naskids/nas-mall/common"
)

type Consignee struct {
	Email string

	StreetAddress string
	City          string
	State         string
	Country       string
	ZipCode       int32
}

type OrderState string

const (
	OrderStatePlaced   OrderState = "placed"
	OrderStatePaid     OrderState = "paid"
	OrderStateCanceled OrderState = "canceled"
)

type Order struct {
	common.Model
	OrderId      string `gorm:"uniqueIndex;size:256"`
	UserId       uint64
	UserCurrency string
	Consignee    Consignee   `gorm:"embedded"`
	OrderItems   []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"`
	OrderState   OrderState
}

func (o Order) TableName() string {
	return "order"
}

func ListOrder(db *gorm.DB, ctx context.Context, userId uint64) (orders []Order, err error) {
	err = db.Model(&Order{}).Where(&Order{UserId: userId}).Preload("OrderItems").Find(&orders).Error
	return
}

func GetOrder(db *gorm.DB, ctx context.Context, userId uint64, orderId string) (order Order, err error) {
	err = db.Where(&Order{UserId: userId, OrderId: orderId}).First(&order).Error
	return
}

func UpdateOrderState(db *gorm.DB, ctx context.Context, userId uint64, orderId string, state OrderState) error {
	return db.Model(&Order{}).Where(&Order{UserId: userId, OrderId: orderId}).Update("order_state", state).Error
}

// DeleteOrder soft deletes an order and its items
func DeleteOrder(db *gorm.DB, ctx context.Context, userId uint64, orderId string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 软删除订单项
		if err := tx.Model(&OrderItem{}).Where("order_id_refer = ?", orderId).Update("deleted_at", tx.NowFunc()).Error; err != nil {
			return err
		}
		
		// 软删除订单
		if err := tx.Model(&Order{}).Where(&Order{UserId: userId, OrderId: orderId}).Update("deleted_at", tx.NowFunc()).Error; err != nil {
			return err
		}
		
		return nil
	})
}
