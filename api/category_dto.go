package api

import (
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

// CategoryDTO 分类数据传输对象
type CategoryDTO struct {
	Id        int32  `json:"id"`
	Name      string `json:"name"`
	IconId    int32  `json:"iconId"`
	Icon      string `json:"icon"` // 图标URL，由IconId生成
	Sort      int    `json:"sort"`
	Status    int    `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ToCategoryDTO 将CategoryModel转换为CategoryDTO
func ToCategoryDTO(category model.CategoryModel, baseURL string) CategoryDTO {
	iconURL := ""
	if category.IconId > 0 {
		iconURL = utils.FormatImageURL(category.IconId, baseURL)
	}
	return CategoryDTO{
		Id:        category.Id,
		Name:      category.Name,
		IconId:    category.IconId,
		Icon:      iconURL,
		Sort:      category.Sort,
		Status:    category.Status,
		CreatedAt: utils.FormatTime(category.CreatedAt),
		UpdatedAt: utils.FormatTime(category.UpdatedAt),
	}
}

// ToCategoryDTOList 将CategoryModel列表转换为CategoryDTO列表
func ToCategoryDTOList(categories []model.CategoryModel, baseURL string) []CategoryDTO {
	result := make([]CategoryDTO, len(categories))
	for i, category := range categories {
		result[i] = ToCategoryDTO(category, baseURL)
	}
	return result
}
