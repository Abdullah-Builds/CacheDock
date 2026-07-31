package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	// Send a command to the server
	fmt.Fprintln(conn, "INFO")

	// Read one line of response
	response, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Println(response)
}
