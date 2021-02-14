package main

import (
	"log"
	"net"
)

func commonServicePort(network, service, domain string) {
	port, err := net.LookupPort(network, service)
	if err != nil {
		log.Fatalf("get service[%s] port error: %s", service, err)
		return
	}
	log.Printf("service[%s] port: %d", service, port)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", domain)
	if err != nil {
		log.Fatalf("test google service port error: %s", err)
		return
	}

	log.Println(tcpAddr.String())
}

