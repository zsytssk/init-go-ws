package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var lastMsg = Message{Type: 0}

type Message struct {
	Type    int
	Content []byte
}

func main() {
	port := 60829
	if !isPortAvailable(port) {
		panic(fmt.Sprintf("端口 %d 已被占用", port))
	}

	http.HandleFunc("/ws", handleWs)
	http.HandleFunc("/send", handleSend)

	go handleMessages()

	fmt.Println("Server started on port:", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}

func isPortAvailable(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// 如果有错误，端口可能已经被占用
		return false
	}
	// 端口可用，关闭监听器
	listener.Close()
	return true
}

func handleSend(w http.ResponseWriter, r *http.Request) {
	respBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// 获取表单参数
	lastMsg.Content = respBody
	lastMsg.Type = 1
	broadcast <- lastMsg

	log.Printf("handleSend: %s\n", respBody)
	// 处理并输出响应
	fmt.Fprintf(w, "Received : %s\n", respBody)
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	if lastMsg.Type != 0 {

	}
	clients[conn] = true

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}
		if string(message) == "getInit" {
			conn.WriteMessage(lastMsg.Type, []byte(lastMsg.Content))
		}
		info := make(map[string]string)
		info["type"] = strconv.Itoa(mt)
		info["content"] = string(message)
		log.Printf("handleWs: %+v", info)
		// broadcast <- lastMsg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteMessage(msg.Type, msg.Content)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
