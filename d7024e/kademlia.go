package d7024e
//Komment
import (
	"sync"
)

type Kademlia struct {
	net   Network
	mutex *sync.Mutex
}

const k = 20
const alpha = 3

func NewKademlia(me Contact) (kademlia *Kademlia){
  kademlia = new(Kademlia)
  kademlia.mutex = &sync.Mutex{}
  kademlia.net = NewNetwork(me, NewRoutingTable(me))
  return kademlia
}

func (kademlia *Kademlia) GetNetwork() *Network {
	return &kademlia.net
}

// func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
// 	// find k closest from local routing table, comes sorted by distance
// 	kclosest := kademlia.net.rt.FindClosestContacts(target, k)
//
//   contacted := []Contact{}
//   result := kclosest
//
//   // contact alpha of k closest to learn about even closer nodes
//   for i := 0; i < alpha; i++ {
//     go kademlia.net.SendFindContactMessage(&kclosest[i])
//     contacted = append(contacted, kclosest[i])
//   }
//
//   // keep contacting unqueried nodes after response/timeout
//   for {
//
// 		// response received
//     if len(kademlia.net.lookupResp) > 0 {
// 			kademlia.net.rt.refreshRT(kademlia.net.lookupResp[0])
//
// 		}
//   }
//
//
//
//
// }

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
