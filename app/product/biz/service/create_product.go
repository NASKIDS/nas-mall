package service

import (
	"context"
	"fmt"
	"github.com/naskids/nas-mall/app/product/biz/model"
	product "github.com/naskids/nas-mall/rpc_gen/kitex_gen/product"

	"github.com/naskids/nas-mall/app/product/biz/dal/mysql"
)

type CreateProductService struct {
	ctx context.Context
} // NewCreateProductService new CreateProductService
func NewCreateProductService(ctx context.Context) *CreateProductService {
	return &CreateProductService{ctx: ctx}
}

// Run create note info
func (s *CreateProductService) Run(req *product.CreateProductReq) (resp *product.CreateProductResp, err error) {
	// 开启事务
	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查所有分类是否存在并获取完整对象
	categories := make([]model.Category, 0, len(req.Categories))
	for _, categoryName := range req.Categories {
		var category model.Category
		if err = tx.Where("name = ?", categoryName).First(&category).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("分类不存在: %s", categoryName)
		}
		categories = append(categories, category)
	}

	// 创建商品记录
	newProduct := &model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Categories:  categories, // 直接关联分类对象，利用GORM的Many2Many关联自动创建关联记录
	}

	// 使用事务创建商品及关联关系
	if err = tx.Create(newProduct).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建商品失败: %v", err)
	}

	// 提交事务
	if err = tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("提交事务失败: %v", err)
	}

	// 转换GORM模型到proto结构
	respProduct := &product.Product{
		Id:          uint64(newProduct.ID),
		Name:        newProduct.Name,
		Description: newProduct.Description,
		Picture:     newProduct.Picture,
		Price:       newProduct.Price,
		Categories:  req.Categories, // 直接使用请求的分类名称
	}

	return &product.CreateProductResp{
		Product: respProduct,
	}, nil
}
