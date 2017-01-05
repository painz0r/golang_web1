package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
	"unicode/utf8"
	"log"
	"sort"
)

func wordFreq() (wordFrq map[string]int) {
	wordFrq = make(map[string]int)
	file, err := os.Open("main.go")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		line := scanner.Text()
		wordFrq[line]++
	}
	return wordFrq
}

func main() {
	words := wordFreq()
	var nums = make([]int, 0, len(words))
	for _, count := range words {
		nums = append(nums, count)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(nums)))
	fmt.Println(nums)
	for _, v := range nums {
		for word, count := range words {
			if count > 5 {
				if v == words[word] {
					fmt.Printf("%s\t%d\n", word, count)
				}
			}

		}

	}

	counts := make(map[rune]int)    // counts of Unicode characters
	var utflen [utf8.UTFMax + 1]int // count of lengths of UTF-8 encodings
	invalid := 0                    // count of invalid UTF-8 characters
	letters := make(map[rune]int)
	digits := make(map[rune]int)
	file, err := os.Open("main.go")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	in := bufio.NewReader(file)
	for {
		r, n, err := in.ReadRune() // returns rune, nbytes, error
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}
		if r == unicode.ReplacementChar && n == 1 {
			invalid++
			continue
		}
		counts[r]++
		utflen[n]++
		if unicode.IsDigit(r) {
			digits[r]++
		}
		if unicode.IsLetter(r) {
			letters[r]++
		}
	}
	fmt.Printf("rune\tcount\n")
	for c, n := range counts {
		fmt.Printf("%q\t%d\n", c, n)
	}
	fmt.Print("\nlen\tcount\n")
	for i, n := range utflen {
		if i > 0 {
			fmt.Printf("%d\t%d\n", i, n)
		}
	}
	if invalid > 0 {
		fmt.Printf("\n%d invalid UTF-8 characters\n", invalid)
	}
	fmt.Print("\nnumber of letters\n")
	for i, n := range letters {
		if n > 0 {
			fmt.Printf("%q\t%d\n", i, n)
		}
	}
	fmt.Print("\nnumber of digits\n")
	for i, n := range digits {
		if n > 0 {
			fmt.Printf("%q\t%d\n", i, n)
		}
	}
}