package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}

	file, err := os.Create("testResults.txt")
	check(err)
	defer file.Close()
	for range os.Args[1:] {
		res := <-ch
		fmt.Println(res)
		file.WriteString(res + "\n")
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
	file.WriteString(fmt.Sprintf("%.2fs elapsed\n", time.Since(start).Seconds()))
}

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}

func fetch(url string, ch chan string) {
	start := time.Now()
	if ok := strings.HasPrefix(url, "http://"); !ok {
		url = "http://" + url
	}
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)

}
