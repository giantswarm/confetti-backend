package websocketutil

import (
	"sync"
	"time"

	"github.com/atreugo/websocket"
	"github.com/giantswarm/microerror"
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

type ClientConfig struct {
	// Hub is the web socket hub that this client belongs to.
	Hub Hub
	// Connection is the websocket connection.
	Connection *websocket.Conn
}

// Client is a middleman between the websocket connection and the hub.
//
// Taken from https://github.com/fasthttp/websocket/blob/master/_examples/chat/fasthttp/client.go
// and modified to match our need.
type Client struct {
	hub  Hub
	conn *websocket.Conn

	// send is the buffered channel of outbound messages.
	send chan []byte

	// wg is a waiting group that synchronizes
	// the writer and reader goroutines.
	// It makes sure that the writer goroutine
	// will be closed before the request goroutine
	// is terminated.
	wg sync.WaitGroup
}

func NewClient(config ClientConfig) (*Client, error) {
	if config.Hub == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Hub must not be empty", config)
	}
	if config.Connection == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Connection must not be empty", config)
	}

	c := &Client{
		hub:  config.Hub,
		conn: config.Connection,
		send: make(chan []byte, 256),
	}

	c.registerToHub()

	return c, nil
}

func (c *Client) Emit(payload []byte) bool {
	select {
	case c.send <- payload:
	default:
		return false
	}

	return true
}

// GetUserValue retrieves a value from the connection context.
func (c *Client) GetUserValue(key string) interface{} {
	return c.conn.UserValue(key)
}

// SaveUserValue saves a value into the connection context.
func (c *Client) SaveUserValue(key string, value interface{}) {
	c.conn.SetUserValue(key, value)
}

func (c *Client) registerToHub() {
	c.wg.Add(1)

	c.hub.RegisterClient(c)

	go c.writePump()
	c.readPump()

	c.wg.Wait()
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	var err error
	var message []byte
	for {
		_, message, err = c.conn.ReadMessage()
		if err != nil {
			return
		}

		clientMessage := ClientMessage{
			Client:  c,
			Payload: message,
		}
		c.hub.SendMessage(clientMessage)
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

		c.wg.Done()
	}()

	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(CloseMessage, []byte{})

				return
			}

			w, err := c.conn.NextWriter(TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(PingMessage, nil); err != nil {
				return
			}
		}
	}
}
