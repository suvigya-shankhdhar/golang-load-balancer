// proxy/proxy.go
package proxy

import (
    "net/http"
    "my-load-balancer/balancer"
    "log"
    "time"
)

// ReverseProxy handles incoming requests and delegates 
// them to the balancer.
type ReverseProxy struct {
    Pool *balancer.ServerPool
}

// Creates a new Reverse proxy instance with a given server pool.
func NewReverseProxy(sp *balancer.ServerPool) *ReverseProxy {
    return &ReverseProxy{
        Pool: sp, 
    }
}

// This is a middleware function
func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
        start := time.Now()
        next.ServeHTTP(w, r)
        duration := time.Since(start)
        log.Printf("%-7s | %-20s | %10s", r.Method, r.URL.Path, duration, )
    })
}

// Main entry point for every request hitting the load balancer.
func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    rp.Pool.ServeHTTP(w, r)
}