package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetCouponList 管理后台获取优惠券列表
func AdminGetCouponList(c *gin.Context) {
	var coupons []model.CouponModel
	db.Get().Order("id DESC").Find(&coupons)

	utils.Success(coupons).WriteJSON(c.Writer)
}

// AdminCreateCoupon 创建优惠券
func AdminCreateCoupon(c *gin.Context) {
	var coupon model.CouponModel
	if err := c.ShouldBindJSON(&coupon); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	if err := db.Get().Create(&coupon).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(coupon).WriteJSON(c.Writer)
}

// AdminUpdateCoupon 更新优惠券
func AdminUpdateCoupon(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var coupon model.CouponModel
	if err := db.Get().Where("id = ?", id).First(&coupon).Error; err != nil {
		utils.Error(404, "优惠券不存在").WriteJSON(c.Writer)
		return
	}

	if err := c.ShouldBindJSON(&coupon); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	coupon.Id = int32(id)
	if err := db.Get().Save(&coupon).Error; err != nil {
		utils.Error(500, "更新失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(coupon).WriteJSON(c.Writer)
}

// AdminDeleteCoupon 删除优惠券
func AdminDeleteCoupon(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	if err := db.Get().Where("id = ?", id).Delete(&model.CouponModel{}).Error; err != nil {
		utils.Error(500, "删除失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

