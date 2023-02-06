package main

import "net"

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

	} else {
		t.server.Broadcast(t, msg)
	}
}
