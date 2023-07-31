package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"site/auth"
	"site/redis"
)

var upgrader = websocket.Upgrader{} // use default options

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

	go redis.Get().ListenForMessages("all", con, username, func(con *websocket.Conn, usr string, msg []byte) error {
		err = con.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return err
		}
		return nil
	})

	defer func(con *websocket.Conn) {
		con.Close()
	}(con)

	for {
		_, message, err := con.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		var msg redis.SingleMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Println("Error:", err)
		}
		var room string
		if msg.Sender == "" || msg.Sender == "0" {
			room = "all"
		} else if msg.Sender < username {
			room = msg.Sender + username
		} else {
			room = username + msg.Sender
		}

		err = redis.Get().SendMessage(room, msg.Sender, msg.Message)
		if err != nil {
			log.Print(err.Error())
		}
	}
	return nil
}
