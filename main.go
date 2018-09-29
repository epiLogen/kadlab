package main

//Komment

import (
	kad "kadlab/d7024e"
	"fmt"
	"time"
)

func main() {

	fmt.Println("Listening... Give command")

	//Skapar allt och startar listen på alla
	id1 := kad.NewRandomKademliaID()
	id2 := kad.NewRandomKademliaID()

	fmt.Println("1")

	contact1 := kad.NewContact(id1, "192.168.0.100:8080")
	contact2 := kad.NewContact(id2, "192.168.0.106:8080")

	fmt.Println("2")

	node1 := kad.NewKademlia(contact1)
	node2 := kad.NewKademlia(contact2)
	fmt.Println("3")
  time.Sleep(10*1000 * time.Millisecond)

	go node1.GetNetwork().Listen(contact1)
	go node2.GetNetwork().Listen(contact2)
  fmt.Println("4")
	time.Sleep(3 * time.Millisecond)
	fmt.Println("5")

	node1.GetNetwork().SendPingMessage(&contact2)
	time.Sleep(10*1000 * time.Millisecond)

}
