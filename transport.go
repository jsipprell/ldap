// This package provide connection transport extended functionality,
// primary for the purposes of sasl gssapi mech security wrapping
package ldap

import (
	"net"
)

type Transporter interface {
	WrapConnection(net.Conn) (t Transport, ok bool)
}

type Transport interface {
	UnwrapConnection() net.Conn
	net.Conn
}

type TransporterFunc func(net.Conn) (Transport, bool)

func (tf TransporterFunc) WrapConnection(c net.Conn) (Transport, bool) {
	return tf(c)
}
