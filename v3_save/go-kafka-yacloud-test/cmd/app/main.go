package main

import (
	"flag"
	"fmt"
	"kafka_yacloud/consumer"
	"kafka_yacloud/producer"
	"os"
)

func main() {
	startPtr := flag.String("start", "none", "lalalushka")
	flag.Parse()

	doneCh := make(chan int)
	switch *startPtr {
	case "p":
		go func() {
			err := producer.StartDS()
			if err != nil {
				fmt.Printf("ERROR: Producer: %s", err.Error())
			}
			doneCh <- 1
		}()
	case "c":
		go func() {
			err := consumer.Start()
			if err != nil {
				fmt.Printf("ERROR: Consumer: %s", err.Error())
			}
			doneCh <- 1
		}()
	case "none":
		fmt.Println("ну че ты...")
		return
	}

	for {
		select {
		case <-doneCh:
			os.Exit(1)
		}
	}
}
