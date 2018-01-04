# self-go-sdk

Golang library to interact with self via the WS APIs

[![GoDoc](https://godoc.org/github.com/open-horizon/self-go-sdk/self?status.svg)](https://godoc.org/github.com/open-horizon/self-go-sdk/self)

## Usage

``` golang
import (
  "github.com/open-horizon/self-go-sdk"
)
```

Initialize a connection to a self instance.

``` golang
conn, err := self.Init("localhost")
```

Subscribe to a list of target.

``` golang
conn.Sub([]self.Target{self.TargetBlackboardStream, self.TargetAgentSociety})
```

Register an event hander. Give it some unique name so you can unregister it latter.

``` golang
handleFunc := func(thing self.Thing){
  // do something with the thing
}
conn.Reg("handle func 1", handleFunc)
```

Once you are done, unregister it.
``` golang
conn.Unreg("handle func 1")
```

handleFunc will be called for all things received. This is often undesirable.

`MakeFilteredHandler()` constructs a ThingHandlerFunc which only calls the main handler func if filterFunc returns true.

``` golang
filterFunc := func(thing self.Thing) bool {
  if thing is desirable {
    return true
  }
  return false
}

handleFunc := func(thing self.Thing){
  // do something with the thing
}
conn.Reg("handle func 1", self.MakeFilteredHandler(filterFunc, handleFunc)))
```

We can combine `ThingFilter`s. Say that we want to print the text of the things whose types are anything but certain types.

``` golang
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
  self.PrintThingText,
))
```
