package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"transfer-plugins/configs"
	"transfer-plugins/internal/kafka"
)

func main() {
	log.Println("Starting a new Sarama consumer")

	configs.Init("dev")

	err := kafka.NewClient()
	if err != nil {
		log.Println(err)
	}

	ctx, cancle := context.WithCancel(context.Background())
	consumer := &kafka.Consumer{
		Ready: make(chan bool),
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
}
