package main

//Komment

import (
	"fmt"
	kad "kadlab/d7024e"
	"time"
)

func main() {

	fmt.Println("Listening... Give command")

	//Skapar allt och startar listen på alla
	id1 := kad.NewRandomKademliaID()
	id2 := kad.NewRandomKademliaID()

	contact1 := kad.NewContact(id1, "192.168.0.100:8080")
	contact2 := kad.NewContact(id2, "192.168.0.106:8080")

	node1 := kad.NewKademlia(contact1)
	node2 := kad.NewKademlia(contact2)
	time.Sleep(10 * 1000 * time.Millisecond)

	go node1.GetNetwork().Listen(contact1)
	go node2.GetNetwork().Listen(contact2)
	time.Sleep(3 * time.Millisecond)

	//node1.GetNetwork().SendPingMessage(&contact2)
	node1.GetNetwork().SendFindContactMessage(&contact2, contact1.ID)
	time.Sleep(10 * 1000 * time.Millisecond)
	node1.GetNetwork().GetRT().PrintRoutingTable()
	node2.GetNetwork().GetRT().PrintRoutingTable()
	time.Sleep(10 * 1000 * time.Millisecond)

	fmt.Println("Startar lookup")
	svar := node1.LookupContact(&contact2)
	fmt.Printf("Svaret blev","%v\n", svar)

}
