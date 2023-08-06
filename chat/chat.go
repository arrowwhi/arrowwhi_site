package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"site/auth"
	"site/database"
)

var upgrader = websocket.Upgrader{} // use default options

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
			newMessage := database.Get().AddMessage(
				rawData["message"].(string),
				username,
				rawData["recipient"].(string))
			if cn, exists := connects[rawData["recipient"].(string)]; exists {
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
			//case "dl_messages":
			//	messages := database.Get().SelectMessages(username, rawData["user"].(string), 10, 0)
			//	d := map[string]interface{}{
			//		"type":    "dl_messages",
			//		"content": messages,
			//	}
			//	jsonData, err := json.Marshal(d)
			//	if err != nil {
			//		fmt.Println("Ошибка при преобразовании массива в JSON:", err)
			//		continue
			//	}
			//
			//	err = con.WriteMessage(websocket.TextMessage, jsonData)
			//	if err != nil {
			//		fmt.Println("Ошибка при отправке данных через веб-сокет:", err)
			//		continue
			//	}
			//break
			//case "close":
			//	return nil
		}
	}
	return nil
}
