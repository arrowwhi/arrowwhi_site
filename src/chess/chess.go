package chess

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
	"site/auth"
)

var InitSetup = []string{
	"rnbqkbnr",
	"pppppppp",
	"........",
	"........",
	"........",
	"........",
	"PPPPPPPP",
	"RNBQKBNR",
}

var upgrader = websocket.Upgrader{} // use default options

var connects = make(map[string]*websocket.Conn)
var chessPlayers = make(map[string]string)

// ChessHandler websocket chess handler
func ChessHandler(c echo.Context) error {
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
			log.Println("error while reading:", err)
			break
		}
		log.Printf("recv: %s", message)
		var rawData map[string]interface{}
		err = json.Unmarshal(message, &rawData)
		if err != nil {
			fmt.Println("Ошибка при анмаршалинге JSON:", err)
			continue
		}
		switch rawData["m_type"] {
		case "select_player":
			break
		case "turn":
			break
		}

	}
	return nil
}

func Play() {

}

//func ()  {
//
//}
