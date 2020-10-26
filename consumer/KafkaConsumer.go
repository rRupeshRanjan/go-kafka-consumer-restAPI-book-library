package consumer

import (
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"go.uber.org/zap"
	"kafka-consumer/config"
	"kafka-consumer/services"
	"time"
)

var log, _ = zap.NewProduction()

func InitConsumer() (*consumergroup.ConsumerGroup, error) {
	consumerConfig := consumergroup.NewConfig()
	consumerConfig.Offsets.Initial = sarama.OffsetOldest
	consumerConfig.Offsets.ProcessingTimeout = 10 * time.Second

	consumerGroup, err := consumergroup.JoinConsumerGroup(
		config.ConsumerGroupName,
		[]string{config.Topic},
		[]string{config.ZkHost},
		consumerConfig)

	return consumerGroup, err
}

func Consume(consumerGroup *consumergroup.ConsumerGroup, topic string) {
	for {
		select {
		// messages coming through channel
		case msg := <- consumerGroup.Messages():
			if msg.Topic == topic {
				log.Info(msg.Topic + " -----> " + string(msg.Value))
				services.ProcessMessage(msg.Value)
				err := consumerGroup.CommitUpto(msg)
				if err != nil {
					log.Error("Error committing offset: " + err.Error())
				}
			}
		}
	}
}
