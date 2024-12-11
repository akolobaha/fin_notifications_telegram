package entity

import (
	"fmt"
	"strconv"
)

type User struct {
	ID       int64
	Name     string
	Email    string
	Telegram string
}

type TargetUser struct {
	Target      Target
	User        User
	ResultValue float64
}

func (tu *TargetUser) GenerateTelegramText() (text string) {
	text = fmt.Sprintf("🏁 Цель %s по %s достигнута:\nцель - %f, последнее значение - %f",
		getRatioText(tu.Target.ValuationRatio),
		tu.Target.Ticker,
		tu.Target.Value,
		tu.ResultValue,
	)

	return text
}

func (u *User) GetTelegramId() (id int64, err error) {
	parseInt, err := strconv.ParseInt(u.Telegram, 10, 64)
	if err != nil {
		return 0, err
	}

	return parseInt, nil
}

func getRatioText(ratio string) string {
	switch ratio {
	case "pbv":
		return "P / Bv"
	case "pe":
		return "P / E"
	case "ps":
		return "P / S"
	case "price":
		return "Цена за акцию"
	default:
		return ""
	}

}
