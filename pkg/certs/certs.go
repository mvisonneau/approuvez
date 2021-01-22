package certs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

// Config ..
type Config struct {
	Disable bool
	CA      string
	Cert    string
	Key     string
}

// CertificateType ..
type CertificateType uint8

const (
	// CertificateTypeServer ..
	CertificateTypeServer CertificateType = iota

	// CertificateTypeClient ..
	CertificateTypeClient
)

// LoadClientCertificates ..
func LoadClientCertificates(cfg Config, serverName string) (credentials.TransportCredentials, error) {
	return LoadCertificates(CertificateTypeClient, cfg, serverName)
}

// LoadServerCertificates ..
func LoadServerCertificates(cfg Config) (credentials.TransportCredentials, error) {
	return LoadCertificates(CertificateTypeServer, cfg, "")
}

// LoadCertificates ..
func LoadCertificates(
	t CertificateType,
	cfg Config,
	serverName string,
) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	bs, err := ioutil.ReadFile(cfg.CA)
	if err != nil {
		return nil, err
	}

	ok := certPool.AppendCertsFromPEM(bs)
	if !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	keyPair, err := tls.LoadX509KeyPair(cfg.Cert, cfg.Key)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{keyPair},
	}

	switch t {
	case CertificateTypeServer:
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsConfig.ClientCAs = certPool
	case CertificateTypeClient:
		tlsConfig.ServerName = serverName
		tlsConfig.RootCAs = certPool
	default:
		return nil, fmt.Errorf("invalid certificate type: %v", t)
	}

	return credentials.NewTLS(tlsConfig), nil
}
