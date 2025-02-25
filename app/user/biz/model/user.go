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
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/naskids/nas-mall/common"
)

type User struct {
	common.Model
	Email          string `gorm:"unique"`
	PasswordHashed string
}

type UserFullMessage struct {
	common.Model
	Id             uint64 `gorm:"primary_key;AUTO_INCREMENT"`
	Email          string `gorm:"unique"`
	PasswordHashed string
	CreatedAt      uint64
	UpdatedAt      uint64
	DeletedAt      uint64
}

func (u User) TableName() string {
	return "user"
}

func GetByEmail(db *gorm.DB, ctx context.Context, email string) (user *User, err error) {
	err = db.WithContext(ctx).Model(&User{}).Where(&User{Email: email}).Where("deleted_at IS NULL").First(&user).Error
	return
}

func GetById(db *gorm.DB, ctx context.Context, id uint64) (*User, error) {
	var user User
	err := db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL"). // 已被删除
		First(&user).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user id:%d not found", id)
	}
	return &user, err
}

func Create(db *gorm.DB, ctx context.Context, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

// 安全更新用户信息
func UpdateUser(db *gorm.DB, ctx context.Context, userID uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return errors.New("更新字段不能为空")
	}

	result := db.WithContext(ctx).
		Model(&User{}).          // 指定模型（自动识别表名）
		Where("id = ?", userID). // 限定更新范围
		Updates(updates)         // 传入更新字段

	if result.Error != nil {
		return fmt.Errorf("数据库更新失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// 安全删除用户信息
func DeleteUser(db *gorm.DB, ctx context.Context, userID uint64) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 验证用户存在性
		if err := tx.First(&User{}).Where("id = ?", userID).Error; err != nil {
			return fmt.Errorf("用户不存在: %w", err)
		}
		// 按依赖顺序删除关联数据
		/*
			associations := []struct {
				Table    interface{}
				Relation string
			}{
				{&Payment{}, "user_id"},
				{&Order{}, "user_id"},
				{&cart{}, "user_id"},
			}

			for _, assoc := range associations {
				if err := tx.Unscoped().Where(assoc.Relation+" = ?", userID).Delete(assoc.Table).Error; err != nil {
					return fmt.Errorf("删除 %T 失败: %w", assoc.Table, err)
				}
			}
		*/

		// 最后删除主表
		if err := tx.Where("id = ?", userID).Delete(&User{}).Error; err != nil {
			return fmt.Errorf("主表删除失败: %w", err)
		}

		return nil
	})
}
