package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net"

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

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

func get_db_data(ctx context.Context ) error  {	
	log.Println("get db data")
	select {
	case <-ctx.Done():
		log.Println("get_db_data canceled.")
	}
	return nil
}

func get_rpc_data(cxt context.Context) error {
	log.Println("get rpc data")
	select {
	case <-cxt.Done():
		log.Println("get_rpc_data canceled.")
	}
	return nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	go get_db_data(ctx)
	go get_rpc_data(ctx)
	select {
	case <-time.After(10 * time.Second):
		w.Write([]byte("hello"))
	case <-ctx.Done():
		log.Println("request canceled.")
	}
}

func (a *App) Run() error{
	eg,ctx := errgroup.WithContext(a.ctx)

	http.HandleFunc("/",helloHandler)
	srv := &http.Server{
		Addr: "localhost:8080",
	}

	// 设置http server的 BaseContext 为 App.ctx
	srv.BaseContext = func(l net.Listener) context.Context {
		return a.ctx
	}

	eg.Go(func() error {
		<-ctx.Done() // 监听App.ctx 取消
		log.Println("Cancelling http server...")
		stopCtx,cancel := context.WithTimeout(ctx,5 * time.Second)
		defer cancel()
		return srv.Shutdown(stopCtx)
	})

	eg.Go(func() error {
		return srv.ListenAndServe()
	})

	c := make(chan os.Signal,1)
	signal.Notify(c,a.sigs...)
	// 处理 App 退出信号
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case signal := <-c:
				log.Printf("accept exit signal: %v\n", signal)
				return a.Stop()
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
		log.Println(err)
	}
	log.Print("Server Exited Properly!")
}