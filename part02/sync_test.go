package part02

import (
	"bytes"
	"io/ioutil"
	"log"
	"runtime"
	"sync"
	"testing"
)

func TestNoMutex(t *testing.T) {
	done := make(chan struct{}, 4)
	counter := 0
	for i := 0; i < 4; i++ {
		go func(i int) {
			log.Println("current index: ", i)
			val := counter
			runtime.Gosched()

			val++
			counter = val

			done <- struct{}{}

		}(i)
	}

	for i := 0; i < 4; i++ {
		<-done
	}

	// counter maybe :1,2,3,4
	log.Println("counter=", counter)
}

/*
2022/03/02 09:11:35 current index:  3
2022/03/02 09:11:35 current index:  0
2022/03/02 09:11:35 current index:  2
2022/03/02 09:11:35 current index:  1
2022/03/02 09:11:35 counter= 2
--- PASS: TestNoMutex (0.00s)
PASS
2022/03/02 09:11:56 current index:  3
2022/03/02 09:11:56 current index:  2
2022/03/02 09:11:56 current index:  1
2022/03/02 09:11:56 current index:  0
2022/03/02 09:11:56 counter= 3
--- PASS: TestNoMutex (0.00s)
PASS
*/
func TestMutex(t *testing.T) {
	done := make(chan struct{}, 4)
	counter := 0
	// Read and write shared variables protected by mutex
	mu := &sync.Mutex{}
	for i := 0; i < 4; i++ {
		go func(i int) {
			mu.Lock()
			log.Println("current index: ", i)
			val := counter
			runtime.Gosched()

			val++
			counter = val
			mu.Unlock()
			done <- struct{}{}

		}(i)
	}

	for i := 0; i < 4; i++ {
		<-done
	}

	// counter = 4
	log.Println("counter=", counter)
}

/*
=== RUN   TestMutex
2022/03/02 09:14:10 current index:  3
2022/03/02 09:14:10 current index:  1
2022/03/02 09:14:10 current index:  0
2022/03/02 09:14:10 current index:  2
2022/03/02 09:14:10 counter= 4
--- PASS: TestMutex (0.00s)
PASS
*/

func TestMutexRwlock(t *testing.T) {
	done := make(chan struct{}, 4)
	counter := 0
	// Read and write shared variables protected by mutex rw lock
	mu := &sync.RWMutex{}

	// write data
	for i := 0; i < 4; i++ {
		go func(i int) {
			log.Println("current index: ", i)
			mu.Lock()
			counter++
			mu.Unlock()
			done <- struct{}{}
		}(i)
	}

	for i := 0; i < 4; i++ {
		<-done
	}

	// read data
	done2 := make(chan struct{})
	go func() {
		mu.RLock()
		log.Println("current counter: ", counter)
		mu.RUnlock()
		close(done2)
	}()

	<-done2

	// counter = 4
	log.Println("counter=", counter)
}

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		log.Println("1111")
	}()
	go func() {
		defer wg.Done()
		log.Println("22222")
	}()

	// wg cannot be copied, if passed to a function,
	// it needs to be passed by pointer
	// fatal error: all goroutines are asleep - deadlock!
	// 'func' passes lock by value: type 'sync.WaitGroup' contains 'interface{}' which is 'sync.Locker'
	//
	// The underlying source of WaitGroup
	// type WaitGroup struct {
	//    noCopy noCopy
	//    state1 [3]uint32
	// }
	// A WaitGroup waits for a collection of goroutines to finish.
	// The main goroutine calls Add to set the number of goroutines to wait for. Then each of the goroutines
	// runs and calls Done when finished. At the same time, Wait can be used to block until
	// all goroutines have finished.
	// A WaitGroup must not be copied after first use.
	go func(wg sync.WaitGroup) {
		defer wg.Done()

		log.Println("33333")
	}(wg)

	wg.Wait()

	log.Println("ok")

}

func TestWaitGroupNoCopy(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		log.Println("1111")
	}()
	go func() {
		defer wg.Done()
		log.Println("22222")
	}()

	// pass wg by pointer
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		log.Println("33333")
	}(&wg)

	wg.Wait()
	log.Println("ok")
}

/*
=== RUN   TestWaitGroupNoCopy
2022/03/02 09:21:54 33333
2022/03/02 09:21:54 1111
2022/03/02 09:21:54 22222
2022/03/02 09:21:54 ok
--- PASS: TestWaitGroupNoCopy (0.00s)
PASS
*/

func TestMap(t *testing.T) {
	m := sync.Map{} // Concurrent safe map operations
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()

		m.Store("abc", "hello")

		value, ok := m.Load("abc")
		log.Println("abc: ", value, "ok: ", ok)
	}()

	go func() {
		defer wg.Done()
		m.Store("name", "daheige")

		// The loaded result is true if the value was loaded, false if stored.
		a, ok := m.LoadOrStore("my-topic", "test")
		log.Println("a: ", a, "ok: ", ok)
	}()
	wg.Wait()
	log.Println("ok")
}

/*
=== RUN   TestMap
2022/03/02 09:36:48 a:  test ok:  false
2022/03/02 09:36:48 abc:  hello ok:  true
2022/03/02 09:36:48 ok
--- PASS: TestMap (0.00s)
PASS
*/

func newValue(i int64) int64 {
	return i + 1
}

func TestPool(t *testing.T) {
	// A Pool must not be copied after first use.
	pool := &sync.Pool{}
	pool.Put(newValue(1))
	pool.Put(newValue(2))
	val := pool.Get().(int64)
	log.Println("value: ", val)

	val = pool.Get().(int64)
	log.Println("value: ", val)
}

/*
2022/03/02 09:39:54 value:  2
2022/03/02 09:39:54 value:  3
*/
type conn struct {
	name string
}

func newConnection(name string) *conn {
	return &conn{
		name: name,
	}
}

func TestPoolConnection(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			return newConnection("redis")
		},
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			conn := pool.Get().(*conn)
			defer pool.Put(conn) // put conn to pool

			log.Println("name: ", conn.name)
		}()
	}

	wg.Wait()
	log.Println("ok")
}

func writeFile(pool *sync.Pool, filename string) error {
	buf := pool.Get().(*bytes.Buffer)

	defer pool.Put(buf)

	// Reset buf
	buf.Reset()

	buf.WriteString("foo")

	return ioutil.WriteFile(filename, buf.Bytes(), 0644)
}

func TestPoolWriteFile(t *testing.T) {
	pool := &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	writeFile(pool, "test.log")
}

func TestOnce(t *testing.T) {
	var once sync.Once
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(func() {
				log.Println("hello")
			})
		}()
	}
	wg.Wait()
	log.Println("ok")
}

/*
=== RUN   TestOnce
2022/03/02 10:13:21 hello
2022/03/02 10:13:21 ok
--- PASS: TestOnce (0.00s)
PASS
*/
