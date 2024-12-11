package telegram

import (
	"fin_notifications_telegram/internal/log"
	"fmt"
	"github.com/mymmrac/telego"
	//th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func SendMessageToUser(token string, message string, tgUserId int64) error {
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Error("Ошибка инициализации телеграм-бота: ", err)
	}

	botUser, _ := bot.GetMe()
	fmt.Printf("Bot User: %+v\n", botUser)
	defer bot.StopLongPolling()

	_, err = bot.SendMessage(tu.Message(tu.ID(tgUserId), message))
	if err != nil {
		return err
	}
	return nil
}
