// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"fmt"

	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	gt "github.com/33cn/plugin/plugin/dapp/game/types"
)

var glog = log.New("module", "execs.game")

var driverName = gt.GameX

// Init register dapp
func Init(name string, cfg *types.Chain33Config, sub []byte) {
	drivers.Register(cfg, GetName(), newGame, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

//InitExecType ...
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&Game{}))
}

// Game the game inherits all the attributes of the driverBase.
type Game struct {
	drivers.DriverBase
}

func newGame() drivers.Driver {
	g := &Game{}
	g.SetChild(g)
	g.SetExecutorType(types.LoadExecutorType(driverName))
	return g
}

// GetName get name
func GetName() string {
	return newGame().GetName()
}

// GetDriverName get driver name
func (g *Game) GetDriverName() string {
	return driverName
}

// update Index
func (g *Game) updateIndex(log *gt.ReceiptGame) (kvs []*types.KeyValue) {
	// save the index generated by this action first.
	kvs = append(kvs, addGameAddrIndex(log.Status, log.GameId, log.Addr, log.Index))
	kvs = append(kvs, addGameStatusIndex(log.Status, log.GameId, log.Index))
	if log.Status == gt.GameActionMatch {
		kvs = append(kvs, addGameAddrIndex(log.Status, log.GameId, log.CreateAddr, log.Index))
		kvs = append(kvs, delGameAddrIndex(gt.GameActionCreate, log.CreateAddr, log.PrevIndex))
		kvs = append(kvs, delGameStatusIndex(gt.GameActionCreate, log.PrevIndex))
	}
	if log.Status == gt.GameActionCancel {
		kvs = append(kvs, delGameAddrIndex(gt.GameActionCreate, log.CreateAddr, log.PrevIndex))
		kvs = append(kvs, delGameStatusIndex(gt.GameActionCreate, log.PrevIndex))
	}

	if log.Status == gt.GameActionClose {
		kvs = append(kvs, addGameAddrIndex(log.Status, log.GameId, log.MatchAddr, log.Index))
		kvs = append(kvs, delGameAddrIndex(gt.GameActionMatch, log.MatchAddr, log.PrevIndex))
		kvs = append(kvs, delGameAddrIndex(gt.GameActionMatch, log.CreateAddr, log.PrevIndex))
		kvs = append(kvs, delGameStatusIndex(gt.GameActionMatch, log.PrevIndex))
	}
	return kvs
}

// rollback Index
func (g *Game) rollbackIndex(log *gt.ReceiptGame) (kvs []*types.KeyValue) {

	kvs = append(kvs, delGameAddrIndex(log.Status, log.Addr, log.Index))
	kvs = append(kvs, delGameStatusIndex(log.Status, log.Index))

	if log.Status == gt.GameActionMatch {
		kvs = append(kvs, delGameAddrIndex(log.Status, log.CreateAddr, log.Index))
		kvs = append(kvs, addGameAddrIndex(gt.GameActionCreate, log.GameId, log.CreateAddr, log.PrevIndex))
		kvs = append(kvs, addGameStatusIndex(gt.GameActionCreate, log.GameId, log.PrevIndex))
	}

	if log.Status == gt.GameActionCancel {
		kvs = append(kvs, addGameAddrIndex(gt.GameActionCreate, log.GameId, log.CreateAddr, log.PrevIndex))
		kvs = append(kvs, addGameStatusIndex(gt.GameActionCreate, log.GameId, log.PrevIndex))
	}

	if log.Status == gt.GameActionClose {
		kvs = append(kvs, delGameAddrIndex(log.Status, log.MatchAddr, log.Index))
		kvs = append(kvs, addGameAddrIndex(gt.GameActionMatch, log.GameId, log.MatchAddr, log.PrevIndex))
		kvs = append(kvs, addGameAddrIndex(gt.GameActionMatch, log.GameId, log.CreateAddr, log.PrevIndex))
		kvs = append(kvs, addGameStatusIndex(gt.GameActionMatch, log.GameId, log.PrevIndex))
	}
	return kvs
}

func calcGameStatusIndexKey(status int32, index int64) []byte {
	key := fmt.Sprintf("LODB-game-status:%d:%018d", status, index)
	return []byte(key)
}

func calcGameStatusIndexPrefix(status int32) []byte {
	key := fmt.Sprintf("LODB-game-status:%d:", status)
	return []byte(key)
}
func calcGameAddrIndexKey(status int32, addr string, index int64) []byte {
	key := fmt.Sprintf("LODB-game-addr:%d:%s:%018d", status, addr, index)
	return []byte(key)
}
func calcGameAddrIndexPrefix(status int32, addr string) []byte {
	key := fmt.Sprintf("LODB-game-addr:%d:%s:", status, addr)
	return []byte(key)
}
func addGameStatusIndex(status int32, gameID string, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcGameStatusIndexKey(status, index)
	record := &gt.GameRecord{
		GameId: gameID,
		Index:  index,
	}
	kv.Value = types.Encode(record)
	return kv
}
func addGameAddrIndex(status int32, gameID, addr string, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcGameAddrIndexKey(status, addr, index)
	record := &gt.GameRecord{
		GameId: gameID,
		Index:  index,
	}
	kv.Value = types.Encode(record)
	return kv
}
func delGameStatusIndex(status int32, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcGameStatusIndexKey(status, index)
	kv.Value = nil
	return kv
}
func delGameAddrIndex(status int32, addr string, index int64) *types.KeyValue {
	kv := &types.KeyValue{}
	kv.Key = calcGameAddrIndexKey(status, addr, index)
	//value置nil,提交时，会自动执行删除操作
	kv.Value = nil
	return kv
}

// ReplyGameList the data structure returned when querying the game list.
type ReplyGameList struct {
	Games []*Game `json:"games"`
}

// ReplyGame the data structure returned when querying a single game.
type ReplyGame struct {
	Game *Game `json:"game"`
}

// GetPayloadValue get payload value
func (g *Game) GetPayloadValue() types.Message {
	return &gt.GameAction{}
}

// GetTypeMap get TypeMap
func (g *Game) GetTypeMap() map[string]int32 {
	return map[string]int32{
		"Create": gt.GameActionCreate,
		"Match":  gt.GameActionMatch,
		"Cancel": gt.GameActionCancel,
		"Close":  gt.GameActionClose,
	}
}

// CheckReceiptExecOk return true to check if receipt ty is ok
func (g *Game) CheckReceiptExecOk() bool {
	return true
}
