package main

import (
	"context"
	"errors"
	"fmt"
	"mint/Data/mysql"
	"mint/Data/redis"
	"mint/controller"
	"mint/logger"
	"mint/pkg/snowflake"
	"mint/routers"
	"mint/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

//Go Web开发较为通用的脚手架模板

func main() {
	//	1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failer, err: %v\n", err)
		return
	}

	//	2. 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failer, err: %v\n", err)
		return
	}
	defer zap.L().Sync()

	//	3. 初始化MySQL连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init mysql failer, err: %v\n", err)
		return
	}
	defer mysql.Close()

	//	4. 初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init redis failer, err: %v\n", err)
		return
	}
	defer redis.Close()

	//	5. 初始化雪花算法节点
	if err := snowflake.Init(); err != nil {
		zap.L().Fatal("init snowflake failer", zap.Error(err))
		return
	}

	//	6. 初始化翻译器
	if err := controller.InitTrans("zh"); err != nil {
		zap.L().Fatal("init translator failed", zap.Error(err))
		return
	}

	//	7. 注册路由
	r := routers.SetUpRouter()

	//	8. 启动服务并实现优雅关机

	//	8.1 声明原生的 HTTP Server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	//	8.2 将服务启动放到一个单独的 Goroutine 中
	go func() {
		zap.L().Info(fmt.Sprintf("Gin 服务器已启动，监听端口: %d", viper.GetInt("app.port")))
		// ErrServerClosed 表示服务器已被正常关闭，属于预期内的行为
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	//	8.3 设置监听操作系统信号的通道
	quit := make(chan os.Signal, 1)
	// kill 默认发送 syscall.SIGTERM
	// kill -2 发送 syscall.SIGINT (Ctrl+C)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	//	8.4 阻塞程序，直到接收到退出信号
	<-quit
	zap.L().Info("接收到退出信号，准备优雅关闭 Gin 服务器...")

	// 8.5 创建一个具有 5 秒超时机制的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 8.6 调用 Shutdown 进行优雅关闭
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Gin 服务器被强制关闭", zap.Error(err))
	}

	zap.L().Info("Gin 服务器已完美退出，所有正在处理的请求已完成！")
}
