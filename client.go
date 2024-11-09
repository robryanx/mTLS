package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	name := flag.String("c", "a", "client name")
	flag.Parse()

	cert, err := os.ReadFile("./certs/ca.crt")
	if err != nil {
		log.Fatalf("could not open certificate file: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	clientCert := fmt.Sprintf("./certs/client.%s.crt", *name)
	clientKey := fmt.Sprintf("./certs/client.%s.key", *name)
	log.Println("Load key pairs - ", clientCert, clientKey)
	certificate, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("could not load certificate: %v", err)
	}

	client := http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS13,
				RootCAs:            caCertPool,
				Certificates:       []tls.Certificate{certificate},
				InsecureSkipVerify: false,
			},
		},
	}

	// Request /hello over port 8443 via the GET method
	// Using curl the verfiy it :
	// curl --trace trace.log -k \
	//   --cacert ./certs/ca.crt  --cert ./certs/client.b.crt --key ./certs/client.b.key  \
	//     https://localhost:8443/hello

	r, err := client.Get("https://localhost:8443/hello")
	if err != nil {
		log.Fatalf("error making get request: %v", err)
	}

	// Read the response body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalf("error reading response: %v", err)
	}

	// Print the response body to stdout
	fmt.Printf("%s\n", body)
}
