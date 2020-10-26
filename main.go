package main

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"kafka-consumer/config"
	"kafka-consumer/consumer"
	"kafka-consumer/repository"
	"kafka-consumer/services"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup
var log, _ = zap.NewProduction()

func main()  {
	config.InitConfig()
	initializeDatabase()

	wg.Add(1)
	go initializeKafkaConsumer()
	wg.Add(1)
	go setupHttpHandles()
	wg.Wait()
}

func initializeKafkaConsumer() {
	defer wg.Done()

	consumerGroup, consumerInitErr := config.InitConsumer()
	errorHandler(consumerInitErr, "consumerGroup")

	defer consumerGroup.Close()
	consumer.Consume(consumerGroup, config.Topic)
}

func setupHttpHandles() {
	defer wg.Done()

	router := mux.NewRouter()
	router.HandleFunc("/books", services.GetAllBooksHandler).Methods("GET")
	router.HandleFunc("/books/{id}", services.GetBookByIdHandler).Methods("GET")
	httpServerInitErr := http.ListenAndServe(config.ServerPort, router)
	errorHandler(httpServerInitErr, "httpServer")
}

func initializeDatabase() {
	databaseInitErr := repository.InitBooksDb()
	errorHandler(databaseInitErr, "database")
}

func errorHandler(err error, entity string) {
	if err != nil {
		log.Error("Error initializing " + entity + " : " + err.Error())
		os.Exit(1);
	}
}
