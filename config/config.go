package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	kafkaAddress    string
	redisAddress    string
	mysqlAddress    string
	infraConfigOnce sync.Once
)

func loadInfraAddress() {
	infraConfigOnce.Do(func() {
		file, err := os.Open("config.json")
		if err != nil {
			kafkaAddress = "error"
			redisAddress = "error"
			mysqlAddress = "error"
			return
		}
		defer file.Close()
		var cfg struct {
			KafkaAddress string `json:"kafka_address"`
			RedisAddress string `json:"redis_address"`
			MySQLAddress string `json:"mysql_address"`
		}
		if err := json.NewDecoder(file).Decode(&cfg); err != nil {
			kafkaAddress = "error"
			redisAddress = "error"
			mysqlAddress = "error"
		} else {
			kafkaAddress = cfg.KafkaAddress
			redisAddress = cfg.RedisAddress
			mysqlAddress = cfg.MySQLAddress
		}
	})
}

func GetKafkaAddress() string {
	loadInfraAddress()
	fmt.Printf("kafkaAddress: %s\n", kafkaAddress)
	return kafkaAddress
}

func GetRedisAddress() string {
	loadInfraAddress()
	fmt.Printf("redisAddress: %s\n", redisAddress)
	return redisAddress
}

func GetMySQLAddress() string {
	loadInfraAddress()
	fmt.Printf("mysqlAddress: %s\n", mysqlAddress)
	return mysqlAddress
}
