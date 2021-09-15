package executor

import (
	"fmt"

	"github.com/33cn/chain33/types"
	echotypes "github.com/33cn/plugin/plugin/dapp/echo/types/echo"
)

func (h *Echo) ExecLocal_Ping(ping *echotypes.Ping, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	var pingLog echotypes.PingLog
	types.Decode(receipt.Logs[0].Log, &pingLog)
	localKey := []byte(fmt.Sprintf(KeyPrefixPingLocal, pingLog.Msg))
	oldValue, err := h.GetLocalDB().Get(localKey)
	if err != nil && err != types.ErrNotFound {
		return nil, err
	}
	if err == nil {
		types.Decode(oldValue, &pingLog)
	}
	pingLog.Count++
	kv := []*types.KeyValue{{Key: localKey, Value: types.Encode(&pingLog)}}
	return &types.LocalDBSet{KV: kv}, nil
}

func (h *Echo) ExecLocal_Pang(ping *echotypes.Pang, tx *types.Transaction, receipt *types.ReceiptData, index int) (*types.LocalDBSet, error) {

	var pangLog echotypes.PangLog
	types.Decode(receipt.Logs[0].Log, &pangLog)
	localKey := []byte(fmt.Sprintf(KeyPrefixPangLocal, pangLog.Msg))
	oldValue, err := h.GetLocalDB().Get(localKey)
	if err != nil && err != types.ErrNotFound {
		return nil, err
	}
	if err == nil {
		types.Decode(oldValue, &pangLog)
	}
	pangLog.Count++
	kv := []*types.KeyValue{{Key: localKey, Value: types.Encode(&pangLog)}}
	return &types.LocalDBSet{KV: kv}, nil
}
