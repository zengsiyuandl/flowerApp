package api

import (
	"strconv"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// parseBannerFormData 解析Banner表单数据
// 返回：title, linkType, linkValue, sort, status
func parseBannerFormData(c *gin.Context) (string, int, string, int, int) {
	title := c.PostForm("title")
	linkTypeStr := c.PostForm("linkType")
	linkValue := c.PostForm("linkValue")
	sortStr := c.PostForm("sort")
	statusStr := c.PostForm("status")

	linkType, _ := strconv.Atoi(linkTypeStr)
	sort, _ := strconv.Atoi(sortStr)
	status, _ := strconv.Atoi(statusStr)

	return title, linkType, linkValue, sort, status
}

// processBannerImage 处理Banner图片上传
// 返回图片ID和错误信息
func processBannerImage(c *gin.Context, existingImageID int32) (int32, string) {
	// 尝试获取上传的图片文件
	if file, err := c.FormFile("image"); err == nil {
		// 有新上传的图片
		imageID, errMsg := UploadImageToDB(file, "banner")
		if errMsg != "" {
			return 0, errMsg
		}
		return imageID, ""
	}

	// 没有新图片，检查是否有imageId参数（编辑时保留原图片）
	imageIdStr := c.PostForm("imageId")
	if imageIdStr != "" {
		if id, err := utils.ParseInt32(imageIdStr); err == nil {
			return id, ""
		}
	}

	// 返回原有图片ID（编辑场景）
	return existingImageID, ""
}

// buildBannerModel 构建Banner模型
func buildBannerModel(id int32, title string, imageID int32, linkType int,
	linkValue string, sort int, status int) model.BannerModel {
	return model.BannerModel{
		Id:        id,
		Title:     title,
		ImageId:   imageID,
		LinkType:  linkType,
		LinkValue: linkValue,
		Sort:      sort,
		Status:    status,
	}
}
