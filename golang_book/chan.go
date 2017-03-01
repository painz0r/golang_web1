package main

import (
	"time"
	"log"
	//"fmt"
	//"sort"
	//"sort"
	"fmt"
	"math/rand"
)

func main() {
	//start := time.Now()
	//h := []int{}
	//y := []int{}
	//doneMain := make(chan struct{})
	//defer close(doneMain)
	//ch := numbers(doneMain)
	//worker1 := sq(doneMain, ch)
	//worker2 := sq(doneMain, ch)
	//worker3 := sq(doneMain, ch)
	//worker4 := sq(doneMain, ch)
	//worker5 := sq(doneMain, ch)
	////worker6 := sq(doneMain, ch)
	////worker7 := sq(doneMain, ch)
	////worker8 := sq(doneMain, ch)
	////worker9 := sq(doneMain, ch)
	////worker10 := sq(doneMain, ch)
	//for n := range merge(doneMain, worker1, worker2, worker3, worker4, worker5) {
	//	//worker6, worker7, worker8, worker9, worker10) {
	//	h = append(h, n)
	//}
	//
	//elapsed := time.Since(start)
	//log.Printf("%s took %s", "First", elapsed)
	//
	//start = time.Now()
	//
	//y = oneThread(y)
	//
	//elapsed = time.Since(start)
	//log.Printf("%s took %s", "Second", elapsed)
	//sort.Ints(h)
	//sort.Ints(y)
	////fmt.Println(h)
	////fmt.Println(y)
	//if testEq(h, y) {
	//	fmt.Println("Equal")
	//} else {
	//	fmt.Println("not equal")
	//}
	//data := []string{"query1", "query2", "query3", "query4"}
	run()
	run2()
}

func numbers(doneMain chan struct{}) <-chan int {
	out := make(chan int, 5)
	go func() {
		defer close(out)
		for i := 0; i <= 20000; i++ {
			select {
			case out <- i:
			case <-doneMain:
				return
			}
		}
	}()
	return out
}

func sq(doneMain chan struct{}, c <-chan int) <-chan int {
	out := make(chan int, 5)
	go func() {
		defer close(out)
		for i := range c {
			select {
			case out <- factorialChan(i):
			case <-doneMain:
				return
			}
		}

	}()

	return out
}

func merge(doneMain <-chan struct{}, cs ...<-chan int) <-chan int {
	out := make(chan int, 1)
	done := make(chan bool)

	for _, c := range cs {
		go func(c <-chan int) {
			for n := range c {
				select {
				case out <- n:
				case <-doneMain:
					return
				}
			}
			done <- true
		}(c)
	}
	go func() {
		for i := 1; i <= len(cs); i++ {
			<-done
		}
		close(out)
	}()

	return out

}

func oneThread(y []int) []int {
	//var res []int
	for i := 0; i <= 20000; i++ {
		//res = append(res, i)
		y = append(y, factorialChan(i))
	}
	//for _, n := range res {
	//	y = append(y, factorialChan(n))
	//}
	return y
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func testEq(a, b []int) bool {

	if a == nil && b == nil {
		return true;
	}

	if a == nil || b == nil {
		return false;
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func factorialChan(n int) int {
	if n == 0 {
		return 1
	}
	return n * factorialChan(n - 1)
}

var (
	Web1 = fakeSearch("web1")
	Image1 = fakeSearch("image1")
	Video1 = fakeSearch("video1")
	Web2 = fakeSearch("web2")
	Image2 = fakeSearch("image2")
	Video2 = fakeSearch("video2")
)

type Result string
type Search func(query string) Result

func fakeSearch(kind ...string) Search {
		return func(query string) Result {
			time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			return Result(fmt.Sprintf("%s result for %q\n", kind, query))
		}
}

func run() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang")
	//result := First("golang",
	//	fakeSearch("replica 1"),
	//	fakeSearch("replica 2"))
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
// Google search 1.0 - no concurrency
//func Google(query string) (results []Result) {
//	results = append(results, Web(query))
//	results = append(results, Image(query))
//	results = append(results, Video(query))
//	return
//}

// Google search 2.0
//func Google(query string) (results []Result) {
//	c := make(chan Result)
//	go func() {
//		c <- Web(query)
//	} ()
//	go func() {
//		c <- Image(query)
//	} ()
//	go func() {
//		c <- Video(query)
//	} ()
//
//	for i := 0; i < 3; i++ {
//		result := <-c
//		results = append(results, result)
//	}
//	return
//}
//Google search 2.1
//func Google(query string) (results []Result) {
//	c := make(chan Result)
//	go func() {
//		c <- Web(query)
//	}()
//	go func() {
//		c <- Image(query)
//	}()
//	go func() {
//		c <- Video(query)
//	}()
//
//	timeout := time.After(80 * time.Millisecond)
//	for i := 0; i < 3; i++ {
//		select {
//		case result := <-c:
//			results = append(results, result)
//		case <-timeout:
//			fmt.Println("timed out")
//			return
//		}
//	}
//	return
//}

//Replicate the servers to avoid discarding of slow results

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	searchReplica := func(i int) {
		c <- replicas[i](query)
	}
	for i := range replicas {
		go searchReplica(i)
	}
	return <-c
}

// Google search 3.0
func Google(query string) (results []Result) {
	c := make(chan Result,3)

	go func() {
		c <- First(query, Web1, Web2)
	}()
	go func() {
		c <- First(query, Image1, Image2)
	}()
	go func() {
		c <- First(query, Video1, Video2)
	}()
	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case result := <-c:
			results = append(results, result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

type Conn struct{
	name string
	waitTime time.Duration
}
var c1 = Conn{name:"req1", waitTime:time.Second}
var c2 = Conn{name:"req2", waitTime:time.Millisecond}

var conns = []Conn{c1, c2}

func Query(conns []Conn, query string) string {
	ch := make(chan string, len(conns))  // buffered
	for _, conn := range conns {
		go func(c Conn) {
			ch <- c.DoQuery(query, c.name, c.waitTime)
		}(conn)
	}
	return <-ch
}

func (c Conn) DoQuery(query string, connName string, waitTime time.Duration) string{
	time.Sleep(waitTime)
	return fmt.Sprintln("Returning search from " + query + " " + connName)
}

func run2() {
	request := Query(conns, "Test string")
	fmt.Println(request)
}