package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"site/auth"
	"site/database"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Настройте настоящую проверку, если это необходимо
	},
}

var connects = make(map[string]*websocket.Conn)

// ChatHandler websocket chat handler
func ChatHandler(c echo.Context) error {
	var username = ""
	cookie, err := c.Cookie("token")
	if err == nil {
		username, err = auth.VerifyAndExtractUsername(cookie.Value)
	}
	con, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Print("upgrade:", err)
		return err
	}
	connects[username] = con
	defer func(con *websocket.Conn) {
		delete(connects, username)
		con.Close()
	}(con)

	for {
		_, message, err := con.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		var rawData map[string]interface{}
		err = json.Unmarshal(message, &rawData)
		if err != nil {
			fmt.Println("Ошибка при анмаршалинге JSON:", err)
			continue
		}

		switch rawData["m_type"].(string) {
		case "message":
			if rawData["user_to"].(string) == "" {
				log.Println("Не указан получатель")
				continue
			} else if rawData["message"].(string) == "" {
				log.Println("Не указано сообщение")
				continue
			}
			newMessage := database.Get().AddMessage(
				rawData["message"].(string),
				username,
				rawData["user_to"].(string))

			cn, exists := connects[newMessage.From]
			if exists {
				callbackId := map[string]interface{}{
					"m_type":   "take_id_from_local",
					"id":       newMessage.ID,
					"local_id": rawData["local_id"],
				}
				msg, err := json.Marshal(callbackId)
				if err != nil {
					log.Print(err)
				}
				err = cn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Print(err)
				}
			}

			cn, exists = connects[rawData["user_to"].(string)]
			if exists {
				msg, err := json.Marshal(newMessage)
				if err != nil {
					log.Print(err)
				}
				err = cn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Print(err)
				}
			}
			break
		case "read_message":
			type readMessage struct {
				ID []int `json:"ids"`
			}
			var id readMessage
			err := json.Unmarshal(message, &id)
			if err != nil {
				log.Print(err)
			}
			ids, err := database.Get().MakeMessagesRead(id.ID)
			if err != nil {
				log.Print(err)
			}
			//TODO переделать жeсткий костыль
			for _, v := range ids {
				msg, err := database.Get().SelectMessageById(v)
				if err != nil {
					log.Print(err.Error())
					continue
				}
				SendRead(msg.From, msg.ID)
				SendRead(msg.To, msg.ID)
			}
		}
	}
	return nil
}

func SendRead(login string, id uint) {
	cn, exists := connects[login]
	if exists {
		callbackId := map[string]interface{}{
			"m_type": "make_read",
			"id":     id,
		}
		msg, err := json.Marshal(callbackId)
		if err != nil {
			log.Print(err)
		}
		err = cn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Print(err)
		}
	}
}
