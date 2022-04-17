package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// func get_db_data(ctx context.Context ) error  {	
// 	log.Println("get db data")
// 	select {
// 	case <-ctx.Done():
// 		log.Println("get_db_data canceled.")
// 	}
// 	return nil
// }

// func get_rpc_data(cxt context.Context) error {
// 	log.Println("get rpc data")
// 	select {
// 	case <-cxt.Done():
// 		log.Println("get_rpc_data canceled.")
// 	}
// 	return nil
// }

// func helloHandler(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()
// 	go get_db_data(ctx)
// 	go get_rpc_data(ctx)
// 	select {
// 	case <-time.After(10 * time.Second):
// 		w.Write([]byte("hello"))
// 	case <-ctx.Done():
// 		log.Println("request canceled.")
// 	}
// }

func server2() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	eg, ctx := errgroup.WithContext(ctx)
	http.HandleFunc("/", helloHandler)
	srv := &http.Server{
		Addr: ":8000",
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGKILL)

	eg.Go(func() error {
		<-ctx.Done()
		defer cancel()
		return srv.Shutdown(ctx)
	})
	eg.Go(func() error {
		return srv.ListenAndServe()
	})

	eg.Go(func() error {
		for {
			select {
			case signal := <-done:
				log.Printf("accept exit signal: %v\n", signal)
				cancel()
				// return srv.Shutdown(ctx)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	})



	log.Println("Server started.")
	
	if err := eg.Wait(); err != nil {
		log.Println(err)
	}

	log.Println("Server Exited Properly.")
}
