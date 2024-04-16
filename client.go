package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = ":3335"
	CONN_TYPE        = "tcp"
	MSG_DISCONNECT   = "Disconnected from the server.\n"
	MSG_HISTORY      = "/history"
	MSG_USER_COUNT   = "/users"
	HISTORY_RESPONSE = "Chat history:\n"
	HISTORY_FILENAME = "chat_history.txt"
)

func Read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf(MSG_DISCONNECT)
			return
		}
		fmt.Print(str)
	}
}

func Write(conn net.Conn, name string) {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if strings.TrimSpace(str) == MSG_HISTORY {
			_, err = writer.WriteString(MSG_HISTORY + "\n")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = writer.Flush()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			continue
		}

		if strings.TrimSpace(str) == MSG_USER_COUNT {
			_, err = writer.WriteString(MSG_USER_COUNT + "\n")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = writer.Flush()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			continue
		}

		_, err = writer.WriteString(name + ": " + str)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func main() {
	fmt.Print("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	conn, err := net.Dial(CONN_TYPE, CONN_HOST+CONN_PORT)
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to the server.")
	go Read(conn)
	Write(conn, name)
}
