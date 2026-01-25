package api

import (
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

// BannerDTO Banner数据传输对象
type BannerDTO struct {
	Id        int32     `json:"id"`
	Title     string    `json:"title"`
	ImageId   int32     `json:"imageId"`
	Image     string    `json:"image"` // 图片URL，由ImageId生成
	LinkType  int       `json:"linkType"`
	LinkValue string    `json:"linkValue"`
	Sort      int       `json:"sort"`
	Status    int       `json:"status"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// ToBannerDTO 将BannerModel转换为BannerDTO
func ToBannerDTO(banner model.BannerModel) BannerDTO {
	imageURL := ""
	if banner.ImageId > 0 {
		imageURL = utils.FormatImageURL(banner.ImageId)
	}
	return BannerDTO{
		Id:        banner.Id,
		Title:     banner.Title,
		ImageId:   banner.ImageId,
		Image:     imageURL,
		LinkType:  banner.LinkType,
		LinkValue: banner.LinkValue,
		Sort:      banner.Sort,
		Status:    banner.Status,
		CreatedAt: utils.FormatTime(banner.CreatedAt),
		UpdatedAt: utils.FormatTime(banner.UpdatedAt),
	}
}

// ToBannerDTOList 将BannerModel列表转换为BannerDTO列表
func ToBannerDTOList(banners []model.BannerModel) []BannerDTO {
	result := make([]BannerDTO, len(banners))
	for i, banner := range banners {
		result[i] = ToBannerDTO(banner)
	}
	return result
}
