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
	Type   string `json:"type"`
	Parent string `json:"parent"`
	Thing  Thing  `json:"thing"`
}

// Thing is a thing
type Thing struct {
	GUID           string            `json:"GUID_"`
	Type           string            `json:"Type_"`
	DataType       string            `json:"m_DataType"`
	Children       []Thing           `json:"m_Children"`
	Confidence     float64           `json:"m_Confidence"`
	CreateTime     float64           `json:"m_CreateTime"`
	Info           string            `json:"m_Info"`
	Name           string            `json:"m_Name"`
	ProxyID        string            `json:"m_ProxyID"`
	State          string            `json:"m_State"`
	Threshold      float64           `json:"m_Threshold"`
	ClassifyIntent bool              `json:"m_ClassifyIntent"`
	Language       string            `json:"m_Language"`
	LocalDialog    bool              `json:"m_LocalDialog"`
	Text           string            `json:"m_Text"`
	ECategory      ThingCategory     `json:"m_eCategory"`
	FImportance    int               `json:"m_fImportance"`
	FLifeSpan      int               `json:"m_fLifeSpan"`
	Data           map[string]string `json:"m_Data"`
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
	TargetConfig           Target = "config"
	TargetDot              Target = "."
)

// ThingCategory are Intu supports different types of things.
// Things belonging to different types are not attached to each other directly but rather a proxy is created and connected via GUID
type ThingCategory int

// Categorys of things
const (
	ThingCategoryINVALID    ThingCategory = -1
	ThingCategoryPERCEPTION ThingCategory = 0
	ThingCategoryAGENCY     ThingCategory = 1
	ThingCategoryMODEL      ThingCategory = 2
)

// ThingEventType Represents different types of an event related to things such as whether a thing has been added, removed, or changed
type ThingEventType int

// Type of thing events
const (
	ThingEventTypeNONE       ThingEventType = 0
	ThingEventTypeADDED      ThingEventType = 1
	ThingEventTypeREMOVED    ThingEventType = 2
	ThingEventTypeSTATE      ThingEventType = 4
	ThingEventTypeIMPORTANCE ThingEventType = 8
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
		logger.Println(err)
		return
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

// Close the conn
func (conn *Conn) Close() {
	conn.conn.Close()
}

// Pub publishes a message to a topic
func (conn *Conn) Pub(target Target, thing Thing) (err error) {
	messageData := msgData{
		Type:  "IThing",
		Event: "add_object",
		Thing: thing,
	}

	dataBytes, err := json.Marshal(messageData)
	if err != nil {
		return
	}
	message := msg{
		Topic:     "blackboard",
		Targets:   []Target{target},
		Data:      string(dataBytes),
		Origin:    conn.selfID,
		Msg:       "publish_at",
		Binary:    false,
		Persisted: true,
		//Type:      "IThing",
	}
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return
	}
	if err = conn.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		return
	}
	return
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
func Init(host string, selfID string) (conn *Conn, err error) {
	u := url.URL{Scheme: "ws", Host: host + ":9443", Path: "/stream"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"selfId": []string{selfID}, "orgId": []string{}, "token": []string{""}})
	if err != nil {
		return
	}
	conn = &Conn{
		conn:     c,
		selfID:   selfID,
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
					//logger.Println("Origin:", message.Origin, "Topic:", message.Topic)
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
