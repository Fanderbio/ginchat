package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FromId   int64  // 消息的发送方
	TargetId int64  // 消息接收者
	Type     int    // 发送类型 群聊 私聊 广播
	Media    int    // 消息类型 文字 图片 音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string // 描述相关的
	Amount   int    // 其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1. TODO: 校验token等合法性
	// token := query.Get("token")

	query := request.URL.Query()
	userId := query.Get("userId")
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")

	isvalid := true // checkToken()

	conn, err := (&websocket.Upgrader{
		// token校验
		CheckOrigin: func(r *http.Request) bool {
			return isvalid
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. 获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// 3.用户关系

	// 4.userid和node绑定，并加锁
	userIdInt, _ := strconv.ParseInt(userId, 10, 64)
	rwLocker.Lock()
	clientMap[int64(userIdInt)] = node
	rwLocker.Unlock()

	// 5.完成发送逻辑
	go sendProcess(node)

	// 6.完成接收的逻辑
	go receiveProcess(node)
	sendMsg(userIdInt, []byte("欢迎进入聊天系统"))
}

func sendProcess(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func receiveProcess(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		boardMsg(data)
		fmt.Println("[ws] <<<<< ", data)
	}
}

var udpsendChan chan []byte = make(chan []byte, 1024)

func boardMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProcess()
	go udpRecvProcess()
}

// 完成upd数据发送协程
func udpSendProcess() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(172, 23, 191, 255),
		Port: 3000,
	})

	defer conn.Close()

	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpsendChan:
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成upd数据接收协程
func udpRecvProcess() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		var buffer [512]byte
		n, err := conn.Read(buffer[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buffer[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: // 私信
		sendMsg(msg.TargetId, data)
		// case 2: // 群发
		// 	sendGroupMsg()
		// case 3: // 广播
		// 	sendAllMsg()
		// case 4:
		//
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()

	if ok {
		node.DataQueue <- msg
	}
}
