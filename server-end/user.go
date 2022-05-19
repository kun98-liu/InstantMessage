package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	Server *Server
}

//创建用户API
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		Server: server,
	}

	go user.ListenMsg()

	return user
}

//监听当前User Channel的方法
func (user *User) ListenMsg() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}

}

func (user *User) Online() {

	user.Server.mapLock.Lock()
	user.Server.OnlineMap[user.Name] = user
	user.Server.mapLock.Unlock()

	user.Server.Broadcast(user, "get on line")
}

func (user *User) Offline() {
	user.Server.mapLock.Lock()
	delete(user.Server.OnlineMap, user.Name)
	user.Server.mapLock.Unlock()

	user.Server.Broadcast(user, "get off line")

}

func (user *User) DoMsg(msg string) {
	//who语句查询在线用户
	if msg == "who" {
		user.Server.mapLock.Lock()

		for _, cur := range user.Server.OnlineMap {
			onlineMsg := "[" + cur.Addr + "]" + cur.Name + ":" + " alive...\n"
			user.SendMsg(onlineMsg)
		}
		user.Server.mapLock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {

		newName := strings.Split(msg, "|")[1]

		_, ok := user.Server.OnlineMap[newName]
		if ok {
			user.SendMsg("当前用户名被使用\n")
		} else {
			user.Server.mapLock.Lock()

			delete(user.Server.OnlineMap, user.Name)
			user.Server.OnlineMap[newName] = user

			user.Server.mapLock.Unlock()

			user.Name = newName

			user.SendMsg("用户名已更新：" + user.Name + "\n")

		}

	} else if len(msg) > 4 && msg[:3] == "to|" {
		//1.获取对方用户名

		toName := strings.Split(msg, "|")[1]
		toMsg := strings.Split(msg, "|")[2]

		if toName == "" {
			user.SendMsg("消息格式不正确")
			return
		}

		//2.根据用户名得到对象
		toUser, ok := user.Server.OnlineMap[toName]

		if !ok {
			user.SendMsg("用户不存在\n")
			return
		}

		if toMsg == "" {
			user.SendMsg("消息为空\n")
			return
		}

		toUser.SendMsg(user.Name + ":" + toMsg)

	} else {

		user.Server.Broadcast(user, msg)
	}
}

func (user *User) SendMsg(msg string) {

	user.conn.Write([]byte(msg))

}
