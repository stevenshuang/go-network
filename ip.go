package main

import (
	"fmt"
	"net"
)

func useIP(name string) error {
	addr := net.ParseIP(name)
	if addr == nil {
		return fmt.Errorf("ivalid address %s", name)
	}
	mask := addr.DefaultMask()
	network := addr.Mask(mask)
	ones, bits := mask.Size()
	fmt.Printf("address=[%s]\ndefault mask length=[%d]\nleading ones=[%d]\nmask=[%s]\nnetwork=[%s]\n",
		addr.String(), bits, ones, mask.String(), network.String())
	return nil
}
