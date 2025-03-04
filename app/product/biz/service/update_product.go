package service

import (
	"context"
	"fmt"
	"github.com/naskids/nas-mall/app/product/biz/model"
	product "github.com/naskids/nas-mall/rpc_gen/kitex_gen/product"

	"github.com/naskids/nas-mall/app/product/biz/dal/mysql"
)

type UpdateProductService struct {
	ctx context.Context
} // NewUpdateProductService new UpdateProductService
func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{ctx: ctx}
}

// 商品模型转proto结构
func convertProductModelToProto(p *model.Product) *product.Product {
	categoryNames := make([]string, 0, len(p.Categories))
	for _, c := range p.Categories {
		categoryNames = append(categoryNames, c.Name)
	}

	return &product.Product{
		Id:          uint64(p.ID),
		Name:        p.Name,
		Description: p.Description,
		Picture:     p.Picture,
		Price:       p.Price,
		Categories:  categoryNames,
	}
}

// Run create note info
// UpdateProductService 修改商品服务
func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	// 开启事务
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 获取现有商品
	var existingProduct model.Product
	// 已被删除 或者 不存在
	if err = tx.First(&existingProduct, req.Id).Where("deleted_at IS NULL").Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("商品不存在")
	}

	// 2. 构建更新字段（处理optional字段）
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Picture != nil {
		updates["picture"] = *req.Picture
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}

	// 3. 执行基础字段更新
	if len(updates) > 0 {
		if err = tx.Model(&existingProduct).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("更新商品信息失败: %v", err)
		}
	}

	// 4. 处理分类更新
	if req.Categories != nil {
		// 获取新分类列表
		newCategories := make([]model.Category, 0, len(req.Categories))
		for _, categoryName := range req.Categories {
			var category model.Category
			if err = tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("分类不存在: %s", categoryName)
			}
			newCategories = append(newCategories, category)
		}

		// 替换关联分类
		if err = tx.Model(&existingProduct).Association("Categories").Replace(newCategories); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("更新分类关联失败: %v", err)
		}
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	// 重新加载完整数据
	if err = mysql.DB.Preload("Categories").First(&existingProduct, req.Id).Error; err != nil {
		return nil, fmt.Errorf("获取更新后数据失败")
	}

	// 转换响应
	return &product.UpdateProductResp{
		Product: convertProductModelToProto(&existingProduct),
	}, nil
}
