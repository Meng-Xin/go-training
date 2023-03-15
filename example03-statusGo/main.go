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
	Start   chan struct{}
	Runing  chan struct{}
	Stop    chan struct{}
}

func NewStatusGo(name string) *StatusGo {
	return &StatusGo{
		Name:   name,
		Start:  make(chan struct{}, 1),
		Runing: make(chan struct{}, 1),
		Stop:   make(chan struct{}, 1),
	}
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
		case <-s.ctx.Done():
			fmt.Println("Work Timeout")
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
	for i := 0; i < 3; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		stateGo.ctx = ctx
		stateGo.ControllerGo(func() {
			Add(1, 2)
		})
		stateGo.ControllerGo(func() {
			SleepTiming(i)

		})
	}
}
