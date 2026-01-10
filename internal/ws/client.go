package ws

import "github.com/gorilla/websocket"

type Client struct {
	Conn  *websocket.Conn
	Token string
	Send  chan []byte
}

func (c *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- c
		c.Conn.Close()
	}()

	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		h.broadcast <- BroadcastMessage{
			Token: c.Token,
			Data:  data,
		}
	}
}

func (c *Client) writePump() {
	for msg := range c.Send {
		c.Conn.WriteMessage(1, msg)
	}
}
