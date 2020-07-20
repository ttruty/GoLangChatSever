package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	s "strings"
)

func main() {

	// get our ca and server certificate
	caPEM, err := ioutil.ReadFile("client.pem")
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM(caPEM)
	clientTLSConf := &tls.Config{
		RootCAs: certpool,
	}

	//dial tcp running on on localhost:8080
	conn, err := tls.Dial("tcp", "127.0.0.1:8000", clientTLSConf)
	if err != nil {
		fmt.Println("failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// write string from console
	fmt.Print("Type STOP to leave chat: \n")
	for {
		go ReadConn(conn)

		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if len(text) == 6 && (s.ToLower(text[0:4]) == "stop") { // len include '/n', probably a better way to handle this but it works
			fmt.Print("--leaving chat--")
			os.Exit(0)
		}
		//write text to the server
		_, err = conn.Write([]byte(text))

	}

}

func ReadConn(conn net.Conn) {
	for {
		//read response
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("failed to read from server:", err)
			os.Exit(1)
		}

		fmt.Println(string(buf[:n]))
	}
}
