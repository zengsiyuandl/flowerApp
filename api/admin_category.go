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

// AdminCreateCategory 创建分类
func AdminCreateCategory(c *gin.Context) {
	var category model.CategoryModel
	if err := c.ShouldBindJSON(&category); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if err := db.Get().Create(&category).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(category).WriteJSON(c.Writer)
}

// AdminUpdateCategory 更新分类
func AdminUpdateCategory(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var category model.CategoryModel
	if err := db.Get().Where("id = ?", id).First(&category).Error; err != nil {
		utils.Error(404, "分类不存在").WriteJSON(c.Writer)
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	category.Id = int32(id)
	if err := db.Get().Save(&category).Error; err != nil {
		utils.Error(500, "更新失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(category).WriteJSON(c.Writer)
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

