package api

import (
	"strconv"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetActivityList 获取活动列表
func GetActivityList(c *gin.Context) {
	now := time.Now()

	var activities []model.ActivityModel
	db.Get().Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now).
		Order("sort DESC, id DESC").Find(&activities)

	// 获取请求的基础URL
	baseURL := getBaseURL(c)

	activityDTOs := ToActivityDTOList(activities, baseURL)
	utils.Success(activityDTOs).WriteJSON(c.Writer)
}

// GetActivityDetail 获取活动详情
func GetActivityDetail(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var activity model.ActivityModel
	if err := db.Get().Where("id = ?", id).First(&activity).Error; err != nil {
		utils.Error(404, "活动不存在").WriteJSON(c.Writer)
		return
	}

	// 获取请求的基础URL
	baseURL := getBaseURL(c)

	activityDTO := ToActivityDTO(activity, baseURL)
	utils.Success(activityDTO).WriteJSON(c.Writer)
}

