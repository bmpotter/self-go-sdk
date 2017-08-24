package self

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

var targets = []Target{TargetBlackboard, TargetAgentSociety, TargetBlackboardStream, TargetGestureManager, TargetSensorManager}

func TestConn(t *testing.T) {
	conn, err := Init("localhost")
	if err != nil {
		t.Fail()
	}

	conn.Sub(targets)
	time.Sleep(1 * time.Second)
	conn.Unsub(targets)
}

func TestPub(t *testing.T) {
	conn, err := Init("localhost")
	if err != nil {
		t.Fail()
	}
	message := msgData{
		Event: "dummy",
		Thing: Thing{
			GUID:       uuid.NewV4().String(),
			Type:       "fooType",
			CreateTime: float64(time.Now().Unix()),
			Text:       "foo bar text",
		},
	}
	conn.Pub(TargetBlackboard, message)
}

func TestReg(t *testing.T) {
	conn, err := Init("localhost")
	if err != nil {
		t.Fail()
	}
	conn.Sub(targets)
  handleFunc := func(thing Thing){
    logger.Println(thing.Info, thing.Text, thing.Type, thing.State)
  }

	conn.Reg("demo1", handleFunc)
	time.Sleep(30 * time.Second)
	conn.Unreg("demo1")
	conn.Unsub(targets)
}
