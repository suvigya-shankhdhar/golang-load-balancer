# Go Load Balancer
A Layer 7 Load Balancer built in Go.

## Features
* **Least Connections Algorithm:** Intelligently routes traffic to the server with the lowest active load.
* **Health Checks:** Automatic TCP-based background health monitoring.
* **Reverse Proxy:** Built using Go's `httputil.ReverseProxy` for high performance.
* **Middleware Logging:** Real-time observability of incoming requests and latency.
* **Graceful Shutdown:** Ensures no connections are dropped during server restarts.

## How to Run
1. Clone the repo.
2. Configure `config.json` with your backend URLs.
3. Run `go run main.go`.
