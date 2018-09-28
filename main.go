package main

import (
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
)

func main() {

	etem := &Person{
		Name: "Etem",
		Age:  14,
	}

	data, err := proto.Marshal(etem)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}

	fmt.Println(data)

	newEtem := &Person{}
	err = proto.Unmarshal(data, newEtem)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}

	fmt.Println(newEtem.GetAge())
	fmt.Println(newEtem.GetName())

}
