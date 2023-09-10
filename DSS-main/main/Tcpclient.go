package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

//Global variable
var cons []net.Conn
var MessagesSent map[string]bool

//Gets the IP address from client
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	hostip, _, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		log.Fatal(err)
	}

	return hostip
}

//Sends the message given from broadcast, to all peers in the network
func sendMessage(conn net.Conn, message *string) {
	enc := gob.NewEncoder(conn)
	err := enc.Encode(message)
	if err != nil {
		log.Fatal(err)
	}
}

//Broadcasts messages to peers in the network
func broadcast(cons []net.Conn, message string) {
	if !MessagesSent[message] {
		for _, conn := range cons {
			sendMessage(conn, &message)
		}
		MessagesSent[message] = true
	}
}

//Handles the connections between clients.
func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		var msg string
		dec := gob.NewDecoder(conn)
		err := dec.Decode(&msg)

		if err != nil {
			fmt.Println("gob decode: " + err.Error())
			return
		}
		if !MessagesSent[msg] {
			fmt.Print("received string through bufio reader: " + msg)
		}
		broadcast(cons, msg)
	}
}

//Server function, handles accepting connections to peer to peer network
func runServer() {
	ln, err := net.Listen("tcp", ":")
	if err != nil {
		log.Fatal(err)
	}

	ipAddress := GetOutboundIP()
	port := strconv.Itoa(ln.Addr().(*net.TCPAddr).Port)
	ipAndPort := ipAddress + ":" + port

	fmt.Println("Listening for connections on IP:port " + ipAndPort)
	for {
		conn, _ := ln.Accept()
		cons = append(cons, conn)
		fmt.Println("Got a connection...")
		go handleConnection(conn)
	}

}

func main() {
	//Initializes a list of connections, aswell as a map of messages
	cons = make([]net.Conn, 0)
	MessagesSent = make(map[string]bool)
	reader := bufio.NewReader(os.Stdin)

	//Asks the user for the ip and port
	fmt.Print("-> Please supply IP address: ")
	ip, _ := reader.ReadString('\n')
	fmt.Print("-> Please supply port: ")
	port, _ := reader.ReadString('\n')

	//Attempts to connect to a network
	conn, _ := net.Dial("tcp", (strings.TrimSpace(ip) + ":" + strings.TrimSpace(port)))
	if conn != nil {
		cons = append(cons, conn)
		go handleConnection(conn)
	}
	go runServer()

	//Message handling
	for {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("User input loop: " + err.Error())
			return
		}
		go broadcast(cons, text)

	}
}
