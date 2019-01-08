package poller

// Event represents a file event
type Event uint8

const (
	// EventNil represents an unset event
	EventNil Event = iota
	// EventCreate represents a create event
	EventCreate
	// EventWrite represents a write event
	EventWrite
	// EventRemove represents a remove event
	EventRemove
	// EventChmod represents a chmod event
	EventChmod
)

func (e Event) String() string {
	switch e {
	case EventCreate:
		return "CREATE"
	case EventWrite:
		return "WRITE"
	case EventRemove:
		return "REMOVE"
	case EventChmod:
		return "CHMOD"

	default:
		return ""
	}
}
