package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	counts := make(map[string]int)
	files := os.Args[1:]
	args := os.Args[1:]
	fmt.Println(args)
	if len(files) == 0 {
		data, _ := fmt.Scanln()
		file, _ := os.OpenFile("test.txt", 1, 0644)
		io.Copy(file, strings.NewReader(string(data)))
		countLines(file, counts)
		for k, v := range counts {
			fmt.Printf("key: %v value: %v", k, v)
		}
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
		for line, n := range counts {
			if n > 1 {
				fmt.Printf("%d\tfile %s\n", n, line)
			}
		}
	}
}

func countLines(f *os.File, counts map[string]int) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		counts[f.Name()+"-"+input.Text()]++
	}
}
