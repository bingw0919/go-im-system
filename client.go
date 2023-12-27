package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.ServerIp, client.ServerPort))
	if err != nil {
		fmt.Println("new client connect failed")
		return nil
	}
	client.conn = conn
	return client
}

var serverIP string
var serverPort int

func init() {
	flag.StringVar(&serverIP, "ip", "127.0.0.1", "set default server ip")
	flag.IntVar(&serverPort, "port", 8888, "set default server port")
}
func (c *Client) selectUsers() {
	sendMsg := "who\n"
	_, err := c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("selectUsers write error: ", err)
		return
	}
}
func (c *Client) privateChat() {
	c.selectUsers()
	var remoteName string
	fmt.Println("please select users:(exit is supported)")
	fmt.Scanln(&remoteName)
	var sendMessage string
	for sendMessage != "exit" {
		fmt.Println("please enter sendMessage")
		fmt.Scanln(&sendMessage)
		for len(sendMessage) != 0 {
			sendMessage = "@" + remoteName + " " + sendMessage
			_, err := c.conn.Write([]byte(sendMessage + "\n"))
			if err != nil {
				fmt.Println("publicChat write error: ", err)
				break
			}
			sendMessage = ""
			fmt.Print("please enter content(exit is supported): ")
			fmt.Scanln(&sendMessage)
		}
		c.selectUsers()
		fmt.Println("please select users:(exit is supported)")
	}
}

func (c *Client) publicChat() {
	var sendMessage string
	fmt.Print("please enter content(exit is supported): ")
	fmt.Scanln(&sendMessage)
	for sendMessage != "exit" {
		if len(sendMessage) != 0 {
			_, err := c.conn.Write([]byte(sendMessage + "\n"))
			if err != nil {
				fmt.Println("publicChat write error: ", err)
				break
			}
		}
		sendMessage = ""
		fmt.Print("please enter content: (exit is supported)")
		fmt.Scanln(&sendMessage)
	}
}

func (c *Client) updateName() bool {
	fmt.Printf(">>>please send name:")
	_, err := fmt.Scanln(&c.Name)
	if err != nil {
		fmt.Println("read name error: ", err)
	}
	sendMsg := "rename|" + c.Name + "\n"
	_, err = c.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("write name error: ", err)
		return false
	}
	return true
}

func (c *Client) Run() {
	for c.flag != 0 {
		for !c.menu() {
		}
		switch c.flag {
		case 1:
			fmt.Println("public mode")
			c.publicChat()
		case 2:
			fmt.Println("private mode")
			c.privateChat()
		case 3:
			fmt.Println("renamed mode")
			c.updateName()
		case 0:
			fmt.Println("exit")
		}
	}
}
func (c *Client) DealResponse() {
	io.Copy(os.Stdout, c.conn)
}

func (c *Client) menu() bool {
	var f int
	fmt.Println("1.public mode")
	fmt.Println("2.private mode")
	fmt.Println("3.update name")
	fmt.Println("0.exit")
	_, err := fmt.Scanln(&f)
	if err != nil {
		fmt.Println("read menu flag error: ", err)
		return false
	}
	if f >= 0 && f <= 3 {
		c.flag = f
		return true
	} else {
		fmt.Println("invalid menu flag")
		return false
	}
}

func main() {
	flag.Parse()
	client := NewClient(serverIP, serverPort)
	if client == nil {
		fmt.Println("client connect failed")
		return
	}
	go client.DealResponse()
	fmt.Println("client connect	")
	client.Run()
}
