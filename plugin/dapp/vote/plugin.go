package types

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/33cn/plugin/plugin/dapp/vote/commands"
	"github.com/33cn/plugin/plugin/dapp/vote/executor"
	"github.com/33cn/plugin/plugin/dapp/vote/rpc"
	votetypes "github.com/33cn/plugin/plugin/dapp/vote/types"
)

/*
 * dap 
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     votetypes.VoteX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
