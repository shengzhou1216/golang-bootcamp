package main

import (
	"log"
	"net/http"
	"time"
)

func context_cancel()  {
	http.ListenAndServe(":8000",http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("processing reqeust")
		select {
		case <-time.After(2 * time.Second):
			w.Write([]byte("request processed"))
		case <-ctx.Done():
			log.Println("request canceled.")
		}
	}))
}