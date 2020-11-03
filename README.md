We are building a kafka consumer in this project, which listens to kafka topic for messages, and based on the data received, 
it creates/updates the records in database. Additionally, it exposes certain endpoints to check the data available in the database.

#### This project demonstrates usage of: 
- Listening from a kafka topic
- Setup and querying SQL db
- goroutine, waitGroup and defer statements
- Exposing Restful APIs for reading data from database
- Config management from external files (Using Viper)
- Logging (Using zap)

#### Tech used:
1. Go - v1.15
2. SQL database
    - "database/sql"
    - "github.com/mattn/go-sqlite3"
3. Mux for request routing
    - "github.com/gorilla/mux"
4. Viper - Configuration management
    - "github.com/spf13/viper"
5. Zap - Logging
    - "go.uber.org/zap" 
6. Kafka consumer
    - "github.com/Shopify/sarama"
    - "github.com/wvanbergen/kafka/consumergroup"