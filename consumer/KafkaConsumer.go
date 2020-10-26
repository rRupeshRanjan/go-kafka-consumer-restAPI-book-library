package consumer

import (
	"github.com/wvanbergen/kafka/consumergroup"
	"go.uber.org/zap"
	"kafka-consumer/services"
)

var log, _ = zap.NewProduction()

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

