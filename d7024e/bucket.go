package d7024e

import (
	"container/list"
	"sync"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
	mtx	*sync.Mutex
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.mtx = &sync.Mutex{}
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	bucket.mtx.Lock()
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
	bucket.mtx.Unlock()
}
//Removes a contact in bucket
func (bucket *bucket) RemoveContact(contact Contact) {
	bucket.mtx.Lock()
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			bucket.list.Remove(e)
			bucket.mtx.Unlock()
			break
		}
	}
	bucket.mtx.Unlock()
}

//Checks if a contact is in bucket
func (bucket *bucket) ContactinBucket(contact Contact) bool {
	bucket.mtx.Lock()
	for e := bucket.list.Front(); e != nil; e = e.Next() {
			nodeID := e.Value.(Contact).ID
			if (contact).ID.Equals(nodeID) {
				bucket.mtx.Unlock()
				return true
		}
	}
	bucket.mtx.Unlock()
	return false

}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	bucket.mtx.Lock()
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}
	
	bucket.mtx.Unlock()
	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
