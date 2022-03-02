package part04

import (
	"log"
	"strconv"
	"testing"
)

func TestChanDone(t *testing.T) {
	done := make(chan struct{}, 1)
	go func() {
		log.Println("abc")
		close(done)
	}()
	<-done
}

func TestMultiChanDone(t *testing.T) {
	done := make(chan struct{}, 3)
	for i := 0; i < 3; i++ {
		go func() {
			log.Println("abc")
			done <- struct{}{}
		}()
	}

	for i := 0; i < 3; i++ {
		<-done
	}
}

// channel messaging
func TestChanMsg(t *testing.T) {
	ch := make(chan string, 3)
	go func() {
		defer close(ch) // finish action close done chan
		for i := 0; i < 3; i++ {
			log.Println("current i := ", i)
			ch <- "hello," + strconv.Itoa(i)
		}
	}()

	for s := range ch {
		log.Println("str: ", s)
	}
}
