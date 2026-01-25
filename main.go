package main

import (
	"fmt"
	"wxcloudrun-golang/config"
	"wxcloudrun-golang/db"
	"wxcloudrun-golang/router"
)

func main() {
	// 初始化配置
	config.Init()

	// 初始化数据库
	if err := db.Init(); err != nil {
		panic(fmt.Sprintf("mysql init failed with %+v", err))
	}

	// 设置路由
	r := router.SetupRouter()

	// 启动服务
	if err := r.Run(":80"); err != nil {
		panic(fmt.Sprintf("server start failed with %+v", err))
	}
}
