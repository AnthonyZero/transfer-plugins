package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
	"transfer-plugins/configs"
	"transfer-plugins/pkg/logger"
)

type Consumer struct {
	Ready chan bool
}

func (consumer *Consumer) Listener(ctx context.Context) {
	topicName := configs.Get().Kafka.TopicAction
	go func() {
		for {
			err := ConsumerGroup().Consume(ctx, []string{topicName}, consumer)
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
	for message := range claim.Messages() {
		logger.Infof("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}
