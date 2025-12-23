package api

import (
	"math/rand"
	"time"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/middleware"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// GetLotteryInfo 获取抽奖活动信息
func GetLotteryInfo(c *gin.Context) {
	userId := middleware.GetUserId(c)

	now := time.Now()

	// 获取当前进行中的抽奖活动
	var lottery model.LotteryModel
	if err := db.Get().Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now).
		Order("id DESC").First(&lottery).Error; err != nil {
		utils.Error(404, "暂无抽奖活动").WriteJSON(c.Writer)
		return
	}

	// 检查今日抽奖次数
	var todayCount int64
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	db.Get().Model(&model.LotteryRecordModel{}).
		Where("lottery_id = ? AND user_id = ? AND created_at >= ?", lottery.Id, userId, todayStart).
		Count(&todayCount)

	// 检查总抽奖次数
	var totalCount int64
	db.Get().Model(&model.LotteryRecordModel{}).
		Where("lottery_id = ? AND user_id = ?", lottery.Id, userId).
		Count(&totalCount)

	canDraw := true
	if lottery.DailyLimit > 0 && int(todayCount) >= lottery.DailyLimit {
		canDraw = false
	}
	if lottery.TotalLimit > 0 && int(totalCount) >= lottery.TotalLimit {
		canDraw = false
	}

	utils.Success(map[string]interface{}{
		"lottery":    lottery,
		"todayCount": todayCount,
		"totalCount": totalCount,
		"canDraw":    canDraw,
	}).WriteJSON(c.Writer)
}

// DrawLottery 执行抽奖
func DrawLottery(c *gin.Context) {
	userId := middleware.GetUserId(c)

	now := time.Now()

	// 获取当前进行中的抽奖活动
	var lottery model.LotteryModel
	if err := db.Get().Where("status = ? AND start_time <= ? AND end_time >= ?", 1, now, now).
		Order("id DESC").First(&lottery).Error; err != nil {
		utils.Error(404, "暂无抽奖活动").WriteJSON(c.Writer)
		return
	}

	// 检查今日抽奖次数
	var todayCount int64
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	db.Get().Model(&model.LotteryRecordModel{}).
		Where("lottery_id = ? AND user_id = ? AND created_at >= ?", lottery.Id, userId, todayStart).
		Count(&todayCount)

	if lottery.DailyLimit > 0 && int(todayCount) >= lottery.DailyLimit {
		utils.Error(400, "今日抽奖次数已用完").WriteJSON(c.Writer)
		return
	}

	// 检查总抽奖次数
	var totalCount int64
	db.Get().Model(&model.LotteryRecordModel{}).
		Where("lottery_id = ? AND user_id = ?", lottery.Id, userId).
		Count(&totalCount)

	if lottery.TotalLimit > 0 && int(totalCount) >= lottery.TotalLimit {
		utils.Error(400, "抽奖次数已用完").WriteJSON(c.Writer)
		return
	}

	// 执行抽奖（简化实现，实际应根据奖品配置和概率）
	rand.Seed(time.Now().UnixNano())
	prizeType := rand.Intn(4) + 1 // 1-优惠券 2-积分 3-余额 4-实物

	prizeName := ""
	prizeValue := ""
	prizeId := int32(0)

	switch prizeType {
	case 1: // 优惠券
		prizeName = "5元优惠券"
		prizeValue = "5"
		// 可以关联实际的优惠券ID
	case 2: // 积分
		prizeName = "100积分"
		prizeValue = "100"
	case 3: // 余额
		prizeName = "10元余额"
		prizeValue = "10"
	case 4: // 实物
		prizeName = "谢谢参与"
		prizeValue = "0"
	}

	// 创建抽奖记录
	record := model.LotteryRecordModel{
		LotteryId: lottery.Id,
		UserId:    userId,
		PrizeType: prizeType,
		PrizeId:   prizeId,
		PrizeName: prizeName,
		PrizeValue: prizeValue,
		Status:    0, // 未领取
	}

	if err := db.Get().Create(&record).Error; err != nil {
		utils.Error(500, "抽奖失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(record).WriteJSON(c.Writer)
}

