package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"
)

func main() {
    server := "localhost:2345"
    cert, err := ioutil.ReadFile("test.zs.org.pem")
    if err != nil {
        log.Fatalf("read ca file error: %s", err)
        return
    }
    certPool := x509.NewCertPool()
    certPool.AppendCertsFromPEM(cert)
    config := tls.Config{
        //InsecureSkipVerify: true,
        RootCAs: certPool,
    }
    conn, err := tls.Dial("tcp", server, &config)
    if err != nil {
        log.Println("dial ", server, " error: ", err)
        os.Exit(1)
    }

    for n := 0; n < 10; n++ {
        log.Printf("writing...%s", string(n+48))
        conn.Write([]byte("Hello: "+string(n+48)))
        var buf [512]byte
        c, err := conn.Read(buf[:])
        if err != nil {
            log.Println("read fatal: ", err)
            continue
        }
        log.Println("read data: ", string(buf[:c]))
    }
}
