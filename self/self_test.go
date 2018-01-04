package self

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

var targets = []Target{TargetAgentSociety, TargetBlackboardStream, TargetGestureManager, TargetSensorManager, TargetBlackboard}

func TestConn(t *testing.T) {
	conn, err := Init("localhost")
	if err != nil {
		t.Fail()
		logger.Fatalln(err)
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

func HandleImage(thing Thing){
	logger.Println(thing)
}

func TestReg(t *testing.T) {
	conn, err := Init("localhost")
	if err != nil {
		t.Fail()
	}
	conn.Sub(targets)
	conn.Reg("misc_types", MakeFilteredHandler(
		Not(
			Or(
				[]ThingFilter{
					MakeThingTypeFilter("Health"),
          MakeThingTypeFilter("IThing"),
          MakeThingTypeFilter("Proxy"),
          MakeThingTypeFilter("Failure"),
					MakeThingTypeFilter("RequestIntent"),
				},
			),
		),
		PrintThingType,
	))
	conn.Unreg("misc_types")
	conn.Reg("func2", MakeFilteredHandler(MakeThingTypeFilter("Text"), PrintThingText))
	conn.Unreg("func2")
	conn.Reg("image type", MakeFilteredHandler(MakeThingTypeFilter("Image"), HandleImage))
	time.Sleep(90 * time.Second)
	conn.Unsub(targets)
}
