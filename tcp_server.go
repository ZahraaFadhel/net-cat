package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

const maxConnections = 10

var (
	connections       = make(map[net.Conn]string)
	connectionsMutex  sync.Mutex
	newConnections    = make(chan net.Conn)
	closedConnections = make(chan net.Conn)
	messages          = make(chan string)
	historyMessages   []string
	historyMutex      sync.Mutex
	names             []string
)

func main() {
	port := "8989"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	listener, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on the port", port)
	go connectionManager()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		newConnections <- conn
	}
}

func connectionManager() {
	activeConnections := 0

	for {
		select {
		case conn := <-newConnections:
			connectionsMutex.Lock()
			if activeConnections < maxConnections {
				activeConnections++
				go handleConnection(conn)
			} else {
				conn.Write([]byte("Server is full. Try again later.\n"))
				conn.Close()
			}
			connectionsMutex.Unlock()

		case conn := <-closedConnections:
			connectionsMutex.Lock()
			// if name, ok := connections[conn]; ok && name != "" {
			// 	broadcastMessage(fmt.Sprintf("%s has left the chat", name))
			// }
			delete(connections, conn)
			activeConnections--
			connectionsMutex.Unlock()
			conn.Close()

		case msg := <-messages:
			historyMutex.Lock()
			historyMessages = append(historyMessages, msg) // Add message to history
			historyMutex.Unlock()
			broadcastMessage(msg)
		}
	}
}
func handleConnection(conn net.Conn) {
	defer func() {
		closedConnections <- conn
	}()

	// Ask for client name
	printLogo(conn)
	conn.Write([]byte("Enter your name: "))
	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	name := scanner.Text()
	name = strings.TrimSpace(name)
	// fmt.Println("name", name, ".")
	for name == "" {
		conn.Write([]byte("Try Again\nEnter your name: "))
		scanner := bufio.NewScanner(conn)
		if !scanner.Scan() {
			return
		}
		name = scanner.Text()
		name = strings.TrimSpace(name)
		continue
	}

	for IsIn(name, names) {
		conn.Write([]byte("Already Taken\nEnter your name: "))
		scanner := bufio.NewScanner(conn)
		if !scanner.Scan() {
			return
		}
		name = scanner.Text()
	}
	names = append(names, name)
	connectionsMutex.Lock()
	connections[conn] = name
	connectionsMutex.Unlock()
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}
	conn.Write([]byte("Welcome to the chat, " + name + "!\n"))
	displayHistory(conn)
	broadcastMessage(fmt.Sprintf("%s has joined the chat", name))

	// Handle messages
	for scanner.Scan() {
		msg := scanner.Text()
		msg = strings.TrimSpace(msg)
		if msg == "" {
			continue
		}
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fullMsg := fmt.Sprintf("[%s][%s]: %s", timestamp, name, msg)
		messages <- fullMsg
	}
	//conn = <-closedConnections
	broadcastMessage(fmt.Sprintf("%s has left the chat", name))
	names = DeleteName(name, names)
	connectionManager()
	//closedConnections <- conn
	//connectionManager()
}
func printLogo(conn net.Conn) {
	filePath := "logo.txt"
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		conn.Write([]byte(scanner.Text() + "\n"))
		//fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file: ", err)
	}
}

func broadcastMessage(msg string) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	for conn := range connections {
		conn.Write([]byte("\033[2K\r" + msg + "\n> "))
	}
}

func displayHistory(c net.Conn) {
	historyMutex.Lock()
	for _, msg := range historyMessages {
		c.Write([]byte(msg + "\n"))
	}
	historyMutex.Unlock()
}
func IsIn(name string, names []string) bool {
	for _, x := range names {
		if x == name {
			return true
		}
	}
	return false
}
func DeleteName(name string, names []string) []string {
	if IsIn(name, names) {
		for i := 0; i < len(names); i++ {
			if names[i] == name {
				copy(names[i:], names[i+1:])
				names = names[:len(names)-1]
				break
			}
		}
	}
	return names
}
