package msgqueue

//This is an interface. contracts/event_created.go fits this interface
type Event interface {
	EventName() string
}
