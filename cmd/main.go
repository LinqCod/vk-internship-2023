package main

import (
	"github.com/linqcod/vk-internship-2023/internal/numbersapi"
	"github.com/linqcod/vk-internship-2023/internal/telegrambot"
	"log"
)

func main() {
	numbersApi := numbersapi.NewApi()

	bot, err := telegrambot.NewNumbersBot(numbersApi)
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}
