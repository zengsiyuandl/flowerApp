package api

import (
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetCouponList 获取可用优惠券列表
func GetCouponList(c *gin.Context) {
	userId := middleware.GetUserId(c)
	status := c.Query("status") // 0-未使用 1-已使用 2-已过期

	query := db.Get().Model(&model.UserCouponModel{}).Where("user_id = ?", userId)

	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		// 默认只返回未使用的
		query = query.Where("status = ?", 0)
	}

	now := time.Now()
	query = query.Where("expire_time > ?", now)

	var userCoupons []model.UserCouponModel
	query.Order("id DESC").Find(&userCoupons)

	// 获取优惠券详情
	type CouponWithDetail struct {
		model.UserCouponModel
		Coupon model.CouponModel `json:"coupon"`
	}

	var result []CouponWithDetail
	for _, uc := range userCoupons {
		var coupon model.CouponModel
		if err := db.Get().Where("id = ?", uc.CouponId).First(&coupon).Error; err == nil {
			result = append(result, CouponWithDetail{
				UserCouponModel: uc,
				Coupon:          coupon,
			})
		}
	}

	utils.Success(result).WriteJSON(c.Writer)
}

// ReceiveCouponRequest 领取优惠券请求
type ReceiveCouponRequest struct {
	CouponId int32 `json:"couponId" binding:"required"`
}

// ReceiveCoupon 领取优惠券
func ReceiveCoupon(c *gin.Context) {
	userId := middleware.GetUserId(c)

	var req ReceiveCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(400, "参数错误").WriteJSON(c.Writer)
		return
	}

	// 查询优惠券
	var coupon model.CouponModel
	if err := db.Get().Where("id = ? AND status = ?", req.CouponId, 1).First(&coupon).Error; err != nil {
		utils.Error(404, "优惠券不存在或已禁用").WriteJSON(c.Writer)
		return
	}

	now := time.Now()
	if now.Before(coupon.StartTime) || now.After(coupon.EndTime) {
		utils.Error(400, "优惠券不在有效期内").WriteJSON(c.Writer)
		return
	}

	// 检查是否已领取
	var count int64
	db.Get().Model(&model.UserCouponModel{}).
		Where("user_id = ? AND coupon_id = ?", userId, req.CouponId).
		Count(&count)
	if count > 0 {
		utils.Error(400, "已领取过该优惠券").WriteJSON(c.Writer)
		return
	}

	// 检查总数量限制
	if coupon.TotalCount > 0 && coupon.UsedCount >= coupon.TotalCount {
		utils.Error(400, "优惠券已领完").WriteJSON(c.Writer)
		return
	}

	// 创建用户优惠券
	userCoupon := model.UserCouponModel{
		UserId:     userId,
		CouponId:   req.CouponId,
		Status:     0, // 未使用
		ExpireTime: coupon.EndTime,
	}

	if err := db.Get().Create(&userCoupon).Error; err != nil {
		utils.Error(500, "领取失败").WriteJSON(c.Writer)
		return
	}

	// 更新优惠券已领取数量
	db.Get().Model(&coupon).Update("used_count", coupon.UsedCount+1)

	utils.Success(userCoupon).WriteJSON(c.Writer)
}

