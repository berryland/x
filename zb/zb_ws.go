package zb

import (
	. "github.com/berryland/x"
	json "github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

const WebSocketServerUrl = "wss://api.zb.com:9999/websocket"

type ZbWebSocketClient struct {
	running   bool
	conn      *websocket.Conn
	decoders  map[string]func([]byte) interface{}
	callbacks map[string]func(interface{})
}

func NewWebSocketClient() *ZbWebSocketClient {
	return &ZbWebSocketClient{running: false, decoders: make(map[string]func([]byte) interface{}), callbacks: make(map[string]func(interface{}))}
}

type eventMessage struct {
	Event   string `json:"event"`
	Channel string `json:"channel"`
}

func (c *ZbWebSocketClient) Connect() {
	if c.running {
		return
	}
	c.running = true

	dialer := &websocket.Dialer{}
	conn, _, err := dialer.Dial(WebSocketServerUrl, nil)
	c.conn = conn
	if err != nil {
		c.Disconnect()
		log.Fatalln("Fail to connect to " + WebSocketServerUrl + ", error: " + err.Error())
	}

	go func() {
		defer c.Disconnect()
		for {
			_, bytes, err := c.conn.ReadMessage()
			if err != nil {
				break
			}

			channel, _ := json.GetString(bytes, "channel")
			if decoder, ok := c.decoders[channel]; ok {
				value := decoder(bytes)
				if callback, ok := c.callbacks[channel]; ok {
					callback(value)
				}
			}
		}
	}()
}

func (c *ZbWebSocketClient) Disconnect() {
	if !c.running {
		return
	}
	c.running = false

	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *ZbWebSocketClient) SubscribeTicker(symbol string, callback func(ticker Ticker)) {
	channel := strings.Replace(symbol, "_", "", 1) + "_ticker"
	c.register(channel, func(value []byte) interface{} {
		return marshalTicker(value)
	}, func(v interface{}) {
		callback(v.(Ticker))
	})
	c.conn.WriteJSON(eventMessage{Event: "addChannel", Channel: channel})
}

func (c *ZbWebSocketClient) register(channel string, decoder func(value []byte) interface{}, callback func(interface{})) {
	c.registerDecoder(channel, decoder)
	c.registerCallback(channel, callback)
}

func (c *ZbWebSocketClient) registerDecoder(channel string, decoder func(value []byte) interface{}) {
	c.decoders[channel] = decoder
}

func (c *ZbWebSocketClient) registerCallback(channel string, callback func(interface{})) {
	c.callbacks[channel] = callback
}
