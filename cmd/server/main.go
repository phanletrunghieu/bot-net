package main

import (
	"log"

	"github.com/phanletrunghieu/botnet/server/service/tcp"
)

func main() {
	errChan := make(chan error)
	go func() {
		err := <-errChan
		log.Println(err)
	}()

	tcpService := tcp.NewTCPService(":8080")
	if tcpService.Error != nil {
		errChan <- tcpService.Error
		panic(tcpService.Error)
	}

	log.Println("Server started!")

	tcpService.Run()
	if tcpService.Error != nil {
		errChan <- tcpService.Error
		panic(tcpService.Error)
	}
}
