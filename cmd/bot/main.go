package main

import (
	"log"

	"github.com/phanletrunghieu/botnet/bot/service/tcp"
)

func main() {
	tcpService := tcp.NewTCPService("127.0.0.1", 8080)
	go func() {
		err := <-tcpService.Error
		log.Println(err)
	}()

	log.Println("Connected!")

	tcpService.Run()
}
