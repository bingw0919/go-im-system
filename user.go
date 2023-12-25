package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser create user
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
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
