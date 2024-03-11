package server

import (
	"eks-inject/internal/discovery"
	"eks-inject/internal/mutate"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

var globalVariables map[string]string

func StartServer() {
	globalVariables = BuildConfig()

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

func BuildConfig() map[string]string {
	variables := make(map[string]string)

	var clusterName string
	var clusterEnvironment string
	var clusterVersion string
	var err error

	clusterName = os.Getenv("CLUSTER_NAME")
	if clusterName == "" {
		clusterName, err = discovery.DiscoverClusterName()
		if err != nil {
			log.Println("Unable to discover cluster name.")
			clusterName = "undefined"
		}
	}
	variables["ClusterName"] = clusterName

	clusterVersion = os.Getenv("CLUSTER_VERSION")
	if clusterVersion == "" {
		clusterVersion, err = discovery.DiscoverKubernetesVersion()
		if err != nil {
			log.Println("Unable to discover cluster version.")
			clusterVersion = "undefined"
		}
	}
	variables["Version"] = clusterVersion

	clusterEnvironment = os.Getenv("CLUSTER_ENVIRONMENT")
	if clusterEnvironment == "" {
		clusterEnvironment, err = discovery.DiscoverEnvironment()
		if err != nil {
			log.Println("Unable to discover cluster environment.")
			clusterEnvironment = "undefined"
		}
	}
	variables["Environment"] = clusterEnvironment

	for k, v := range variables {
		log.Printf("Global variable [%s] = '%s'", k, v)
	}

	return variables
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

	responseBody, err := mutate.ProcessAdmissionReview(body, globalVariables)
	if err != nil {
		log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseBody)
}
