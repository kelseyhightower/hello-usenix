package main

import (
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/braintree/manners"
)

func main() {
    log.Println("Starting up...")

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Printf("%s %s - %s", r.Method, r.URL.Path, r.UserAgent())
        w.Write([]byte("Hello USENIX!\n"))
    })

    // Catch signals in the background using an anonymous go routine.
    go func() {
        signalChan := make(chan os.Signal, 1)
        signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

        // This blocks until this process receives a SIGINT or SIGTERM signal.
        <-signalChan
        log.Println("Shutting down...")

        // stop accepting new requests and begin shutting down.
        manners.Close()
    }()

    // The call to ListenAndServe() blocks until an error is returned
    // or the manners.Close() function is called.
    err := manners.ListenAndServe(":8080", http.DefaultServeMux)
    if err != nil {
       log.Fatal(err)
    }
    log.Println("Shutdown complete.")
}
