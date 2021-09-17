// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

//game action ty
const (
	GameActionCreate = iota + 1
	GameActionMatch
	GameActionCancel
	GameActionClose

	//log for game
	TyLogCreateGame = 711
	TyLogMatchGame  = 712
	TyLogCancleGame = 713
	TyLogCloseGame  = 714
)


var (
	GameX      = "game"
	ExecerGame = []byte(GameX)
)

// action name
const (
	ActionCreateGame = "createGame"
	ActionMatchGame  = "matchGame"
	ActionCancelGame = "cancelGame"
	ActionCloseGame  = "closeGame"
)

// query func name
const (
	FuncNameQueryGameListByIds           = "QueryGameListByIds"
	FuncNameQueryGameListCount           = "QueryGameListCount"
	FuncNameQueryGameListByStatusAndAddr = "QueryGameListByStatusAndAddr"
	FuncNameQueryGameByID                = "QueryGameById"
)
