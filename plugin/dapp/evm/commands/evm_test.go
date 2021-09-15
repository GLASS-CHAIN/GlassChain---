package commands

import (
	"testing"

	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util/testnode"
	"github.com/stretchr/testify/assert"

	evm "github.com/33cn/plugin/plugin/dapp/evm/executor"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"

	"github.com/33cn/chain33/client/mocks"
	_ "github.com/33cn/chain33/system"
	"github.com/stretchr/testify/mock"
)

func TestQueryDebug(t *testing.T) {
	var cfg = types.NewChain33Config(types.GetDefaultCfgstring())
	evm.Init(evmtypes.ExecutorName, cfg, nil)
	var debugReq = evmtypes.EvmDebugReq{Optype: 1}
	js, err := types.PBToJSON(&debugReq)
	assert.Nil(t, err)
	in := &rpctypes.Query4Jrpc{
		Execer:   "evm",
		FuncName: "EvmDebug",
		Payload:  js,
	}

	var mockResp = evmtypes.EvmDebugResp{DebugStatus: "on"}

	mockapi := &mocks.QueueProtocolAPI{}
	mockapi.On("Close").Return()
	mockapi.On("Query", "evm", "EvmDebug", &debugReq).Return(&mockResp, nil)
	mockapi.On("GetConfig", mock.Anything).Return(cfg, nil)

	mock33 := testnode.New("", mockapi)
	defer mock33.Close()
	rpcCfg := mock33.GetCfg().RPC
	rpcCfg.JrpcBindAddr = "127.0.0.1:8899"
	mock33.GetRPC().Listen()

	jsonClient, err := jsonclient.NewJSONClient("http://" + rpcCfg.JrpcBindAddr + "/")
	assert.Nil(t, err)
	assert.NotNil(t, jsonClient)

	var debugResp evmtypes.EvmDebugResp
	err = jsonClient.Call("Chain33.Query", in, &debugResp)
	assert.Nil(t, err)
	assert.Equal(t, "on", debugResp.DebugStatus)
}
