package main

import (
	"flag"
	"github.com/sunriseex/tgbot-notifier/storage/files"
	"log"

	tgClient "github.com/sunriseex/tgbot-notifier/clients/telegram"
	event_consumer "github.com/sunriseex/tgbot-notifier/consumer/event-consumer"

	telegram "github.com/sunriseex/tgbot-notifier/events/telegram"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)
	log.Print("Server started")
	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
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
