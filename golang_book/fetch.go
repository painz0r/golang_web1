package main

import (
	"fmt"
	"net/http"
	"os"
	//"io/ioutil"
	"io"
	"log"
	"strings"
)

func main() {
	for _, url := range os.Args[1:] {
		if ok := strings.HasPrefix(url, "http://"); !ok {
			url = "http://" + url
		}
		log.Println(url)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		//b, err := ioutil.ReadAll(resp.Body)
		b, err := io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}
		fmt.Println("\n***********")
		fmt.Println(resp.Status)
		fmt.Println("***********")
		resp.Body.Close()
		fmt.Printf("type: %T value: %v", b, b)
	}

}
