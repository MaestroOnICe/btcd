package scion

import (
	"github.com/netsec-ethz/scion-apps/pkg/pan"
)

type ScionAddr struct {
	udpAddr pan.UDPAddr
	addr    string
}

type PanConn struct {
	conn pan.ListenConn
}

func (sa *ScionAddr) Network() string {
	return "scion"
}
func (sa *ScionAddr) String() string {
	return sa.addr
}

// parses scion string in the format AS:IP:Port to an ScionAddr
func ParseAddr(address string) (pan.UDPAddr, error) {
	return pan.ResolveUDPAddr(address)
}

func (sa *ScionAddr) ParseAddr() error {
	udpAddr, err := pan.ResolveUDPAddr(sa.addr)
	if err != nil {
		return err
	}
	sa.udpAddr = udpAddr
	return nil
}

func main() {

}
