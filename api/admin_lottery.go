package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetLotteryList 管理后台获取抽奖活动列表
func AdminGetLotteryList(c *gin.Context) {
	var lotteries []model.LotteryModel
	db.Get().Order("id DESC").Find(&lotteries)

	utils.Success(lotteries).WriteJSON(c.Writer)
}

// AdminCreateLottery 创建抽奖活动
func AdminCreateLottery(c *gin.Context) {
	var lottery model.LotteryModel
	if err := c.ShouldBindJSON(&lottery); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	// 检查时间是否有效（避免零值时间导致数据库错误）
	if lottery.StartTime.IsZero() || lottery.EndTime.IsZero() {
		utils.Error(400, "开始时间和结束时间不能为空").WriteJSON(c.Writer)
		return
	}

	if err := db.Get().Create(&lottery).Error; err != nil {
		utils.Error(500, "创建失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(lottery).WriteJSON(c.Writer)
}

// AdminUpdateLottery 更新抽奖活动
func AdminUpdateLottery(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var lottery model.LotteryModel
	if err := db.Get().Where("id = ?", id).First(&lottery).Error; err != nil {
		utils.Error(404, "抽奖活动不存在").WriteJSON(c.Writer)
		return
	}

	if err := c.ShouldBindJSON(&lottery); err != nil {
		utils.Error(400, "参数错误: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	lottery.Id = int32(id)
	if err := db.Get().Save(&lottery).Error; err != nil {
		utils.Error(500, "更新失败: "+err.Error()).WriteJSON(c.Writer)
		return
	}

	utils.Success(lottery).WriteJSON(c.Writer)
}

// AdminGetLotteryRecords 获取抽奖记录
func AdminGetLotteryRecords(c *gin.Context) {
	lotteryId := c.Query("lotteryId")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	query := db.Get().Model(&model.LotteryRecordModel{})
	if lotteryId != "" {
		query = query.Where("lottery_id = ?", lotteryId)
	}

	var total int64
	query.Count(&total)

	var records []model.LotteryRecordModel
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&records)

	utils.Success(map[string]interface{}{
		"list":     records,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}).WriteJSON(c.Writer)
}

