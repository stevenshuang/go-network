package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)


func main() {
    genPublicKey()

    /*
    log.Println("load file........")
    var rsaPrivatreKey rsa.PrivateKey
    loadKey("private.key", &rsaPrivatreKey)
    log.Println("Private key primes: ", rsaPrivatreKey.Primes[0].String(), rsaPrivatreKey.Primes[1].String())
    log.Println("Private key exponent: ", rsaPrivatreKey.D.String())

    var pKey rsa.PublicKey
    loadKey("public.key", &pKey)
    log.Println("Public key modules: ", pKey.N.String())
    log.Println("Public key exponent: ", pKey.E)
    */
    log.Println("gen ca ...")
    startTlsServer()
}


func genPublicKey() {
    reader := rand.Reader
    bitSize := 2048
    key, err := rsa.GenerateKey(reader, bitSize)
    if err != nil {
        log.Fatalf("generate key error %s", err)
        os.Exit(1)
    }

    log.Println("Private key primes: ", key.Primes[0].String(), key.Primes[1].String())
    log.Println("Private key exponent: ", key.D.String())

    publicKey := key.PublicKey
    log.Println("Public key modulees: ", publicKey.N.String())
    log.Println("Public key exponent: ", publicKey.E)

    saveGobKey("private.key", key)
    saveGobKey("public.key", key.PublicKey)
    savePEMKey("private.pem", key)
}


func saveGobKey(fileName string, key interface{}) {
    file, err := os.Create(fileName)
    if err != nil {
        log.Fatalf("create file[%s] error: %s", fileName, err)
        os.Exit(1)
    }
    defer file.Close()
    encoder := gob.NewEncoder(file)
    if err = encoder.Encode(key); err != nil {
        log.Fatalf("gob encode key error: %s", err)
        os.Exit(1)
    }
}

func savePEMKey(fileName string, key *rsa.PrivateKey) {
    file, err := os.Create(fileName)
    if err != nil {
        log.Fatalf("create file[%s] error: %s", fileName, err)
        os.Exit(1)
    }
    defer file.Close()

    privateKey := &pem.Block{
        Type: "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(key),
    }
    pem.Encode(file, privateKey)
}

func loadKey(fileName string, key interface{}) {
    file, err := os.Open(fileName)
    if err != nil {
        log.Fatalf("os open file[%s] error: %s", fileName, err)
        os.Exit(1)
    }
    defer file.Close()

    decoder := gob.NewDecoder(file)
    if err = decoder.Decode(key); err != nil {
        log.Fatalf("decode %s error: %s", fileName, err)
        os.Exit(1)
    }
}


func x509CA() {
    random := rand.Reader
    var key rsa.PrivateKey
    loadKey("private.key", &key)
    
    now := time.Now()
    then := now.Add(365*24*time.Hour)
    template := x509.Certificate{
        SerialNumber: big.NewInt(1),
        Subject: pkix.Name{
            CommonName: "test.zs.org",
            Organization: []string{"ZS"},
        },
        NotBefore: now,
        NotAfter: then,
        SubjectKeyId: []byte{1, 2, 3, 4},
        KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
        BasicConstraintsValid: true,
        IsCA: true,
        DNSNames: []string{"localhost","test.zs.org"},
    }
    derBytes, err := x509.CreateCertificate(random, &template, &template, &key.PublicKey, &key)
    if err != nil {
        log.Fatalf("create ca error: %s", err)
        os.Exit(1)
    }

    caFile, err := os.Create("ca.cer")
    if err != nil {
        log.Fatalf("create ca file error: %s", err)
        os.Exit(1)
    }
    defer caFile.Close()
    caFile.Write(derBytes)

    caPEMFile, err := os.Create("test.zs.org.pem")
    if err != nil {
        log.Fatalf("create file test.zs.org.pem error: %s", err)
        os.Exit(1)
    }
    defer caPEMFile.Close()
    pem.Encode(caPEMFile, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
    
    keyPEMFile, err := os.Create("private.pem")
    if err != nil {
        log.Fatalf("create private.pem error: %s", err)
        os.Exit(1)
    }
    defer keyPEMFile.Close()
    pem.Encode(keyPEMFile, &pem.Block{Type: "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(&key)})
}

func startTlsServer() {
    x509CA()
    cert, err := tls.LoadX509KeyPair("test.zs.org.pem", "private.pem")
    if err != nil {
        log.Fatalf("load x509 key pair error: %s", err)
        os.Exit(1)
    }

    config := tls.Config{Certificates: []tls.Certificate{cert}}
    config.Time = func() time.Time {return time.Now()}
    config.Rand = rand.Reader

    server := ":2345"
    listener, err := tls.Listen("tcp", server, &config)
    if err != nil {
        log.Fatalf("start tls server error: %s", err)
        os.Exit(1)
    }

    log.Println("start listening...")
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err)
            continue
        }

        go handle(conn)
    }
}

func handle(conn net.Conn) {
    defer conn.Close()
    var buf [512]byte
    for {
        n, err := conn.Read(buf[:])
        if err != nil {
            log.Println(err)
            return
        }
        _, err = conn.Write(buf[:n])
        if err != nil {
            return
        }
    }
}
