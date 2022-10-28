package scion

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsec-ethz/scion-apps/pkg/quicutil"
	"inet.af/netaddr"
)

// ScionAddr implements the net.Addr interface and represents a tor address.
type ScionAddr struct {
	udpAddr pan.UDPAddr
	addr    string
}

// Network returns "scion".
//
// This is part of the net.Addr interface.
func (sa ScionAddr) Network() string {
	return "scion"
}

// String returns the scion address.
//
// This is part of the net.Addr interface.
func (sa ScionAddr) String() string {
	return sa.addr
}

// Ensure scionAddr implements the net.Addr interface.
var _ net.Addr = (*ScionAddr)(nil)

// Wrapper for the function SplitHostPort from the pan libary
// This is analogous to net.SplitHostPort, which however refuses to handle SCION addresses.
// The address can be of the form of a SCION address (i.e. of the form "ISD-AS,[IP]:port") or in the form of "hostname:port".
func SplitHostPort(hostport string) (host, port string, err error) {
	if _, ok := IsValidAddress(hostport); ok {
		return pan.SplitHostPort(hostport)
	}
	return net.SplitHostPort(hostport)
}

// parse scion string into an ScionAddr with udp addresse and string address
func ParseAddr(address string) (ScionAddr, error) {
	saUDPAddress, err := pan.ResolveUDPAddr(address)
	if err != nil {
		fmt.Println(address, saUDPAddress)
		return ScionAddr{}, errors.New("could not parse scion address")
	}
	sa := ScionAddr{
		udpAddr: saUDPAddress,
		addr:    address,
	}
	return sa, err
}

// RegEx for identifying a scion address
// found in pan/addr.go
//var addrRegexp = regexp.MustCompile(`^(\d+-[\d:A-Fa-f]+),(\[[^\]]+\]|[^\[\]]+)$`)

// ParseUDPAddr converts an address string to a SCION address.
// The supported formats are:
//
// Recommended:
//  - isd-as,ipv4:port        (e.g., 1-ff00:0:300,192.168.1.1:8080)
//  - isd-as,[ipv6]:port      (e.g., 1-ff00:0:300,[f00d::1337]:8080)
//  - isd-as,[ipv6%zone]:port (e.g., 1-ff00:0:300,[f00d::1337%zone]:8080)
//
// Others:
//  - isd-as,[ipv4]:port (e.g., 1-ff00:0:300,[192.168.1.1]:8080)
//  - isd-as,[ipv4]      (e.g., 1-ff00:0:300,[192.168.1.1])
//  - isd-as,[ipv6]      (e.g., 1-ff00:0:300,[f00d::1337])
//  - isd-as,[ipv6%zone] (e.g., 1-ff00:0:300,[f00d::1337%zone])
//  - isd-as,ipv4        (e.g., 1-ff00:0:300,192.168.1.1)
//  - isd-as,ipv6        (e.g., 1-ff00:0:300,f00d::1337)
//  - isd-as,ipv6%zone   (e.g., 1-ff00:0:300,f00d::1337%zone)
//
// Not supported:
//  - isd-as,ipv6:port    (caveat if ipv6:port builds a valid ipv6 address,
//                         it will successfully parse as ipv6 without error)

func IsValidAddress(address string) (pan.UDPAddr, bool) {
	sa, err := pan.ParseUDPAddr(address)
	return sa, err == nil
	// parts := addrRegexp.FindStringSubmatch(address)
	// if parts == nil {
	// 	return false
	// }
	// _, err := pan.ParseIA(parts[1])
	// if err != nil {
	// 	return false
	// }
	// l3Trimmed := strings.Trim(parts[2], "[]")
	// _, err2 := netaddr.ParseIP(l3Trimmed)
	// if err2 != nil {
	// 	return false
	// }
	// return true
}

// return a connection for a scion
// using quicutil from scion-apps
// SingleStream implements an opaque, bi-directional data stream using QUIC, intending to be a drop-in replacement for TCP
func Dial(address string) (net.Conn, error) {
	fmt.Printf("SCION: dialing to %v\n", address)
	// Parse addr into a scionAddr
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	fmt.Printf("SCION: parsed %v\n", address)
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"hello-quic"},
	}
	fmt.Printf("SCION: tls set %v\n", addr)
	// Set Pinging Selector with active probing on two paths
	// selector := &pan.PingingSelector{
	// 	Interval: 2 * time.Second,
	// 	Timeout:  time.Second,
	// }
	// selector.SetActive(2)
	selector := pan.NewDefaultSelector()
	fmt.Printf("SCION: selector set %v\n", addr)
	session, err := pan.DialQUIC(context.Background(), netaddr.IPPort{}, addr, nil, selector, "", tlsCfg, nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("SCION: dialed to %v\n", addr)

	ss, err := quicutil.NewSingleStream(session)
	if err != nil {
		return nil, err
	}
	return ss, nil
}

// return a listener for a scion connection
// using quicutil from scion-apps
// SingleStreamListener is a wrapper for a quic.Listener, returning SingleStream connections from Accept
// This allows to use quic in contexts where a (TCP-)net.Listener is expected.
func Listen(addr string) (net.Listener, error) {
	fmt.Printf("SCION: listening to %v\n", addr)
	tlsCfg := &tls.Config{
		Certificates: quicutil.MustGenerateSelfSignedCert(),
		NextProtos:   []string{"hello-quic"},
	}
	fmt.Printf("SCION: tls set %v\n", addr)
	udpAddr, err := pan.ParseUDPAddr(addr)
	if err != nil {
		return nil, err
	}
	fmt.Printf("SCION: parsed %v\n", addr)
	quicListener, err := pan.ListenQUIC(context.Background(), netaddr.IPPortFrom(udpAddr.IP, udpAddr.Port), nil, tlsCfg, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("SCION: listened to %v\n", addr)
	return quicutil.SingleStreamListener{Listener: quicListener}, nil
}
