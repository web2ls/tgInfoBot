package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type exactTime struct {
	Time int `json:"time"`
}

var buttonsMarkup = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Точное время (московское)"),
	),
)

var timeURL = "https://yandex.com/time/sync.json?geo=213"

const msInSec = 1000

func checkOnAdditionalZero(dateValue string) string {
	if len(dateValue) < 2 {
		return "0" + dateValue
	}

	return dateValue
}

func getTimeAsString(timestamp int) string {
	fullDate := time.Unix(int64(timestamp/msInSec), 0)
	location, _ := time.LoadLocation("Europe/Moscow")
	fullDateByLocation := fullDate.In(location)
	hour := checkOnAdditionalZero(strconv.Itoa(fullDateByLocation.Hour()))
	minute := checkOnAdditionalZero(strconv.Itoa(fullDateByLocation.Minute()))
	second := checkOnAdditionalZero(strconv.Itoa(fullDateByLocation.Second()))

	return fmt.Sprintf("%s : %s : %s", string(hour), string(minute), string(second))
}

func main() {
	fmt.Println("Bot has been started...")

	TOKEN := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		panic(err)
	}

	for update := range updates {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg.Text = "You are wellcome"

			default:
				msg.Text = "Unknown command, sorry"
			}
		} else {
			switch update.Message.Text {
			case "Точное время (московское)":
				data := exactTime{}

				res, err := http.Get(timeURL)
				if err != nil {
					log.Println(err)
				}

				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Println(err)
				}

				err = json.Unmarshal(body, &data)
				if err != nil {
					log.Println(err)
				}

				msg.Text = getTimeAsString(data.Time)

			default:
				msg.Text = "I don't know what is this..."
			}
		}

		msg.ReplyMarkup = buttonsMarkup
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
