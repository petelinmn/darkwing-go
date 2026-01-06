package main

import (
    "log"
    "os"
    "time"
)

func main() {
    interval := 5 * time.Second
    if v := os.Getenv("INTERVAL_SECONDS"); v != "" {
        if d, err := time.ParseDuration(v + "s"); err == nil {
            interval = d
        }
    }

    log.Printf("worker starting; tick interval: %s", interval)
    ticker := time.NewTicker(interval)
    defer ticker.Stop()

    for i := 1; ; i++ {
        <-ticker.C
        log.Printf("worker tick %d", i)
    }
}
