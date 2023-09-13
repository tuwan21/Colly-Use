package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"net/http"
	"strconv"
	"sync"
)

// Node 当前用户节点 userId和Node的映射关系
type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	// 群组的消息分发
	GroupSet set.Interface
}

// 映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(ctx *gin.Context) {
	userId := ctx.DefaultQuery("userId", "")
	token := ctx.DefaultQuery("token", "")
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	isValied := checkToken(userIdInt, token)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isValied
		},
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

	if err != nil {
		fmt.Println("WebSocket连接错误")
		return
	}

	// 绑定到当前节点
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSet:  set.New(set.ThreadSafe),
	}
	// 映射关系的绑定
	rwLocker.Lock()
	clientMap[userIdInt] = node
	rwLocker.Unlock()
	// 根据当前用户的id来加入群组
	AddGroupByAccountId(userIdInt)
	// 连接成功就给当前用户发送一个Hello world
	sendMsg(userIdInt, []byte("hello world!"))
	// 发送数据给客户端
	go senProc(node)
	// 接收客户端的数据
	go recvProc(node)
}

// 内部使用判断token合法
func checkToken(userId int64, token string) bool {
	//...
	return true
}

// 将数据推送到管道中,然后让管道中拿数据发送出去
func sendMsg(userId int64, message []byte) {
	rwLocker.RLock()
	node, isOk := clientMap[userId]
	rwLocker.RUnlock()
	if isOk {
		node.DataQueue <- message
	}
}

// 从管道中获取数据发送出去
func senProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println("发送消息失败")
				return
			}
		}
	}
}

// 接收客户端数据
func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println("接收数据失败", err)
			return
		}
		// 将数据处理转发给对应的人
		dispatch(data)
	}
}

// 分发数据
func dispatch(data []byte) {
	type Message struct {
		UserId int64  `json:"userId"`
		Msg    string `json:"msg"`
	}
	fmt.Println("接收到的数据", string(data))
	// 解析出数据
	message := Message{}
	err := json.Unmarshal(data, &message)
	if err != nil {
		fmt.Println("解析数据失败:", err.Error())
		return
	}
	fmt.Println("解析的数据为:", message)
	// 发送数据
	sendMsg(message.UserId, data)
}

// SendMessage 定义一个对外的方法(比如在别的接口中要发送数据到WebSocket当中)
func SendMessage(userId int64, message interface{}) {
	str, _ := json.Marshal(message)
	sendMsg(userId, str)
}

// addGroupByAccountId 加入群聊
func AddGroupByAccountId(userId, groupId int64) {
	fmt.Println(userId, "加入群聊")
	rwLocker.Lock()
	node, isOk := FrontClientMap[userId]
	if isOk {
		node.GroupSet.Add(groupId)
		fmt.Println("加入群聊成功")
	}
	rwLocker.Unlock()
}

func main() {

}
