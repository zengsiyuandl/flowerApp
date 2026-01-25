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

// AdminCreateActivity 创建活动（支持multipart/form-data，可同时上传图片）
func AdminCreateActivity(c *gin.Context) {
	title, content, startTime, endTime, sort, status := parseActivityFormData(c)

	if title == "" {
		utils.Error(400, "活动标题不能为空").WriteJSON(c.Writer)
		return
	}

	imageID, errMsg := processActivityImage(c, 0)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	activity := buildActivityModel(
		0, title, content, imageID, startTime, endTime, sort, status,
	)

	if err := db.Get().Create(&activity).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	baseURL := getBaseURL(c)
	activityDTO := ToActivityDTO(activity, baseURL)
	utils.Success(activityDTO).WriteJSON(c.Writer)
}

// AdminUpdateActivity 更新活动（支持multipart/form-data，可同时上传图片）
func AdminUpdateActivity(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var existingActivity model.ActivityModel
	if err := db.Get().Where("id = ?", id).First(&existingActivity).Error; err != nil {
		utils.Error(404, "活动不存在").WriteJSON(c.Writer)
		return
	}

	title, content, startTime, endTime, sort, status := parseActivityFormData(c)

	if title == "" {
		utils.Error(400, "活动标题不能为空").WriteJSON(c.Writer)
		return
	}

	imageID, errMsg := processActivityImage(c, existingActivity.ImageId)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	activity := buildActivityModel(
		int32(id), title, content, imageID, startTime, endTime, sort, status,
	)

	if err := db.Get().Save(&activity).Error; err != nil {
		utils.Error(500, "更新失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	baseURL := getBaseURL(c)
	activityDTO := ToActivityDTO(activity, baseURL)
	utils.Success(activityDTO).WriteJSON(c.Writer)
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

