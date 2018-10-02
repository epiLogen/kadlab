package main

//Komment

import (
	"fmt"
	kad "kadlab/d7024e"
	"time"
	"log"
	"net"
)

func main() {
	//dockermain()
	mymain()

}

func dockermain(){
	//Fixa IP
	myIP := GetOutboundIP()
	fmt.Println(myIP.String() + ":8080")

	//Skapa main kontakt
	mainID := kad.NewKademliaID("FFFFFFFFFFF11111111111111111111111111005")
	maincontact := kad.NewContact(mainID,"172.19.0.2:8080")

	if myIP.String() == "172.19.0.2" { //Jag är main
		node := kad.NewKademlia(maincontact)
		go node.GetNetwork().Listen(maincontact) //Starta main listen
		node.GetNetwork().GetRT().PrintRoutingTable() //Printa min RT
	}else{
		time.Sleep(10 * 1000 * time.Millisecond) //Chilla
		id1 := kad.NewRandomKademliaID()
		contact1 := kad.NewContact(id1, myIP.String()+":8080")
		node1 := kad.NewKademlia(contact1)
		go node1.GetNetwork().Listen(contact1) //Starta min listen
		node1.GetNetwork().SendFindContactMessage(&maincontact, contact1.ID) //Informera main om att jag finns
		time.Sleep(10 * 1000 * time.Millisecond) //Chilla
		node1.GetNetwork().GetRT().PrintRoutingTable() //Printa min RT
	}



}

func mymain() {
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

func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}
