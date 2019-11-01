package client

import (
	"fmt"
	"test-log-emitter/config"
	"time"

	"code.cloudfoundry.org/go-loggregator"
	"google.golang.org/grpc"
)

const CertsDir = "certs/loggregator"
const CACertPath = CertsDir + "/ca_cert"
const CertPath = CertsDir + "/cert"
const KeyPath = CertsDir + "/key"

func NewLoggregatorIngressClient(config config.LoggregatorConfig) (*loggregator.IngressClient, error) {
	tlsConfig, err := loggregator.NewIngressTLSConfig(
		config.CA,
		config.Cert,
		config.Key,
	)
	if err != nil {
		return nil, err
	}

	opts := []loggregator.IngressOption{
		loggregator.WithAddr(fmt.Sprintf("127.0.0.1:%d", config.Port)),
	}

	opts = append(opts, loggregator.WithDialOptions(grpc.WithBlock(), grpc.WithTimeout(time.Second)))

	return loggregator.NewIngressClient(tlsConfig, opts...)
}
