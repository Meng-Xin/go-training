package main

import (
	"fmt"
	"time"
)

type StateMachine struct {
	currentState string
	stateMap     map[string]map[string]string
	eventChMap   map[string]chan struct{}
}

func NewStateMachine(initialState string, stateMap map[string]map[string]string) *StateMachine {
	stateMachine := &StateMachine{
		currentState: initialState,
		stateMap:     stateMap,
		eventChMap:   make(map[string]chan struct{}),
	}

	for state := range stateMap {
		stateMachine.eventChMap[state] = make(chan struct{})
	}

	go stateMachine.run()

	return stateMachine
}

func (sm *StateMachine) run() {
	for {
		select {
		case <-sm.eventChMap[sm.currentState]:
			nextState := sm.stateMap[sm.currentState]["event"]
			fmt.Printf("State transition from %s to %s\n", sm.currentState, nextState)
			sm.currentState = nextState
		}
	}
}

func (sm *StateMachine) SendEvent(event string) {
	sm.eventChMap[sm.currentState] <- struct{}{}
}

func main() {
	stateMap := map[string]map[string]string{
		"start": {
			"event": "pause",
		},
		"pause": {
			"event": "resume",
		},
		"resume": {
			"event": "stop",
		},
	}

	stateMachine := NewStateMachine("start", stateMap)

	time.Sleep(1 * time.Second)
	stateMachine.SendEvent("event")

	time.Sleep(1 * time.Second)
	stateMachine.SendEvent("event")

	time.Sleep(1 * time.Second)
	stateMachine.SendEvent("event")
}
