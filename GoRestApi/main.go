package main

import (
	tgClient "GoRestApi/clients/telegram"
	event_consumer "GoRestApi/consumer/event-consumer"
	"GoRestApi/events/telegram"
	"GoRestApi/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Printf("service started")

	consumer := event_consumer.New(eventProcessor, eventProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("consumer stopped")
	}
	//fetcher = fetcher.New() - общается с апи тг. Фетчер будет отправлять запросы для получения новых событий

	//processor = processor.New() - будет отправлять новыен сообщения (ссылки на статьи)

	//consumer.Start(fetcher, p	kekekke := files.New("storage")rocessor)
}

func mustToken() string { //функции must - это те, которые падают если возникает ошибка
	token := flag.String("bot-token", "", "token for tg bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is  not specified")
	}

	return *token
}

//t.me/My_Saved_Links_Bot
