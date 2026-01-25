package api

import (
	"strconv"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/db/model"
	"wxcloudrun-golang/utils"

	"github.com/gin-gonic/gin"
)

// AdminGetOrderList 管理后台获取订单列表
func AdminGetOrderList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	query := db.Get().Model(&model.OrderModel{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("order_no LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var orders []model.OrderModel
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&orders)

	utils.Success(map[string]interface{}{
		"list":     orders,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	}).WriteJSON(c.Writer)
}

// AdminShipOrder 发货
func AdminShipOrder(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	var order model.OrderModel
	if err := db.Get().Where("id = ?", id).First(&order).Error; err != nil {
		utils.Error(404, "订单不存在").WriteJSON(c.Writer)
		return
	}

	if order.Status != 1 {
		utils.Error(400, "订单状态不允许发货").WriteJSON(c.Writer)
		return
	}

	// 更新订单状态为配送中
	if err := db.Get().Model(&order).Update("status", 2).Error; err != nil {
		utils.Error(500, "发货失败").WriteJSON(c.Writer)
		return
	}

	// 创建配送信息
	delivery := model.DeliveryModel{
		OrderId:      order.Id,
		OrderNo:      order.OrderNo,
		DeliveryType: 1, // 自营配送
		Status:       1, // 已分配
	}
	// 使用 Select 明确指定要插入的字段，排除 start_time 和 complete_time，避免 MySQL 报错
	fieldsToSelect := []string{
		"order_id", "order_no", "delivery_type", "delivery_man_id", "delivery_man_name",
		"delivery_man_phone", "third_party_order_no", "status",
		"created_at", "updated_at",
	}
	if err := db.Get().Select(fieldsToSelect).Create(&delivery).Error; err != nil {
		utils.Error(500, "创建配送记录失败").WriteJSON(c.Writer)
		return
	}

	utils.Success(nil).WriteJSON(c.Writer)
}

