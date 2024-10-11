package main

import (
	"flag"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/slack-go/slack"
)

// Константи для доступу до API месенджерів
const (
	discordToken  = "DISCORD_TOKEN"
	telegramToken = "TELEGRAM_TOKEN"
	slackToken    = "SLACK_TOKEN"
)

func sendToSlack(channel, message string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer handlePanic()

	go func() {
		api := slack.New(slackToken)
		_, _, err := api.PostMessage(channel, slack.MsgOptionText(message, false))
		if err != nil {
			panic(fmt.Sprintf("Помилка при відправленні повідомлення до Slack: %v", err))
		}
		fmt.Println("Повідомлення надіслано до Slack")
	}()
}

func sendToDiscord(channelID, message string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer handlePanic()

	go func() {
		dg, err := discordgo.New("Bot " + discordToken)
		if err != nil {
			panic(fmt.Sprintf("Помилка при підключенні до Discord: %v", err))
		}
		defer dg.Close()

		_, err = dg.ChannelMessageSend(channelID, message)
		if err != nil {
			panic(fmt.Sprintf("Помилка при відправленні повідомлення до Discord: %v", err))
		}
		fmt.Println("Повідомлення надіслано до Discord")
	}()
}

func sendToTelegram(chatID int64, message string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer handlePanic()

	go func() {
		bot, err := tgbotapi.NewBotAPI(telegramToken)
		if err != nil {
			panic(fmt.Sprintf("Помилка при підключенні до Telegram: %v", err))
		}

		msg := tgbotapi.NewMessage(chatID, message)
		_, err = bot.Send(msg)
		if err != nil {
			panic(fmt.Sprintf("Помилка при відправленні повідомлення до Telegram: %v", err))
		}
		fmt.Println("Повідомлення надіслано до Telegram")
	}()
}

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Printf("Виникла помилка: %v\n", r)
	}
}

func main() {
	// Прапорці для вибору месенджера
	slackFlag := flag.Bool("slack", false, "Надіслати повідомлення до Slack")
	telegramFlag := flag.Bool("telegram", false, "Надіслати повідомлення до Telegram")
	discordFlag := flag.Bool("discord", false, "Надіслати повідомлення до Discord")

	// Параметри для повідомлення
	message := flag.String("message", "", "Текст повідомлення")
	channel := flag.String("channel", "", "Канал або чат ID для повідомлення (для Telegram - Chat ID)")

	flag.Parse()

	if *message == "" || *channel == "" {
		log.Fatal("Необхідно вказати повідомлення і канал/чат ID")
	}

	var wg sync.WaitGroup

	// Відправка повідомлення у вибраний месенджер
	if *slackFlag {
		wg.Add(1)
		sendToSlack(*channel, *message, &wg)
	} else if *telegramFlag {
		wg.Add(1)
		var chatID int64
		fmt.Sscanf(*channel, "%d", &chatID)
		sendToTelegram(chatID, *message, &wg)
	} else if *discordFlag {
		wg.Add(1)
		sendToDiscord(*channel, *message, &wg)
	} else {
		log.Fatal("Необхідно вибрати хоча б один месенджер для відправки (-slack, -telegram, -discord)")
	}

	wg.Wait()

	fmt.Println("Повідомлення надіслано")
}
