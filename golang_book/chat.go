// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 254.
//!+

// Chat is a server that lets clients chat with each other.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

//!+broadcaster
type client chan string // an outgoing message channel

var (
	entering     = make(chan client)
	leaving      = make(chan client)
	messages     = make(chan string) // all incoming client messages
	enteringName = make(chan string)
	leavingName  = make(chan string)
)

func parseName(s string) string {
	var name string
	name = strings.Trim(s, "You are ")
	return name
}

func broadcaster(clientNames map[string]bool) {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true
		case name := <-enteringName:
			clientNames[name] = true
		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		case name := <-leavingName:
			delete(clientNames, name)

		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn, clientNames map[string]bool) {
	var name string
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	conn.Write([]byte("Hello\n"))
	fmt.Fprintln(conn, "Please enter your name:")
	fmt.Fscanln(conn, &name)
	ch <- "You are " + name

	if len(clientNames) != 0 {
		fmt.Fprintln(conn, "Here is who's in the chat room at the moment:")
		for n := range clientNames {
			fmt.Fprintln(conn, n)
		}
	}
	messages <- name + " has arrived"
	entering <- ch
	enteringName <- name
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + "_(" + who + ")" + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- ch
	messages <- name + " has left"
	leavingName <- name
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	clientNames := make(map[string]bool)
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster(clientNames)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, clientNames)
	}
}

//!-main
