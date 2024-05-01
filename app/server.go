package main

import (
	"flag"
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
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection : ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}
func handleRequest(conn net.Conn) {
	var dirPath string
	flag.StringVar(&dirPath, "directory", ".", "Directory path")
	flag.Parse()

	requestData := make([]byte, 1024)
	_, err := conn.Read(requestData)
	if err != nil {
		fmt.Println("Error while reading from the connection:", err.Error())
	}
	requestString := string(requestData)
	httpMethod := strings.Split(requestString, "\r\n")[0]
	pathString := strings.Split(httpMethod, " ")[1]
	if pathString == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(pathString, "/files/") && len(dirPath) > 0 {
		contents, err := os.ReadFile(dirPath)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}
		content := string(contents)
		conn.Write([]byte("HTTP/1.1 200 OK\r\n" + "Content-Type: application/octet-stream\r\n" + "Content-Length:" + strconv.Itoa(len(content)) + "\r\n\r\n" + content))
	} else if strings.HasPrefix(pathString, "/echo/") {
		content := strings.TrimSpace(pathString[6:])
		conn.Write([]byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\n" + "Content-Length:" + strconv.Itoa(len(content)) + "\r\n\r\n" + content))
	} else if strings.HasPrefix(pathString, "/user-agent") {
		userAgent := strings.Split(requestString, "\r\n")[2]
		conn.Write([]byte("HTTP/1.1 200 OK\r\n" + "Content-Type: text/plain\r\n" + "Content-Length:" + strconv.Itoa(len(userAgent[12:])) + "\r\n\r\n" + userAgent[12:]))
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
	return
}
