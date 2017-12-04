package main

import "bufio"
import "fmt"
import "os"

import "gitlab.com/michael.golfi/appinsights/insights"
import "github.com/docker/docker/daemon/logger"

const version string = "1.0"

func processDiagnosticMessage(message string) {
	fmt.Println(message)
	fmt.Printf("> ")
}

const (
	url  = ""
	iKey = ""
)

func main() {
	count := 0

	config := make(map[string]string)
	config["insights-url"] = url
	config["insights-key"] = iKey
	info := logger.Info{
		Config: config,
	}

	client, err := insights.New(info)
	if err != nil {
		panic(err)
	}

	fmt.Println()
	fmt.Printf("Sending telemetry with iKey '%s'.", iKey)

	fmt.Println()

	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, _, err := in.ReadLine()
		if err != nil {
			panic(err)
		}

		input := string(line)

		if input == "exit" {
			return
		}

		if input == "" {
			continue
		}

		msg := logger.NewMessage()
		msg.Line = []byte(fmt.Sprintf("%s %d", input, count))
		client.Log(msg)
	}
}
