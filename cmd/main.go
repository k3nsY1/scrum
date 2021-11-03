package main

import (
	"log"
	"time"

	bot "github.com/k3nsY1/scrumbot/pkg/bot"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	//Соединение с ботом
	b, err := tb.NewBot(tb.Settings{
		Token:  "1939233113:AAG1ieW_rrxf9XVy9C5Lv7_Gxg62Bn63nDs",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	bt := bot.CreateBot(b)

	bt.Init()

}
