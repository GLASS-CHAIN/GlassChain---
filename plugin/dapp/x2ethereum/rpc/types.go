package rpc

import (
	rpctypes "github.com/33cn/chain33/rpc/types"
	x2ethereumtypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

/*
 * rp 
 */

// grp servic 
type channelClient struct {
	rpctypes.ChannelClient
}

// Jrpc json rp 
type Jrpc struct {
	cli *channelClient
}

// Grpc grpc
type Grpc struct {
	*channelClient
}

// Init init rpc
func Init(name string, s rpctypes.RPCServer) {
	cli := &channelClient{}
	grpc := &Grpc{channelClient: cli}
	cli.Init(name, s, &Jrpc{cli: cli}, grpc)
	/ grpc servic grpc server pb.g 
	x2ethereumtypes.RegisterX2EthereumServer(s.GRPC(), grpc)
}
