package main

import (
	"bytes"
	"fmt"

	"golang.org/x/crypto/blowfish"
)

func main() {
    // symmetric key encryption
    key := []byte("this is key")
    cipher ,err := blowfish.NewCipher(key)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    src := []byte("hello\n\n\n")
    var enc [512]byte
    cipher.Encrypt(enc[:], src)

    var decrypt [8]byte
    cipher.Decrypt(decrypt[:], enc[:])
    res := bytes.NewBuffer(nil)
    res.Write(decrypt[:8])
    fmt.Println(string(res.Bytes()))
}
