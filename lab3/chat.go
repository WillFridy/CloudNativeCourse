// Demonstration of channels with a chat application
// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// Chat is a server that lets clients chat with each other.

package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

// type client chan<- string // an outgoing message channel
type client struct {
	channel chan<- string
	name    string
}

//Goble var
var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

func main() {
	//Network codeing
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

// Ask about *hint

func broadcaster() {
	clients := make(map[client]bool) // all connected clients

	for {
		select {

		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			// Reads thru all clients
			for cli := range clients {
				cli.channel <- msg
			}

		case cli := <-entering:
			//Throws cli into map
			clients[cli] = true
			//
			cli.channel <- "Currently Online:"
			// Go Through Client list
			for c := range clients {
				// Doesn't Print Self
				if c != cli {
					cli.channel <- c.name
				}
			}
		case cli := <-leaving:
			delete(clients, cli)
			close(cli.channel)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // outgoing client messages
	go clientWriter(conn, ch)

	//Add Line to input user name pass to who cin, cout
	var cli client
	var who string
	cli.channel = ch

	// Entering Name
	ch <- "Enter User Name"
	reader := bufio.NewReader(conn)      // waits for input
	data, _ := reader.ReadString('\n')   // stop reading after \n is read
	who = strings.TrimSuffix(data, "\n") // removes \n from user input

	cli.name = who
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- cli

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
	}
	// NOTE: ignoring potential errors from input.Err()

	leaving <- cli
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

/*
Have to store the name, currently only stores the channal
use map to store storage
*/
