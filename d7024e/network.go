package d7024e

//Komment
import (
	"container/list"
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
}

func NewNetwork(me Contact, rt *RoutingTable) Network {
	network := Network{}
	network.me = me
	network.rt = rt
	network.mtx = &sync.Mutex{}
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
			LookupId:   *proto.String(input[3]),
			LookupResp: *proto.String(input[4]),
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

func (network *Network) SendPingMessage(contact *Contact) {
	message := buildMessage([]string{"ping", network.me.ID.String(), network.me.Address})
	sendMessage(contact.Address, message)
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
	message := buildMessage([]string{"lookupdata", network.me.ID.String(), network.me.Address, "", "", hash})
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
	//fmt.Println("Refreshing routing table")
	bucket := network.rt.buckets[network.rt.getBucketIndex(contact.ID)]

	// go through bucket to see if contact already exists
	var element *list.Element
	bucket.mtx.Lock()
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if contact.ID.Equals(nodeID) {
			fmt.Println("Equals triggad")
			element = e
		}
	}

	if element == nil { // contact not in bucket
		//fmt.Println("Contact not in bucket")

		if bucket.list.Len() < bucketSize { // bucket not full -> add contact in front
			//fmt.Println("bucket not full -> add contact in front")
			bucket.list.PushFront(contact)
		} else { // bucket is full -> ping oldest contact
			//fmt.Println("bucket is full -> ping oldest contact")
			oldestContact := bucket.list.Back().Value.(Contact)
//
//			bucket.RemoveContact(oldestContact)
//			bucket.AddContact(contact)
//
			bucket.mtx.Unlock()
			network.SendPingMessage(&oldestContact)

			time.Sleep(100 * time.Millisecond) // give 1 sec to respond
			bucket.mtx.Lock()
			if !network.inPingResp(&oldestContact) {
				fmt.Println("Han är inte alive")
				bucket.RemoveContact(oldestContact)
				bucket.AddContact(contact)
			}	else{
				fmt.Println("Han är alive")
			}
		}
	} else { // contact is in bucket -> move it to front
		//fmt.Println("contact is in bucket -> move it to front")
		bucket.list.MoveToFront(element)
	}
	bucket.mtx.Unlock()
	//network.rt.PrintRoutingTable()

}

// checks if contact responded to ping, removes the response if so
func (network *Network) inPingResp(c *Contact) bool {
	for i := 0; i < len(network.pingResp); i++ {
		if c.ID.Equals(network.pingResp[i].ID) {
			network.mtx.Lock()
			if i == 0 && len(network.pingResp) == 1{
				network.pingResp = []Contact{}
			} else if i == 0 {
				network.pingResp = network.pingResp[i+1:]
			} else if i == len(network.pingResp) {
				network.pingResp = network.pingResp[:i-1]
			} else {
				network.pingResp = append(network.pingResp[:i-1], network.pingResp[i+1:]...)
			}
			network.mtx.Unlock()
			fmt.Println("Jag har hittat pingen")
			return true
		}
	}
	return false
}
