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

func NewClient(serverIp string, port int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: port,
		flag:       999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, port))

	if err != nil {
		fmt.Println("net.Dial error", err)
		return nil
	}

	client.conn = conn

	return client

}

var serverIp string
var port int

func init() {

	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器ip地址")
	flag.IntVar(&port, "port", 9000, "设置服务器port")

}

func main() {
	flag.Parse()
	client := NewClient(serverIp, port)

	if client == nil {
		fmt.Println(">>>>>连接失败<<<<<< ")
		return
	}
	go client.DealResponse()
	fmt.Println(">>>>>连接成功<<<<<< ")

	client.Run()

}

func (client *Client) Run() {
	for client.flag != 0 {
		//循环调用menu。直到输入合法数字
		for client.menu() != true {
		}

		switch client.flag {

		case 1:
			client.PublicChat()
			break

		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))

	if err != nil {
		fmt.Println("conn.Write error:", err)
		return
	}

}
func (client *Client) PrivateChat() {

	var toName string

	client.SelectUsers()
	fmt.Println("-----请输入聊天对象的用户名，exit退出-------")

	fmt.Scanln(&toName)

	for toName != "exit" {
		var chatMsg string
		fmt.Println("-----输入聊天内容，exit退出------")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if chatMsg != "" {
				sendMsg := "to|" + toName + "|" + chatMsg + "\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write error:", err)
					break
				}
			}

			chatMsg = ""
			fmt.Println("-----输入聊天内容，exit退出------")
			fmt.Scanln(&chatMsg)

		}

		client.SelectUsers()
		fmt.Println("-----请输入聊天对象的用户名，exit退出-------")

		fmt.Scanln(&toName)
	}

}

func (client *Client) PublicChat() {
	fmt.Println("-----输入聊天内容，exit退出------")

	var chatMsg string

	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if chatMsg != "" {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write error:", err)
				break
			}
		}

		chatMsg = ""
		fmt.Println("-----输入聊天内容，exit退出------")
		fmt.Scanln(&chatMsg)

	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write error : ", err)
		return false
	}

	return true
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.改名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("请输入合法数字")
		return false
	}
}
