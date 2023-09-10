package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"testing"
)

func TestAConnection(t *testing.T) {
	message := "Something written!\n"
	server, client := net.Pipe()

	go func() {
		defer server.Close()
		server.Write([]byte(message))
	}()

	defer client.Close()
	b, err := io.ReadAll(client)
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Println(string(b))

}

func TestSendMessage(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:0")
	message := "Is this delivered!"
	if err != nil {
		log.Fatal()
	}
	sendMessage(conn, &message)
	handleConnection(conn)
	fmt.Printf("Ved ikke hvad der skal testes")
}
