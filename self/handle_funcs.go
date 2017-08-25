package self

// ThingHandlerFunc handels an incoming Thing
type ThingHandlerFunc func(Thing)

// PrintThingHandler just prints the thing
func PrintThingHandler(thing Thing) {
	logger.Println(thing)
}

// PrintThingText prints only the text of the thing.
func PrintThingText(thing Thing) {
	logger.Println("thing.Text:", thing.Text)
}

// PrintThingType prints only the Type of the thing.
func PrintThingType(thing Thing) {
	logger.Println("thing.Type:", thing.Type)
}

// ThingFilter desides if a Thing needs to be handled.
type ThingFilter func(Thing) bool

// DummyFilterTrue always returns true
func DummyFilterTrue(Thing) bool {
	return true
}

// DummyFilterFalse always returns false
func DummyFilterFalse(Thing) bool {
	return false
}

// MakeFilterNth returns a filter that returns true every n times.
func MakeFilterNth(n int) ThingFilter {
	var i int
	return func(thing Thing) bool {
		i++
		if i > n {
			i = 0
			return true
		}
		return false
	}
}

// MakeThingTypeFilter returns a Thing which returns true if the things type matches the given typeName
func MakeThingTypeFilter(typeName string) ThingFilter {
	return func(thing Thing) bool {
		return typeName == thing.Type
	}
}

// And returns a filter which returns true if none of the filters return false.
func And(filters []ThingFilter) ThingFilter {
	return func(thing Thing) bool {
		for _, filter := range filters {
			if !filter(thing) {
				return false
			}
		}
		return false
	}
}

// Or returns a filter which returns true if any of the filters return true.
func Or(filters []ThingFilter) ThingFilter {
	return func(thing Thing) bool {
		for _, filter := range filters {
			if filter(thing) {
				return true
			}
		}
		return false
	}
}

// Xor returns a filter which returns true if filter 1 xor filter2 return true.
func Xor(filter1, filter2 ThingFilter) ThingFilter {
	return func(thing Thing) bool {
		return filter1(thing) != filter2(thing)
	}
}

// Not returns a filter which inverts the provided filter.
func Not(filter ThingFilter) ThingFilter {
	return func(thing Thing) bool {
		return !filter(thing)
	}
}

// HasText returns true if the things text is non empty.
func HasText(thing Thing) bool {
	return len(thing.Text) > 0
}

// MakeFilteredHandler makes a MsgHandlerFunc that only called the given MsgHandlerFunc is filter returns true.
func MakeFilteredHandler(filter func(Thing) bool, handle ThingHandlerFunc) ThingHandlerFunc {
	return func(thing Thing) {
		if filter(thing) {
			handle(thing)
		}
	}
}
