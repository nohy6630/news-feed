package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var (
	infraAddress    string
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
		}
		if err := json.NewDecoder(file).Decode(&cfg); err != nil || cfg.InfraAddress == "" {
			infraAddress = "error"
		} else {
			infraAddress = cfg.InfraAddress
		}
	})
}

func GetInfraAddress() string {
	loadInfraAddress()
	fmt.Printf("infraAddress: %s\n", infraAddress)
	return infraAddress
}
