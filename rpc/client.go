package main

import (
	"log"
	"net/rpc"
)

type Args1 struct {
	A, B int
}

type Quotient1 struct {
	Quo, Rem int
}

func main() {
	server := "127.0.0.1:2345"
	client, err := rpc.DialHTTP("tcp", server)
	if err != nil {
		log.Fatal(err)
		return
	}
	a := Args1{17, 8}
	var reply int
	if err = client.Call("Arith.Multiply", &a, &reply); err != nil {
		log.Println(err)
		return
	}
	log.Println("result: ", reply)

	var quo Quotient1
	if err = client.Call("Arith.Divide", &a, &quo); err != nil {
		log.Println(err)
		return
	}

	log.Println("quo: ", quo.Quo, " rem: ", quo.Rem)
}
