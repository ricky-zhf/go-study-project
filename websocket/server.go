package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// 定义一个升级器，用于将 HTTP 连接升级为 WebSocket 连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许跨域请求，在实际生产环境中应根据具体情况设置
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocket 处理函数
func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 尝试将 HTTP 连接升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket 升级失败:", err)
		return
	}
	// 确保在函数结束时关闭连接
	defer conn.Close()

	for {
		// 读取客户端发送的消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息出错:", err)
			break
		}
		// 将消息原样返回给客户端
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			log.Println("写入消息出错:", err)
			break
		}
	}
}

func runServer() {
	// 注册 WebSocket 处理函数
	http.HandleFunc("/ws", wsHandler)
	// 启动 HTTP 服务器，监听 8080 端口
	log.Println("服务器启动，监听端口 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
