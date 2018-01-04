package main

import (
	"os"
	"time"

	"github.com/open-horizon/self-go-sdk/self"
	"github.com/satori/go.uuid"
)

func main() {
	conn, err := self.Init("localhost", "PubAgent")
	if err != nil {
		os.Exit(1)
	}
	time.Sleep(3 * time.Second)
	thing := self.Thing{
		GUID:        uuid.NewV4().String(),
		Type:        "IThing",
		CreateTime:  float64(time.Now().Unix()),
		Text:        "intent_text",
		Confidence:  0.9,
		Info:        "some_info",
		Name:        "some_name",
		State:       "ADDED",
		ECategory:   self.ThingCategoryPERCEPTION,
		FImportance: 1,
		FLifeSpan:   3600,
	}
	conn.Pub(self.Dot, thing)
	time.Sleep(1 * time.Second)
	conn.Close()
}
