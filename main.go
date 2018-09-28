package main

import (
	"d7024e"
	"container/list"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {

	fmt.Println("Listening... Give command")

	//Skapar allt och startar listen på alla
	id1 := NewKademliaID("kuk")
	id2 := NewKademliaID("fitta")

	contact1 := NewContact(id1, "172.18.0.2:8080")
	contact2 := NewContact(id2, "172.18.0.3:8080")

	node1 := NewKademlia(contact1)
	node2 := NewKademlia(contact2)

  //time.Sleep(5*60*1000 * time.Millisecond)

	go node1.GetNetwork().Listen(contact1)
	go node2.GetNetwork().Listen(contact2)

	node1.GetNetwork().SendPingMessage(&contact2)

}
