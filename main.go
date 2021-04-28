package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"transfer-plugins/configs"
	"transfer-plugins/internal/influxdb"
	"transfer-plugins/internal/kafka"
	"transfer-plugins/pkg/logger"
	"transfer-plugins/utils"
)

var (
	config = flag.String("env", "", "请输入运行环境:\n dev:开发环境\n fat:测试环境\n uat:预上线环境\n pro:正式环境\n")
)

func main() {
	flag.Parse()
	if *config == "" {
		configs.Init("dev")
		log.Println("Warning: '-env' cannot be found, The default 'dev' will be used.")
	} else {
		envs := []string{"dev", "fat", "uat", "pro"}
		if !utils.Contain(*config, envs) {
			log.Println("Error: '-env' it is illegal")
			flag.Usage()
			os.Exit(1)
		}
		configs.Init(*config)
	}

	//日志初始化
	logger.InitZapLogger(configs.Get().Base.LogPath, logger.ToLevel("info"))

	//kafka消费组
	if err := kafka.NewClient(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//influxdb2
	if err := influxdb.NewClient(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ctx, cancle := context.WithCancel(context.Background())
	consumer := &kafka.Consumer{
		Ready:   make(chan bool),
		Service: influxdb.NewService(influxdb.WriteApi()),
	}
	consumer.Listener(ctx)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("Terminating: context cancelled")
	case <-sigterm:
		log.Println("Terminating: via signal")
	}
	cancle()
	kafka.Close()
	influxdb.Close()
}
