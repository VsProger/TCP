package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	CONN_PORT           = ":3335"
	CONN_TYPE           = "tcp"
	MSG_HISTORY_REQUEST = "/history"
	MSG_USER_COUNT      = "/users"
	HISTORY_FILENAME    = "chat_history.txt"
)

var (
	clients     = make(map[net.Conn]bool)
	broadcast   = make(chan string)
	joinChannel = make(chan net.Conn)
)

func main() {
	listener, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Println("Listening on " + CONN_PORT)

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}
		clients[conn] = true
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	joinChannel <- conn
	reader := bufio.NewReader(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Client disconnected.")
			delete(clients, conn)
			return
		}

		str = strings.TrimSpace(str)

		switch str {
		case MSG_HISTORY_REQUEST:
			sendChatHistory(conn)
		case MSG_USER_COUNT:
			sendUserCount(conn)
		default:
			log.Println("Received:", str)
			saveToChatHistory(str)

			time := time.Now().Format(time.ANSIC)
			responseStr := fmt.Sprintf("[%v] %v", time, str)
			broadcast <- responseStr
		}
	}
}

func broadcaster() {
	for {
		select {
		case msg := <-broadcast:
			for client := range clients {
				_, err := client.Write([]byte(msg + "\n"))
				if err != nil {
					log.Println("Error broadcasting message:", err)
					delete(clients, client)
					client.Close()
				}
			}
		case <-joinChannel:
			log.Println("New client joined.")
		}
	}
}

func sendChatHistory(conn net.Conn) {
	file, err := os.Open(HISTORY_FILENAME)
	if err != nil {
		log.Println("Error opening history file:", err)
		conn.Write([]byte("No chat history available.\n"))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		conn.Write([]byte(line + "\n"))
	}
}

func sendUserCount(conn net.Conn) {
	count := len(clients)
	conn.Write([]byte(fmt.Sprintf("Number of connected users: %d\n", count)))
}

func saveToChatHistory(message string) {
	file, err := os.OpenFile(HISTORY_FILENAME, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error opening history file:", err)
		return
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s\n", message)
	if err != nil {
		log.Println("Error writing to history file:", err)
		return
	}
}
