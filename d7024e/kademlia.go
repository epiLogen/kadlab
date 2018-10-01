package d7024e
//Komment
import (
	"sync"
	"time"
)

type Kademlia struct {
	net   Network
	mutex *sync.Mutex
}

type ContactTime struct {
	contact Contact
	ct time.Time
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

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	// find k closest from local routing table, comes sorted by distance
	kclosest := kademlia.net.rt.FindClosestContacts(target, k)

	contacted := []Contact{}
	unresponded := []Contact{}
	contacttimes := []ContactTime{}
	//result := []Contact{}
	//result := kclosest
	fewerthanalpha := false
	currentalpha := alpha

	if(len(kclosest) < alpha){
		currentalpha = len(kclosest)
		fewerthanalpha = true
	}

	// contact alpha of k closest to learn about even closer nodes
	for i := 0; i < currentalpha; i++ {
		go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)
		contacted = append(contacted, kclosest[i])
		unresponded = append(unresponded, kclosest[i])
		contacttime = append(contacttimes, ContactTime{kclosest[i], time.Now()})
	}

	// keep contacting unqueried nodes after response/timeout
	for {
		time.Sleep(500 * time.Millisecond)

		// response received
		if len(kademlia.net.lookupResp) > 0 && len(kademlia.net.lookupResponder) > 0 {
			unresponded = deleteContact(kademlia.net.lookupResponder[0], unresponded)
			contacttime = deleteTime(kademlia.net.lookupResponder[0], contacttime)
			//Uppdaterar routing table
			net.mtx.Lock()
			for i := 0; i < len(kademlia.net.lookupResp); i++ {
				kademlia.net.rt.refreshRT(kademlia.net.lookupResp[i]) //Todo
			}
			kademlia.net.lookupResp = kademlia.net.lookupResp[1:]
			kademlia.net.lookupResponder = kademlia.net.lookupResponder[1:]
			net.mtx.Unlock()

			//Uppdatera kclosest
			kclosest = kademlia.net.rt.FindClosestContacts(target, k)

			if len(kclosest) < alpha{ //Sätt alpha
				currentalpha = len(kclosest)
			}
			else{
				currentalpha = alpha
			}

			for i := 0; i < currentalpha; i++ { //Rekursiv lookup
				if !isElementof(kclosest[i], contacted) && !isElementof(kclosest[i], unresponded) {
					go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)
					contacted = append(contacted, kclosest[i])
					unresponded = append(unresponded, kclosest[i])
					contacttime = append(contacttimes, ContactTime{kclosest[i], time.Now()})
				}
			}

		}

		//Kolla om nån är sen
		for i := 0; i<len(contacttime); i++ {
			if time.Now().Sub(contacttime[i].ct).Nanoseconds() > 10000000000 {
				for { //Rekursiv lookup
					if !isElementof(kclosest[i], contacted) && !isElementof(kclosest[i], unresponded) {
						go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)
						contacted = append(contacted, kclosest[i])
						unresponded = append(unresponded, kclosest[i])
						contacttime = append(contacttimes, ContactTime{kclosest[i], time.Now()})
						break
						}
					}
				}
			}

		}
	}

func (kademlia *Kademlia) deleteContact(target Contact, contacts []Contact) []Contact {
	for i := 0 < len(contacts); i++ {
		if(target.String() == contacts[i].String()){
			result := append(contacts[:i-1], contacts[i+1:])
			return result
		}
	}

	// for index, element := range contacts {
	// 	if(element.String() == target.String()){
	// 		result := append(contacts[:index-1], contacts[index+1:])
	// 	}
	}

}

func (kademlia *Kademlia) deleteTime(target Contact, contacttimes []ContactTime) []ContactTime {
	for i := 0 < len(contacts); i++ {
		if(target.String() == contacttimes[i].contact.String()){
			result := append(contacttimes[:i-1], contacttimes[i+1:])
			return result
		}
	}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func (kademlia *Kademlia) isElementof(target Contact, contacts []Contact) bool {
	for i := 0 < len(contacts); i++ {
		if(target.String() == contacts[i].String()){
			return true
		}
	}
	return false
}
