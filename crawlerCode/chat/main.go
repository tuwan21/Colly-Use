package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源连接，实际生产环境需要根据需要进行跨域控制
	},
}

func main() {
	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("Error upgrading connection:", err.Error())
			return
		}

		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err.Error())
				break
			}

			fmt.Printf("Received message: %s\n", msg)

			// 处理接收到的消息（你可以在这里加入自定义逻辑）

			// 发送响应消息给客户端
			err = conn.WriteMessage(websocket.TextMessage, []byte("Server received: "+string(msg)))
			if err != nil {
				fmt.Println("Error writing message:", err.Error())
				break
			}
		}
	})

	r.Run(":8080")
}
