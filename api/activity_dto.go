package api

import (
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

// ActivityDTO 活动数据传输对象
type ActivityDTO struct {
	Id        int32  `json:"id"`
	Type      int    `json:"type"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	ImageId   int32  `json:"imageId"`
	Image     string `json:"image"` // 活动图片URL，由ImageId生成
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Status    int    `json:"status"`
	Sort      int    `json:"sort"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ToActivityDTO 将ActivityModel转换为ActivityDTO
func ToActivityDTO(activity model.ActivityModel, baseURL string) ActivityDTO {
	imageURL := ""
	if activity.ImageId > 0 {
		imageURL = utils.FormatImageURL(activity.ImageId, baseURL)
	}
	return ActivityDTO{
		Id:        activity.Id,
		Type:      activity.Type,
		Title:     activity.Title,
		Content:   activity.Content,
		ImageId:   activity.ImageId,
		Image:     imageURL,
		StartTime:  utils.FormatTime(activity.StartTime),
		EndTime:    utils.FormatTime(activity.EndTime),
		Status:    activity.Status,
		Sort:      activity.Sort,
		CreatedAt: utils.FormatTime(activity.CreatedAt),
		UpdatedAt: utils.FormatTime(activity.UpdatedAt),
	}
}

// ToActivityDTOList 将ActivityModel列表转换为ActivityDTO列表
func ToActivityDTOList(activities []model.ActivityModel, baseURL string) []ActivityDTO {
	result := make([]ActivityDTO, len(activities))
	for i, activity := range activities {
		result[i] = ToActivityDTO(activity, baseURL)
	}
	return result
}
