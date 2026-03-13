package healthcheck

import (
	"log"
	"net"
	// "net/http"
	"net/url"
	"my-load-balancer/balancer"
	"time"
)

// Pings a server to check if it's reachable.
func isServerAlive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Printf("Site uncreachable, error: %v", err)
		return false
	}
	conn.Close()
	return true
}

// RunHealthCheck starts a 
func RunHealthCheck(servers []*balancer.ServerInfo, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C{
			log.Println("Starting a health check...")
			for _, s := range servers {
				status := isServerAlive(s.URL)
				s.SetAlive(status)
			}
		}
	} ()
}

