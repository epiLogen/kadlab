package main

//Komment

import (
	"fmt"
	kad "kadlab/d7024e"
	"log"
	"net"
	"time"
)

func main() {
	//mymain()
	//time.Sleep(30 * 1000 * time.Millisecond)
	fmt.Println("hello world")
	dockermain()

}

func dockermain() {
	//Fixa IP
	myIP := GetOutboundIP()
	fmt.Println(myIP.String() + ":8080")

	//Skapa main kontakt
	mainID := kad.NewKademliaID("FFFFFFFFFFF11111111111111111111111111005")
	maincontact := kad.NewContact(mainID, "172.19.0.2:8080")

	if myIP.String() == "172.19.0.2" { //Jag är main
		fmt.Println("Jag vet att jag är main")
		node := kad.NewKademlia(maincontact)
		time.Sleep(5 * 1000 * time.Millisecond)
		go node.GetNetwork().Listen(maincontact)
		time.Sleep(40 * 1000 * time.Millisecond)      //Starta main listen
		node.GetNetwork().GetRT().PrintRoutingTable() //Printa min RT
		time.Sleep(40 * 1000 * time.Millisecond)
	} else {
		time.Sleep(30 * 1000 * time.Millisecond) //Chilla
		id1 := kad.NewRandomKademliaID()
		contact1 := kad.NewContact(&id1, myIP.String()+":8080")
		fmt.Println(contact1.String())
		node1 := kad.NewKademlia(contact1)
		go node1.GetNetwork().Listen(contact1)                               //Starta min listen
		node1.GetNetwork().SendFindContactMessage(&maincontact, contact1.ID) //Informera main om att jag finns
		time.Sleep(15 * 1000 * time.Millisecond)                             //Chilla
		node1.GetNetwork().GetRT().PrintRoutingTable()  //Print RT
		time.Sleep(15 * 1000 * time.Millisecond)

		if myIP.String() == "172.19.0.3" { //Denna nod får göra en lookup hos mainen av mainen och få 20 noder
			fmt.Println("Startar lookup")
			svar := node1.LookupContact(&maincontact)
			fmt.Printf("Svaret blev", "%v\n", svar)
			node1.GetNetwork().GetRT().PrintRoutingTable()  //Print RT
		}else{
			time.Sleep(10 * 1000 * time.Millisecond)
		}
	}

}

func mymain() {
	fmt.Println("Listening... Give command")

	//Skapar allt och startar listen på alla
	id1 := kad.NewRandomKademliaID()
	id2 := kad.NewRandomKademliaID()

	contact1 := kad.NewContact(&id1, "192.168.0.100:8080")
	contact2 := kad.NewContact(&id2, "192.168.0.106:8080")

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
	fmt.Printf("Svaret blev", "%v\n", svar)
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
