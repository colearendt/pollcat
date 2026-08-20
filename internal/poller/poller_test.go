package poller

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/colearendt/pollcat/internal/model"
	"github.com/stretchr/testify/assert"
)

type fakeDialer struct {
	conn net.Conn
	err  error
}

func (f *fakeDialer) DialContext(_ context.Context, _, _ string) (net.Conn, error) {
	return f.conn, f.err
}

type fakeConn struct{}

func (fakeConn) Read(_ []byte) (n int, err error)   { return 0, nil }
func (fakeConn) Write(_ []byte) (n int, err error)  { return 0, nil }
func (fakeConn) Close() error                       { return nil }
func (fakeConn) LocalAddr() net.Addr                { return nil }
func (fakeConn) RemoteAddr() net.Addr               { return nil }
func (fakeConn) SetDeadline(_ time.Time) error      { return nil }
func (fakeConn) SetReadDeadline(_ time.Time) error  { return nil }
func (fakeConn) SetWriteDeadline(_ time.Time) error { return nil }

type fakeResolver struct {
	ips []string
	err error
}

func (f *fakeResolver) LookupHost(_ context.Context, _ string) ([]string, error) {
	return f.ips, f.err
}

type fakePinger struct {
	latency  time.Duration
	response string
	err      error
}

func (f *fakePinger) Ping(_ context.Context, _ string) (time.Duration, string, error) {
	return f.latency, f.response, f.err
}

func TestTCPPoller_Poll_Success(t *testing.T) {
	t.Parallel()
	p := &TCPPoller{Dialer: &fakeDialer{conn: fakeConn{}}}
	res := p.Poll(context.Background(), model.Target{Address: "1.2.3.4:80"})

	assert.True(t, res.Success)
	assert.Equal(t, model.PollTypeTCP, res.Type)
	assert.Equal(t, "connected", res.Response)
	assert.Zero(t, res.Error)
	assert.GreaterOrEqual(t, res.Latency, time.Duration(0))
}

func TestTCPPoller_Poll_Failure(t *testing.T) {
	t.Parallel()
	p := &TCPPoller{Dialer: &fakeDialer{err: errors.New("connection refused")}}
	res := p.Poll(context.Background(), model.Target{Address: "1.2.3.4:80"})

	assert.False(t, res.Success)
	assert.Equal(t, "connection refused", res.Error)
}

func TestUDPPoller_Poll_Success(t *testing.T) {
	t.Parallel()
	p := &UDPPoller{Dialer: &fakeDialer{conn: fakeConn{}}}
	res := p.Poll(context.Background(), model.Target{Address: "1.2.3.4:53"})

	assert.True(t, res.Success)
	assert.Equal(t, model.PollTypeUDP, res.Type)
	assert.Equal(t, "socket ready", res.Response)
	assert.Zero(t, res.Error)
	assert.GreaterOrEqual(t, res.Latency, time.Duration(0))
}

func TestUDPPoller_Poll_Failure(t *testing.T) {
	t.Parallel()
	p := &UDPPoller{Dialer: &fakeDialer{err: errors.New("network unreachable")}}
	res := p.Poll(context.Background(), model.Target{Address: "1.2.3.4:53"})

	assert.False(t, res.Success)
	assert.Equal(t, "network unreachable", res.Error)
}

func TestDNSPoller_Poll_Success(t *testing.T) {
	t.Parallel()
	p := &DNSPoller{Resolver: &fakeResolver{ips: []string{"93.184.216.34"}}}
	res := p.Poll(context.Background(), model.Target{Address: "example.com"})

	assert.True(t, res.Success)
	assert.Equal(t, model.PollTypeDNS, res.Type)
	assert.Contains(t, res.Response, "93.184.216.34")
	assert.Zero(t, res.Error)
}

func TestDNSPoller_Poll_Failure(t *testing.T) {
	t.Parallel()
	p := &DNSPoller{Resolver: &fakeResolver{err: errors.New("NXDOMAIN")}}
	res := p.Poll(context.Background(), model.Target{Address: "bad.example.com"})

	assert.False(t, res.Success)
	assert.Equal(t, "NXDOMAIN", res.Error)
}

func TestICMPPoller_Poll_Success(t *testing.T) {
	t.Parallel()
	p := &ICMPPoller{Pinger: &fakePinger{latency: 10 * time.Millisecond, response: "reply from 1.2.3.4"}}
	res := p.Poll(context.Background(), model.Target{Address: "1.2.3.4"})

	assert.True(t, res.Success)
	assert.Equal(t, model.PollTypeICMP, res.Type)
	assert.Equal(t, "reply from 1.2.3.4", res.Response)
	assert.Zero(t, res.Error)
	assert.Equal(t, 10*time.Millisecond, res.Latency)
}

func TestICMPPoller_Poll_Failure(t *testing.T) {
	t.Parallel()
	p := &ICMPPoller{Pinger: &fakePinger{err: errors.New("permission denied")}}
	res := p.Poll(context.Background(), model.Target{Address: "1.2.3.4"})

	assert.False(t, res.Success)
	assert.Equal(t, "permission denied", res.Error)
}
