package main

import (
	"context"
	"news-feed/listener"
	"news-feed/manager"
)

func main() {
	km, _ := manager.GetKafkaManager()
	go km.Consume(context.Background())

	rest := listener.GetRestListener()
	err := rest.Start(":8081")
	if err != nil {
		return
	}
}
