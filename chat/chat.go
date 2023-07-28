package chat

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log"
)

var upgrader = websocket.Upgrader{} // use default options
var conns []*websocket.Conn

// ChatHandler websocket chat handler
func ChatHandler(c echo.Context) error {
	con, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Print("upgrade:", err)
		return err
	}
	conns = append(conns, con)

	defer func(con *websocket.Conn) {
		removeFirstElement(con)
		con.Close()
	}(con)

	for {
		mt, message, err := con.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		for _, conn := range conns {
			if con == conn {
				continue
			}
			err = conn.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	}
	return nil
}

func removeFirstElement(c *websocket.Conn) {
	for i, value := range conns {
		if value == c {
			// Удаление элемента из среза
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}
}
