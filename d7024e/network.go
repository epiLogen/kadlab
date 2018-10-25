package d7024e

//Komment
import (
	"fmt"
	"kadlab/protobuf"
	"net"
	"strings"
	"sync"
	"time"
	"crypto/sha1"

	"github.com/golang/protobuf/proto"
)

type Network struct {
	me              Contact
	rt              *RoutingTable
	lookupResp      [][]Contact
	lookupResponder []Contact
	pingResp        []Contact
	mtx             *sync.RWMutex
	fs              FileSystem
	data 						string
}

func NewNetwork(me Contact, rt *RoutingTable) Network {
	network := Network{}
	network.me = me
	network.rt = rt
	network.mtx = &sync.RWMutex{}
	network.fs = NewFileSystem()
	network.data = ""
	return network
}

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
		resp := buildMessage([]string{"pingresp", me.ID.String(), me.Address})
		sendMessage(message.GetSenderAddr(), resp)

	case "pingresp":
		//net.mtx.Lock()
		//fmt.Println("I got pingresponse from:", message.GetSenderId(), message.GetSenderAddr())
		id := NewKademliaID(message.GetSenderId())
		responder := NewContact(id, message.GetSenderAddr())
		net.pingResp = append(net.pingResp, responder)
		//net.mtx.Unlock()
		//fmt.Printf("Utskrift av pingresp", "%v,\n", net.pingResp, "\n")

	case "lookup":
		//net.mtx.Lock()
		fmt.Println("I got a lookupreq from:", message.GetSenderId(), message.GetSenderAddr())
		id := NewKademliaID(message.GetLookupId())
		kclosest := net.rt.FindClosestContacts(id, 20)
		//net.mtx.Unlock()
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
		//fmt.Println("The response in string is:", s)

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
		time.Sleep(10 * time.Millisecond)

		//Uppdaterar routing table ebola bugfix
		breakern := false
		for i := 0; i < len(net.lookupResp); i++ {
			//net.mtx.RLock()
			if breakern == true {
				break
			} else if i > len(net.lookupResp) -1 {
				//net.mtx.RUnlock()
				break
			}
				for j := 0; j < len(net.lookupResp[i]); j++ {
					if i > len(net.lookupResp) -1 {
						breakern = true
						break
					}
					if j > len(net.lookupResp[i]){
						breakern = true
						break
					}
					net.RefreshRT(net.lookupResp[i][j])
					if i > len(net.lookupResp) -1 {
						breakern = true
						break
					}
				}
			//net.mtx.RUnlock()
			time.Sleep(10 * time.Millisecond)
		}

	case "lookupdata":
		fmt.Println("Just got lookupdata request key is: ", message.GetKey())
		key := NewKademliaIDnp(message.GetKey())
		file := net.fs.GetFile(key)

		if file == "" { // file not found -> send back kclosest
			fmt.Println("File not found")
			kclosest := net.rt.FindClosestContacts(&key, 20)

			s := ""
			for i := 0; i < len(kclosest); i++ {
				s = s + kclosest[i].String() + "\n"
			}
			response := buildMessage([]string{"lookupresp", me.ID.String(), me.Address, key.String(), s})
			sendMessage(message.GetSenderAddr(), response)
		} else { // file found -> send back file
			fmt.Println("File found")
			response := buildMessage([]string{"lookupdataresp", me.ID.String(), me.Address, file})
			sendMessage(message.GetSenderAddr(), response)

		}
	case "lookupdataresp":
		//fmt.Println("Sending my file :)")
		net.mtx.Lock()
		net.data = string(message.Data)
		net.mtx.Unlock()
		time.Sleep(10 * time.Millisecond)

	case "store":
		fmt.Println("Got a store request")
		hash := []byte(message.GetData())
		key := KademliaID(sha1.Sum(hash))

		//key := NewKademliaID(message.GetKey())	//Provar regen
		file := message.GetData()
		publisher := message.GetSenderId()
		net.fs.Store(key, file, publisher)
		time.Sleep(50 * time.Millisecond)

	case "pin":
		fmt.Println("Got a pin request")
		key := NewKademliaIDnp(message.GetKey())
		file := net.fs.GetFile(key)

		if file != "" {
			net.fs.Pin(key)
		}

	case "unpin":
		fmt.Println("Got an unpin request")
		key := NewKademliaIDnp(message.GetKey())
		file := net.fs.GetFile(key)

		if file != "" {
			net.fs.Unpin(key)
		}


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
		//fmt.Println("Building Ping")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
		}
		return message
	case "pingresp":
		//fmt.Println("Bulding pingresp")
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
	case "lookupdataresp":
		fmt.Println("Building lookupdataresp")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
			Data:       *proto.String(input[3]),
		}
		return message
	case "store":
		fmt.Println("Building store")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
			Key:        *proto.String(input[3]),
			Data:       *proto.String(input[4]),
		}
		return message

	case "pin":
		fmt.Println("Building pin")
		message := &protobuf.Kmessage{
			Label:      *proto.String(input[0]),
			SenderId:   *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
			Key:        *proto.String(input[3]),
		}
		return message
	case "unpin":
		fmt.Println("Building unpin")
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
	//network.pingResp = []Contact{}
	message := buildMessage([]string{"ping", network.me.ID.String(), network.me.Address})
	sendMessage(contact.Address, message)

	//fmt.Println("Skickat ping väntar på svar")
	time.Sleep(time.Second * 2)
	//fmt.Println("Väntat klart")
	if network.inPingResp(contact) {
		//fmt.Println("Fick svar")
		time.Sleep(10 * time.Millisecond)
		return true
	} else {
		//fmt.Println("Fick inte ett svar")
		time.Sleep(10 * time.Millisecond)
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
			//fmt.Println("Write Error: ", err)
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
		//fmt.Println(me.Address, "Väntar")
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
	message := buildMessage([]string{"lookup", network.me.ID.String(), network.me.Address, targetid.String()})
	sendMessage(contact.Address, message)
}

//Find value RPC
func (network *Network) SendFindDataMessage(contact *Contact, hash string) {
	message := buildMessage([]string{"lookupdata", network.me.ID.String(), network.me.Address, hash})
	sendMessage(contact.Address, message)
}

func (network *Network) SendPinMessage(contact *Contact, key string) {
	message := buildMessage([]string{"pin", network.me.ID.String(), network.me.Address, key})
	sendMessage(contact.Address, message)
}

func (network *Network) SendUnPinMessage(contact *Contact, key string) {
	message := buildMessage([]string{"unpin", network.me.ID.String(), network.me.Address, key})
	sendMessage(contact.Address, message)
}

//Store RPC
func (network *Network) SendStoreMessage(contact *Contact, key string, data string) {
	//fmt.Println("Storen kallad jao")
	message := buildMessage([]string{"store", network.me.ID.String(), network.me.Address, key, data})
	sendMessage(contact.Address, message)
}

func (network *Network) RefreshRT(contact Contact) {
	//If me do nothing
	if contact.String() == network.me.String() {
		return
	}

	// find the bucket contact should be placed in
	bucket := network.rt.buckets[network.rt.getBucketIndex(contact.ID)]

	if bucket.ContactinBucket(contact) {
		bucket.AddContact(contact)
	} else {
		if bucket.list.Len() < bucketSize { // bucket not full -> add contact in front
			bucket.AddContact(contact)
		} else {
			//fmt.Println("Bucket full")
			oldestContact := bucket.list.Back().Value.(Contact)
			svar := network.SendPingMessage(&oldestContact)
			if !svar {
				//fmt.Println("Contact nere")
				bucket.RemoveContact(oldestContact)
				bucket.AddContact(contact)
			} else {
				//fmt.Println("Contact Levande")
			}
		}
	}
}

// checks if contact responded to ping and emties pingresp
func (network *Network) inPingResp(c *Contact) bool {
	if len(network.pingResp) > 1 {
		network.pingResp = []Contact{}
		return true
	}
	network.pingResp = []Contact{}
	return false
}
