package main

import (
	"net/http"
	"os"
	"log"
	"encoding/json"
	"fmt"
)

const ImdbUrl = "https://omdbapi.com?plot=short&r=json&t="

type Movie struct {
	Title, Poster string
	//Rest interface{}
	//Year, Rated, Released, Runtime, Genre, Director,
	//Writer, Actors, Plot, Language, Country,
	//Awards, Metascore string
	//ImdbRating                string `json:"imdbRating"`
	//ImdbVotes                 string `json:"imdbVotes"`
	//ImdbID                    string `json:"imdbID"`
	//Type                      string
	//Response                  string
}

func main() {
	var title string
	title = os.Args[len(os.Args) - 1]
	if title == "" {
		log.Fatalln("You need a movie name")
	}
	resp, err := http.Get(ImdbUrl + title)
	if err != nil {
		log.Fatalln(err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		log.Fatalln(err)
	}
	var result Movie
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		log.Fatalln(err)
	}
	resp.Body.Close()
	fmt.Printf("%s\t%s\n", result.Title, result.Poster)

	fmt.Printf("%s\n", result)
	file, err := os.Create("Poster_URL.html")
	if err != nil {
		log.Panicln(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	var htmlRes = 	`<html>
				<img src="` + result.Poster + `"/>
			</html>`
	b1 := []byte(htmlRes)
	if _, err := file.Write(b1); err != nil {
		log.Fatal(err)
	}
}