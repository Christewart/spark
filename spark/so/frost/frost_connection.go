package frost

import (
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lightsparkdev/spark/common"
	sparkgrpc "github.com/lightsparkdev/spark/common/grpc"
)

type FrostGRPCConnectionFactory interface {
	NewFrostGRPCConnection(signerAddress string) (*grpc.ClientConn, error)
	SetTimeoutProvider(timeoutProvider sparkgrpc.TimeoutProvider)
}

type frostGRPCConnectionFactorySecure struct {
	ClientTimeoutConfig *common.ClientTimeoutConfig
}

func NewFrostGRPCConnectionFactorySecure() *frostGRPCConnectionFactorySecure {
	return &frostGRPCConnectionFactorySecure{
		ClientTimeoutConfig: nil,
	}
}

func (f *frostGRPCConnectionFactorySecure) SetTimeoutProvider(timeoutProvider sparkgrpc.TimeoutProvider) {
	f.ClientTimeoutConfig = &common.ClientTimeoutConfig{
		TimeoutProvider: timeoutProvider,
	}
}

func (f *frostGRPCConnectionFactorySecure) NewFrostGRPCConnection(signerAddress string) (*grpc.ClientConn, error) {
	if shouldDialFrostOverTCP(signerAddress) {
		clientOpts := common.BasicClientOptions(signerAddress, nil, f.ClientTimeoutConfig)
		clientOpts = append(clientOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		return grpc.NewClient(signerAddress, clientOpts...)
	}

	return common.NewGRPCConnectionUnixDomainSocket(signerAddress, nil, f.ClientTimeoutConfig)
}

func shouldDialFrostOverTCP(signerAddress string) bool {
	if strings.HasPrefix(signerAddress, "unix:///") || strings.HasPrefix(signerAddress, "unix:/") {
		return false
	}
	_, _, err := net.SplitHostPort(signerAddress)
	return err == nil
}
