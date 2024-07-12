package main

import (
	"context"
	"flag"
	"github.com/sunriseex/tgbot-notifier/storage/sqlite"
	"log"

	tgClient "github.com/sunriseex/tgbot-notifier/clients/telegram"
	"github.com/sunriseex/tgbot-notifier/consumer/event-consumer"

	"github.com/sunriseex/tgbot-notifier/events/telegram"
)

const (
	tgBotHost = "api.telegram.org"
	//storagePath       = "storage"
	sqliteStoragePath = "storage/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	//s := files.New(storagePath)
	stor, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect storage ", err)
	}

	if err := stor.Init(context.TODO()); err != nil {
		log.Fatal("can't initialize storage ", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		stor,
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
