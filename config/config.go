package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	infraAddress    string
	mysqlAddress    string
	infraConfigOnce sync.Once
)

func loadInfraAddress() {
	infraConfigOnce.Do(func() {
		file, err := os.Open("config.json")
		if err != nil {
			infraAddress = "error"
			return
		}
		defer file.Close()
		var cfg struct {
			InfraAddress string `json:"infra_address"`
			MySQLAddress string `json:"mysql_address"`
		}
		if err := json.NewDecoder(file).Decode(&cfg); err != nil {
			infraAddress = "error"
			mysqlAddress = "error"
		} else {
			infraAddress = cfg.InfraAddress
			mysqlAddress = cfg.MySQLAddress
		}
	})
}

func GetInfraAddress() string {
	loadInfraAddress()
	fmt.Printf("infraAddress: %s\n", infraAddress)
	return infraAddress
}

func GetMySQLAddress() string {
	loadInfraAddress()
	fmt.Printf("mysqlAddress: %s\n", mysqlAddress)
	return mysqlAddress
}
