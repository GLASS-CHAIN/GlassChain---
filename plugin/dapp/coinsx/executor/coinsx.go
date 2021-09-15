// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor



// package none execer for unknow execer
// all none transaction exec ok, execept nofee
// nofee transaction will not pack into block

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	coinsTy "github.com/33cn/plugin/plugin/dapp/coinsx/types"
)

// var clog = log.New("module", "execs.coins")
var (
	driverName = coinsTy.CoinsxX
	clog       = log.New("module", "execs.paracross")
)

// Init defines a register function
func Init(name string, cfg *types.Chain33Config, sub []byte) {
	if name != driverName {
		panic("system dapp can't be rename")
	}
	drivers.Register(cfg, driverName, newCoinsx, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
	setPrefix()
}

//InitExecType the initialization process is relatively heavyweight, lots of reflect, so it's global
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Coinsx{}))
}

// GetName return name string
func GetName() string {
	return newCoinsx().GetName()
}

// Coins defines coins
type Coinsx struct {
	drivers.DriverBase
}

func newCoinsx() drivers.Driver {
	c := &Coinsx{}
	c.SetChild(c)
	c.SetExecutorType(types.LoadExecutorType(driverName))
	return c
}

// GetDriverName get drive name
func (c *Coinsx) GetDriverName() string {
	return driverName
}

func (c *Coinsx) CheckTx(tx *types.Transaction, index int) error {
	ety := c.GetExecutorType()
	amount, err := ety.Amount(tx)
	if err != nil {
		return err
	}
	if amount < 0 {
		return types.ErrAmount
	}
	return nil
}

// IsFriend Coinsx contract  the mining transaction that runs the ticket contract
func (c *Coinsx) IsFriend(myexec, writekey []byte, othertx *types.Transaction) bool {

	if !c.AllowIsSame(myexec) {
		return false
	}
	types.AssertConfig(c.GetAPI())
	types := c.GetAPI().GetConfig()
	if othertx.ActionName() == "miner" {
		for _, exec := range types.GetMinerExecs() {
			if types.ExecName(exec) == string(othertx.Execer) {
				return true
			}
		}
	}

	return false
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (c *Coinsx) CheckReceiptExecOk() bool {
	return true
}
