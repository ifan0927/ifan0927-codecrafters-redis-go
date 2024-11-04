package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

type EventType int

const (
	EventRead EventType = iota
	EventWrite
)

type Event struct {
	Type EventType
	Conn net.Conn
}

type EventLoop struct {
	Running bool
	Queue   chan Event
	handler *EventHandler
}
type EventHandler struct {
	Loop *EventLoop
}

func (h *EventHandler) handleEvent(event Event) {
	if event.Type == EventRead {
		write_event := Event{
			Type: EventWrite,
			Conn: event.Conn,
		}
		h.Loop.AddEvent(write_event)

		read_event := Event{
			Type: EventRead,
			Conn: event.Conn,
		}
		h.Loop.AddEvent(read_event)

	} else if event.Type == EventWrite {
		event.Conn.Write([]byte("+PONG\r\n"))
		return
	}
}
func NewEventLoop() *EventLoop {
	el := &EventLoop{
		Running: false,
		Queue:   make(chan Event),
	}
	handler := &EventHandler{
		Loop: el,
	}
	el.handler = handler
	return el
}
func (el *EventLoop) Start() {
	if el.Running {
		return
	}
	el.Running = true

	go func() {
		for {
			event := <-el.Queue
			el.handler.handleEvent(event)
		}
	}()
}

func (el *EventLoop) AddEvent(e Event) {
	select {
	case el.Queue <- e:
		fmt.Printf("Added event type: %d, queue length: %d\n", e.Type, len(el.Queue))
	default:
		fmt.Println("Queue is full, dropping event")
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests. 123
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	el := NewEventLoop()
	el.Start()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		e := Event{
			Type: EventRead,
			Conn: c,
		}
		el.AddEvent(e)
	}

}
