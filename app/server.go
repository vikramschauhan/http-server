package main

import (
	"flag"
	"fmt"
	"github.com/codecrafters-io/http-server-starter-go/constants"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

var dirPath string

func main() {
	flag.StringVar(&dirPath, "directory", ".", "Directory path")
	flag.Parse()
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
	defer conn.Close()
	requestData := make([]byte, 1024)
	_, err := conn.Read(requestData)

	if err != nil {
		fmt.Println("Error while reading from the connection:", err.Error())
	}
	requestString := string(requestData)
	if strings.HasPrefix(requestString, "POST") {
		handlePostRequest(conn, requestString)
	} else {
		handleGetRequest(conn, requestString)
	}
}

func handlePostRequest(conn net.Conn, requestString string) {
	pathString := getPathString(requestString)
	requestBody := strings.Split(requestString, "\r\n")[6]
	if strings.HasPrefix(pathString, "/files/") && len(dirPath) > 0 {
		err := writeFileContents(pathString, requestBody)
		if err != nil {
			writeResponse(conn, constants.NotFoundResponse, "", "")
		}
		writeResponse(conn, constants.CreatedResponse, "", "")
	} else {
		writeResponse(conn, constants.NotFoundResponse, "", "")
	}
}

func handleGetRequest(conn net.Conn, requestString string) {
	pathString := getPathString(requestString)
	if pathString == "/" {
		conn.Write([]byte(constants.OKResponse + "\r\n\r\n"))
	} else if strings.HasPrefix(pathString, "/files/") && len(dirPath) > 0 {
		content, err := getFileContents(pathString)
		if err != nil {
			writeResponse(conn, constants.NotFoundResponse, "", "")
		}
		writeResponse(conn, constants.OKResponse, constants.OctetStream, content)
	} else if strings.HasPrefix(pathString, "/echo/") {
		writeResponse(conn, constants.OKResponse, constants.TextPlain, strings.TrimSpace(pathString[6:]))
	} else if strings.HasPrefix(pathString, "/user-agent") {
		userAgent := strings.Split(requestString, "\r\n")[2]
		writeResponse(conn, constants.OKResponse, constants.TextPlain, userAgent[12:])
	} else {
		writeResponse(conn, constants.NotFoundResponse, "", "")
	}
}

func getPathString(requestString string) string {
	httpMethod := strings.Split(requestString, "\r\n")[0]
	return strings.Split(httpMethod, " ")[1]
}

func getFileContents(pathString string) (string, error) {
	contents, err := ioutil.ReadFile(getFilePath(pathString))
	return string(contents), err
}

func writeFileContents(pathString string, content string) error {
	return ioutil.WriteFile(getFilePath(pathString), []byte(strings.TrimRight(content, "\x00")), 0644)
}

func getFilePath(pathString string) string {
	fileName := strings.TrimPrefix(pathString, "/files/")
	return dirPath + string(os.PathSeparator) + fileName
}
func writeResponse(conn net.Conn, response string, contentType string, content string) {
	if response == constants.OKResponse {
		conn.Write([]byte(response + "\r\n" + constants.ContentTypeKey + contentType + "\r\n" + constants.ContentLengthKey + strconv.Itoa(len(content)) + "\r\n\r\n" + content))
	} else {
		conn.Write([]byte(response + "\r\n\r\n"))
	}
}
