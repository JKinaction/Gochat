package ws

import (
	"GoChat/model"
	session2 "GoChat/session"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub      *Hub
	username string
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}
type Msg struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func MesHandler(status int, clientMes *Msg) {

}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		fmt.Println(c.conn.RemoteAddr(), "read:", string(message))
		var clientMsg Msg
		json.Unmarshal(message, &clientMsg)
		var globalMes []byte
		if clientMsg.Data != nil {
			if clientMsg.Status == 1 {
				clientMsg.Data.(map[string]interface{})["user"] = c.username
				clientMsg.Data.(map[string]interface{})["time"] = time.Now().Format("[15:04:05] ")
				_, err = model.AddMsgtoDB(c.username, clientMsg.Data.(map[string]interface{})["content"].(string), clientMsg.Data.(map[string]interface{})["time"].(string))
				if err != nil {
					log.Fatalln("Unable to add message to database")
				}
				marshal, err := json.Marshal(clientMsg)
				if err != nil {
					log.Fatalln(err)
				}
				globalMes = marshal
				c.hub.broadcast <- globalMes
			} else if clientMsg.Status == 2 {
				clientMsg.Data.(map[string]interface{})["user"] = c.username
				clientMsg.Data.(map[string]interface{})["time"] = time.Now().Format("[15:04:05] ")
				_, err = model.AddPriMsgtoDB(c.username, clientMsg.Data.(map[string]interface{})["content"].(string),
					clientMsg.Data.(map[string]interface{})["toUser"].(string),
					clientMsg.Data.(map[string]interface{})["time"].(string))
				if err != nil {
					log.Fatalln("Unable to add message to database")
				}
				marshal, err := json.Marshal(clientMsg)
				if err != nil {
					log.Fatalln(err)
				}
				c.hub.priMes <- marshal
			}
		}
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			fmt.Println(c.conn.RemoteAddr(), "write:", string(message))
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err = w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	username := session2.GetSession(c, "user")
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), username: username.(string)}
	client.hub.register <- client
	//len := len(client.hub.clients)

	//var onlineClients []string
	//for cl := range client.hub.clients {
	//	onlineClients = append(onlineClients, cl.username)
	//}
	//marshal, err := json.Marshal(onlineClients)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//conn.WriteMessage(1, marshal)
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
