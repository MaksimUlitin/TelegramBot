package main

import (
	"context"
	"flag"
	"log"

	tgClient "github.com/MaksimUlitin/internal/cliens"
	"github.com/MaksimUlitin/internal/consumer/eventconsumer"
	"github.com/MaksimUlitin/internal/storage/sqlite"
	"github.com/MaksimUlitin/internal/telegram"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "internal/data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can`t connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can`t init storage: ", err)
	}
	eventsProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), s)

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
		"token for access to telegram-bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}
	return *token
}
