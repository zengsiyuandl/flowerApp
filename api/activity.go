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

	activityDTOs := ToActivityDTOList(activities)
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

	activityDTO := ToActivityDTO(activity)
	utils.Success(activityDTO).WriteJSON(c.Writer)
}

