package main

import (
	"flag"
	"log"

	tgClient "github.com/MaksimUlitin/cliens"
	"github.com/MaksimUlitin/consumer/eventconsumer"

	"github.com/MaksimUlitin/events/telegram"
	"github.com/MaksimUlitin/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), files.New(storagePath))

	log.Print("servec started")

	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {

		log.Fatal("servec is stopped", err)
	}

}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for acces to telegram-bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
