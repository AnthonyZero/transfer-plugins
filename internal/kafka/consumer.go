package kafka

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"transfer-plugins/configs"
	"transfer-plugins/internal/influxdb"
	"transfer-plugins/internal/models"
	"transfer-plugins/pkg/logger"
)

type Consumer struct {
	Ready   chan bool
	Service influxdb.Service
}

func (consumer *Consumer) Listener(ctx context.Context) {
	topicName := configs.Get().Kafka.TopicAction
	go func() {
		for {
			err := (*ConsumerGroup()).Consume(ctx, []string{topicName}, consumer)
			if err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				log.Printf("Cancle : %v", ctx.Err())
				return
			}
			//consumer.Ready = make(chan bool)
		}
	}()
	<-consumer.Ready
	log.Println("Sarama consumer up and running!...")
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	log.Println("Mark the consumer as ready")
	close(consumer.Ready)
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//var points []models.UserAction
	for message := range claim.Messages() {
		value := string(message.Value)

		var userAction models.UserAction
		json.Unmarshal([]byte(value), &userAction)
		consumer.Service.WritePoint(userAction)
		session.MarkMessage(message, "") //commit offset

		logger.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", value, message.Timestamp, message.Topic)
	}
	//log.Println("here no execute")
	return nil
}
