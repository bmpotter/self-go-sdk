package main

import (
	"fmt"
	"time"

	"github.com/open-horizon/self-go-sdk/self"
)

func main() {
	conn, err := self.Init("localhost", "")
	if err != nil {
		panic(err)
	}
	targets := []self.Target{self.TargetBlackboardStream, self.TargetAgentSociety}
	conn.Sub(targets)
	conn.Reg("misc_types", self.MakeFilteredHandler(
		self.MakeThingTypeFilter("IThing"),
		func(thing self.Thing) {
			fmt.Println("type:", thing.Type, "name:", thing.Name, "text:", thing.Text, "conf:", thing.Confidence, "data:", thing.Data, "dataType:", thing.DataType)
		},
	))
	time.Sleep(2000 * time.Minute)
	conn.Unreg("handle func 1")
	conn.Unreg("print_text")
	conn.Unsub(targets)
}
