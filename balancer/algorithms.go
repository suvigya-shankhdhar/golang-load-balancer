package balancer

import (
	"math"
	"sync/atomic"
)

// Balancer is an interface for different load balancing algorithms.
type Balancer interface {
	Next(servers []*ServerInfo) *ServerInfo
}

// RoundRobinBalancer implements Balancer interface with Round Robin algorithm.
type RoundRobinBalancer struct {
	current uint64
}

// Returns a new Round Robin Balancer struct.
func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{
		current: 0,
	}
}

func (rr *RoundRobinBalancer) Next(servers []*ServerInfo) *ServerInfo {
	var aliveServers []*ServerInfo
	for _, s := range servers {
		if s.IsAlive() {
			aliveServers = append(aliveServers, s)
		}
	}
	if len(aliveServers) == 0 {
		return nil
	}

	nextIndex := atomic.AddUint64(&rr.current, 1)
	return aliveServers[nextIndex%uint64(len(aliveServers))]
}

// Implements the Balancer interface for Least Connections Algorithm.
type LeastConnectionsBalancer struct{}

func NewLeastConnectionsBalancer() *LeastConnectionsBalancer {
	return &LeastConnectionsBalancer{}
}

func (l *LeastConnectionsBalancer) Next(servers []*ServerInfo) *ServerInfo {
	var bestServer *ServerInfo
	minConnections := int64(math.MaxInt64)

	for _, s := range servers {
		if s.IsAlive() {
			count := s.Connections
			if count < minConnections {
				minConnections = count
				bestServer = s
			}
		}
	}
	return bestServer
}
