package consumer

import (
	"github.com/Shopify/sarama"
	"github.com/wvanbergen/kafka/consumergroup"
	"go-kafka-consumer-restAPI-book-library/config"
	"go-kafka-consumer-restAPI-book-library/services"
	"go.uber.org/zap"
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
		case msg := <-consumerGroup.Messages():
			if msg.Topic == topic {
				log.Info(msg.Topic + " -----> " + string(msg.Value))
				_, processErr := services.ProcessMessage(msg.Value)
				if processErr == nil || processErr.Error() == config.InvalidDataMessage {
					err := consumerGroup.CommitUpto(msg)
					if err != nil {
						log.Error("Error committing offset: " + err.Error())
					}
				}
			}
		}
	}
}
