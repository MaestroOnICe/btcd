package scion

import (
	"context"
	"crypto/tls"
	"errors"
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

// Wrapper for the function SplitHostPort
// The pan libary function is analogous to net.SplitHostPort, which however refuses to handle SCION addresses.
// The address can be of the form of a SCION address (i.e. of the form "ISD-AS,[IP]:port") or in the form of "hostname:port".
// We use this as a universal function and decide if pan or net is needed
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
		log.Debugf("Parsed: ", address, "into: ", saUDPAddress)
		return ScionAddr{}, errors.New("could not parse scion address")
	}
	sa := ScionAddr{
		udpAddr: saUDPAddress,
		addr:    address,
	}
	return sa, err
}

// naive way to check for a valid scion address
// TODO: regex for checking
func IsValidAddress(address string) (pan.UDPAddr, bool) {
	sa, err := pan.ParseUDPAddr(address)
	return sa, err == nil
}

// return a connection for a scion
// using quicutil from scion-apps
// SingleStream implements an opaque, bi-directional data stream using QUIC, intending to be a drop-in replacement for TCP
func Dial(address string) (net.Conn, error) {
	log.Debugf("SCION: dialing to %v\n", address)

	// Parse addr into a scionAddr
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}
	log.Debugf("SCION: parsed %v\n", addr)

	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"hello-quic"},
	}
	log.Debugf("SCION: tls set %v\n", tlsCfg)

	// create default selector TODO: create specific selector
	selector := pan.NewDefaultSelector()
	log.Debugf("SCION: selector set %v\n", selector)

	//dial
	session, err := pan.DialQUIC(context.Background(), netaddr.IPPort{}, addr, nil, selector, "", tlsCfg, nil)
	if err != nil {
		return nil, err
	}
	log.Debugf("SCION: dialed to %v\n", session)

	//return drop-in replacement stream
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
	log.Debugf("SCION: listening to %v\n", addr)
	tlsCfg := &tls.Config{
		Certificates: quicutil.MustGenerateSelfSignedCert(),
		NextProtos:   []string{"hello-quic"},
	}
	log.Debugf("SCION: tls set %v\n", addr)

	udpAddr, err := pan.ParseUDPAddr(addr)
	if err != nil {
		return nil, err
	}
	log.Debugf("SCION: parsed %v\n", addr)

	quicListener, err := pan.ListenQUIC(context.Background(), netaddr.IPPortFrom(udpAddr.IP, udpAddr.Port), nil, tlsCfg, nil)
	if err != nil {
		return nil, err
	}
	log.Debugf("SCION: listened to %v\n", addr)
	return quicutil.SingleStreamListener{Listener: quicListener}, nil
}
