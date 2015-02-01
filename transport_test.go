package ldap

import (
	"fmt"
	"net"
	"testing"
	"time"
)

var (
	tw_base_dn string   = "dc=lab,ou=nodes,ou=hdb"
	tw_filter  []string = []string{
		"(puppetclass=*)"}
	tw_attributes []string = []string{
		"cn",
		"description"}
	test *testing.T
)

type testTransport struct {
	T    *testing.T
	conn net.Conn
}

func (tt *testTransport) UnwrapConnection() net.Conn {
	return tt.conn
}

func (tt *testTransport) Read(b []byte) (n int, err error) {
	if n, err = tt.conn.Read(b); err == nil {
		tt.T.Logf("%v read %v octets from %v\n", tt.LocalAddr(), n, tt.RemoteAddr())
	}
	return
}

func (tt *testTransport) Write(b []byte) (n int, err error) {
	if n, err = tt.conn.Write(b); err == nil {
		tt.T.Logf("%v wrote %v octets to %v\n", tt.LocalAddr(), n, tt.RemoteAddr())
	}
	return
}

func (tt *testTransport) Close() error {
	return tt.conn.Close()
}

func (tt *testTransport) LocalAddr() net.Addr {
	return tt.conn.LocalAddr()
}

func (tt *testTransport) RemoteAddr() net.Addr {
	return tt.conn.RemoteAddr()
}

func (tt *testTransport) SetDeadline(tm time.Time) error {
	return tt.conn.SetDeadline(tm)
}

func (tt *testTransport) SetReadDeadline(tm time.Time) error {
	return tt.conn.SetReadDeadline(tm)
}

func (tt *testTransport) SetWriteDeadline(tm time.Time) error {
	return tt.conn.SetWriteDeadline(tm)
}

func makeWrapper(c net.Conn) (tr Transport, ok bool) {
	if c == nil {
		ok = false
	} else {
		tr = &testTransport{test, c}
		ok = true
	}
	return
}

func TestWrapper(t *testing.T) {
	test = t
	defer func() { test = nil }()

	fmt.Printf("TestWrapper: starting...\n")
	l := NewLDAPConnection(ldap_server, ldap_port)
	l.TransportWrapper = TransporterFunc(makeWrapper)
	if err := l.Connect(); err != nil {
		t.Error(err)
		return
	}
	defer l.Close()

	req := NewSearchRequest(tw_base_dn, ScopeWholeSubtree, DerefAlways, 0, 0, false, tw_filter[0], tw_attributes, nil)
	sr, err := l.Search(req)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("TestWrapper: %v -> num of entries = %d", req.Filter, len(sr.Entries))
}
