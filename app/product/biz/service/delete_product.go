package service

import (
	"context"
	"fmt"

	"github.com/naskids/nas-mall/app/product/biz/model"
	product "github.com/naskids/nas-mall/rpc_gen/kitex_gen/product"

	"github.com/naskids/nas-mall/app/product/biz/dal/mysql"
)

type DeleteProductService struct {
	ctx context.Context
} // NewDeleteProductService new DeleteProductService
func NewDeleteProductService(ctx context.Context) *DeleteProductService {
	return &DeleteProductService{ctx: ctx}
}

// Run create note info
// DeleteProductService 删除商品服务
func (s *DeleteProductService) Run(req *product.DeleteProductReq) (resp *product.DeleteProductResp, err error) {
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 获取商品（包含软删除的记录）
	var productToDelete model.Product
	if err = tx.Unscoped().Preload("Categories").First(&productToDelete, req.Id).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("商品不存在")
	}

	// 2. 执行软删除（自动设置deleted_at）
	if err = tx.Delete(&productToDelete).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("删除失败: %v", err)
	}

	// 3. 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交失败: %v", err)
	}

	return &product.DeleteProductResp{
		Product: convertProductModelToProto(&productToDelete),
	}, nil
}
