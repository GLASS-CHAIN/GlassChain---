package testnode

import (
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util/testnode"
)

/*
1. solo   
2.    
*/

//ParaNode 
type ParaNode struct {
	Main *testnode.Chain33Mock
	Para *testnode.Chain33Mock
}

//NewParaNode 
func NewParaNode(main *testnode.Chain33Mock, para *testnode.Chain33Mock) *ParaNode {
	if main == nil {
		main = testnode.New("", nil)
		main.Listen()
	}
	if para == nil {
		cfg := types.NewChain33Config(DefaultConfig)
		testnode.ModifyParaClient(cfg, main.GetCfg().RPC.GrpcBindAddr)
		para = testnode.NewWithConfig(cfg, nil)
		para.Listen()
	}
	return &ParaNode{Main: main, Para: para}
}

//Close 
func (node *ParaNode) Close() {
	node.Para.Close()
	node.Main.Close()
}
