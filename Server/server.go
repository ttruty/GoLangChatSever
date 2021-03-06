package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/ttruty/Homework/week-7/MakeCerts"
)

// Resource: https://gist.github.com/drewolson/3950226 - go chat app
// Resource: https://golangforall.com/en/post/golang-tcp-server-chat.html
// Resource: https://stackoverflow.com/questions/36417199/how-to-broadcast-message-using-channel

type ClientJob struct {
	name    string
	message string
	conn    net.Conn
}

func generateResponses(clientJobs chan ClientJob) {
	for {
		// Wait for the next job to come off the queue.
		clientJob := <-clientJobs

		// Send back the response.
		clientJob.conn.Write([]byte(clientJob.name + ">" + clientJob.message))
	}
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

	// add client to map in struct
	// Using sync.Map to store map off connected clients
	var connMap = &sync.Map{} //https://golang.org/pkg/sync/#Map

	// run forever, keep listening for connections
	for {
		//accept an incoming connection and create a handle to the connection (conn)
		//TODO: need to save conns in the final project, so can write to
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		//unique id for connected client
		id := uuid.New().String()
		connMap.Store(id, conn)

		// concurrently handle connection
		go handleConnection(id, conn, connMap)
	}
}

func handleConnection(id string, c net.Conn, connMap *sync.Map) {
	//defer closing connection and deleting the connection map
	defer func() {
		connMap.Delete(id)
		c.Close()
	}()

	clientJobs := make(chan ClientJob)
	go generateResponses(clientJobs)

	remoteAddr := c.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)
	// scanner for input
	scanner := bufio.NewScanner(c)

	for {
		// check if able to receive in scan
		ok := scanner.Scan()
		if !ok {
			break
		}

		// for each connection in the connection map
		connMap.Range(func(key, value interface{}) bool {
			if conn, ok := value.(net.Conn); ok {
				if c.RemoteAddr() != conn.RemoteAddr() { //only send to other client, not slef
					clientJobs <- ClientJob{remoteAddr, scanner.Text(), conn}
				} else {
					fmt.Println(remoteAddr, ">", scanner.Text())
				}
			}
			return true
		})

	}
	fmt.Println("Client at " + remoteAddr + " disconnected.")
}
