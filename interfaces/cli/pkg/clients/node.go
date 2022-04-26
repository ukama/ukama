package clients

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/pkg/errors"
	"github.com/ukama/ukamaX/cli/pkg"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type NodeClient struct {
	log pkg.Logger
}

func NewNodeClient(log pkg.Logger) *NodeClient {
	return &NodeClient{log: log}
}

// sends file over tls to the ip with certificate. If client cert and key are missing then only server cert is checked
// host - ip or host of the server
// caCert -  CA cert that signed the server's certificate
// clientCertFile - the name the clients's certificate file (optioinal)
// clientKeyFile - the name the client's key certificate file (optioinal)
func (c *NodeClient) SendFile(host string, caCertFile string, clientCertFile string, clientKeyFile string, r io.Reader) error {

	if caCertFile == "" {
		return fmt.Errorf("caCert is required but missing")
	}

	var cert tls.Certificate
	var err error
	if clientCertFile != "" && clientKeyFile != "" {
		cert, err = tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return errors.Wrap(err, "error creating x509 keypair from client cert file")
		}
	}

	caCert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return errors.Wrap(err, "error opening cert file")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	t := &http.Transport{
		TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		},
	}

	client := http.Client{Transport: t, Timeout: 15 * time.Second}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://%s/config", host), r)
	if err != nil {
		return errors.Wrap(err, "unable to create http request")
	}

	resp, err := client.Do(req)
	if err != nil {
		switch e := err.(type) {
		case *url.Error:
			return errors.Wrap(e, "url.Error received on http request")
		default:
			return errors.Wrap(err, "unexpected error received")
		}
	}

	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return errors.Wrap(err, "unexpected error reading response body")
	}

	c.log.Printf("Response: %s", string(b))
	return nil
}
