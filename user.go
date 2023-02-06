package main

import "net"

// User is a struct for user data
type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

// NewUser is User constructor
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

// ListenMessage is a listener for User
func (user *User) ListenMessage() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\n"))
	}
}
