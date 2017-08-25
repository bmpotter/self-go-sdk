// Package self lets you connect to a self instance via the WS API.
// To run a self instance, see https://github.com/watson-intu/self, or https://github.com/michaeldye/self-docker-builder
package self

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var logger = log.New(os.Stdout, "self: ", log.Lshortfile)

type msg struct {
	Binary    bool     `json:"binary"`
	Data      string   `json:"data"` // json msgData
	Msg       string   `json:"msg"`
	Origin    string   `json:"origin"`
	Persisted bool     `json:"persisted"`
	Targets   []Target `json:"targets"`
	Topic     string   `json:"topic"`
	Request   string   `json:"request"`
	Type      string   `json:"type"`
}

type msgData struct {
	Event  string `json:"event"`
	Parent string `json:"parent"`
	Thing  Thing  `json:"thing"`
}

// Thing is a thing
type Thing struct {
	GUID           string  `json:"GUID_"`
	Type           string  `json:"Type_"`
	Children       []Thing `json:"m_Children"`
	Confidence     float64 `json:"m_Confidence"`
	CreateTime     float64 `json:"m_CreateTime"`
	Info           string  `json:"m_Info"`
	Name           string  `json:"m_Name"`
	ProxyID        string  `json:"m_ProxyID"`
	State          string  `json:"m_State"`
	Threshold      float64 `json:"m_Threshold"`
	ClassifyIntent bool    `json:"m_ClassifyIntent"`
	Language       string  `json:"m_Language"`
	LocalDialog    bool    `json:"m_LocalDialog"`
	Text           string  `json:"m_Text"`
	ECategory      int     `json:"m_eCategory"`
	FImportance    int     `json:"m_fImportance"`
	FLifeSpan      int     `json:"m_fLifeSpan"`
}

// Target is a target that can be subscribed to
type Target string

// Targets
const (
	TargetBlackboard       Target = "blackboard"
	TargetAgentSociety     Target = "agent-society"
	TargetBlackboardStream Target = "blackboard-stream"
	TargetGestureManager   Target = "gesture-manager"
	TargetSensorManager    Target = "sensor-manager"
	TargetModels           Target = "sensors"
)

// Conn is a ws conn + metadata
type Conn struct {
	mutex    *sync.Mutex
	handlers map[string]ThingHandlerFunc
	conn     *websocket.Conn
	selfID   string
}

// Sub subscribes to a topic
func (conn *Conn) Sub(targets []Target) {
	subTopic := msg{
		Targets: targets,
		Msg:     "subscribe",
		Origin:  "/.",
	}
	msg, err := json.Marshal(subTopic)
	if err != nil {
		panic(err)
	}
	if err = conn.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		logger.Println(err)
	}
}

// Unsub unsubscribes from a topic
func (conn *Conn) Unsub(targets []Target) {
	subTopic := msg{
		Targets: targets,
		Msg:     "unsubscribe",
		Origin:  "/.",
	}
	msg, err := json.Marshal(subTopic)
	if err != nil {
		panic(err)
	}
	if err = conn.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		logger.Println(err)
	}
}

// Pub publishes a message to a topic
func (conn *Conn) Pub(target Target, data msgData) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	message := msg{
		Targets: []Target{target},
		Data:    string(dataBytes),
		Origin:  conn.selfID,
		Msg:     "publish",
	}
	msgBytes, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	if err = conn.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		logger.Println(err)
	}
}

// Reg registers a function to the connection
func (conn *Conn) Reg(name string, handlerFunc ThingHandlerFunc) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.handlers[name] = handlerFunc
	return
}

// Unreg removes a handler from the connection
func (conn *Conn) Unreg(name string) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	delete(conn.handlers, name)
	return
}

// Init the connection
func Init(host string) (conn *Conn, err error) {
	u := url.URL{Scheme: "ws", Host: host + ":9443", Path: "/stream"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"selfId": []string{""}, "orgId": []string{}, "token": []string{""}})
	if err != nil {
		return
	}
	conn = &Conn{
		conn:     c,
		selfID:   "/.",
		handlers: make(map[string]ThingHandlerFunc),
		mutex:    &sync.Mutex{},
	}
	go func() {
		for {
			_, msgBytes, err := c.ReadMessage()
			if err != nil {
				logger.Println(err)
			} else {
				var message msg
				if err := json.Unmarshal(msgBytes, &message); err != nil {
					logger.Println(err)
				} else {
					var messageData msgData
					if err := json.Unmarshal([]byte(message.Data), &messageData); err != nil {
					} else {
						conn.mutex.Lock()
						for _, handle := range conn.handlers {
							go handle(messageData.Thing)
						}
						conn.mutex.Unlock()
					}
				}
			}
		}
	}()
	return
}
