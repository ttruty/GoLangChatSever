package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"

	"github.com/ttruty/Homework/week-7/MakeCerts"
)

// Resource: https://gist.github.com/drewolson/3950226 - go chat app
type ConnectedClient struct {
	conn     net.Conn
	incoming chan string
	outgoing chan string
}

func main() {

	// get our ca and server certificate
	serverTLSConf, _, err := MakeCerts.Certsetup()
	if err != nil {
		log.Fatalf("server: load cert error %s", err)
	}

	//start tcp server secured with TSL listening on localhost:8080
	serverTLSConf.Rand = rand.Reader
	service := "0.0.0.0:8000"
	listener, err := tls.Listen("tcp", service, serverTLSConf)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}

	defer listener.Close()

	//struct to hold connected clients
	//connectedClient := make(chan ConnectedClient)

	// run forever, keep listening for connections
	for {
		//accept an incoming connection and create a handle to the connection (conn)
		//TODO: need to save conns in the final project, so can write to
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// concurrently handle connection
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)
	scanner := bufio.NewScanner(conn)

	for {
		ok := scanner.Scan()

		if !ok {
			break
		}

		handleMessage(scanner.Text(), conn)
	}

	fmt.Println("Client at " + remoteAddr + " disconnected.")
}

func handleMessage(message string, conn net.Conn) {
	fmt.Println(conn.RemoteAddr().String() + "> " + message)
}
