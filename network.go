package d7024e

type Network struct {
  me  Contact
}

type RpcHandler struct {
  network *Network
  mutex   *sync.Mutex
}

func NewRpcHandler(n *Network) *RpcHandler {
  rpc := &RpcHandler{}
  rpc.network = n
  rpc.mutex = &sync.Mutex{}
  return rpc
}

func (network *Network) Listen(me Contact) {
	Addr, err1 := net.ResolveUDPAddr("udp", me.Address)
	Conn, err2 := net.ListenUDP("udp", Addr)
  rpchandler := NewRpcHandler(network)


	if (err1 != nil) || (err2 != nil) {
		fmt.Println("Resolve error:", err1)
		fmt.Println("Listen error:", err2)
	}

	buffer := make([]byte, 4096)
	defer Conn.Close()

	for {
		n, _, err1 := Conn.ReadFromUDP(buffer)

		if err1 != nil {
			fmt.Println("Read Error:", err1)
		}

		rawdata := buffer[:n]
    go rpchandler.handleRPC(rawdata, &me, network)

		if err2 != nil {
			fmt.Println("Unmarshal error:", err2)
		} else {
			fmt.Println("Received data:")
			fmt.Println("type:", message.GetType())
			fmt.Println("senderID:", message.GetSenderId())
			fmt.Println("senderAddr:", message.GetSenderAddr())
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (rpc *RpcHandler) handleRPC(rawdata []byte, me *Contact, net *Network){
  message := &protobuf.Kmessage{}
  err := proto.Unmarshal(rawdata, message)
  if err != nil {
    fmt.Println(err)
  }
  switch message.GetType() {
  case "ping":
    fmt.Println("Just got pinged omg")
    fmt.Println("type:", message.GetType())
    fmt.Println("senderID:", message.GetSenderId())
    fmt.Println("senderAddr:", message.GetSenderAddr())
  case default:
    fmt.Println("Wrong message")
  }

}

func buildMessage(input []string) *Kmessage {
  switch input[0] {
  case "ping":
    message = &Kmessage{
      type:     proto.String(input[0]),
      senderId: proto.String(input[1]),
      senderAddr: proto.String(input[2])
    }
    return message
  }
}

func (network *Network) SendPingMessage(contact *Contact) {
  message := buildMessage([]string{"ping", network.me.ID.String(), network.me.Address})
  sendMessage(contact.Address, message)
}

func sendMessage(Address string, message *protobuf.KademliaMessage) {
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

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
