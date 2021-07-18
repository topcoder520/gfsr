package httpclient

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"net/http"
	"sync"
)

//go:embed cert/ca.pem
var caPem string

//go:embed cert/client.key
var clientKey string

//go:embed cert/client.pem
var clientPem string

var httpClientOnce sync.Once
var httpTLSClientOnce sync.Once
var client *http.Client

func Once() *http.Client {
	httpClientOnce.Do(func() {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	})
	return client
}

func OnceTLS() *http.Client {
	httpTLSClientOnce.Do(func() {
		client = createClient()
	})
	return client
}

func createClient() *http.Client {
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(caPem))
	//客户端证书与密钥
	cliCrt, err := tls.X509KeyPair([]byte(clientPem), []byte(clientKey))
	if err != nil {
		return nil
	}
	cl := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			},
		},
	}
	return cl
}
