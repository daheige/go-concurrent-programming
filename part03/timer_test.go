package part03

import (
	"log"
	"testing"
	"time"
)

/*
time.Timer timers must be created with the time.NewTimer,
time.AfterFunc or time.After functions. When the timer expires,
the expiration time will be sent to the channel held by the timer,
and the goroutine that subscribes to the channel will receive
the timer expiration time.
*/

func TestTimer(t *testing.T) {
	done := make(chan struct{}, 1)
	go func() {
		defer close(done)

		log.Println("2222")
		time.Sleep(10 * time.Second)
	}()

	// select {
	// case <-time.After(2 * time.Second):
	// 	log.Println("timeout")
	// case <-done:
	// 	log.Println("done")
	// }

	// timer
	tn := time.NewTimer(2 * time.Second)
	select {
	case <-tn.C:
		tn.Stop()
		log.Println("timeout")
	case <-done:
		log.Println("done")
	}

	// ticker
	// Ticker can trigger time events periodically,
	// each time the specified time interval is reached,
	// the event will be triggered.
	// time.Ticker needs to be created by
	// time.NewTicker or time.Tick.
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for tc := range ticker.C {
			log.Println("ticker at: ", tc)
		}
	}()

	time.Sleep(5 * time.Second)
	ticker.Stop()
	log.Println("ticker stop")
}
