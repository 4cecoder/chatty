package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("Starting Chat Server...")
	listener, err := net.Listen("tcp", ":7777")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Waiting for connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("New connection:", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	username := conn.RemoteAddr().String()
	defer conn.Close()
	defer fmt.Println("Connection closed:", username)

	scanner := bufio.NewScanner(conn)

	conn.Write([]byte("Enter your name: "))

	scanner.Scan()
	username = scanner.Text()

	// add the new connection to the connections map
	connections[username] = conn

	for {
		scanner.Scan()
		msg := scanner.Text()

		if msg == "" {
			// ignore empty messages
			continue
		}

		if msg == "/quit" {
			// remove the connection from the connections map
			delete(connections, username)
			return
		}

		broadcast(username, msg)
	}
}

func broadcast(sender, msg string) {
	for _, conn := range connections {
		if conn != nil && conn.RemoteAddr().String() != sender {
			conn.Write([]byte(fmt.Sprintf("[%s]: %s\n", sender, msg)))
		}
	}
}

var connections = make(map[string]net.Conn)
