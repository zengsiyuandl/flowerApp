package api

import (
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

// ProductDTO 商品数据传输对象
type ProductDTO struct {
	Id            int32   `json:"id"`
	CategoryId    int32   `json:"categoryId"`
	Name          string  `json:"name"`
	Subtitle      string  `json:"subtitle"`
	MainImageId   int32   `json:"mainImageId"`
	MainImage     string  `json:"mainImage"` // 主图URL，由MainImageId生成
	Price         float64 `json:"price"`
	OriginalPrice float64 `json:"originalPrice"`
	Stock         int     `json:"stock"`
	Sales         int     `json:"sales"`
	Description   string  `json:"description"`
	Detail        string  `json:"detail"`
	Sort          int     `json:"sort"`
	Status        int     `json:"status"`
	IsHot         int     `json:"isHot"`
	IsNew         int     `json:"isNew"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt"`
}

// ToProductDTO 将ProductModel转换为ProductDTO
func ToProductDTO(product model.ProductModel, baseURL string) ProductDTO {
	mainImageURL := ""
	if product.MainImageId > 0 {
		mainImageURL = utils.FormatImageURL(product.MainImageId, baseURL)
	}
	return ProductDTO{
		Id:            product.Id,
		CategoryId:    product.CategoryId,
		Name:          product.Name,
		Subtitle:      product.Subtitle,
		MainImageId:   product.MainImageId,
		MainImage:     mainImageURL,
		Price:         product.Price,
		OriginalPrice: product.OriginalPrice,
		Stock:         product.Stock,
		Sales:         product.Sales,
		Description:   product.Description,
		Detail:        product.Detail,
		Sort:          product.Sort,
		Status:        product.Status,
		IsHot:         product.IsHot,
		IsNew:         product.IsNew,
		CreatedAt:     utils.FormatTime(product.CreatedAt),
		UpdatedAt:     utils.FormatTime(product.UpdatedAt),
	}
}

// ToProductDTOList 将ProductModel列表转换为ProductDTO列表
func ToProductDTOList(products []model.ProductModel, baseURL string) []ProductDTO {
	result := make([]ProductDTO, len(products))
	for i, product := range products {
		result[i] = ToProductDTO(product, baseURL)
	}
	return result
}

// ProductImageDTO 商品图片数据传输对象
type ProductImageDTO struct {
	Id        int32  `json:"id"`
	ProductId int32  `json:"productId"`
	ImageId   int32  `json:"imageId"`
	ImageUrl  string `json:"imageUrl"` // 图片URL，由ImageId生成
	Sort      int    `json:"sort"`
	CreatedAt string `json:"createdAt"`
}

// ToProductImageDTO 将ProductImageModel转换为ProductImageDTO
func ToProductImageDTO(image model.ProductImageModel, baseURL string) ProductImageDTO {
	imageURL := ""
	if image.ImageId > 0 {
		imageURL = utils.FormatImageURL(image.ImageId, baseURL)
	}
	return ProductImageDTO{
		Id:        image.Id,
		ProductId: image.ProductId,
		ImageId:   image.ImageId,
		ImageUrl:  imageURL,
		Sort:      image.Sort,
		CreatedAt: utils.FormatTime(image.CreatedAt),
	}
}

// ToProductImageDTOList 将ProductImageModel列表转换为ProductImageDTO列表
func ToProductImageDTOList(images []model.ProductImageModel, baseURL string) []ProductImageDTO {
	result := make([]ProductImageDTO, len(images))
	for i, image := range images {
		result[i] = ToProductImageDTO(image, baseURL)
	}
	return result
}
