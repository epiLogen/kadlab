package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {

	address := "172.19.0.3:8081"
  //mainaddress := "172.19.0.2:8081"

	// Verify that a subcommand and one argument has been provided
	if len(os.Args) < 2 {
		fmt.Println("Missing one command or argument. Avalible commands: \n-store, \n-cat, \n-pin, \n-unpin")
		os.Exit(1)
	}
	if len(os.Args) < 3 {
		fmt.Println("Missing argument. Add a filename (for store) or a hash")
		os.Exit(1)
	}
	if len(os.Args) > 3 {
		fmt.Println("To many arguments")
		os.Exit(1)
	}

	// Switch on the subcommand
	switch os.Args[1] {
	case "store":
		content, err := ioutil.ReadFile(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		inputstring := string(content)
		fmt.Printf("Command received: %s argument: %s \n", os.Args[1], inputstring)
		send(os.Args[1], inputstring, address)
	case "cat":
		fmt.Printf("Command received: %s argument: %s \n", os.Args[1], os.Args[2])
		send(os.Args[1], os.Args[2], address)

	case "pin":
		fmt.Printf("Command received: %s argument: %s \n", os.Args[1], os.Args[2])
		send(os.Args[1], os.Args[2], address)

	case "unpin":
		fmt.Printf("Command received: %s argument: %s \n", os.Args[1], os.Args[2])
		send(os.Args[1], os.Args[2], address)

	default:
		fmt.Printf("Unavalible command: %s \n", os.Args[1])
		os.Exit(1)
	}
}

func Listener(Address string) {
	Addr, err1 := net.ResolveUDPAddr("udp", Address) //Lite weird (utan udp)
	Conn, err2 := net.ListenUDP("udp", Addr)

	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Connection Error Listen: ", err1, "\n", err2)
	}
	//read connection
	defer Conn.Close()

	buf := make([]byte, 4096)
	for {
		n, _, err := Conn.ReadFromUDP(buf)
		//fmt.Println("Received ", string(buf[0:n]), " from ", addr)
		fmt.Println(string(buf[:n]))

		if err != nil {
			fmt.Println("Read Error: ", err)
		}
		break
	}
}

func send(command string, args string, address string) {
  myIP := GetOutboundIP()

	ServerAddr, err1 := net.ResolveUDPAddr("udp", address)
	LocalAddr, err2 := net.ResolveUDPAddr("udp", myIP.String() + ":8080")
	Conn, err3 := net.DialUDP("udp", LocalAddr, ServerAddr)
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Println("UDP-Error: ", err1, err2, err3)
	}
	buf := []byte(command + "," + args)
	_, err := Conn.Write(buf)
	if err != nil {
		fmt.Println("Write Error: ", err)
	}
	Conn.Close()
	if command == "store" || command == "cat" { //wait for response
		Listener(LocalAddr.String())
	}
}
