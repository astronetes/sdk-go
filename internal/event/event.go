package event

import "context"

// All custom events names must be of this type.
type Name string

// All custom event types must satisfy this interface.
type Event interface {
	Handle(ctx context.Context)
}

type Dispatcher[T ~string] struct {
	events map[T]Listener[T]
}

func NewDispatcher[T ~string]() *Dispatcher[T] {
	return &Dispatcher[T]{
		events: make(map[T]Listener[T]),
	}
}

type Listener[T any] interface {
	Listen(ctx context.Context, event T)
}
