package commands

import (
	"context"
	"encoding/json"
	"fin_notifications_telegram/internal/config"
	"fin_notifications_telegram/internal/entity"
	"fin_notifications_telegram/internal/log"
	"fin_notifications_telegram/internal/telegram"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func ReadFromQueue(ctx context.Context, cfg *config.Config) error {
	conn, err := amqp.Dial(cfg.GetRabbitDSN())
	if err != nil {
		log.Error("Failed to connect to RabbitMQ: ", err)
		return err
	}
	defer conn.Close()

	chNotification, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a channel: ", err)
		return err
	}
	defer chNotification.Close()

	chEmailConfirm, err := conn.Channel()
	if err != nil {
		log.Error("Failed to open a channel: ", err)
		return err
	}
	defer chEmailConfirm.Close()

	msgs, err := chNotification.Consume(
		cfg.RabbitQueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	for {
		select {
		case msg := <-msgs:
			var targetUser entity.TargetUser
			err := json.Unmarshal(msg.Body, &targetUser)
			if err != nil {
				return err
			}

			text := targetUser.GenerateTelegramText()
			tgId, err := targetUser.User.GetTelegramId()
			if err != nil {
				log.Error("Failed to get telegram id: ", err)
			}

			err = telegram.SendMessageToUser(cfg.TelegramToken, text, tgId)
			if err != nil {
				log.Error("Failed to send message: ", err)
			}

			// Подтверждение сообщения после успешной обработки
			if err := msg.Ack(false); err != nil {
				log.Error("Failed to ack message: %s", err)
			}

		case <-ctx.Done():
			log.Info("Сервис обработки данных остановлен")
			time.Sleep(5 * time.Second)
			return nil
		}
	}
}
