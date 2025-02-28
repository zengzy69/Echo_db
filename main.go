package main

import (
	"echoDB/config"
	"echoDB/db"
	_ "echoDB/docs"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"time"
)

func main() {
	// 加载配置文件

	config, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// 创建EchoDB实例
	echoDB := db.NewEchoDB(config)

	engine := db.NewGossipEngine(config.Gossip.Peers, config.Gossip.Port, config.Gossip.NodeID)

	go engine.StartGossipServer()

	// 启动Gossip传播，每5秒一次
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for range ticker.C {
			engine.Gossip()
		}
	}()

	// 根据配置选择一致性算法并启动服务
	if config.ConsistencyAlgorithm == "Gossip" {
		go echoDB.Gossip.StartGossipServer()
	}

	// 从MySQL加载数据
	items, err := db.LoadDataFromMySQL(config)
	if err != nil {
		log.Fatalf("Error loading data from MySQL: %v", err)
	}

	// 打印加载的数据
	fmt.Println("Loaded items from MySQL:", items)

	//
	router := gin.Default()

	router.GET("/check-update", db.NeedsUpdate)

	router.POST("/update-version", db.UpdateVersion)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务
	router.Run(":8080")
}
