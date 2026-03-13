package balancer

import (
	"fmt"
	"net/http"
)

// ServerPool holds all backend Servers and the load balancing strategy.
type ServerPool struct {
	Servers []*ServerInfo
	Algorithm Balancer
}


// Creates a new pool of Servers with a chosen load balancing strategy.
func NewServerPool(algo Balancer) *ServerPool {
	return &ServerPool {
		Servers: make([]*ServerInfo, 0), 
		Algorithm: algo,
	}
}

// Adds a new Server to an existing Server pool.
func (sp *ServerPool) AddServer(s *ServerInfo) {
	sp.Servers = append(sp.Servers, s)
}

// GetNextServer returns the next avaible server 
// according to the chosen algorithm. 
func (sp *ServerPool) GetNextServer() *ServerInfo {
	return sp.Algorithm.Next(sp.Servers)
}

// Selects the next server from Server Pool and 
// calls the ServeHTTP() method on it. 
func (sp *ServerPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	peer := sp.GetNextServer() 
	if peer != nil {
		peer.ServeHTTP(w, r)
		return 
	}

	http.Error(w, "Service not available", http.StatusServiceUnavailable)
}

// Returns a summary of the server pool for loggin or monitoring.
func (sp *ServerPool) GetStatus() {
	for _, s := range sp.Servers {
		fmt.Printf("Backends: %s | Alive: %v | Active Conns: %d\n", s.URL, s.IsAlive(), s.Connections)
	}
}
