package executor

import (
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	echotypes "github.com/33cn/plugin/plugin/dapp/echo/types/echo"
)

var (
	KeyPrefixPing = "mavl-echo-ping:%s"
	KeyPrefixPang = "mavl-echo-pang:%s"

	KeyPrefixPingLocal = "LODB-echo-ping:%s"
	KeyPrefixPangLocal = "LODB-echo-pang:%s"
)

func Init(name string, cfg *types.Chain33Config, sub []byte) {
	dapp.Register(cfg, echotypes.EchoX, newEcho, 0)
	InitExecType()
}

func InitExecType() {
	ety := types.LoadExecutorType(echotypes.EchoX)
	ety.InitFuncList(types.ListMethod(&Echo{}))
}

type Echo struct {
	dapp.DriverBase
}

func newEcho() dapp.Driver {
	c := &Echo{}
	c.SetChild(c)
	c.SetExecutorType(types.LoadExecutorType(echotypes.EchoX))
	return c
}

func (h *Echo) GetDriverName() string {
	return echotypes.EchoX
}
