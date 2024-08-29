package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)
var lastmsg = Message{Type: 0}

type Message struct {
	Type    int
	Content string
}

func main() {
	port := 60829
	if !isPortAvailable(port) {
		panic(fmt.Sprintf("端口 %d 已被占用", port))
	}

	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/send", sendPage)

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

func sendPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// 获取表单参数
	message := r.FormValue("msg")
	lastmsg.Content = message
	lastmsg.Type = 1
	broadcast <- lastmsg

	// 处理并输出响应
	fmt.Fprintf(w, "Received msg: %s\n", message)
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	if lastmsg.Type != 0 {
		conn.WriteMessage(lastmsg.Type, []byte(lastmsg.Content))
	}
	clients[conn] = true

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}
		lastmsg.Type = mt
		lastmsg.Content = string(message)
		log.Printf("recv: %+v", lastmsg)
		broadcast <- lastmsg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteMessage(msg.Type, []byte(msg.Content))
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
