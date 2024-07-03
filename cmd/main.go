package main

import (
	"flag"
	"log"
)

func main() {
	t := mustToken()

	//tgClient = telegram,New(token)

	//fetcher = fetcher.New(tgClient)

	//processor = processor.New(tgClient)

	//consumer = consumer.Start(fetcher,processor)

}

func mustToken() string {
	token := flag.String(
		"token",
		"",
		"token for authentication",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("Please provide a token")
	}

	return *token
}
