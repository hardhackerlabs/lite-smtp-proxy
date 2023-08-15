package main

import (
	"crypto/tls"
	"log"
	"os"
	"time"

	"github.com/emersion/go-smtp"
)

func main() {
	// Create a smtp sever instance
	b := &Backend{}
	if err := b.init(); err != nil {
		log.Fatalln(err)
		return
	}
	srv := smtp.NewServer(b)

	// Setup TLS if needed
	certFile := os.Getenv("SMTP_PROXY_CERT")
	keyFile := os.Getenv("SMTP_PROXY_KEY")

	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalln(err)
			return
		}
		srv.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	// Setup Port
	port := os.Getenv("SMTP_PROXY_PORT")
	if port == "" {
		if srv.TLSConfig != nil {
			port = "587"
		} else {
			port = "25"
		}
	}

	srv.Addr = ":" + port
	srv.Domain = ""
	srv.ReadTimeout = 10 * time.Second
	srv.WriteTimeout = 10 * time.Second
	srv.MaxMessageBytes = 1024 * 1024
	srv.MaxLineLength = 4000
	srv.MaxRecipients = 50
	srv.AllowInsecureAuth = true

	log.Printf("Start to listen on \"%s\"", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
