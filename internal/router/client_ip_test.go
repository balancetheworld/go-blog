package router

import (
	"net"
	"testing"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/test/mock"
)

type clientIPTestConn struct {
	*mock.Conn
	remoteAddr net.Addr
}

func (c *clientIPTestConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func newClientIPTestContext(remoteIP string, forwardedIP string) *app.RequestContext {
	ctx := app.NewContext(0)
	ctx.SetConn(&clientIPTestConn{
		Conn: mock.NewConn(""),
		remoteAddr: &net.TCPAddr{
			IP:   net.ParseIP(remoteIP),
			Port: 8080,
		},
	})
	ctx.Request.Header.Set("X-Forwarded-For", forwardedIP)
	ctx.Request.Header.Set("X-Real-IP", forwardedIP)
	return ctx
}

func TestClientIPIgnoresForwardedHeadersByDefault(t *testing.T) {
	clientIP, err := buildClientIPFunc("")
	if err != nil {
		t.Fatal(err)
	}

	ctx := newClientIPTestContext("203.0.113.10", "198.51.100.20")
	if value := clientIP(ctx); value != "203.0.113.10" {
		t.Fatalf("expected remote IP, got %q", value)
	}
}

func TestClientIPUsesHeadersOnlyFromTrustedProxy(t *testing.T) {
	clientIP, err := buildClientIPFunc("127.0.0.1/32, ::1/128")
	if err != nil {
		t.Fatal(err)
	}

	trusted := newClientIPTestContext("127.0.0.1", "198.51.100.20")
	if value := clientIP(trusted); value != "198.51.100.20" {
		t.Fatalf("expected forwarded IP, got %q", value)
	}

	untrusted := newClientIPTestContext("203.0.113.10", "198.51.100.20")
	if value := clientIP(untrusted); value != "203.0.113.10" {
		t.Fatalf("expected remote IP, got %q", value)
	}
}

func TestClientIPRejectsInvalidTrustedProxyCIDR(t *testing.T) {
	if _, err := buildClientIPFunc("not-a-cidr"); err == nil {
		t.Fatal("expected invalid CIDR error")
	}
}
