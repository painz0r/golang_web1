package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

//!+broadcaster
type client chan string // an outgoing message channel
type clientAndName struct {
	client
	name string
}

var (
	entering = make(chan clientAndName)
	leaving  = make(chan clientAndName)
	messages = make(chan string, 5) // all incoming client messages
	//enteringName = make(chan string)
	//leavingName  = make(chan string)
)

func broadcaster(clientNames map[string]bool) {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				select {
				case cli <- msg:
				default:
					fmt.Println("The message was dropped")
				}
			}

		case cli := <-entering:
			clients[cli.client] = true
			clientNames[cli.name] = true
		//case name := <-enteringName:
		//	clientNames[name] = true
		case cli := <-leaving:
			delete(clients, cli.client)
			delete(clientNames, cli.name)
			close(cli.client)
			//case name := <-leavingName:
			//	delete(clientNames, name)

		}
	}
}

func idleCheck(conn net.Conn, out *chan bool, timeout *<-chan time.Time) {
	for {
		select {
		case <-*out:
			fmt.Println("hello from inside")
			*timeout = time.After(10 * time.Minute)
		case <-*timeout:
			fmt.Println("hello from timeout")
			fmt.Fprintln(conn, "\t", strings.ToUpper("GoodBye, you've been too idle"))
			conn.Close()
			return
		default:
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn, clientNames map[string]bool, out *chan bool, timeout *<-chan time.Time) {
	var name string
	ch := make(chan string) // outgoing client messages

	go func() {
		*out <- true
	}()
	go idleCheck(conn, out, timeout)
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
	entering <- clientAndName{ch, name}
	//enteringName <- name
	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- name + "_(" + who + ")" + ": " + input.Text()
		*out <- true
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- clientAndName{ch, name}
	messages <- name + " has left"
	//leavingName <- name
	conn.Close()
}

func clientWriter(conn net.Conn, ch chan string) {
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
		go func() {
			out := make(chan bool)
			timeout := make(<-chan time.Time)
			handleConn(conn, clientNames, &out, &timeout)
		}()
	}
}

//!-main
