package main

import (
	"context"
	"errors"
	"log"
	"time"
)

func operation1(ctx context.Context) error {
	time.Sleep(10 * time.Millisecond)
	return errors.New("failed")
}

func operation2(ctx context.Context) {
	select {
	case <-time.After(500 * time.Millisecond):
		log.Println("done")
	case <-ctx.Done():
		log.Println("halted operation2")

	}
}

func emmit_cancel() {
	ctx := context.Background()

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		err := operation1(ctx)
		if err != nil {
			cancel()
		}
	}()

	operation2(ctx)
}
