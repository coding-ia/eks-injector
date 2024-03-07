package main

import (
	"eks-inject/internal/mutate"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("Starting server...")

	mux := http.NewServeMux()

	mux.HandleFunc("/inject", handleWebhook)

	s := &http.Server{
		Addr:           ":8443",
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1048576
	}

	log.Fatal(s.ListenAndServeTLS("/etc/ssl/tls.crt", "/etc/ssl/tls.key"))
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
	}

	log.Println(string(body))

	clusterName := os.Getenv("EKS_CLUSTER_NAME")
	responseBody, err := mutate.ProcessAdmissionReview(body, clusterName)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}