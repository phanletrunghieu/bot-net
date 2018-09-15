package main

import (
	"log"

	"github.com/phanletrunghieu/bot-net/boss/service/cli"
	"github.com/phanletrunghieu/bot-net/boss/service/tcp"
)

func main() {
	bossService := tcp.NewTCPService("127.0.0.1", 8081)
	go func() {
		err := <-bossService.Error
		log.Println(err)
	}()

	go bossService.Run()

	cli.ExecuteMain(bossService)
}
