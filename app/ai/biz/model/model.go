package model

import (
	"context"
	"fmt"
	"regexp"

	"gorm.io/gorm"

	"github.com/pkg/errors"

	"github.com/naskids/nas-mall/app/order/biz/model"
)

func GetOrderFromRawSQL(ctx context.Context, db *gorm.DB, rawSql *string) (orders []model.Order, err error) {
	// 参数校验
	if rawSql == nil || *rawSql == "" {
		return executeFallback(ctx, db)
	}

	// 1. 校验SQL合法性
	if valid, err := validateSQLSyntax(ctx, db, *rawSql); !valid || err != nil {
		return executeFallback(ctx, db)
	}

	// 2. 尝试执行原始SQL
	if err := executeAndScan(ctx, db, *rawSql, &orders); err != nil {
		return executeFallback(ctx, db)
	}

	// 3. 校验结果是否为空（根据业务需求可选）
	if len(orders) == 0 {
		return executeFallback(ctx, db)
	}

	return orders, nil
}

// 使用EXPLAIN验证SQL语法
func validateSQLSyntax(ctx context.Context, db *gorm.DB, sql string) (bool, error) {
	// 防止非SELECT语句
	if matched, _ := regexp.MatchString(`^SELECT\s+`, sql); !matched {
		return false, errors.New("only select statements are allowed")
	}

	explainSQL := fmt.Sprintf("EXPLAIN %s", sql)
	if err := db.WithContext(ctx).Exec(explainSQL).Error; err != nil {
		return false, errors.Wrap(err, "EXPLAIN validation failed")
	}
	return true, nil
}

// 执行SQL并扫描结果
func executeAndScan(ctx context.Context, db *gorm.DB, sql string, dest interface{}) error {
	result := db.WithContext(ctx).Raw(sql).Scan(dest)
	if result.Error != nil || result.RowsAffected == 0 {
		return errors.Wrap(result.Error, "execute raw SQL failed")
	}
	return nil
}

// 执行兜底方案
func executeFallback(ctx context.Context, db *gorm.DB) (orders []model.Order, err error) {
	err = db.Model(&model.Order{}).Limit(5).Preload("OrderItems").Find(&orders).WithContext(ctx).Error
	return
}
