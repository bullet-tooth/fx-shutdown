package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Manager struct {
}

func New(lc fx.Lifecycle) *Manager {
	m := &Manager{}

	ctx, cancel := context.WithCancel(context.Background())
	bg, bgCtx := errgroup.WithContext(ctx)

	lc.Append(fx.Hook{
		OnStart: func(appCtx context.Context) error {
			m.StartBackgroundProcess(bgCtx, bg)
			return nil
		},
		OnStop: func(_ context.Context) error {
			cancel()
			return bg.Wait()
		},
	})
	return m
}

func (m *Manager) StartBackgroundProcess(bgCtx context.Context, bg *errgroup.Group) {
	bg.Go(func() error {
		defer func() {
			log.Println("Exit coroutine")
		}()
		for {
			select {
			case <-bgCtx.Done():
				log.Println("received context done")
				return nil
			case <-time.After(1 * time.Second):
				log.Println("Working...")
			}
		}
	})
}

func main() {
	fx.New(
		fx.Provide(New),
		fx.Invoke(func(*Manager) {}), // just to invoke the manager
	).Run()
}
