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

type weather struct {
	Temp int    `json:"temp"`
	Icon string `json:"icon"`
	Link string `json:"link"`
}

type parent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type zoneID struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Offset            int      `json:"offset"`
	OffsetString      string   `json:"offsetString"`
	ShowSunriseSunset bool     `json:"showSunriseSunset"`
	Sunrise           string   `json:"sunrise"`
	Sunset            string   `json:"sunset"`
	IsNight           bool     `json:"isNight"`
	SkyColor          string   `json:"skyColor"`
	Weather           weather  `json:"weather"`
	Parents           []parent `json:"parents"`
}

type clocks struct {
	RegionNumber zoneID `json:"213"`
}

type exactTime struct {
	Time   int    `json:"time"`
	Clocks clocks `json:"clocks"`
}

var buttonsMarkup = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Курсы валют"),
		tgbotapi.NewKeyboardButton("Точное время (московское)"),
	),
)

var testURL = "https://yandex.com/time/sync.json?geo=213"

const msInSec = 1000

func getTimeAsString(timestamp int) string {
	fullDate := time.Unix(int64(timestamp/msInSec), 0)
	hour := strconv.Itoa(fullDate.Hour())
	minute := strconv.Itoa(fullDate.Minute())
	second := strconv.Itoa(fullDate.Second())
	return fmt.Sprintf("%s : %s : %s", string(hour), string(minute), string(second))
}

func main() {
	fmt.Println("Bot has been started...")

	text := `{"time":1587963941623,"clocks":{"213":{"id":213,"name":"Москва","offset":10800000,"offsetString":"UTC+3:00","showSunriseSunset":true,"sunrise":"04:55","sunset":"20:00","isNight":false,"skyColor":"#a0cdff","weather":{"temp":6,"icon":"skc-d","link":"https://yandex.ru/pogoda/moscow"},"parents":[{"id":225,"name":"Россия"}]}}}`

	TOKEN := os.Getenv("TOKEN")
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
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
		switch update.Message.Command() {
		case "start":
			msg.Text = "You are wellcome"

		default:
			msg.Text = "Неизвестная команда. Пардоньте"
		}

		switch update.Message.Text {
		case "Курсы валют":
			testData := exactTime{}
			data := exactTime{}

			textBytes := []byte(text)
			err := json.Unmarshal(textBytes, &testData)
			if err != nil {
				log.Println(err)
			}

			fmt.Println("message before make asunc call")
			res, err := http.Get(testURL)
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

			getTimeAsString(data.Time)

			msg.Text = getTimeAsString(data.Time)

		default:
			msg.Text = "Боюсь, что не понял вас..."
		}

		msg.ReplyMarkup = buttonsMarkup
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
