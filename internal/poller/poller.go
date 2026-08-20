package poller

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"github.com/colearendt/pollcat/internal/model"
)

// Dialer abstracts net.Dial for testability.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// Resolver abstracts DNS resolution for testability.
type Resolver interface {
	LookupHost(ctx context.Context, host string) ([]string, error)
}

// Pinger abstracts ICMP ping for testability.
type Pinger interface {
	Ping(ctx context.Context, address string) (latency time.Duration, response string, err error)
}

// Poller performs a single probe and returns a Result.
type Poller interface {
	Poll(ctx context.Context, target model.Target) model.Result
}

// TCPPoller probes TCP connectivity.
type TCPPoller struct {
	Dialer Dialer
}

func (p *TCPPoller) Poll(ctx context.Context, target model.Target) model.Result {
	start := time.Now()
	conn, err := p.Dialer.DialContext(ctx, "tcp", target.Address)
	latency := time.Since(start)

	res := model.Result{
		Timestamp: start,
		Type:      model.PollTypeTCP,
		Target:    target.Address,
		Latency:   latency,
	}

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return res
	}
	defer conn.Close()

	res.Success = true
	res.Response = "connected"
	return res
}

// UDPPoller probes UDP socket creation (connectionless, but resolves address and creates socket).
type UDPPoller struct {
	Dialer Dialer
}

func (p *UDPPoller) Poll(ctx context.Context, target model.Target) model.Result {
	start := time.Now()
	conn, err := p.Dialer.DialContext(ctx, "udp", target.Address)
	latency := time.Since(start)

	res := model.Result{
		Timestamp: start,
		Type:      model.PollTypeUDP,
		Target:    target.Address,
		Latency:   latency,
	}

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return res
	}
	defer conn.Close()

	res.Success = true
	res.Response = "socket ready"
	return res
}

// DNSPoller probes DNS resolution.
type DNSPoller struct {
	Resolver Resolver
	Timeout  time.Duration
}

func (p *DNSPoller) Poll(ctx context.Context, target model.Target) model.Result {
	start := time.Now()
	if p.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.Timeout)
		defer cancel()
	}
	ips, err := p.Resolver.LookupHost(ctx, target.Address)
	latency := time.Since(start)

	res := model.Result{
		Timestamp: start,
		Type:      model.PollTypeDNS,
		Target:    target.Address,
		Latency:   latency,
	}

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return res
	}

	res.Success = true
	res.Response = fmt.Sprintf("%v", ips)
	return res
}

// ICMPPoller probes ICMP echo (ping).
type ICMPPoller struct {
	Pinger Pinger
}

func (p *ICMPPoller) Poll(ctx context.Context, target model.Target) model.Result {
	start := time.Now()
	latency, response, err := p.Pinger.Ping(ctx, target.Address)

	res := model.Result{
		Timestamp: start,
		Type:      model.PollTypeICMP,
		Target:    target.Address,
		Latency:   latency,
	}

	if err != nil {
		res.Success = false
		res.Error = err.Error()
		return res
	}

	res.Success = true
	res.Response = response
	return res
}

// DefaultPinger is a Pinger implementation using golang.org/x/net/icmp.
// Requires root/admin privileges on most systems.
type DefaultPinger struct {
	Timeout time.Duration
}

func (p *DefaultPinger) Ping(ctx context.Context, address string) (time.Duration, string, error) {
	start := time.Now()

	ipAddr, err := net.ResolveIPAddr("ip4", address)
	if err != nil {
		return 0, "", fmt.Errorf("resolve: %w", err)
	}

	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, "", fmt.Errorf("icmp listen: %w (may require root/admin privileges)", err)
	}
	defer c.Close()

	msg := &icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("pollcat"),
		},
	}
	data, err := msg.Marshal(nil)
	if err != nil {
		return 0, "", fmt.Errorf("marshal icmp: %w", err)
	}

	_, err = c.WriteTo(data, ipAddr)
	if err != nil {
		return 0, "", fmt.Errorf("send icmp: %w", err)
	}

	reply := make([]byte, 1500)
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(p.Timeout)
		if p.Timeout == 0 {
			deadline = time.Now().Add(5 * time.Second)
		}
	}
	c.SetReadDeadline(deadline)

	n, peer, err := c.ReadFrom(reply)
	if err != nil {
		return 0, "", fmt.Errorf("read reply: %w", err)
	}

	latency := time.Since(start)

	rm, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return latency, "", fmt.Errorf("parse reply: %w", err)
	}

	if rm.Type != ipv4.ICMPTypeEchoReply {
		return latency, "", fmt.Errorf("unexpected icmp type: %v", rm.Type)
	}

	return latency, fmt.Sprintf("reply from %v", peer), nil
}
