package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetCategoryList 管理后台获取分类列表（包含所有状态）
func AdminGetCategoryList(c *gin.Context) {
	var categories []model.CategoryModel
	db.Get().Order("sort DESC, id ASC").Find(&categories)

	utils.Success(categories).WriteJSON(c.Writer)
}

// AdminCreateCategory 创建分类（支持multipart/form-data，可同时上传图片）
func AdminCreateCategory(c *gin.Context) {
	name, sort, status := parseCategoryFormData(c)

	if name == "" {
		utils.Error(400, "分类名称不能为空").WriteJSON(c.Writer)
		return
	}

	iconID, errMsg := processCategoryIcon(c, 0)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	category := buildCategoryModel(0, name, iconID, sort, status)

	if err := db.Get().Create(&category).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	baseURL := getBaseURL(c)
	categoryDTO := ToCategoryDTO(category, baseURL)
	utils.Success(categoryDTO).WriteJSON(c.Writer)
}

// AdminUpdateCategory 更新分类（支持multipart/form-data，可同时上传图片）
func AdminUpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var existingCategory model.CategoryModel
	if err := db.Get().Where("id = ?", id).First(&existingCategory).Error; err != nil {
		utils.Error(404, "分类不存在").WriteJSON(c.Writer)
		return
	}

	name, sort, status := parseCategoryFormData(c)

	if name == "" {
		utils.Error(400, "分类名称不能为空").WriteJSON(c.Writer)
		return
	}

	iconID, errMsg := processCategoryIcon(c, existingCategory.IconId)
	if errMsg != "" {
		utils.Error(400, errMsg).WriteJSON(c.Writer)
		return
	}

	category := buildCategoryModel(int32(id), name, iconID, sort, status)

	if err := db.Get().Save(&category).Error; err != nil {
		utils.Error(500, "更新失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	baseURL := getBaseURL(c)
	categoryDTO := ToCategoryDTO(category, baseURL)
	utils.Success(categoryDTO).WriteJSON(c.Writer)
}

// AdminDeleteCategory 删除分类
func AdminDeleteCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	if err := db.Get().Where("id = ?", id).Delete(&model.CategoryModel{}).Error; err != nil {
		utils.Error(500, "删除失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

