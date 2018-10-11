package d7024e
//Komment
import (
	"sync"
	"time"
	"fmt"
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

func (net *Network) GetRT() *RoutingTable {
	return net.rt
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	// find k closest from local routing table, comes sorted by distance
	kclosest := kademlia.net.rt.FindClosestContacts(target.ID, k)

	contacted := []Contact{}
	unresponded := []Contact{}
	contacttimes := []ContactTime{}
	missedtime := []Contact{}
	currentcon := alpha

	if len(kclosest) == 0 {
		fmt.Println("Tom routing tabell")
		return []Contact{}
	} else if len(kclosest) < alpha{
		currentcon = len(kclosest)
	}

	// contact alpha of k closest to learn about even closer nodes
	for i := 0; i < currentcon; i++ {
		go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)
		contacted = append(contacted, kclosest[i])
		unresponded = append(unresponded, kclosest[i])
		contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()})
	}

	// keep contacting unqueried nodes after response/timeout
	for {
		time.Sleep(500 * time.Millisecond)

		// response received
		if len(kademlia.net.lookupResp) > 0 && len(kademlia.net.lookupResponder) > 0 {
			fmt.Println("Jag fick ett svar!!")
			//fmt.Printf("Fick svar","%v, %v, %v, %v\n", kademlia.net.lookupResp[0], kademlia.net.lookupResponder[0], unresponded, contacttimes, "\n")

			// if reponder missed time
			if isElementof(kademlia.net.lookupResponder[0], missedtime){
				missedtime = kademlia.deleteContact(kademlia.net.lookupResponder[0], missedtime)
			} else {
				currentcon = currentcon -1
			}

			//Delete from unresponded
			kademlia.net.mtx.Lock()
			unresponded = kademlia.deleteContact(kademlia.net.lookupResponder[0], unresponded)
			contacttimes = kademlia.deleteTime(kademlia.net.lookupResponder[0], contacttimes)

			//Tar bort responsen och respondern från nätverket
			if len(kademlia.net.lookupResp) <= 1 {
				kademlia.net.lookupResp = [][]Contact{}
				fmt.Println(len(kademlia.net.lookupResp))
			} else{
				kademlia.net.lookupResp = kademlia.net.lookupResp[1:]
			}
			if len(kademlia.net.lookupResponder) <= 1 {
				kademlia.net.lookupResponder = []Contact{}
				fmt.Printf("%v\n", kademlia.net.lookupResponder)
			} else{
				kademlia.net.lookupResponder = kademlia.net.lookupResponder[1:]
			}
			kademlia.net.mtx.Unlock()

			//Uppdatera kclosest
			kclosest = kademlia.net.rt.FindClosestContacts(target.ID, k)

			//Start the missing connections (to alpha)
			for i := 0; i < len(kclosest); i++ { //För varje element i kclosest (Så man ej går out of bound)
				if currentcon < alpha {            //Om currentcon är mindre än alpha
					if !isElementof(kclosest[i], contacted) && !isElementof(kclosest[i], unresponded) {  //Om nuvarande element inte kontaktad och inte väntande. Skicka RPC
						//fmt.Println("Jag ska inte köras, ny go, currentalpha är ", currentcon, len(contacted), len(unresponded))
						fmt.Println("Jag startar en missing connection")
						go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)
						contacted = append(contacted, kclosest[i])
						unresponded = append(unresponded, kclosest[i])
						contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()})
						currentcon = currentcon + 1 //Öka current con
					}
				}
			}
		}

		//Kolla om nån är sen
		for i := 0; i<len(contacttimes); i++ {
			fmt.Println("Kollar times")
			if time.Now().Sub(contacttimes[i].ct).Nanoseconds() > 5000000000000 { //nån är sen
					//lägg till i missed time och minska currentconnections (väntar ej längre på han)
				missedtime = append(missedtime, contacttimes[i].contact)
				currentcon = currentcon -1
				fmt.Println("Nån har missat tiden")

				for j := 0; j < len(kclosest); j++ { // Gå igenom kclosest

					if !isElementof(kclosest[j], contacted) && !isElementof(kclosest[j], unresponded) { //om nån ej blivit kontaktad
						//fmt.Println("sista ifen, jag ska definitivt inte köras", string(len(contacttimes)))
						go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)  //Kontakta han
						contacted = append(contacted, kclosest[i])
						unresponded = append(unresponded, kclosest[i])
						contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()})
						currentcon = currentcon + 1 //Öka current con
						break							//Endast en anslutning per missad tid
					}
				}
			}
		}

		//Kolla om det finns nån som inte har blivit kontaktad och inte missat tiden isf fortsätt
		terminate := true
		for i := 0; i < len(kclosest); i++ {
			if !isElementof(kclosest[i], contacted) || (!isElementof(kclosest[i], missedtime) && isElementof(kclosest[i], unresponded)){
				fmt.Println("Det finns ännu nån att kontakta")
				terminate = false
			}
		}

		if terminate { //Om inte avsluta och returnera kclosest
			fmt.Println("\n")
			fmt.Printf("Lookup avslutad","%v, %v, %v, %v", len(kademlia.net.lookupResp), kademlia.net.lookupResponder, unresponded, contacttimes, contacted, currentcon, "\n")
			fmt.Println("Terminerar")
			return kclosest
		}

	}
}

	//Tar bort en kontakt i array av kontakter
	func (kademlia *Kademlia) deleteContact(target Contact, contacts []Contact) []Contact {
		result := []Contact{}
		for i := 0; i < len(contacts); i++ {
			if(target.String() == contacts[i].String()){
				if i == 0 && len(contacts) == 1{
					result = []Contact{}
				} else if i == 0 {
					result = contacts[i+1:]
				} else if i == len(contacts) {
					result = contacts[:i-1]
				} else {
					result = append(contacts[:i-1], contacts[i+1:]...)
				}
				return result
			}
		}
		return contacts
	}

	//Tar bort en kontakt och tid från en array av Contacttimes
	func (kademlia *Kademlia) deleteTime(target Contact, contacttimes []ContactTime) []ContactTime {
		result := []ContactTime{}
		for i := 0; i < len(contacttimes); i++ {
				if(target.String() == contacttimes[i].contact.String()){
					if i == 0 && len(contacttimes) == 1{
						result = []ContactTime{}
					} else if i == 0 {
						result = contacttimes[i+1:]
					} else if i == len(contacttimes) {
						result = contacttimes[:i-1]
					} else {
						result = append(contacttimes[:i-1], contacttimes[i+1:]...)
					}
				return result
			}
		}
		return contacttimes
	}

	func (kademlia *Kademlia) LookupData(hash string) {
		// TODO
	}

	func (kademlia *Kademlia) Store(data []byte) {
		// TODO
	}

	//Kollar om en kontakt är ett element av en kontaktarray
	func isElementof(target Contact, contacts []Contact) bool {
		svar := false
		for i := 0; i < len(contacts); i++ {
			if(target.String() == contacts[i].String()){
				svar = true
			}
		}
		return svar
	}
