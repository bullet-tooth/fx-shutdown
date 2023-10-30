package main

import (
	"context"
	"log"
	"sync"
	"time"

	"go.uber.org/fx"
)

type Manager struct {
	shutdown chan struct{}
	wg       *sync.WaitGroup
}

func New(lc fx.Lifecycle) *Manager {
	m := &Manager{
		shutdown: make(chan struct{}, 1),
		wg:       &sync.WaitGroup{},
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			m.StartBackgroundProcess()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			m.Shutdown()
			return nil
		},
	})
	return m
}

func (m *Manager) StartBackgroundProcess() {
	m.wg.Add(1)
	go func() {
		defer func() {
			log.Println("Exit coroutine")
			m.wg.Done()
		}()

		for {
			select {
			case <-m.shutdown:
				return
			case <-time.After(1 * time.Second):
				log.Println("Working...")
			}
		}
	}()
}

func (m *Manager) Shutdown() {
	m.shutdown <- struct{}{}
	log.Println("Shutdown")
	m.wg.Wait()
}

func main() {
	fx.New(
		fx.Provide(New),
		fx.Invoke(func(*Manager) {}),
	).Run()
}
