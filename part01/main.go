/*
The cancellation function needs to be implemented from two aspects to complete:
	 Listen for cancel events

	 emit cancel event
*/
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// The Context type of the Go language context standard library provides a Done() method,
// which returns a channel of type <-chan struct{}. This channel will receive a value of
// type struct{} every time the context receives a cancellation event.
// So listening to the cancellation event in the Go language is to
// wait to receive <-ctx.Done()

func httpServer(address string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.Println("processing request...")
		select {
		case <-time.After(2 * time.Second): // request handler in 2 second
			w.Write([]byte("request processed"))
		case <-ctx.Done():
			// When the browser is closed,
			// the following operations are performed
			// 2022/03/01 21:05:08 request cancelled
			log.Println("request cancelled")
		}
	})

	server := http.Server{
		Addr:         fmt.Sprintf(address),
		Handler:      mux,
		WriteTimeout: 3 * time.Second,
		ReadTimeout:  3 * time.Second,
	}

	log.Fatalln(server.ListenAndServe())
}

// emit cancel event
// If you have an action that can be canceled, you must emit a cancel event through the context.
// This can be done through the cancel function returned by the WithCancel function of
// the context package (withCancel also returns a context object that supports cancellation).
// This function takes no arguments and returns nothing. It is called when the context
// needs to be canceled, and a cancel event is emitted.

func task1(ctx context.Context) error {
	time.Sleep(600 * time.Millisecond)
	return errors.New("failed")
}

func task2(ctx context.Context) {
	// get value from ctx
	log.Println("name: ", ctx.Value("name"))
	select {
	case <-time.After(500 * time.Millisecond):
		log.Println("timeout")
	case <-ctx.Done():
		log.Println("handler task2")
	}
}

func doSomething() {
	// create a ctx
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	// run task2 in goroutine
	go func() {
		task2(ctx)
	}()

	err := task1(ctx)
	if err != nil {
		cancel() // cancel task2 run
	}
}

func doSomething2() {
	ctx := context.Background()
	// create a ctx with timeout
	ctx, cancel := context.WithTimeout(ctx, 600*time.Millisecond)
	ctx = context.WithValue(ctx, "name", "daheige")
	defer cancel()
	// run task2 in goroutine
	go func() {
		task2(ctx)
	}()

	err := task1(ctx)
	if err != nil {
		log.Println("err: ", err)
	}
}

func main() {
	// httpServer(":8080")

	// doSomething()

	doSomething2()

	time.Sleep(2 * time.Second)
}

// This context will be cancelled after 3 seconds
// If you need to cancel before expiration you can use the cancel
// function as in the previous example
// ctx, cancel := context.WithTimeout(ctx, 3*time.Second)

// The context will be cancelled at 2009-11-10 23:00:00
// ctx, cancel := context.WithDeadline(ctx, time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
