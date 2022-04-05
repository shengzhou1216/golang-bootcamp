package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type App struct {
	ctx context.Context
	cancel func()
	sigs []os.Signal
}

func New() *App {
	var (
		ctx context.Context
		cancel context.CancelFunc
	)
	ctx = context.Background()
	ctx,cancel = context.WithCancel(ctx)
	return &App{
		ctx: ctx,
		cancel: cancel,
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
}

func (a *App) Run() error{
	wg := sync.WaitGroup{}
	eg,ctx := errgroup.WithContext(a.ctx)

	srv := &http.Server{
		Addr: "localhost:8080",
	}

	wg.Add(1)
	eg.Go(func() error {
		wg.Done()
		return srv.ListenAndServe()
	})
	wg.Wait()

	c := make(chan os.Signal,1)
	signal.Notify(c,a.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				if err := srv.Shutdown(ctx); err != nil {
					fmt.Printf("Failed to stop app %v.\n",err)
					return err
				}
			}
		}
	})
	log.Print("Server started")
	if err := eg.Wait(); err != nil && errors.Is(err,context.Canceled) {
		return err
	}
	return nil
}

func main()  {
	app := New()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}