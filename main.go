package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

var networkType = flag.String("net", "tcp", "tcp, udp")
var domain = flag.String("domain", "www.bing.com", "ip address")
var service = flag.String("service", "telnet", "service such as telnet")
var addrWithPort = flag.String("addr-with-port","www.bing.com:80", "domain with port")

func main() {
	flag.Parse()
    log.Println("resolve ip addr...")
	addr, err := net.ResolveIPAddr("ip", *domain)
	if err != nil {
		fmt.Fprintf(os.Stderr, "resolve ip addr error: %s", err)
		os.Exit(1)
	}

    log.Println("test ip...")
	if err := useIP(addr.String()); err != nil {
		fmt.Fprintf(os.Stderr, "test ip error: %s", err)
		os.Exit(1)
	}

    log.Println("resolve tcp addr...")
	commonServicePort(*networkType, *service, *addrWithPort)
    
    log.Println("dial tcp...")
    tcpHeadReq(*networkType, *addrWithPort)
}
