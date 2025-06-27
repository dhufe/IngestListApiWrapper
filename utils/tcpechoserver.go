package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp4", ":21757")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(c)
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	packet := make([]byte, 4096)
	defer c.Close()

	n, _ := c.Read(packet)
	fmt.Printf("Received %d bytes from client.\n", n)
	packet = packet[:n]
	n, _ = c.Write(packet)
	fmt.Printf("Send %d bytes\n", n)
}
