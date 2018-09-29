package d7024e
//Komment
import (
	"kadlab/protobuf"
	"fmt"
	"net"
	"sync"
	"time"
	"github.com/golang/protobuf/proto"
)

type Network struct {
  me  Contact
  rt  *RoutingTable
  lookupResp [][]Contact
  pingResp []*Contact
  mtx *sync.Mutex
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

func handleRPC(ch chan []byte, me *Contact, net *Network){
  rawdata := <-ch
  message := &protobuf.Kmessage{}

  err := proto.Unmarshal(rawdata, message)
  if err != nil {
    fmt.Println(err)
  }
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
    net.pingResp = append(net.pingResp, me)
    net.mtx.Unlock()


  default:
    fmt.Println("Wrong message")
		fmt.Println("label:", message.GetLabel())
		fmt.Println("senderID:", message.GetSenderId())
		fmt.Println("senderAddr:", message.GetSenderAddr())
  }

}

func buildMessage(input []string) *protobuf.Kmessage {
	fmt.Println(input[0])
  switch input[0] {
  case "ping":
		fmt.Println("Building Ping")
    message := &protobuf.Kmessage{
      Label:     *proto.String(input[0]),
      SenderId: *proto.String(input[1]),
      SenderAddr: *proto.String(input[2]),
    }
    return message
	case "pingresp":
		fmt.Println("Bulding pingresp")
		message := &protobuf.Kmessage{
			Label:     *proto.String(input[0]),
			SenderId: *proto.String(input[1]),
			SenderAddr: *proto.String(input[2]),
		}
		return message
	default:
		fmt.Println("Building Error message")
		message := &protobuf.Kmessage{
			Label:     *proto.String("Error"),
			SenderId: *proto.String(input[1]),
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
	fmt.Println("1")

	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Resolve error:", err1)
		fmt.Println("Listen error:", err2)
	}

  ch := make(chan []byte)
	buffer := make([]byte, 4096)
	fmt.Println("2")

	for {
    time.Sleep(10 * time.Millisecond)
		fmt.Println("Väntar")
		n, _, err1 := Conn.ReadFromUDP(buffer)

		fmt.Println("3")

		if err1 != nil {
			fmt.Println("Read Error:", err1)
		}

		rawdata := buffer[:n]
    go handleRPC(ch, &me, network)
		fmt.Println("4")
    ch <- rawdata
		fmt.Println("5")
	}
}



func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}

func (network *Network) refreshRT(resp []Contact) {
  // should try to refresh rt using contacts from a lookupresponse
}
