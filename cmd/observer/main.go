package main

import (
	"log"

	"flag"
	"github.com/Darkpenguin1/k8-observer/internal/observer"
	"github.com/Darkpenguin1/k8-observer/internal/server"

	"context"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	inCluster := flag.Bool("in-cluster", false, "Run inside a Kubernetes cluster")
	flag.Parse()

	o, err := observer.New(*inCluster)
	if err != nil {
		log.Printf("create observer: %v", err)
		return
	}

	type result struct {
		name string
		err  error
	}
	results := make(chan result, 2)

	go func() { results <- result{"observer", o.Run(ctx)} }()
	go func() { results <- result{"server", server.Run(ctx)} }()

	completed := 0
	select {
	case <-ctx.Done(): // Ctrl+C or SIGTERM
	case r := <-results: // one component stopped first
		completed++
		if r.err != nil {
			log.Printf("%s: %v", r.name, r.err)
		}
	}

	stop() // tell the other component to stop

	for completed < 2 {
		r := <-results
		completed++
		if r.err != nil {
			log.Printf("%s: %v", r.name, r.err)
		}
	}
}
