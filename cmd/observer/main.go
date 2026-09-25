package main

import (
	"fmt"
	"log"
	"net/http"
	"flag"
	"github.com/Darkpenguin1/k8-observer/internal/observer"
)

func main() {
	
	inCluster := flag.Bool("in-cluster", false, "Run inside a Kubernetes cluster")
	flag.Parse()

	newObserver, err := observer.New(*inCluster)
	if err != nil {
		log.Fatalf("Failed to create observer: %v", err)
	}
	fmt.Println(newObserver)
	
	
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world from observer!")
	})
	port := ":8080"
	log.Fatal(http.ListenAndServe(port, nil))
}