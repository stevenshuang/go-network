package main

import (
	"io/ioutil"
	"log"
	"net"
)

func tcpHeadReq(networkType, domain string) {
    tcpAddr, err := net.ResolveTCPAddr(networkType, domain)
    if err != nil {
        log.Fatalf("parse tcp address error: %s", err)
        return
    }


    log.Println("tcp address: ", tcpAddr.String())
    conn, err := net.DialTCP(networkType, nil, tcpAddr)
    if err != nil {
        log.Fatalf("dial %s error: %s", tcpAddr.String(), err)
        return
    }
    
    _, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
    if err != nil {
        log.Fatalf("head %s error: %s", tcpAddr.String(), err)
        return
    }
    r, err := ioutil.ReadAll(conn)
    if err != nil {
        log.Fatalf("read response error: %s", err)
        return
    }
    log.Println("result: ", string(r))
}
