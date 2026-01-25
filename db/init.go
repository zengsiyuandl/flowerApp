package db

import (
	"fmt"
	"os"
	"time"
	"wxcloudrun-golang/db/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var dbInstance *gorm.DB

// Init 初始化数据库
func Init() error {
	user := os.Getenv("MYSQL_USERNAME")
	pwd := os.Getenv("MYSQL_PASSWORD")
	addr := os.Getenv("MYSQL_ADDRESS")
	dataBase := os.Getenv("MYSQL_DATABASE")
	if dataBase == "" {
		dataBase = "golang_demo"
	}

	// 修复连接字符串格式：使用单个&而不是&&
	source := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		user, pwd, addr, dataBase)
	
	// 隐藏密码打印（安全考虑）
	sourceForLog := fmt.Sprintf("%s:***@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local&timeout=10s&readTimeout=30s&writeTimeout=30s",
		user, addr, dataBase)
	fmt.Println("start init mysql with ", sourceForLog)

	// 重试连接（最多3次）
	var db *gorm.DB
	var err error
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(source), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err == nil {
			break
		}
		fmt.Printf("DB Open error (attempt %d/%d), err=%s\n", i+1, maxRetries, err.Error())
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after %d attempts: %v", maxRetries, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)  // 减少空闲连接数
	sqlDB.SetMaxOpenConns(50)  // 减少最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour)

	dbInstance = db

	// 自动迁移所有表
	if err := AutoMigrate(); err != nil {
		fmt.Printf("AutoMigrate error,err=%s\n", err.Error())
		// 迁移失败不阻止启动，只记录错误
	}

	fmt.Println("finish init mysql successfully")
	return nil
}

// AutoMigrate 自动迁移表结构
func AutoMigrate() error {
	return dbInstance.AutoMigrate(
		&model.UserModel{},
		&model.AddressModel{},
		&model.RechargeModel{},
		&model.CategoryModel{},
		&model.ProductModel{},
		&model.ProductImageModel{},
		&model.OrderModel{},
		&model.OrderItemModel{},
		&model.CartModel{},
		&model.PaymentModel{},
		&model.DeliveryModel{},
		&model.DeliveryTrackModel{},
		&model.ActivityModel{},
		&model.CouponModel{},
		&model.UserCouponModel{},
		&model.LotteryModel{},
		&model.LotteryRecordModel{},
		&model.BannerModel{},
	)
}

// Get ...
func Get() *gorm.DB {
	return dbInstance
}
