package self

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"os"
)

var logger = log.New(os.Stdout, "self: ", log.Lshortfile)

type msg struct {
	Binary    bool     `json:"binary"`
	Data      string   `json:"data"`
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
	CreateType     float64 `json:"m_CreateTime"`
	Info           string  `json:"m_Info"`
	Name           string  `json:"m_Name"`
	ProxyID        string  `json:"m_ProxyID"`
	State          string  `json:"m_State"`
	Threshold      float64 `json:"m_Threshold"`
	ClassifyIntent bool    `json:"m_ClassifyIntent"`
	CreateTime     float64 `json:"m_CreateTime"`
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
	conn   *websocket.Conn
	selfID string
}

// Sub subscribes to a topic
func (conn *Conn) Sub(target Target) {
	subTopic := msg{
		Targets: []Target{target},
		Msg:     "subscribe",
		Origin:  conn.selfID,
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
func (conn *Conn) Unsub(target Target) {
	subTopic := msg{
		Targets: []Target{target},
		Msg:     "unsubscribe",
		Origin:  conn.selfID,
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
func (conn *Conn) Pub(target Target, Msg string) {
	subTopic := msg{
		Targets: []Target{target},
		Msg:     Msg,
		Origin:  conn.selfID,
	}
	msg, err := json.Marshal(subTopic)
	if err != nil {
		panic(err)
	}
	if err = conn.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
		logger.Println(err)
	}
}

// Init the connection
func Init(host string) (conn *Conn, err error) {
	u := url.URL{Scheme: "ws", Host: host + ":9443", Path: "/stream"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"selfId": []string{""}, "orgId": []string{}, "token": []string{""}})
	if err != nil {
		return
	}
	conn = &Conn{
		conn:   c,
		selfID: "/.",
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
					//logger.Println(message.Topic)
					var messageData msgData
					if err := json.Unmarshal([]byte(message.Data), &messageData); err != nil {
						logger.Println(err)
						logger.Println(message.Data)
					} else {
						logger.Println(messageData.Thing.Text)
					}
				}
			}
		}
	}()
	return
}
