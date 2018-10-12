package d7024e

//Komment
import (
	"fmt"
	"kadlab/protobuf"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type Network struct {
	me              Contact
	rt              *RoutingTable
	lookupResp      [][]Contact
	lookupResponder []Contact
	pingResp        []Contact
	mtx             *sync.Mutex
	fs              FileSystem
}

func NewNetwork(me Contact, rt *RoutingTable) Network {
	network := Network{}
	network.me = me
	network.rt = rt
	network.mtx = &sync.Mutex{}
	network.fs = NewFileSystem()
	return network
}

// type RpcHandler struct {
//   network *Network
//   mutex   *sync.Mutex
// }
//
// func NewRpcHandler(n *Network) *RpcHandler {
//   rpc := &RpcHandler{}
//   rpc.network = n
//   rpc.mutex = &sync.Mutex{}
//   return rpc
// }

func handleRPC(ch chan []byte, me *Contact, net *Network) {
	rawdata := <-ch
	message := &protobuf.Kmessage{}

	err := proto.Unmarshal(rawdata, message)
	if err != nil {
		fmt.Println(err)
	}

	//Add new contact to RT if needed
	newContact := NewContact(NewKademliaID(message.GetSenderId()), message.GetSenderAddr())
	if newContact.String() == me.String() {
		return
	}
	//fmt.Println("refreshar med ny contact", message.GetSenderAddr())
	net.RefreshRT(newContact)

	//Ping findnode findvalue store
	switch message.GetLabel() {
	case "ping":
		fmt.Println("Just got pinged omg")
		fmt.Println("label:", message.GetLabel())
		fmt.Println("senderID:", message.GetSenderId())
		fmt.Println("senderAddr:", message.GetSenderAddr())

		resp := buildMessage([]string{"pingresp", me.ID.String(), me.Address})
		sendMessage(message.GetSenderAddr(), resp)

	case "pingresp":
		net.mtx.Lock()
		fmt.Println("I got pingresponse from:", message.GetSenderId(), message.GetSenderAddr())
		id := NewKademliaID(message.GetSenderId())
		responder := NewContact(id, message.GetSenderAddr())
		net.pingResp = append(net.pingResp, responder)
		net.mtx.Unlock()
		fmt.Printf("Utskrift av pingresp", "%v,\n", net.pingResp, "\n")

	case "lookup":
		net.mtx.Lock()
		fmt.Println("I got a lookupreq from:", message.GetSenderId(), message.GetSenderAddr())
		id := NewKademliaID(message.GetLookupId())
		kclosest := net.rt.FindClosestContacts(id, 20)
		net.mtx.Unlock()
		t := ""
		for i := 0; i < len(kclosest); i++ {
			t = t + kclosest[i].String() + "\n"
		}
		resp := buildMessage([]string{"lookupresp", me.ID.String(), me.Address, message.GetLookupId(), t})
		sendMessage(message.GetSenderAddr(), resp)

	case "lookupresp":
		fmt.Println("I got a lookuprespons from:", message.GetSenderId(), message.GetSenderAddr())
		t := string(message.GetLookupResp())

		contacts := []Contact{}
		s := strings.Split(t, "\n")
		fmt.Println("The response in string is:", s)

		for i := 0; i < len(s)-1; i++ {
			row := strings.Split(s[i], "\"")
			contacts = append(contacts, NewContact(NewKademliaID(row[1]), row[3]))
		}

		net.mtx.Lock()
		net.lookupResp = append(net.lookupResp, [][]Contact{contacts}...)
		id := NewKademliaID(message.GetSenderId())
		responder := NewContact(id, message.GetSenderAddr())
		net.lookupResponder = append(net.lookupResponder, responder)

		//Uppdaterar routing table
		for i := 0; i < len(net.lookupResp); i++ {
			for j := 0; j < len(net.lookupResp[i]); j++ {
				time.Sleep(50 * time.Millisecond)
				net.RefreshRT(net.lookupResp[i][j])
			}
		}
		net.mtx.Unlock()

	case "lookupdata":
		// return file if node has it, otherwise return kclosest

	default:
		fmt.Println("Wrong message")
		fmt.Println("label:", message.GetLabel())
		fmt.Println("senderID:", message.GetSenderId())
		fmt.Println("senderAddr:", message.GetSenderAddr())
	}

}

func buildMessage(input []string) *protobuf.Kmessage {
	switch input[0] {
	case "ping":
		fmt.Println("Building Ping")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
		}
		return message
	case "pingresp":
		fmt.Println("Bulding pingresp")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
		}
		return message
	case "lookup":
		fmt.Println("Building lookup")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
			LookupId:   *proto.String(input[3]),
		}
		return message
	case "lookupresp":
		fmt.Println("Building lookupresp")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
			LookupId:   *proto.String(input[3]),
			LookupResp: *proto.String(input[4]),
		}
		return message
	case "lookupdata":
		fmt.Println("Building lookupdata")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
			Key:        *proto.String(input[3]),
		}
		return message
	default:
		fmt.Println("Building Error message")
		message := &protobuf.Kmessage{
			Label:      *proto.String("Error"),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
		}
		return message
	}
}

func (network *Network) SendPingMessage(contact *Contact) bool {
	network.pingResp = []Contact{}
	message := buildMessage([]string{"ping", network.me.ID.String(), network.me.Address})
	sendMessage(contact.Address, message)

	fmt.Println("Skickat ping väntar på svar")
	time.Sleep(time.Second * 2)
	fmt.Println("Väntat klart")
	network.mtx.Lock()
	if network.inPingResp(contact) {
		fmt.Println("Fick svar")
		network.mtx.Unlock()
		return true
	} else {
		fmt.Println("Fick inte ett svar")
		network.mtx.Unlock()
		return false
	}

}

func sendMessage(Address string, message *protobuf.Kmessage) {
	if len(Address) >= 14 {
		//fmt.Println("send to anddress: ", Address)
		data, err := proto.Marshal(message)
		if err != nil {
			fmt.Println("Marshal Error: ", err)
		}

		Conn, err := net.Dial("udp", Address)
		if err != nil {
			fmt.Println("UDP-Error: ", err)
		}
		defer Conn.Close()
		_, err = Conn.Write(data)
		if err != nil {
			fmt.Println("Write Error: ", err)
		}
	}

}

func (network *Network) Listen(me Contact) {
	fmt.Println("Lyssnar")
	Addr, err1 := net.ResolveUDPAddr("udp", me.Address)
	Conn, err2 := net.ListenUDP("udp", Addr)
	defer Conn.Close()

	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Resolve error:", err1)
		fmt.Println("Listen error:", err2)
	}

	ch := make(chan []byte)
	buffer := make([]byte, 4096)

	for {
		time.Sleep(10 * time.Millisecond)
		fmt.Println(me.Address, "Väntar")
		n, _, err1 := Conn.ReadFromUDP(buffer)

		if err1 != nil {
			fmt.Println("Read Error:", err1)
		}

		rawdata := buffer[:n]
		go handleRPC(ch, &me, network)
		ch <- rawdata
	}
}

//Find node RPC
func (network *Network) SendFindContactMessage(contact *Contact, targetid *KademliaID) {
	network.lookupResp = [][]Contact{}
	network.lookupResponder = []Contact{}
	message := buildMessage([]string{"lookup", network.me.ID.String(), network.me.Address, targetid.String()})
	sendMessage(contact.Address, message)
}

//Find value RPC
func (network *Network) SendFindDataMessage(contact *Contact, hash string) {
	message := buildMessage([]string{"lookupdata", network.me.ID.String(), network.me.Address, hash})
	sendMessage(contact.Address, message)
}

//Store RPC
func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) RefreshRT(contact Contact) {
	//If me do nothing
	if contact.String() == network.me.String() {
		return
	}

	// find bucket contact should be placed in
	bucket := network.rt.buckets[network.rt.getBucketIndex(contact.ID)]

	if bucket.ContactinBucket(contact) {
		bucket.AddContact(contact)
	} else {
		if bucket.list.Len() < bucketSize { // bucket not full -> add contact in front
			bucket.AddContact(contact)
		} else {
			fmt.Println("Bucket full")
			oldestContact := bucket.list.Back().Value.(Contact)
			svar := network.SendPingMessage(&oldestContact)
			if !svar {
				fmt.Println("Contact DÖD")
				bucket.RemoveContact(oldestContact)
				bucket.AddContact(contact)
			} else {
				fmt.Println("Contact Levande")
			}
		}
	}
}

// checks if contact responded to ping, removes the response if so
func (network *Network) inPingResp(c *Contact) bool {
	for i := 0; i < len(network.pingResp); i++ {
		if c.String() == network.pingResp[i].String() {
			if (i == 0) && (len(network.pingResp) == 1) {
				network.pingResp = []Contact{}
			} else if i == 0 {
				network.pingResp = network.pingResp[i+1:]
			} else if i == len(network.pingResp) {
				network.pingResp = network.pingResp[:i-1]
			} else {
				network.pingResp = append(network.pingResp[:i-1], network.pingResp[i+1:]...)
			}
			fmt.Println("Jag har hittat pingen")
			return true
		} else {
		}
	}
	return false
}
