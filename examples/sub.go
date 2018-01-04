package main

import (
	"fmt"
	"time"

	"github.com/open-horizon/self-go-sdk/self"
)

func main() {
	conn, err := self.Init("localhost", "SubAgent")
	if err != nil {
		panic(err)
	}
	targets := []self.Target{self.TargetBlackboard, self.TargetAgentSociety, self.TargetBlackboardStream, self.TargetSensorManager}
	conn.Sub(targets)
	conn.Reg("emotion",
		func(thing self.Thing) {
			fmt.Println("type:", thing.Type, "name:", thing.Name)
		},
	)
	fmt.Println("subbed")
	time.Sleep(90 * time.Second)
	conn.Unsub(targets)
}
