package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fundnet/backend/internal/config"
	"fundnet/backend/internal/handlers"
	"fundnet/backend/internal/models"
	"fundnet/backend/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	if err := models.InitDB(cfg.Database.Path); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer models.CloseDB()

	// 初始化服务
	fundService := services.NewFundService()
	估值Service := services.NewEstimateService()

	// 设置 Gin 模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 路由
	router := gin.Default()

	// 启用 CORS
	router.Use(corsMiddleware())

	// 注册路由
	handlers.RegisterRoutes(router, fundService, 估值Service)

	// 启动定时任务
	go startScheduler(fundService, 估值Service, cfg)

	// 创建 HTTP 服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// 优雅关闭
	go func() {
		log.Printf("Server starting on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 关闭数据库
	models.CloseDB()
	log.Println("Server stopped")
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func startScheduler(fundService *services.FundService, 估值Service *services.EstimateService, cfg *config.Config) {
	ticker := time.NewTicker(time.Duration(cfg.App.RefreshInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Executing scheduled estimation update...")
		fundService.UpdateAllFundData()
		估值Service.RefreshAllEstimates()
	}
}
