package rpc

import (
	"context"

	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	echotypes "github.com/33cn/plugin/plugin/dapp/echo/types/echo"
)

type Jrpc struct {
	cli *channelClient
}

type channelClient struct {
	rpctypes.ChannelClient
}

func Init(name string, s rpctypes.RPCServer) {
	cli := &channelClient{}
	cli.Init(name, s, &Jrpc{cli: cli}, nil)
}

func (c *Jrpc) QueryPing(param *echotypes.Query, result *interface{}) error {
	if param == nil {
		return types.ErrInvalidParam
	}
	reply, err := c.cli.QueryPing(context.Background(), param)
	if err != nil {
		return err
	}
	*result = reply
	return nil
}

func (c *channelClient) QueryPing(ctx context.Context, queryParam *echotypes.Query) (types.Message, error) {
	return c.Query(echotypes.EchoX, "GetPing", queryParam)
}
