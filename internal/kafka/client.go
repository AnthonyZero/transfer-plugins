package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"strings"
	"transfer-plugins/configs"
)

var client *Client

type Client struct {
	client sarama.ConsumerGroup
}

func NewClient() error {
	conf := configs.Get().Kafka
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	switch conf.Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		return errors.New(fmt.Sprintf("Unrecognized consumer group partition assignor: %s", conf.Assignor))
	}
	if conf.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	consumerGroup, err := sarama.NewConsumerGroup(strings.Split(conf.Addr, ","), conf.Group, config)
	if err != nil {
		log.Printf("Error creating consumer group client: %v", err)
		return err
	}
	client = &Client{
		client: consumerGroup,
	}
	return nil
}

func ConsumerGroup() sarama.ConsumerGroup {
	return client.client
}

func Close() {
	if err := client.client.Close(); err != nil {
		log.Printf("Error closing client err:%v\n", err)
	}
	log.Printf("Kafka consumerGroup stopped\n")
}
