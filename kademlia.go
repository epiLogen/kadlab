package d7024e

type Kademlia struct {
  net Network
  mutex *sync.Mutex
}

func (kademlia *Kademlia) GetNetwork() *Network {
  return &kademlia.net
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
