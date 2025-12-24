package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetActivityList 管理后台获取活动列表
func AdminGetActivityList(c *gin.Context) {
	var activities []model.ActivityModel
	db.Get().Order("id DESC").Find(&activities)

	utils.Success(activities).WriteJSON(c.Writer)
}

// AdminCreateActivity 创建活动
func AdminCreateActivity(c *gin.Context) {
	var activity model.ActivityModel
	if err := c.ShouldBindJSON(&activity); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if err := db.Get().Create(&activity).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(activity).WriteJSON(c.Writer)
}

// AdminUpdateActivity 更新活动
func AdminUpdateActivity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var activity model.ActivityModel
	if err := db.Get().Where("id = ?", id).First(&activity).Error; err != nil {
		utils.Error(404, "活动不存在").WriteJSON(c.Writer)
		return
	}

	if err := c.ShouldBindJSON(&activity); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	activity.Id = int32(id)
	if err := db.Get().Save(&activity).Error; err != nil {
		utils.Error(500, "更新失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(activity).WriteJSON(c.Writer)
}

// AdminDeleteActivity 删除活动
func AdminDeleteActivity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	if err := db.Get().Where("id = ?", id).Delete(&model.ActivityModel{}).Error; err != nil {
		utils.Error(500, "删除失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

