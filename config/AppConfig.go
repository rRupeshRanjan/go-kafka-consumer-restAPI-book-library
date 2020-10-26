package config

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var ConsumerGroupName string
var ZkHost string
var Topic string
var ServerPort string
var log zap.Logger

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err !=nil {
		log.Error("Error while reading config file " + err.Error())
	} else {
		ConsumerGroupName = viper.GetString("consumer.consumerGroup")
		ZkHost = viper.GetString("consumer.zkHost")
		Topic = viper.GetString("consumer.topic")
		ServerPort = viper.GetString("server.port")
	}
}