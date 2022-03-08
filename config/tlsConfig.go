package config

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	"google.golang.org/grpc/credentials"
)

func CreateTlsConfig() credentials.TransportCredentials {
	caPem, err := ioutil.ReadFile("tls/ca-cert.pem")
	if err != nil {
		log.Fatalf("error while loading TLS ca cert %v", err)
	}

	certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(caPem) {
		log.Fatal(err)
	}

	serverCert, err := tls.LoadX509KeyPair("tls/server-cert.pem", "tls/server-key.pem")
	if err != nil {
		log.Fatalf("error while loading server cert %v", err)
	}

	conf := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs: certPool,
	}
	tlsCredentials := credentials.NewTLS(conf)
	return tlsCredentials
}