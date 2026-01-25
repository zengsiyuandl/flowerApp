package api

import (
	"strconv"
	"time"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// parseActivityFormData 解析活动表单数据
func parseActivityFormData(c *gin.Context) (
	title, content string,
	startTime, endTime time.Time,
	sort, status int,
) {
	title = c.PostForm("title")
	content = c.PostForm("content")
	startTimeStr := c.PostForm("startTime")
	endTimeStr := c.PostForm("endTime")
	sortStr := c.PostForm("sort")
	statusStr := c.PostForm("status")

	startTime, _ = time.Parse("2006-01-02 15:04:05", startTimeStr)
	endTime, _ = time.Parse("2006-01-02 15:04:05", endTimeStr)
	sort, _ = strconv.Atoi(sortStr)
	status, _ = strconv.Atoi(statusStr)

	return
}

// processActivityImage 处理活动封面图上传
func processActivityImage(c *gin.Context, existingImageID int32) (int32, string) {
	if file, err := c.FormFile("image"); err == nil {
		imageID, errMsg := UploadImageToDB(file, "activity")
		if errMsg != "" {
			return 0, errMsg
		}
		return imageID, ""
	}

	imageIdStr := c.PostForm("imageId")
	if imageIdStr != "" {
		if id, err := utils.ParseInt32(imageIdStr); err == nil {
			return id, ""
		}
	}

	return existingImageID, ""
}

// buildActivityModel 构建活动模型
func buildActivityModel(
	id int32,
	title, content string,
	imageId int32,
	startTime, endTime time.Time,
	sort, status int,
) model.ActivityModel {
	return model.ActivityModel{
		Id:        id,
		Title:     title,
		Content:   content,
		ImageId:   imageId,
		StartTime: startTime,
		EndTime:   endTime,
		Sort:      sort,
		Status:    status,
	}
}
