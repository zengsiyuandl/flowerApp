package main

import (
	"fmt"
	"os"
	"wxcloudrun-golang/config"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/router"
	"wxcloudrun-golang/utils"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	// 确保存储目录存在（支持环境变量配置的持久化路径）
	storagePath := utils.GetStoragePath()
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		fmt.Printf("Warning: failed to create storage directory (%s): %v\n", storagePath, err)
	} else {
		fmt.Printf("Storage directory initialized: %s\n", storagePath)
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务
	if err := r.Run(":80"); err != nil {
		panic(fmt.Sprintf("server start failed with %+v", err))
	}
}
