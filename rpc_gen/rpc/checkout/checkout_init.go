package checkout

import (
	"sync"

	"github.com/cloudwego/kitex/client"
)

var (
	// todo edit custom config
	defaultClient     RPCClient
	defaultDstService = "checkout"
	defaultClientOpts []client.Option
	once              sync.Once
)

func init() {
	DefaultClient()
}

func DefaultClient() RPCClient {
	once.Do(func() {
		defaultClientOpts = append(defaultClientOpts, client.WithHostPorts("10.1.2.133:8885"))
		defaultClient = newClient(defaultDstService, defaultClientOpts...)
	})
	return defaultClient
}

func newClient(dstService string, opts ...client.Option) RPCClient {
	c, err := NewRPCClient(dstService, opts...)
	if err != nil {
		panic("failed to init client: " + err.Error())
	}
	return c
}

func InitClient(dstService string, opts ...client.Option) {
	defaultClient = newClient(dstService, opts...)
}
