package main

import (
	"time"

	"github.com/open-horizon/self-go-sdk/self"
)

func main() {
	conn, err := self.Init("localhost")
	if err != nil {
		panic(err)
	}
	targets := []self.Target{self.TargetBlackboardStream, self.TargetAgentSociety}
	conn.Sub(targets)
	conn.Reg("misc_types", self.MakeFilteredHandler(
		self.Not(
			self.Or(
				[]self.ThingFilter{
					self.MakeThingTypeFilter("Health"),
					self.MakeThingTypeFilter("IThing"),
					self.MakeThingTypeFilter("Proxy"),
					self.MakeThingTypeFilter("Failure"),
					self.MakeThingTypeFilter("RequestIntent"),
				},
			),
		),
		self.PrintThingType,
	))
	handleFunc := func(thing self.Thing) {
		// do something with the thing
	}
	conn.Reg("handle func 1", handleFunc)
	conn.Reg("print_text", self.MakeFilteredHandler(self.MakeThingTypeFilter("Text"), self.PrintThingText))
	time.Sleep(20 * time.Second)
	conn.Unreg("handle func 1")
	conn.Unreg("print_text")
	conn.Unsub(targets)
}
