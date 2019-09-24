package event

import "github.com/kumparan/kumnats"

// NatsHelloMessage :nodoc:
type NatsHelloMessage kumnats.NatsMessage

const (
	// TypeHello :nodoc:
	TypeHello = kumnats.EventType("hello")

	// NatsHelloChannel :nodoc:
	NatsHelloChannel = "hello-channel"
)
