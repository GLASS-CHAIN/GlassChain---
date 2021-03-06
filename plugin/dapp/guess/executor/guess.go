// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	gty "github.com/33cn/plugin/plugin/dapp/guess/types"
)

var logger = log.New("module", "execs.guess")

var driverName = gty.GuessX

// Init Guess
func Init(name string, cfg *types.Chain33Config, sub []byte) {
	driverName := GetName()
	if name != driverName {
		panic("system dapp can't be rename")
	}

	drivers.Register(cfg, driverName, newGuessGame, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Guess{}))
}

type Guess struct {
	drivers.DriverBase
}

func newGuessGame() drivers.Driver {
	t := &Guess{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

func GetName() string {
	return newGuessGame().GetName()
}

func (g *Guess) ExecutorOrder() int64 {
	return drivers.ExecLocalSameTime
}

func (g *Guess) GetDriverName() string {
	return gty.GuessX
}

/*
// GetPayloadValue GuessAction
func (g *Guess) GetPayloadValue() types.Message {
	return &pkt.GuessGameAction{}
}*/

// CheckReceiptExecOk return true to check if receipt ty is ok
func (g *Guess) CheckReceiptExecOk() bool {
	return true
}
