package chat

import (
	"net/url"
	"sync"

	"github.com/apex/log"
	"github.com/gorilla/websocket"
)

type WebsocketMessage struct {
	MessageType int
	Bytes       []byte
}

type WebsocketConnection struct {
	url        *url.URL
	connection *websocket.Conn
	connecting sync.Mutex
	incoming   chan WebsocketMessage
	outgoing   chan WebsocketMessage
	log        log.Interface
}

func NewWebsocketConnection(url *url.URL) *WebsocketConnection {
	return &WebsocketConnection{
		url: url,
		log: log.WithField("connection", url.Hostname()),
	}
}

func (w *WebsocketConnection) GetMessages() (chan WebsocketMessage, error) {
	err := w.ensureConnection()
	return w.incoming, err
}

func (w *WebsocketConnection) SendMessage(message WebsocketMessage) error {
	err := w.ensureConnection()
	if err != nil {
		return err
	}
	w.outgoing <- message
	return nil
}

func (w *WebsocketConnection) ensureConnection() error {
	w.connecting.Lock()
	defer w.connecting.Unlock()
	if w.connection != nil {
		return nil
	}
	var err error
	w.log.Debug("Connecting")
	w.connection, _, err = websocket.DefaultDialer.Dial(w.url.String(), nil)
	if w.connection != nil {
		w.log.Debug("Connected!")
		w.incoming = make(chan WebsocketMessage)
		w.outgoing = make(chan WebsocketMessage)
		go w.readMessages()
		go w.writeMessages()
	}
	return err
}

func (w *WebsocketConnection) readMessages() {
	for {
		messageType, message, err := w.connection.ReadMessage()
		if err != nil {
			w.log.Error(err.Error())
			break
		}
		w.incoming <- WebsocketMessage{MessageType: messageType, Bytes: message}
	}
	close(w.incoming)
	close(w.outgoing)
	w.connection.Close()
	w.connection = nil
}

func (w *WebsocketConnection) writeMessages() {
	for message := range w.outgoing {
		w.connection.WriteMessage(message.MessageType, message.Bytes)
	}
}
