package config

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"github.com/wvanbergen/kafka/consumergroup"
	"go.uber.org/zap"
	"time"
)

var consumerGroupName string
var zkHost string
var Topic string
var ServerPort string
var log zap.Logger

func InitConsumer() (*consumergroup.ConsumerGroup, error) {
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetOldest
	config.Offsets.ProcessingTimeout = 10 * time.Second

	consumerGroup, err := consumergroup.JoinConsumerGroup(
		consumerGroupName,
		[]string{Topic},
		[]string{zkHost},
		config)

	return consumerGroup, err
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err !=nil {
		log.Error("Error while reading config file " + err.Error())
	} else {
		consumerGroupName = viper.GetString("consumer.consumerGroup")
		zkHost = viper.GetString("consumer.zkHost")
		Topic = viper.GetString("consumer.topic")
		ServerPort = viper.GetString("server.port")
	}
}