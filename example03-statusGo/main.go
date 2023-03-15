package main

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type StatusGo struct {
	Name    string
	running bool
	mu      sync.Mutex
	ctx     context.Context
	TimeOut int64
	Start   chan struct{}
	Runing  chan struct{}
	Stop    chan struct{}
}

func NewStatusGo(name string) *StatusGo {
	return &StatusGo{
		Name:   name,
		ctx:    context.Background(),
		Start:  make(chan struct{}, 1),
		Runing: make(chan struct{}, 1),
		Stop:   make(chan struct{}, 1),
	}
}

func (s *StatusGo) SetTimeout(timing int64) {
	s.TimeOut = timing
}

func (s *StatusGo) start() {
	s.Start <- struct{}{}
}

func (s *StatusGo) setRunning(state bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = state
}
func (s *StatusGo) getRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *StatusGo) startGo(f func()) {
	ctx, _ := context.WithTimeout(s.ctx, time.Duration(s.TimeOut)*time.Second)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				// 处理panic
				fmt.Println("Recovered from panic:", r)
			}
			s.setRunning(false)
			s.Stop <- struct{}{}
		}()
		f()
		select {
		case <-ctx.Done():
			fmt.Println("Context Canceled")
			return
		}
		// 动态调用函数
		// reflect.ValueOf(f).Call(convertArgs(args))
	}()
}

// convertArgs 动态函数调用参数解析Args
func convertArgs(args []interface{}) []reflect.Value {
	values := make([]reflect.Value, len(args))
	for i, arg := range args {
		values[i] = reflect.ValueOf(arg)
	}
	return values
}

func (s *StatusGo) ControllerGo(f func()) {
	s.start()
	for {
		select {
		case <-s.Start:
			fmt.Println("Work Start")
			if !s.getRunning() {
				s.setRunning(true)
				s.startGo(f)
				s.Runing <- struct{}{}
			}
		case <-s.Runing:
			fmt.Println("Work Running")
		case <-s.Stop:
			fmt.Println("Work End")
			return
		}
	}
}

func Add(x, y int) int {
	fmt.Println(x + y)
	return x + y
}

func SleepTiming(mark int) {
	fmt.Printf("开摆-%d-号\n", mark)
	time.Sleep(time.Second * 15)
	fmt.Printf("摆烂结束,%d号\n", mark)
}

func main() {
	stateGo := NewStatusGo("Test")
	stateGo.SetTimeout(3)
	for i := 0; i < 3; i++ {
		stateGo.ControllerGo(func() {
			Add(1, 2)
		})
		stateGo.ControllerGo(func() {
			SleepTiming(i)
		})
	}
}
