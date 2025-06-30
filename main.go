package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

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
var handleListeners = make(map[uint32]chan string)

type Message struct {
	Type    int
	ID      uint32
	Content string
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

var id uint32 = 0

func genId() uint32 {
	id = id + 1
	return id
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
	ID := genId()
	lastMsg.Content = string(respBody)
	lastMsg.Type = 1
	lastMsg.ID = ID
	broadcast <- lastMsg

	log.Printf("handleSend: %s\n", respBody)
	// 创建用于接收响应的 channel
	respChan := make(chan string, 1) // 缓冲防止 goroutine 泄漏

	handleListeners[ID] = respChan

	// 等待响应或超时
	select {
	case response := <-respChan: // 同步写入响应
		w.Write([]byte(response))
	case <-time.After(3 * time.Second): // 设置超时时间
		http.Error(w, "Response timeout", http.StatusGatewayTimeout)
	case <-r.Context().Done(): // 客户端提前断开连接
	}
	delete(handleListeners, ID)
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	clients[conn] = true

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}
		if string(message) == "getInit" {
			clientSendMsg(conn, lastMsg)
			continue
		}
		msg := Message{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			continue
		}
		if respChan, ok := handleListeners[msg.ID]; ok {
			respChan <- msg.Content // 发送响应
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			clientSendMsg(client, msg)
		}
	}
}

func clientSendMsg(client *websocket.Conn, msg Message) error {
	if msg.ID == 0 {
		msg.ID = genId()
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return client.WriteMessage(msg.Type, data)
}
