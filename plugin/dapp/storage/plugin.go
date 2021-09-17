package types

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/33cn/plugin/plugin/dapp/storage/commands"
	"github.com/33cn/plugin/plugin/dapp/storage/executor"
	"github.com/33cn/plugin/plugin/dapp/storage/rpc"
	storagetypes "github.com/33cn/plugin/plugin/dapp/storage/types"
)

/*
 * dap 
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     storagetypes.StorageX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
