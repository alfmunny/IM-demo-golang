package main

import (
	"net"
	"strings"
)

// User is a struct for user data
type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

// NewUser is User constructor
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

// ListenMessage is a listener for User
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}

func (t *User) GoOnline() {
	t.server.mapLock.Lock()
	t.server.OnlineMap[t.Name] = t
	t.server.mapLock.Unlock()

	// Broadcast message for other users
	t.server.Broadcast(t, "is online")

}

func (t *User) GoOffline() {
	t.server.mapLock.Lock()
	delete(t.server.OnlineMap, t.Name)
	t.server.mapLock.Unlock()

	// Broadcast message for other users
	t.server.Broadcast(t, "is offline")

}

func (t *User) SendMsg(msg string) {
	t.conn.Write([]byte(msg))
}

func (t *User) DoMessage(msg string) {
	if msg == "who" {
		t.server.mapLock.Lock()
		for _, user := range t.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "online\n"
			t.SendMsg(onlineMsg)
		}
		t.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {

		newName := strings.Split(msg, "|")[1]
		_, ok := t.server.OnlineMap[newName]
		if ok {
			t.SendMsg("Name already exists")
		} else {
			t.server.mapLock.Lock()
			delete(t.server.OnlineMap, t.Name)
			t.server.OnlineMap[newName] = t
			t.server.mapLock.Unlock()
			t.Name = newName
			t.SendMsg("User is renamed: " + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]

		if remoteName == "" {
			t.SendMsg("Message format not correct, please use to|name|message.\n")
			return
		}

		remoteUser, ok := t.server.OnlineMap[remoteName]
		if !ok {
			t.SendMsg("No such user")
			return
		}

		content := strings.Split(msg, "|")[2]

		if content == "" {
			t.SendMsg("No message content, please seond again")
			return
		}

		remoteUser.SendMsg(t.Name + " said to You: " + content + "\n")

	} else {
		t.server.Broadcast(t, msg)
	}
}
