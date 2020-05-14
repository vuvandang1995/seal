package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/bitnami-labs/sealed-secrets/pkg/crypto"
	"go.uber.org/zap"
)

const (
	CertificateBlockType = "CERTIFICATE"
)

var (
	logger, _ = zap.NewProduction()
	sugar     = logger.Sugar()
)

func (s *Server) handleCheckHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	}
}

func (s *Server) handleEncrypt() http.HandlerFunc {

	parseLabel := func(r *http.Request) []byte {
		ns := r.URL.Query().Get("namespace")
		sec := r.URL.Query().Get("name")

		// consider as cluster-wide
		if len(ns) == 0 {
			return []byte{}
		}

		// consider as namespace
		if len(sec) == 0 {
			return []byte(fmt.Sprintf("%s", ns))
		}

		return []byte(fmt.Sprintf("%s/%s", ns, sec))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(base64.NewDecoder(base64.StdEncoding, r.Body))

		if err != nil {
			sugar.Errorf("%s", err)
			w.WriteHeader(500)
			return
		}

		sugar.Debugf("%s", data)

		f, err := os.Open(fmt.Sprintf("%s/%s.pub", s.config.PublicKeysDirectory, r.URL.Query().Get("env")))
		

		if err != nil {
			sugar.Errorf("%s", err)
			w.WriteHeader(500)
			return
		}

		pubKey, err := parsePublicKey(f)

		if err != nil {
			sugar.Errorf("%s", err)
			w.WriteHeader(500)
			fmt.Println(s.config.PublicKeysDirectory)
			return
		}

		label := parseLabel(r)

		out, err := crypto.HybridEncrypt(rand.Reader, pubKey, data, label)

		if err != nil {
			sugar.Errorf("%s", err)
			w.WriteHeader(500)
			return
		}

		sugar.Infof("payload '%s********' was encrypted %s", data[:len(data)/3], label)

		fmt.Fprint(w, base64.StdEncoding.EncodeToString(out))
	}
}

func parsePublicKey(r io.Reader) (*rsa.PublicKey, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	certs, err := parseCertsPEM(data)
	if err != nil {
		return nil, err
	}

	// ParseCertsPem returns error if len(certs) == 0, but best to be sure...
	if len(certs) == 0 {
		return nil, errors.New("Failed to read any certificates")
	}

	cert, ok := certs[0].PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("Expected RSA public key but found %v", certs[0].PublicKey)
	}

	return cert, nil
}

func parseCertsPEM(pemCerts []byte) ([]*x509.Certificate, error) {
	ok := false
	certs := []*x509.Certificate{}
	for len(pemCerts) > 0 {
		var block *pem.Block
		block, pemCerts = pem.Decode(pemCerts)
		if block == nil {
			break
		}
		// Only use PEM "CERTIFICATE" blocks without extra headers
		if block.Type != CertificateBlockType || len(block.Headers) != 0 {
			continue
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return certs, err
		}

		certs = append(certs, cert)
		ok = true
	}

	if !ok {
		return certs, errors.New("data does not contain any valid RSA or ECDSA certificates")
	}
	return certs, nil
}

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, s.config.HomeContent)
	}
}
