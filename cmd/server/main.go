package main

import (
	"log"
	"net/http"

	"github.com/phanletrunghieu/bot-net/server/service/boss"
	"github.com/phanletrunghieu/bot-net/server/service/client"
)

func main() {
	// create a file server
	http.Handle("/", http.FileServer(http.Dir("./public")))
	go func() {
		http.ListenAndServe(":8082", nil)
	}()

	//client
	serviceClient := client.NewClientService(8080)
	go func() {
		err := <-serviceClient.Error
		log.Println(err)
	}()

	go serviceClient.Run()

	// boss
	serviceBoss := boss.NewBossService(8081, serviceClient)
	go func() {
		err := <-serviceBoss.Error
		log.Println(err)
	}()

	serviceBoss.Run()
}
