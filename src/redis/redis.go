package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"log"
)

type Json struct {
	MType string
	Body  struct{}
}

type SingleMessage struct {
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	Recipient string `json:"user_to"`
}

type Redis struct {
	ctx    context.Context
	client *redis.Client
	Stream string
}

var v *Redis

func Get() *Redis {
	if v == nil {
		v = &Redis{}
		err := v.NewRedis()
		if err != nil {
			log.Fatal("Ошибка при подключении к Redis:", err)
		}
	}
	return v
}

func (r *Redis) NewRedis() error {
	r.ctx = context.Background()
	r.client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	r.Stream = "chat"
	// Проверка поддержки Redis Stream
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) SendMessage(chatRoom, sender, message string) error {
	// Добавление сообщения в Redis Stream
	_, err := r.client.XAdd(r.ctx, &redis.XAddArgs{
		Stream: chatRoom,
		Values: map[string]interface{}{"sender": sender, "content": message},
	}).Result()

	return err
}

func (r *Redis) ListenForMessages(chatRoom *string, con *websocket.Conn, usr string, fn func(*websocket.Conn, []byte) error, stopChan chan struct{}) {
	for {

		select {
		case <-stopChan:
			fmt.Println("Worker received stop signal. Exiting...")
			return
		default:
			// Получение сообщений из Redis Stream с определенными ограничениями (получим только новые сообщения)
			messages, err := r.client.XRead(r.ctx, &redis.XReadArgs{
				Streams: []string{*chatRoom, "$"}, // "$" указывает на получение только новых сообщений
				Count:   0,
				Block:   0,
			}).Result()

			if err != nil {
				log.Println("Ошибка при получении сообщений из Redis Stream:", err)
				return
			}

			// Выводим сообщения в консоль
			for _, message := range messages {
				streamMessages := message.Messages
				for _, streamMessage := range streamMessages {
					msg := SingleMessage{
						Message:   streamMessage.Values["content"].(string),
						Sender:    streamMessage.Values["sender"].(string),
						Recipient: usr,
					}

					jsonData, err := json.Marshal(msg)
					if err != nil {
						fmt.Println("Ошибка при маршалинге в JSON:", err)
						return
					}

					err = fn(con, jsonData)
					if err != nil {
						log.Println("Ошибка при отправке сообщения клиенту:", err)
					}
				}
			}
		}
	}
}
