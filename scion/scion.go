package scion

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"time"

	"github.com/netsec-ethz/scion-apps/pkg/pan"
	"github.com/netsec-ethz/scion-apps/pkg/quicutil"
	"inet.af/netaddr"
)

// ScionAddr implements the net.Addr interface
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

// Wrapper for the function JoinHostPort
// Preserves the functionalty of the net.JoinHostPort function
// We decided, that scion runs on port 8666, if not otherwise specified
func JoinHostPort(host, port string) string {
	if _, ok := IsValidAddress(host); ok {
		return host + ":8666"
	}
	return net.JoinHostPort(host, port)
}

// parse scion string into an ScionAddr with udp addresse and string address
func ParseAddr(address string) (ScionAddr, error) {
	saUDPAddress, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return ScionAddr{}, errors.New("could not parse scion address")
	}
	sa := ScionAddr{
		udpAddr: saUDPAddress,
		addr:    address,
	}
	return sa, err
}

// naive way to check for a valid scion address
// TODO: regex for checking, regex seems to thow erros, so for now we use this, which uses regex under the hood
func IsValidAddress(address string) (pan.UDPAddr, bool) {
	sa, err := pan.ParseUDPAddr(address)
	return sa, err == nil
}

// return a connection for a scion
// using quicutil from scion-apps
// SingleStream implements an opaque, bi-directional data stream using QUIC, intending to be a drop-in replacement for TCP
// We use SingleStream because it is a working drop-in which requires no futher code changes
func Dial(address string) (net.Conn, error) {

	// parse addr into a scionAddr
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}

	// default tls, as in examples
	tlsCfg := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"hello-quic"},
	}

	// create pinging selector, which pings two paths
	selector := &pan.PingingSelector{
		// Interval for pinging. Must be positive.
		Interval: 500 * time.Millisecond,
		// Timeout for the individual pings. Must be positive and less than Interval.
		Timeout: 400 * time.Millisecond,
	}
	// enables active pinging on at most numActive paths.
	selector.SetActive(2)

	//dial
	session, err := pan.DialQUIC(context.Background(), netaddr.IPPort{}, addr, nil, selector, "", tlsCfg, nil)
	if err != nil {
		return nil, err
	}

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
func Listen(address string) (net.Listener, error) {

	//check for valid address
	addr, err := pan.ResolveUDPAddr(address)
	if err != nil {
		return nil, err
	}

	// default tls, as in examples
	tlsCfg := &tls.Config{
		Certificates: quicutil.MustGenerateSelfSignedCert(),
		NextProtos:   []string{"hello-quic"},
	}

	// listen
	quicListener, err := pan.ListenQUIC(context.Background(), netaddr.IPPortFrom(addr.IP, addr.Port), nil, tlsCfg, nil)
	if err != nil {
		return nil, err
	}
	return quicutil.SingleStreamListener{Listener: quicListener}, nil
}
