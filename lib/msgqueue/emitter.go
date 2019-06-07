package msgqueue

//Event Emiiter is the interface that can be imported anywhere and it will
//actually emit events that are passed to it.
type EventEmitter interface {
	Emit(event Event) error
}
