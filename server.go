package main

import (
	"fmt"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	Msg chan string
}

//创建一个Server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Msg:       make(chan string),
	}

	return server
}

//监听Msg Channel， 有消息就发送给在线User
func (server *Server) ListenMsg() {
	for {
		msg := <-server.Msg

		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.mapLock.Unlock()

	}
}

/**
建立连接后
*/
func (server *Server) Handler(conn net.Conn) {

	user := NewUser(conn, server)

	user.Online()

	//监听用户是否活跃
	isLive := make(chan bool)

	//接收客户端传来的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err:", err)
				return
			}
			msg := string(buf[:n-1])

			user.DoMsg(msg)
			isLive <- true
		}
	}()

	for {

		select {
		case <-isLive:
			//重置定时器
			//不做任何事，就是为了激活select
		case <-time.After(time.Second * 30):
			//chaoshi
			user.SendMsg("待机时间过长\n")

			close(user.C)

			conn.Close()
			runtime.Goexit()
		}
	}

}

//广播消息
func (server *Server) Broadcast(user *User, msg string) {

	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg

	server.Msg <- sendMsg

}

//启动服务器的接口
func (server *Server) Start() {
	//listen

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))

	if err != nil {
		fmt.Println("listener error")
		return
	}

	//close
	defer listener.Close()

	//启动ListenMsg的goroutine
	go server.ListenMsg()

	for {

		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accpet error")
			continue
		}

		go server.Handler(conn)
		//handler
	}

}
