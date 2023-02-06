package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

func (t *Server) ListenMessage() {
	for {
		msg := <-t.Message

		t.mapLock.Lock()
		for _, cli := range t.OnlineMap {
			cli.C <- msg
		}
		t.mapLock.Unlock()
	}

}

func (t *Server) BroadcastUser(user *User, msg string) {
	sendMsg := "[" + user.Addr + "] " + user.Name + ": " + msg
	t.Message <- sendMsg
}

func (t *Server) Handler(conn net.Conn) {
	// User goes online
	user := NewUser(conn)

	t.mapLock.Lock()
	t.OnlineMap[user.Name] = user
	t.mapLock.Unlock()

	// Broadcast message for other users
	t.BroadcastUser(user, "is online")
}

func (t *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.IP, t.Port))
	if err != nil {
		fmt.Println("net.Listen err")
		return
	}

	go t.ListenMessage()

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("listener accept err:", err)
		}

		go t.Handler(conn)
	}

}
