package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// Online 用户上线功能
func (u *User) Online() {
	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	u.server.BroadCast(u, "user online...")
}

// Offline 用户下线功能
func (u *User) Offline() {
	u.server.mapLock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.mapLock.Unlock()
	u.server.BroadCast(u, "user offline...")
}

func (u *User) SendMessage(msg string) {
	_, err := u.conn.Write([]byte(msg))
	if err != nil {
		fmt.Println("SendMessage Error")
	}
}

// DoMessage 用户处理消息
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前所有用户
		u.server.mapLock.Lock()
		onlineCount := fmt.Sprintf("there is  %d person online ....", len(u.server.OnlineMap))
		u.SendMessage(onlineCount)
		for _, cli := range u.server.OnlineMap {
			sendMsg := "[" + cli.Addr + "]" + cli.Name + "online、、、、、、"
			u.SendMessage(sendMsg)
		}
		u.server.mapLock.Unlock()
	} else {
		u.server.BroadCast(u, msg)
	}
}

// NewUser create user
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

// ListenMessage 监听chan,有消息立刻发送给服务端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		fmt.Println("write msg:" + msg)
		_, err := u.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("ListenMessage write error: ", err)
		}
	}
}
