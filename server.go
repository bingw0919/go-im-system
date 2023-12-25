package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// NewServer 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}
func (s *Server) Handler(conn net.Conn) {
	//当前链接的业务
	fmt.Println("Conn Succeed。。。")
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
	for {
		//accept
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// do handler
		go s.Handler(conn)
	}
}
