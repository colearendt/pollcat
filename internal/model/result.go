package model

import "time"

// PollType distinguishes between TCP and DNS probes.
type PollType string

const (
	PollTypeTCP  PollType = "tcp"
	PollTypeUDP  PollType = "udp"
	PollTypeICMP PollType = "icmp"
	PollTypeDNS  PollType = "dns"
)

// Target describes a single thing to poll.
type Target struct {
	Type     PollType
	Address  string // host:port for TCP/UDP, hostname for DNS/ICMP
	Interval time.Duration
}

// Result is the outcome of a single poll.
type Result struct {
	Timestamp time.Time
	Type      PollType
	Target    string
	Success   bool
	Latency   time.Duration
	Response  string // e.g. resolved IP or TCP dial result
	Error     string
}
