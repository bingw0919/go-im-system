package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	Message   chan string
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// BroadCast 广播消息
func (s *Server) BroadCast(user *User, message string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + message
	s.Message <- sendMsg
}

// ListenMessage 监听Message
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		fmt.Println("server transmit msg:" + msg)
		s.mapLock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.mapLock.Unlock()
	}
}

func (s *Server) Handler(conn net.Conn) {
	//当前链接的业务
	fmt.Println("Conn Succeed。。。")
	//用户上线，将用户加入到onlineMap
	user := NewUser(conn)
	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()
	s.BroadCast(user, "online...")

	//当前handler阻塞
	select {}
}

// Start 启动服务器的接口
func (s *Server) Start() {
	//socket start
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Printf("net.Listen error: %v", err)
		return
	}
	//close listener socket
	defer func(listen net.Listener) {
		err := listen.Close()
		if err != nil {
			fmt.Printf("net.Listen close error: %v", err)
		}
	}(listen)
	fmt.Println("Server Started~~~")

	go s.ListenMessage()

	for {
		//accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		// do handler
		go s.Handler(conn)
	}
}
