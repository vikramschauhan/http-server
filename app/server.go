package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection : ", err.Error())
		os.Exit(1)
	}
	requestData := make([]byte, 1024)
	_, err = conn.Read(requestData)
	if err != nil {
		fmt.Println("Error while reading from the connection:", err.Error())
	}
	requestString := string(requestData)
	httpMethod := strings.Split(requestString, "\r\n")[0]
	pathString := strings.Split(httpMethod, " ")[1]
	if strings.HasPrefix(pathString, "/user-agent") {
		content := strings.TrimSpace(pathString[0:])
		userAgent := strings.Split(requestString, "\r\n")[2]
		conn.Write([]byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\n" + "Content-Length:" + strconv.Itoa(len(content)) + "\r\n\r\n" + userAgent[12:]))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
