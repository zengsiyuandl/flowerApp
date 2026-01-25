package api

import (
	"strconv"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// parseCategoryFormData 解析分类表单数据
func parseCategoryFormData(c *gin.Context) (name string, sort, status int) {
	name = c.PostForm("name")
	sortStr := c.PostForm("sort")
	statusStr := c.PostForm("status")

	sort, _ = strconv.Atoi(sortStr)
	status, _ = strconv.Atoi(statusStr)

	return
}

// processCategoryIcon 处理分类图标上传
func processCategoryIcon(c *gin.Context, existingIconID int32) (int32, string) {
	if file, err := c.FormFile("icon"); err == nil {
		iconID, errMsg := UploadImageToDB(file, "category")
		if errMsg != "" {
			return 0, errMsg
		}
		return iconID, ""
	}

	iconIdStr := c.PostForm("iconId")
	if iconIdStr != "" {
		if id, err := utils.ParseInt32(iconIdStr); err == nil {
			return id, ""
		}
	}

	return existingIconID, ""
}

// buildCategoryModel 构建分类模型
func buildCategoryModel(
	id int32,
	name string,
	iconId int32,
	sort, status int,
) model.CategoryModel {
	return model.CategoryModel{
		Id:     id,
		Name:   name,
		IconId: iconId,
		Sort:   sort,
		Status: status,
	}
}
