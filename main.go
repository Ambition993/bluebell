package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web_app_base/controller"
	"web_app_base/dao/mysql"
	"web_app_base/dao/redis"
	"web_app_base/logger"
	"web_app_base/pkg/snowflake"
	"web_app_base/router"
	"web_app_base/settings"
)

func main() {
	// init config
	if err := settings.Init(); err != nil {
		fmt.Printf("init config failed , err:%v\n", err)
		return
	}
	// init logger
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed , err :%v\n", err)
		return
	}
	zap.L().Debug("logger init ok")
	defer zap.L().Sync()
	// init mysql
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init logger failed , err :%v\n", err)
		return
	}
	defer mysql.Close()
	// init redis
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init logger failed , err :%v\n", err)
		return
	}
	defer redis.Close()
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed , err :%v\n", err)
		return
	}

	// init routes
	r := router.SetupRouter()

	//初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed , err :%v\n", err)
		return
	}

	//  start sever gracefully
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	log.Println("Server exiting")
}
