package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"context"
	"os"
)

func main() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)

	req, err := http.ReadRequest(r)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	fmt.Println("Received request:", req.Method, req.URL)

	ip := req.RemoteAddr

	// Get the headers of the request
	headers := req.Header

	// Write the headers and IP to a file
	file, err := os.Create("request_info.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString("IP: " + ip + "\n")
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	for key, value := range headers {
		_, err := file.WriteString(key + ": " + value[0] + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}

	// Modify the request URL to point to the backend server

	cloneReq := req.Clone(context.Background())
	cloneReq.URL = &url.URL{
		Scheme: "http",
		Host: "localhost:5000",
		Path: req.URL.Path,
	}

	var request_url = "http://localhost:5000"+req.URL.Path
	// Send the request to the backend server
	resp, err := http.DefaultClient.Get(request_url, )
	if err != nil {
		fmt.Println("first", err)
		return
	}
	defer resp.Body.Close()

	// Write the response back to the client
	buf := bufio.NewWriter(conn)
	err = resp.Write(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = buf.Flush()
	if err != nil {
		fmt.Println(err)
		return
	}






}