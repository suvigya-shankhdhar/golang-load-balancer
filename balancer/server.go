package balancer

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
)

// ServerInfo holds the data about a specific server 
// the proxy forwards to. 
type ServerInfo struct {
	URL					*url.URL
	Alive				bool	
	ReverseProxy		*httputil.ReverseProxy
	mu 					sync.RWMutex
	Connections			int64
}

// SetAlive updates the status of a Backend Server to alive. 
func (s *ServerInfo) SetAlive(alive bool) {
	s.mu.Lock()
	defer s.mu.Unlock() 
	s.Alive = alive
}

// IsAlive returns the current status of a Backend Server.
func (s *ServerInfo) IsAlive() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Alive
}

// NewBackendServer creates a new Backend Server instance 
// and initializes its ReverseProxy
func NewBackendServer(serverURL *url.URL) *ServerInfo {
	return &ServerInfo {
		URL: serverURL, 
		Alive: true,
		ReverseProxy: httputil.NewSingleHostReverseProxy(serverURL), 
	}
}

// ServerHTTP 
func (s *ServerInfo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&s.Connections, 1)
	defer atomic.AddInt64(&s.Connections, -1)
	s.ReverseProxy.ServeHTTP(w, r)
}