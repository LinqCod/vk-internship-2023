package main

import (
	"github.com/linqcod/vk-internship-2023/internal/telegrambot"
	"github.com/linqcod/vk-internship-2023/pkg/config"
	"log"
)

func init() {
	config.LoadConfig()
}

func main() {
	bot, err := telegrambot.NewNumbersBot()
	if err != nil {
		log.Fatal(err)
	}

	bot.Start()
}
