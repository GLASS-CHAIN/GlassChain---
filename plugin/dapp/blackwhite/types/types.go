// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import (
	"encoding/json"

	"github.com/33cn/chain33/types"
)

// blackwhite action type
const (
	// BlackwhiteActionCreate blackwhite create action
	BlackwhiteActionCreate = iota
	// BlackwhiteActionPlay blackwhite play action
	BlackwhiteActionPlay
	// BlackwhiteActionShow blackwhite show action
	BlackwhiteActionShow
	// BlackwhiteActionTimeoutDone blackwhite timeout action
	BlackwhiteActionTimeoutDone
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, ExecerBlackwhite)
	types.RegFork(BlackwhiteX, InitFork)
	types.RegExec(BlackwhiteX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(BlackwhiteX, "ForkBlackWhiteV2", 900000)
	cfg.RegisterDappFork(BlackwhiteX, "Enable", 850000)
}

//InitExecutor ...
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(BlackwhiteX, NewType(cfg))
}

type BlackwhiteType struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *BlackwhiteType {
	c := &BlackwhiteType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

func (b *BlackwhiteType) GetPayload() types.Message {
	return &BlackwhiteAction{}
}

func (b *BlackwhiteType) GetName() string {
	return BlackwhiteX
}

func (b *BlackwhiteType) GetLogMap() map[int64]*types.LogInfo {
	return logInfo
}

func (b *BlackwhiteType) GetTypeMap() map[string]int32 {
	return actionName
}

func (b BlackwhiteType) ActionName(tx *types.Transaction) string {
	var g BlackwhiteAction
	err := types.Decode(tx.Payload, &g)
	if err != nil {
		return "unknown-Blackwhite-action-err"
	}
	if g.Ty == BlackwhiteActionCreate && g.GetCreate() != nil {
		return "BlackwhiteCreate"
	} else if g.Ty == BlackwhiteActionShow && g.GetShow() != nil {
		return "BlackwhiteShow"
	} else if g.Ty == BlackwhiteActionPlay && g.GetPlay() != nil {
		return "BlackwhitePlay"
	} else if g.Ty == BlackwhiteActionTimeoutDone && g.GetTimeoutDone() != nil {
		return "BlackwhiteTimeoutDone"
	}
	return "unknown"
}

// Amount ...
func (b BlackwhiteType) Amount(tx *types.Transaction) (int64, error) {
	return 0, nil
}

func (b BlackwhiteType) CreateTx(action string, message json.RawMessage) (*types.Transaction, error) {
	glog.Debug("Blackwhite.CreateTx", "action", action)
	var tx *types.Transaction
	return tx, nil
}

// BlackwhiteRoundInfo ...
type BlackwhiteRoundInfo struct {
}

// Input for convert struct
func (t *BlackwhiteRoundInfo) Input(message json.RawMessage) ([]byte, error) {
	var req ReqBlackwhiteRoundInfo
	err := json.Unmarshal(message, &req)
	if err != nil {
		return nil, err
	}
	return types.Encode(&req), nil
}

// Output for convert struct
func (t *BlackwhiteRoundInfo) Output(reply interface{}) (interface{}, error) {
	return reply, nil
}

// BlackwhiteByStatusAndAddr ...
type BlackwhiteByStatusAndAddr struct {
}

// Input for convert struct
func (t *BlackwhiteByStatusAndAddr) Input(message json.RawMessage) ([]byte, error) {
	var req ReqBlackwhiteRoundList
	err := json.Unmarshal(message, &req)
	if err != nil {
		return nil, err
	}
	return types.Encode(&req), nil
}

// Output for convert struct
func (t *BlackwhiteByStatusAndAddr) Output(reply interface{}) (interface{}, error) {
	return reply, nil
}

// BlackwhiteloopResult ...
type BlackwhiteloopResult struct {
}

// Input for convert struct
func (t *BlackwhiteloopResult) Input(message json.RawMessage) ([]byte, error) {
	var req ReqLoopResult
	err := json.Unmarshal(message, &req)
	if err != nil {
		return nil, err
	}
	return types.Encode(&req), nil
}

// Output for convert struct
func (t *BlackwhiteloopResult) Output(reply interface{}) (interface{}, error) {
	return reply, nil
}
