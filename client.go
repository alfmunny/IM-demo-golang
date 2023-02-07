package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	IP   string
	Port int
	Name string
	conn net.Conn
	flag int
}

func NewClient(ip string, port int) *Client {
	client := &Client{
		IP:   ip,
		Port: port,
		flag: 999,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", client.IP, client.Port))

	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.conn = conn
	return client
}

func (t *Client) menu() bool {

	var flag int

	fmt.Println("1. Public Chat")
	fmt.Println("2. Private Chat")
	fmt.Println("3. Rename")
	fmt.Println("0. Exit")

	_, err := fmt.Scanln(&flag)
	if err != nil {
		fmt.Println("Scanln err:", err)
	}

	if flag >= 0 && flag <= 3 {
		t.flag = flag
		return true
	}

	fmt.Println(">>>>>>>> Please choose a valid number <<<<<<<<<<")
	return false
}

var ip string
var port int

func init() {
	flag.StringVar(&ip, "ip", "127.0.0.1", "IP address")
	flag.IntVar(&port, "port", 8888, "Port")

}

func (t *Client) UpdateName() bool {
	fmt.Println(">>>>>> Please input your name:")
	fmt.Scanln(&t.Name)
	_, err := t.conn.Write([]byte("rename|" + t.Name + "\n"))
	if err != nil {
		fmt.Println("Write err,", err)
		return false
	}

	return true
}

func (t *Client) ListUsers() {
	t.conn.Write([]byte("who\n"))
}

func (t *Client) PrivateChat() {
	t.ListUsers()

	for {
		var slectedUser string = ""
		fmt.Println(">>>>>>>>>>Please select user or type exit:")
		fmt.Scanln(&slectedUser)
		if slectedUser != "exit" {
			var message string = ""
			for {
				fmt.Println(">>>>>>>>>>Please chat or exit:")
				fmt.Scanln(&message)
				if message != "exit" {
					_, err := t.conn.Write([]byte("to|" + slectedUser + "|" + message + "\n"))
					if err != nil {
						fmt.Println("conn Write err:", err)
					}
				} else {
					break
				}
			}
		} else {
			break
		}
	}

}

func (t *Client) PublicChat() {
	chatMsg := ""
	fmt.Println(">>>>>> Please start input:")

	for chatMsg != "exit" {
		_, err := fmt.Scanln(&chatMsg)
		if err != nil {
			fmt.Println("Scanln err, ", err)
		}

		if len(chatMsg) != 0 {
			_, err := t.conn.Write([]byte(chatMsg + "\n"))
			if err != nil {
				fmt.Println("conn Write err, ", err)
				break
			}
		}
	}
}

func (t *Client) ResponseHandle() {
	io.Copy(os.Stdout, t.conn)
}

func (t *Client) Run() {
	for t.flag != 0 {
		for t.menu() != true {

		}

		switch t.flag {
		case 1:
			fmt.Println("Public chat mode")
			t.PublicChat()
			break
		case 2:
			fmt.Println("Private chat mode")
			t.PrivateChat()
			break
		case 3:
			fmt.Println("Rename mode")
			t.UpdateName()
			break
		}
	}

}

func main() {
	flag.Parse()
	client := NewClient(ip, port)
	if client == nil {
		fmt.Println(">>>>>>>>>>> conntetion failed...")
		return
	}

	fmt.Println(">>>>>>>>>>> Connection succeed...")

	go client.ResponseHandle()

	client.Run()
}
