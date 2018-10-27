package d7024e
//Komment
import (
	"sync"
	"time"
	"fmt"
	"crypto/sha1"
	"math/rand"
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
const republishmin = 60
const republishforpub = 24*60

func NewKademlia(me Contact) (kademlia *Kademlia){
	kademlia = new(Kademlia)
	kademlia.mutex = &sync.Mutex{}
	kademlia.net = NewNetwork(me, NewRoutingTable(me))
	go kademlia.StartRepublish()
	go kademlia.RemoveExpired()
	return kademlia
}

func (kademlia *Kademlia) StartRepublish() { //Republish that is run on every node once an hour
	fmt.Println("StartRepublish initierad")
	rand := rand.Intn(20)
	time.Sleep(time.Duration(republishmin) * time.Duration(60-rand) * 1000 * time.Millisecond)
	//fmt.Println("StartRepublish startad")
	republishfiles := kademlia.getFileSystem().GetRepublish(republishmin)

	data := ""
	hash := []byte("data")
	key := KademliaID(sha1.Sum(hash))
	keystring := key.String()
	contact := NewContact(&key, "1234567")
	//fmt.Println("Kör lookupen i minirepublishen")
	kclosest := kademlia.LookupContact(&contact)

	if len(republishfiles) == 0 || len(kclosest) == 0 {
		fmt.Println("StartRepublish: tom routing table/inge filer att republisha")
		go kademlia.StartRepublish()
		return
	}
	fmt.Println("StartRepublish kör igång o republishar")
	for i := 0; i < len(republishfiles); i++ {
		key = republishfiles[i].key
		keystring = key.String()
		contact = NewContact(&key, "1234567")
		data = republishfiles[i].data
		kclosest = kademlia.LookupContact(&contact)
		for i := 0; i < len(kclosest); i++ {
			go kademlia.net.SendStoreMessage(&kclosest[i], keystring, data)
			time.Sleep(3 * time.Millisecond)
		}
	}
	go kademlia.StartRepublish()
}

func (kademlia *Kademlia) GetNetwork() *Network {
	return &kademlia.net
}

func (kademlia *Kademlia) getFileSystem() *FileSystem {
	return &kademlia.net.fs
}

func (net *Network) GetRT() *RoutingTable {
	return net.rt
}

func (net *Network) GetFS() FileSystem {
	return net.fs
}

func (kademlia *Kademlia) RemoveExpired(){ //Removes expired files once each 24 hours
	fmt.Println("Remove Expired started in network")
	time.Sleep(24 * 60 * 60 * 1000 * time.Millisecond)
	fs := kademlia.getFileSystem()
	fs.mtx.Lock()
	fmt.Println("length of fs is ", len(fs.files))
	for i := 0; i<len(fs.files); i++ {
		fmt.Println("exp: looprunda", i)
		if fs.Expired(fs.files[i].key) {
			fmt.Println("A file has expired")
			fs.Delete(fs.files[i].key)
			fmt.Println("A file has been deleted")
		}
	}
	fs.mtx.Unlock()
	go kademlia.RemoveExpired()
}

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	// find k closest from local routing table, comes sorted by distance
	kclosest := kademlia.net.rt.FindClosestContacts(target.ID, k)

	contacted := []Contact{} //Dom som blivit kontaktade, kan ej tas bort
	contacttimes := []ContactTime{}  //tider, vi går igenom dom sen tar bort då vi hittar nån
	missedtime := []Contact{} //Dom som missat tiden kan ej tas bort
	currentcon := alpha //Antal aktiva connections
	currenttime := time.Now()


	if len(kclosest) == 0 { //Om kclosest inte räcker till
		fmt.Println("Tom routing tabell")
		return []Contact{}
	} else if len(kclosest) < alpha{
		currentcon = len(kclosest) //Sätt connections
	}

	// contact alpha of k closest to learn about even closer nodes
	for i := 0; i < currentcon; i++ {
		//fmt.Println("Startar en connection")
		contacted = append(contacted, kclosest[i])
		contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()})
		go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)
	}

	prevtime := time.Now()
	// keep contacting unqueried nodes after response/timeout
	for {
		time.Sleep(100 * time.Millisecond)
		currenttime = time.Now()
		if currenttime.Sub(prevtime).Nanoseconds() > 5000000000 {
			fmt.Printf("Lookup avslutad","längd: lookupresp", len(kademlia.net.lookupResp), "lookupresponder", kademlia.net.lookupResponder, "len contacttimes", len(contacttimes), "len contacted", len(contacted), "currentcon", currentcon, "len kclosest", len(kclosest), "\n")
			return kclosest //time out
		}


		// response received
		kademlia.net.mtx.RLock()
		if len(kademlia.net.lookupResp) > 0 && len(kademlia.net.lookupResponder) > 0 {
			kademlia.net.mtx.RUnlock()
			time.Sleep(10 * time.Millisecond)
			prevtime = time.Now()
			fmt.Println("Jag fick ett svar")


			if !isElementof(kademlia.net.lookupResponder[0], missedtime){ //Om kontakten inte missad tiden dra av connection
				currentcon = currentcon -1
			}


			//Delete from contacttimes
			contacttimes = kademlia.deleteTime(kademlia.net.lookupResponder[0], contacttimes) //Ta bort från contecttimes så vi inte går igenom han
			kademlia.net.mtx.Lock()
			//Tar bort responsen och respondern från nätverket
			if len(kademlia.net.lookupResp) <= 1 {
				//fmt.Println("Tömmer lookupResp")
				kademlia.net.lookupResp = [][]Contact{}
			} else{
				kademlia.net.lookupResp = kademlia.net.lookupResp[1:]
			}
			if len(kademlia.net.lookupResponder) <= 1 {
				kademlia.net.lookupResponder = []Contact{}
			} else{
				kademlia.net.lookupResponder = kademlia.net.lookupResponder[1:]
			}
			kademlia.net.mtx.Unlock()
			kclosest = kademlia.net.rt.FindClosestContacts(target.ID, k) //Uppdatera k closest med nya routing tablen (från responsen)

			//Start the missing connections (to alpha)
			for i := 0; i < len(kclosest); i++ { //För varje element i kclosest (Så man ej går out of bound)
				if currentcon < alpha {            //Om currentcon är mindre än alpha
					if !isElementof(kclosest[i], contacted) && !isElementof(kclosest[i], missedtime) {  //Om nuvarande element inte kontaktad och inte har missat tiden, kontakta.
						fmt.Println("Startar en missing connection")
						go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)
						contacted = append(contacted, kclosest[i]) //Lägg till i contacted
						contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()}) //Lägg till tid
						currentcon = currentcon + 1 //Öka current con
						time.Sleep(100 * time.Millisecond)
					}
				}else{ //Om det finns tillräckligt med cons avbryt ökningen direkt
					break
				}
			}
		} else{
			kademlia.net.mtx.RUnlock()
			time.Sleep(250 * time.Millisecond)
		}
		time.Sleep(250 * time.Millisecond)
		//Kolla om nån är sen. Gå igenom alla
		for i := 0; i<len(contacttimes)-1; i++ {
			//fmt.Println("Kollar times")
			if time.Now().Sub(contacttimes[i].ct).Nanoseconds() > 500000000 { //nån är sen, ta bort dom
					//lägg till i missed time och minska currentconnections (väntar ej längre på han)

				missedtime = append(missedtime, contacttimes[i].contact)
				currentcon = currentcon -1
				//fmt.Println("Nån har missat tiden") //Tas bort från contact times (Vi ska inte gå igenom dom igen)

				//Ta bort från routing tablen så dom inte kommer med i kclosest igen. Om detta är en nod.
				if contacttimes[i].contact.Address != "1234567"{
					bucket := kademlia.net.rt.buckets[kademlia.net.rt.getBucketIndex(contacttimes[i].contact.ID)]
					bucket.RemoveContact(contacttimes[i].contact)
				}
				kclosest = kademlia.net.rt.FindClosestContacts(target.ID, k) //Updatera kclosest

				for j := 0; j < len(kclosest); j++ { // Gå igenom kclosest

					if !isElementof(kclosest[j], contacted) && !isElementof(kclosest[j], missedtime) { //om nån ej blivit kontaktad och inte missat tiden
						fmt.Println("Startar en missing connection")
						go kademlia.net.SendFindContactMessage(&kclosest[i], target.ID)  //Kontakta han
						contacted = append(contacted, kclosest[i])
						contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()})
						currentcon = currentcon + 1 //Öka current con
						break							//Endast en anslutning per missad tid
					}
					if currentcon >= alpha { //Om tillräckligt med cons avbryt dirr
						break
					}
				}
				contacttimes = kademlia.deleteTime(contacttimes[i].contact, contacttimes) //Ta bort från contacttimes i slutet så man undviker indexfel
			}
		}

		terminate := true
		for i := 0; i < len(kclosest); i++ {
			if !isElementof(kclosest[i], contacted){ // Om nån inte blivit kontaktad
				terminate = false
			} else if !isElementof(kclosest[i], missedtime) { //eller om nån har det men inte har missat tiden
				terminate = false
			}
		}

		if currentcon <= 0 ||  currentcon > alpha || len(contacted) > 100 { //If we dont have any waiting connections terminate (timeout)
			terminate = true
		}


		if terminate { //Om terminate avsluta och returnera kclosest
			fmt.Println("\n")
			fmt.Printf("Lookup avslutad","längd: lookupresp", len(kademlia.net.lookupResp), "lookupresponder", kademlia.net.lookupResponder, "len contacttimes", len(contacttimes), "len contacted", len(contacted), "currentcon", currentcon, "len kclosest", len(kclosest), "\n")
			//fmt.Println("Terminerar")
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

	func (kademlia *Kademlia) LookupData(hash string) string {

		key := NewKademliaIDnp(hash)
		fs := kademlia.GetNetwork().GetFS()
		if fs.GetFile(key) != ""{
			fmt.Println("LookupData: I have the file")
			//return kademlia.net.fs.GetFile(key)
		}

		kclosest := kademlia.net.rt.FindClosestContacts(&key, k)

		contacted := []Contact{}
		contacttimes := []ContactTime{}
		missedtime := []Contact{}
		currentcon := alpha
		currenttime := time.Now()
		prevtime := time.Now()

		if len(kclosest) == 0 {
			fmt.Println("Tom routing tabell")
			return "Tom routing tabell"
		} else if len(kclosest) < alpha{
			currentcon = len(kclosest)
		}

		// contact alpha of k closest to learn about even closer nodes
		for i := 0; i < currentcon; i++ {
			go kademlia.net.SendFindDataMessage(&kclosest[i], hash)
			contacted = append(contacted, kclosest[i])
			contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()})
		}

		// keep contacting unqueried nodes after response/timeout
		for {
			time.Sleep(100 * time.Millisecond)
			currenttime = time.Now()
			if currenttime.Sub(prevtime).Nanoseconds() > 50000000000 {
				return kademlia.LookupDataD(hash)
			}

			if kademlia.net.data != "" { //File found
				time.Sleep(10 * time.Millisecond)
				fmt.Println("LookupData: FILE WAS FOUND")
				data := kademlia.net.data
				kademlia.net.data = ""
				time.Sleep(10 * time.Millisecond)
				return data
			}	else if len(kademlia.net.lookupResp) > 0 && len(kademlia.net.lookupResponder) > 0 {
				time.Sleep(10 * time.Millisecond)
				fmt.Println("LookupData: LOOKUPRESPONSE WAS received")
				// if reponder missed time
				if !isElementof(kademlia.net.lookupResponder[0], missedtime){
					currentcon = currentcon -1
				}

				//Delete from ct
				contacttimes = kademlia.deleteTime(kademlia.net.lookupResponder[0], contacttimes)

				kademlia.net.mtx.Lock()
				//Tar bort responsen och respondern från nätverket
				if len(kademlia.net.lookupResp) <= 1 {
					kademlia.net.lookupResp = [][]Contact{}
				} else{
					kademlia.net.lookupResp = kademlia.net.lookupResp[1:]
				}
				if len(kademlia.net.lookupResponder) <= 1 {
					kademlia.net.lookupResponder = []Contact{}
				} else{
					kademlia.net.lookupResponder = kademlia.net.lookupResponder[1:]
				}
				kademlia.net.mtx.Unlock()

				//Uppdatera kclosest
				kclosest = kademlia.net.rt.FindClosestContacts(&key, k)

				//Start the missing connections (to alpha)
				for i := 0; i < len(kclosest); i++ {
					if currentcon < alpha {
						if !isElementof(kclosest[i], contacted) && !isElementof(kclosest[i], missedtime) {
							fmt.Println("Jag startar en missing connection")
							go kademlia.net.SendFindDataMessage(&kclosest[i], hash)
							contacted = append(contacted, kclosest[i])
							contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()})
							currentcon = currentcon + 1
							//time.Sleep(100 * time.Millisecond)
						}
					}else{
						break
					}
				}
			}	else{
				time.Sleep(250 * time.Millisecond)
			}
			time.Sleep(250 * time.Millisecond)

			for i := 0; i<len(contacttimes)-1; i++ {
				//fmt.Println("Kollar times")
				if time.Now().Sub(contacttimes[i].ct).Nanoseconds() > 5000000000 { //nån är sen, ta bort dom
						//lägg till i missed time och minska currentconnections (väntar ej längre på han)

					missedtime = append(missedtime, contacttimes[i].contact)
					currentcon = currentcon -1
					//fmt.Println("Nån har missat tiden") //Tas bort från contact times (Vi ska inte gå igenom dom igen)

					kclosest = kademlia.net.rt.FindClosestContacts(&key, k) //Updatera kclosest

					for j := 0; j < len(kclosest); j++ { // Gå igenom kclosest

						if !isElementof(kclosest[j], contacted) && !isElementof(kclosest[j], missedtime) { //om nån ej blivit kontaktad och inte missat tiden
							fmt.Println("Jag startar en missing connection")
							go kademlia.net.SendFindDataMessage(&kclosest[i], hash)
							contacted = append(contacted, kclosest[i])
							contacttimes = append(contacttimes, ContactTime{kclosest[i], time.Now()})
							currentcon = currentcon + 1 //Öka current con
							break							//Endast en anslutning per missad tid
						}
						if currentcon >= alpha { //Om tillräckligt med cons avbryt dirr
							break
						}
					}
					contacttimes = kademlia.deleteTime(contacttimes[i].contact, contacttimes)
				}
			}

			terminate := true
			for i := 0; i < len(kclosest); i++ {
				if !isElementof(kclosest[i], contacted){ // Om nån inte blivit kontaktad
					terminate = false
				} else if !isElementof(kclosest[i], missedtime) { //eller om nån har det men inte har missat tiden
					terminate = false
				}
			}

			if currentcon <= 0 || len(contacted) > 100 { //If we dont have any waiting connections terminate (timeout)
				terminate = true
			}

			if terminate {
				return kademlia.LookupDataD(hash)
			}

		}
	}

	func (kademlia *Kademlia) Store(data string) {
		fmt.Println("Storen startad i kademlia")
		hash := []byte(data)
		key := KademliaID(sha1.Sum(hash))
		keystring := key.String()

		contact := NewContact(&key, "1234567")
		//fmt.Println("kallar lookup")
		fs := kademlia.getFileSystem()
		fs.Store(key, data, kademlia.GetNetwork().me.ID.String())
		kclosest := kademlia.LookupContact(&contact)

		fmt.Printf("Lookup i store blev klar Svaret och key blev", "%v\n", kclosest, keystring)

		for i := 0; i < len(kclosest); i++ {
			go kademlia.net.SendStoreMessage(&kclosest[i], keystring, data)
			time.Sleep(50 * time.Millisecond)
		}
		go kademlia.Republish(data)
	}

	func (kademlia *Kademlia) Republish(data string) { //Republish that the file publisher makes once every 24 hours
		fmt.Println("Republish rutinen startad från Storen")
		time.Sleep(republishforpub*60*1000 * time.Millisecond)
		fmt.Println("Storen startad från Republish")
		kademlia.Store(data)
	}

	func (kademlia *Kademlia) Pin(key KademliaID) {
		contact := NewContact(&key, "1234567")
		kclosest := kademlia.LookupContact(&contact)
		keystring := key.String()
		fmt.Println("Pin börjar skicka pinreqs")
		for i := 0; i < len(kclosest); i++ {
			go kademlia.net.SendPinMessage(&kclosest[i], keystring)
		}
	}

	func (kademlia *Kademlia) UnPin(key KademliaID) {
		contact := NewContact(&key, "1234567")
		kclosest := kademlia.LookupContact(&contact)
		keystring := key.String()
		fmt.Println("Pin börjar skicka unpinreqs")
		for i := 0; i < len(kclosest); i++ {
			go kademlia.net.SendUnPinMessage(&kclosest[i], keystring)
		}
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

	func (kademlia *Kademlia) LookupDataD(hash string) string { //Backup lookup (20 nodes call)
		key := NewKademliaIDnp(hash)

		contacted := []Contact{}
		kclosest := kademlia.net.rt.FindClosestContacts(&key, k)
		for i := 0; i < len(kclosest); i++ { //20 full k closest
			if kademlia.net.data != "" {
				return kademlia.net.data
			} else if !isElementof(kclosest[i], contacted) {
				go kademlia.net.SendFindDataMessage(&kclosest[i], hash)
				contacted = append(contacted, kclosest[i])
				time.Sleep(1000 * time.Millisecond)
				kclosest = kademlia.net.rt.FindClosestContacts(&key, k)
			}
		}

		kclosest = kademlia.net.rt.FindClosestContacts(&key, k)
		index := 0

		if kademlia.net.data != "" {
			return kademlia.net.data
		}

		starttime := time.Now()
		currenttime := starttime
		for {
			currenttime = time.Now()
			kclosest = kademlia.net.rt.FindClosestContacts(&key, k)
			if kademlia.net.data != "" {
				return kademlia.net.data
			}
			if currenttime.Sub(starttime).Nanoseconds() > 1000000000 { //Does not crash anymore
				return "Data not found"
			}
			time.Sleep(15 * 1000 * time.Millisecond)
			if index <= len(kclosest) && !isElementof(kclosest[index], contacted) {
				go kademlia.net.SendFindDataMessage(&kclosest[index], hash)
				contacted = append(contacted, kclosest[index])
				index = index + 1
			}
		}
	}
