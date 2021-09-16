// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package state

import (
	"fmt"
	"strings"

	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/client"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/model"
	evmtypes "github.com/33cn/plugin/plugin/dapp/evm/types"
)

type MemoryStateDB struct {

	StateDB db.KV

	LocalDB db.KVDB

	CoinsAccount *account.DB

	evmPlatformAddr string

	accounts map[string]*ContractAccount

	refund uint64

	logs    map[common.Hash][]*model.ContractLog
	logSize uint

	snapshots  []*Snapshot
	currentVer *Snapshot
	versionID  int

	preimages map[common.Hash][]byte

	txHash  common.Hash
	txIndex int

	blockHeight int64

	stateDirty map[string]interface{}
	dataDirty  map[string]interface{}
	api        client.QueueProtocolAPI
}

func NewMemoryStateDB(StateDB db.KV, LocalDB db.KVDB, CoinsAccount *account.DB, blockHeight int64, api client.QueueProtocolAPI) *MemoryStateDB {
	mdb := &MemoryStateDB{
		StateDB:         StateDB,
		LocalDB:         LocalDB,
		CoinsAccount:    CoinsAccount,
		evmPlatformAddr: address.GetExecAddress(api.GetConfig().ExecName("evm")).String(),
		accounts:        make(map[string]*ContractAccount),
		logs:            make(map[common.Hash][]*model.ContractLog),
		logSize:         0,
		preimages:       make(map[common.Hash][]byte),
		stateDirty:      make(map[string]interface{}),
		dataDirty:       make(map[string]interface{}),
		blockHeight:     blockHeight,
		refund:          0,
		txIndex:         0,
		api:             api,
	}
	return mdb
}

func (mdb *MemoryStateDB) Prepare(txHash common.Hash, txIndex int) {
	mdb.txHash = txHash
	mdb.txIndex = txIndex
	log15.Info("MemoryStateDB::Prepare", "txHash", txHash.Hex(), "txIndex", txIndex, "logSize", mdb.logSize)
}

func (mdb *MemoryStateDB) CreateAccount(addr, creator string, execName, alias string) {
	acc := mdb.GetAccount(addr)
	if acc == nil {

		acc := NewContractAccount(addr, mdb)
		acc.SetCreator(creator)
		acc.SetExecName(execName)
		acc.SetAliasName(alias)
		mdb.accounts[addr] = acc
		mdb.addChange(createAccountChange{baseChange: baseChange{}, account: addr})
	}
}

func (mdb *MemoryStateDB) addChange(entry DataChange) {
	if mdb.currentVer != nil {
		mdb.currentVer.append(entry)
	}
}

func (mdb *MemoryStateDB) SubBalance(addr, caddr string, value uint64) {
	res := mdb.Transfer(addr, caddr, value)
	log15.Debug("transfer result", "from", addr, "to", caddr, "amount", value, "result", res)
}

func (mdb *MemoryStateDB) AddBalance(addr, caddr string, value uint64) {
	res := mdb.Transfer(caddr, addr, value)
	log15.Debug("transfer result", "from", addr, "to", caddr, "amount", value, "result", res)
}

// GetBalance
func (mdb *MemoryStateDB) GetBalance(addr string) uint64 {
	ac := mdb.CoinsAccount.LoadExecAccount(addr, mdb.evmPlatformAddr)
	return uint64(ac.Balance)
}


func (mdb *MemoryStateDB) GetNonce(addr string) uint64 {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.GetNonce()
	}
	return 0
}

func (mdb *MemoryStateDB) SetNonce(addr string, nonce uint64) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		acc.SetNonce(nonce)
	}
}

func (mdb *MemoryStateDB) GetCodeHash(addr string) common.Hash {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return common.BytesToHash(acc.Data.GetCodeHash())
	}
	return common.Hash{}
}

func (mdb *MemoryStateDB) GetCode(addr string) []byte {
	//if "15wDXJKYxTq3FvfL5uK4MCZXdQfcSTGUtt" == addr {
	//	panic("MemoryStateDB::debugCall::GetCode")
	//}
	log15.Debug("MemoryStateDB::debugCall::GetCode", "addr", addr)
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.Data.GetCode()
	}
	return nil
}

func (mdb *MemoryStateDB) SetCode(addr string, code []byte) {
	log15.Debug("MemoryStateDB::debugCall::SetCode", "addr", addr)
	acc := mdb.GetAccount(addr)
	if acc != nil {
		mdb.dataDirty[addr] = true
		acc.SetCode(code)
	}
}

func (mdb *MemoryStateDB) SetAbi(addr, abi string) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		mdb.dataDirty[addr] = true
		acc.SetAbi(abi)
	}
}

func (mdb *MemoryStateDB) GetAbi(addr string) string {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.Data.GetAbi()
	}
	return ""
}

func (mdb *MemoryStateDB) GetCodeSize(addr string) int {
	code := mdb.GetCode(addr)
	if code != nil {
		return len(code)
	}
	return 0
}

func (mdb *MemoryStateDB) AddRefund(gas uint64) {
	mdb.addChange(refundChange{baseChange: baseChange{}, prev: mdb.refund})
	mdb.refund += gas
}

func (mdb *MemoryStateDB) GetRefund() uint64 {
	return mdb.refund
}

func (mdb *MemoryStateDB) GetAccount(addr string) *ContractAccount {
	if acc, ok := mdb.accounts[addr]; ok {
		return acc
	}

	contract := NewContractAccount(addr, mdb)
	contract.LoadContract(mdb.StateDB)
	if contract.Empty() {
		return nil
	}
	mdb.accounts[addr] = contract
	return contract
}

func (mdb *MemoryStateDB) GetState(addr string, key common.Hash) common.Hash {

	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.GetState(key)
	}
	return common.Hash{}
}

func (mdb *MemoryStateDB) SetState(addr string, key common.Hash, value common.Hash) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		acc.SetState(key, value)

		cfg := mdb.api.GetConfig()
		if !cfg.IsDappFork(mdb.blockHeight, "evm", evmtypes.ForkEVMState) {
			mdb.stateDirty[addr] = true
		}
	}
}

func (mdb *MemoryStateDB) TransferStateData(addr string) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		acc.TransferState()
	}
}

func (mdb *MemoryStateDB) UpdateState(addr string) {
	mdb.stateDirty[addr] = true
}

func (mdb *MemoryStateDB) Suicide(addr string) bool {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		mdb.addChange(suicideChange{
			baseChange: baseChange{},
			account:    addr,
			prev:       acc.State.GetSuicided(),
		})
		mdb.stateDirty[addr] = true
		return acc.Suicide()
	}
	return false
}

func (mdb *MemoryStateDB) HasSuicided(addr string) bool {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		return acc.HasSuicided()
	}
	return false
}

func (mdb *MemoryStateDB) Exist(addr string) bool {
	return mdb.GetAccount(addr) != nil
}

func (mdb *MemoryStateDB) Empty(addr string) bool {
	acc := mdb.GetAccount(addr)

	if acc != nil && !acc.Empty() {
		return false
	}

	if mdb.GetBalance(addr) != 0 {
		return false
	}
	return true
}

func (mdb *MemoryStateDB) RevertToSnapshot(version int) {
	if version >= len(mdb.snapshots) {
		return
	}

	ver := mdb.snapshots[version]

	if ver == nil || ver.id != version {
		log15.Crit(fmt.Errorf("Snapshot id %v cannot be reverted", version).Error())
		return
	}

	for index := len(mdb.snapshots) - 1; index >= version; index-- {
		mdb.snapshots[index].revert()
	}

	mdb.snapshots = mdb.snapshots[:version]
	mdb.versionID = version
	if version == 0 {
		mdb.currentVer = nil
	} else {
		mdb.currentVer = mdb.snapshots[version-1]
	}

}

func (mdb *MemoryStateDB) Snapshot() int {
	id := mdb.versionID
	mdb.versionID++
	mdb.currentVer = &Snapshot{id: id, statedb: mdb}
	mdb.snapshots = append(mdb.snapshots, mdb.currentVer)
	log15.Debug("MemoryStateDB::Snapshot", "mdb.versionID", mdb.versionID)
	return id
}

func (mdb *MemoryStateDB) GetLastSnapshot() *Snapshot {
	if mdb.versionID == 0 {
		return nil
	}
	return mdb.snapshots[mdb.versionID-1]
}

func (mdb *MemoryStateDB) GetReceiptLogs(addr string) (logs []*types.ReceiptLog) {
	acc := mdb.GetAccount(addr)
	if acc != nil {
		if mdb.stateDirty[addr] != nil {
			stateLog := acc.BuildStateLog()
			if stateLog != nil {
				logs = append(logs, stateLog)
			}
		}

		if mdb.dataDirty[addr] != nil {
			logs = append(logs, acc.BuildDataLog())
		}
		return
	}
	return
}

func (mdb *MemoryStateDB) GetChangedData(version int) (kvSet []*types.KeyValue, logs []*types.ReceiptLog) {
	if version < 0 || version >= len(mdb.snapshots) {
		return
	}

	for _, snapshot := range mdb.snapshots {
		kv, log := snapshot.getData()
		if kv != nil {
			kvSet = append(kvSet, kv...)
		}
		if log != nil {
			logs = append(logs, log...)
		}
	}

	return
}

func (mdb *MemoryStateDB) CanTransfer(sender string, amount uint64) bool {
	senderAcc := mdb.CoinsAccount.LoadExecAccount(sender, mdb.evmPlatformAddr)
	log15.Info("CanTransfer", "balance", senderAcc.Balance, "sender", sender, "evmPlatformAddr", mdb.evmPlatformAddr,
		"mdb.CoinsAccount", mdb.CoinsAccount)

	return senderAcc.Balance >= int64(amount)
}

type TransferType int

const (
	_ TransferType = iota

	NoNeed

	ToExec

	FromExec

	Error
)

func (mdb *MemoryStateDB) Transfer(sender, recipient string, amount uint64) bool {
	log15.Debug("transfer from contract to external(contract)", "sender", sender, "recipient", recipient, "amount", amount)
	var (
		ret *types.Receipt
		err error
	)

	value := int64(amount)
	if value < 0 {
		return false
	}

	if 0 == value {
		return true
	}

	ret, err = mdb.CoinsAccount.ExecTransfer(sender, recipient, mdb.evmPlatformAddr, int64(amount))
	if err != nil {
		log15.Error("transfer error", "sender", sender, "recipient", recipient, "amount", amount, "err info", err)
		return false
	}
	if ret != nil {
		mdb.addChange(transferChange{
			baseChange: baseChange{},
			amount:     value,
			data:       ret.KV,
			logs:       ret.Logs,
		})
	}
	log15.Info("transfer successful", "balance", mdb.CoinsAccount.LoadExecAccount(recipient, mdb.evmPlatformAddr).Balance,
		"mdb.CoinsAccount", mdb.CoinsAccount)

	return true
}


func (mdb *MemoryStateDB) transfer2Contract(sender, recipient string, amount int64) (ret *types.Receipt, err error) {

	contract := mdb.GetAccount(recipient)
	if contract == nil {
		return nil, model.ErrAddrNotExists
	}
	creator := contract.GetCreator()
	if len(creator) == 0 {
		return nil, model.ErrNoCreator
	}
	execAddr := recipient

	ret = &types.Receipt{}

	cfg := mdb.api.GetConfig()
	if cfg.IsDappFork(mdb.GetBlockHeight(), "evm", evmtypes.ForkEVMFrozen) {
		rs, err := mdb.CoinsAccount.ExecTransfer(sender, execAddr, execAddr, amount)
		if err != nil {
			return nil, err
		}

		ret.KV = append(ret.KV, rs.KV...)
		ret.Logs = append(ret.Logs, rs.Logs...)
	} else {
		if strings.Compare(sender, creator) != 0 {

			rs, err := mdb.CoinsAccount.ExecTransfer(sender, creator, execAddr, amount)
			if err != nil {
				return nil, err
			}

			ret.KV = append(ret.KV, rs.KV...)
			ret.Logs = append(ret.Logs, rs.Logs...)
		}
	}

	return ret, nil
}


func (mdb *MemoryStateDB) transfer2External(sender, recipient string, amount int64) (ret *types.Receipt, err error) {

	contract := mdb.GetAccount(sender)
	if contract == nil {
		return nil, model.ErrAddrNotExists
	}
	creator := contract.GetCreator()
	if len(creator) == 0 {
		return nil, model.ErrNoCreator
	}

	execAddr := sender

	cfg := mdb.api.GetConfig()
	if cfg.IsDappFork(mdb.GetBlockHeight(), "evm", evmtypes.ForkEVMFrozen) {
		ret, err = mdb.CoinsAccount.ExecTransfer(execAddr, recipient, execAddr, amount)
		if err != nil {
			return nil, err
		}
	} else {

		if strings.Compare(creator, recipient) != 0 {
			ret, err = mdb.CoinsAccount.ExecTransfer(creator, recipient, execAddr, amount)
			if err != nil {
				return nil, err
			}
		}
	}
	return ret, nil
}

func (mdb *MemoryStateDB) mergeResult(one, two *types.Receipt) (ret *types.Receipt) {
	ret = one
	if ret == nil {
		ret = two
	} else if two != nil {
		ret.KV = append(ret.KV, two.KV...)
		ret.Logs = append(ret.Logs, two.Logs...)
	}
	return
}

func (mdb *MemoryStateDB) AddLog(log *model.ContractLog) {
	newEvmLog := &types.EVMLog{
		Topic: [][]byte{log.Topics[0].Bytes()},
		Data:  log.Data,
	}
	if len(log.Topics) > 0 {
		for i := 1; i < len(log.Topics); i++ {
			newEvmLog.Topic = append(newEvmLog.Topic, log.Topics[i].Bytes())
		}
	}
	receiptLog := &types.ReceiptLog{
		Ty:  evmtypes.TyLogEVMEventData,
		Log: types.Encode(newEvmLog),
	}

	mdb.addChange(addLogChange{
		txhash: mdb.txHash,
		logs:   []*types.ReceiptLog{receiptLog}})

	log.TxHash = mdb.txHash
	log.Index = int(mdb.logSize)
	mdb.logs[mdb.txHash] = append(mdb.logs[mdb.txHash], log)
	mdb.logSize++
	log15.Info("MemoryStateDB::AddLog", "txhash", mdb.txHash.Hex(), "blockHeight", mdb.blockHeight, "txIndex", mdb.txIndex,
		"mdb.logSize", mdb.logSize, "topic", log.Topics[0].Hex())
}

func (mdb *MemoryStateDB) AddPreimage(hash common.Hash, data []byte) {

	if _, ok := mdb.preimages[hash]; !ok {
		mdb.addChange(addPreimageChange{hash: hash})
		pi := make([]byte, len(data))
		copy(pi, data)
		mdb.preimages[hash] = pi
	}
}

func (mdb *MemoryStateDB) PrintLogs() {
	items := mdb.logs[mdb.txHash]
	log15.Debug("PrintLogs", "item number:", len(items), "txhash", mdb.txHash.Hex())
	for _, item := range items {
		item.PrintLog()
	}
}

func (mdb *MemoryStateDB) WritePreimages(number int64) {
	for k, v := range mdb.preimages {
		log15.Debug("Contract preimages ", "key:", k.Str(), "value:", common.Bytes2Hex(v), "block height:", number)
	}
}

func (mdb *MemoryStateDB) ResetDatas() {
	mdb.currentVer = nil
	mdb.snapshots = mdb.snapshots[:0]
}

func (mdb *MemoryStateDB) GetBlockHeight() int64 {
	return mdb.blockHeight
}

func (mdb *MemoryStateDB) GetConfig() *types.Chain33Config {
	return mdb.api.GetConfig()
}
