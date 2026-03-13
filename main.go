package main

import (
    "context"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/signal"
    "syscall"
    "time"
    "my-load-balancer/balancer"
    "my-load-balancer/config"
    "my-load-balancer/healthcheck"
    "my-load-balancer/proxy"
)

func main() {
    // Load Configuration
    cfg, err := config.LoadConfig(`C:\Users\shank\OneDrive\Documents\go programming\reverse_proxy\config.json`)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // Initialize the Algorithm and Pool
    algo := balancer.NewLeastConnectionsBalancer()
    pool := balancer.NewServerPool(algo)

    // Register Backends from Config
    for _, sURL := range cfg.Backends {
        serverURL, err := url.Parse(sURL)
        if err != nil {
            log.Fatalf("Invalid backend URL %s: %v", sURL, err)
        }

        backend := balancer.NewBackendServer(serverURL)
        pool.AddServer(backend)
        log.Printf("Registered backend: %s", sURL)
    }

    // Start Health Checker in the background
    healthcheck.RunHealthCheck(pool.Servers, 10 * time.Second)

    // Setup the Reverse Proxy
    lbProxy := proxy.NewReverseProxy(pool)
    loggedProxy := proxy.Logger(lbProxy)

    server := &http.Server{
        Addr:   cfg.ListenPort, 
        Handler: loggedProxy, 
    }

    // Listens for a "stop" signal
    go func() {
        log.Printf("Load Balncer started on %s", cfg.ListenPort)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Listen: %s\n", err)
        }
    }()

    // Wait for interrupt signal to shutdown
    quit := make(chan os.Signal, 1) 
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Println("Shutting down Load Balancer....")

    // Give the server 5 seconds to finish current requests
    ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forces to shutdown: ", err)
    }
    log.Println("Load Balancer exiting")

}