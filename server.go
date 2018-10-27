package main

import (
  kademlia "kadlab/d7024e"
	"crypto/sha1"
  "strings"
	"fmt"
	"net"
)

type Server struct {
	kad *kademlia.Kademlia
}

func NewServer(address string, kad *kademlia.Kademlia) {
	server := &Server{}
	server.kad = kad
	go server.SListen(address)
}

func (server *Server) SListen(Address string) {
	Addr, err1 := net.ResolveUDPAddr("udp", Address)
	Conn, err2 := net.ListenUDP("udp", Addr)
	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Connection Error Listen: ", err1, "\n", err2)
	}
	//read connection
	defer Conn.Close()

	buf := make([]byte, 4096)
	for {
		n, addr, err := Conn.ReadFromUDP(buf)
		go handler(buf[:n], server, addr.String())

		if err != nil {
			fmt.Println("Read Error: ", err)
		}
	}
}

func handler(traffic []byte, server *Server, addr string) {
	s := strings.Split(string(traffic), ",")
	switch s[0] {
	case "store":
    hash := []byte(s[1])
		key := kademlia.KademliaID(sha1.Sum(hash))
		go server.kad.Store(s[1])
		SendCMD(key.String(), addr)
	case "cat":
		filehash := server.kad.LookupData(s[1])
		if filehash == "" {
			filehash = "Data not found"
		}
		SendCMD(filehash, addr)
	case "pin":
		server.kad.Pin(kademlia.NewKademliaIDnp(s[1]))
	case "unpin":
		server.kad.UnPin(kademlia.NewKademliaIDnp(s[1]))
	case "default":
    fmt.Println("Wrong input")
	}
}

func SendCMD(data string, address string) {
	Conn, err := net.Dial("udp", address)
	if err != nil {
		fmt.Println("UDP-Error: ", err)
	}
	buf := []byte(data)
	defer Conn.Close()
	_, err = Conn.Write(buf)
	if err != nil {
		fmt.Println("Write Error: ", err)
	}
}
