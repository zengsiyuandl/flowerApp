package api

import (
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"
)

// LotteryDTO 抽奖活动数据传输对象
type LotteryDTO struct {
	Id         int32  `json:"id"`
	Name       string `json:"name"`
	ImageId    int32  `json:"imageId"`
	Image      string `json:"image"` // 活动图片URL，由ImageId生成
	Rule       string `json:"rule"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	DailyLimit int    `json:"dailyLimit"`
	TotalLimit int    `json:"totalLimit"`
	Status     int    `json:"status"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

// ToLotteryDTO 将LotteryModel转换为LotteryDTO
func ToLotteryDTO(lottery model.LotteryModel, baseURL string) LotteryDTO {
	imageURL := ""
	if lottery.ImageId > 0 {
		imageURL = utils.FormatImageURL(lottery.ImageId, baseURL)
	}
	return LotteryDTO{
		Id:         lottery.Id,
		Name:       lottery.Name,
		ImageId:    lottery.ImageId,
		Image:      imageURL,
		Rule:       lottery.Rule,
		StartTime:  utils.FormatTime(lottery.StartTime),
		EndTime:    utils.FormatTime(lottery.EndTime),
		DailyLimit: lottery.DailyLimit,
		TotalLimit: lottery.TotalLimit,
		Status:     lottery.Status,
		CreatedAt:  utils.FormatTime(lottery.CreatedAt),
		UpdatedAt:  utils.FormatTime(lottery.UpdatedAt),
	}
}
