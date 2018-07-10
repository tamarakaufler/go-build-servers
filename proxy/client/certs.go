package main

import (
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
)

// getCertPool appends our cert to the system pool of trusted certificates
func getCertPool(certPath string) (*x509.CertPool, error) {
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	certs, err := ioutil.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to append %q to RootCAs: %v", err)
	}

	if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
		return nil, errors.New("Could not append a cert: using system certs only")
	}

	return rootCAs, nil
}
