package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Starting Chat Client...\n")

	conn, err := net.Dial("tcp", "localhost:7777")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	go receiveMessages(conn)

	scanner := bufio.NewScanner(os.Stdin)

	// get the user's name
	fmt.Println("Enter your name: ")
	scanner.Scan()
	username := scanner.Text()
	// Send the username to the server
	conn.Write([]byte(username + "\n"))

	for {
		scanner.Scan()
		msg := scanner.Text()

		if msg == "/quit" {
			return
		}

		broadcast(conn, username, msg)
	}
}

func receiveMessages(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	for {
		if scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}
}

func broadcast(conn net.Conn, sender, msg string) {
	// check if the message already has the sender's name
	if !strings.HasPrefix(msg, sender+": ") {
		conn.Write([]byte(msg + "\n"))
	} else {
		conn.Write([]byte(sender + msg + "\n"))
	}
}
