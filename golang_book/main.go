package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"time"
)

func Echo(ws *websocket.Conn) {
	var err error

	//ws.SetReadDeadline(time.Now().Add(10 * time.Second))

	//fmt.Println(timer)
	times_up := make(chan bool)

	go func() {
		time.NewTimer(time.Second * 2)
		times_up <- true
	}()

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
		if <-times_up {
			fmt.Println("Can't send")
			break
		}

	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
