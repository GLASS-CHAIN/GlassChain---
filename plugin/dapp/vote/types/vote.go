package types

import (
	"reflect"

	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
)

/*
 * 
 * actio lo  
 * actio lo i nam 
 */

// actio i name 
const (
	TyUnknowAction = iota + 100
	TyCreateGroupAction
	TyUpdateGroupAction
	TyCreateVoteAction
	TyCommitVoteAction
	TyCloseVoteAction
	TyUpdateMemberAction

	NameCreateGroupAction  = "CreateGroup"
	NameUpdateGroupAction  = "UpdateGroup"
	NameCreateVoteAction   = "CreateVote"
	NameCommitVoteAction   = "CommitVote"
	NameCloseVoteAction    = "CloseVote"
	NameUpdateMemberAction = "UpdateMember"
)

// lo i 
const (
	TyUnknownLog = iota + 100
	TyCreateGroupLog
	TyUpdateGroupLog
	TyCreateVoteLog
	TyCommitVoteLog
	TyCloseVoteLog
	TyUpdateMemberLog

	NameCreateGroupLog  = "CreateGroupLog"
	NameUpdateGroupLog  = "UpdateGroupLog"
	NameCreateVoteLog   = "CreateVoteLog"
	NameCommitVoteLog   = "CommitVoteLog"
	NameCloseVoteLog    = "CloseVoteLog"
	NameUpdateMemberLog = "UpdateMemberLog"
)

var (
	//VoteX 
	VoteX = "vote"
	/ actionMap
	actionMap = map[string]int32{
		NameCreateGroupAction:  TyCreateGroupAction,
		NameUpdateGroupAction:  TyUpdateGroupAction,
		NameCreateVoteAction:   TyCreateVoteAction,
		NameCommitVoteAction:   TyCommitVoteAction,
		NameCloseVoteAction:    TyCloseVoteAction,
		NameUpdateMemberAction: TyUpdateMemberAction,
	}
	/ lo i lo  lo 
	logMap = map[int64]*types.LogInfo{
		TyCreateGroupLog:  {Ty: reflect.TypeOf(GroupInfo{}), Name: NameCreateGroupLog},
		TyUpdateGroupLog:  {Ty: reflect.TypeOf(GroupInfo{}), Name: NameUpdateGroupLog},
		TyCreateVoteLog:   {Ty: reflect.TypeOf(VoteInfo{}), Name: NameCreateVoteLog},
		TyCommitVoteLog:   {Ty: reflect.TypeOf(CommitInfo{}), Name: NameCommitVoteLog},
		TyCloseVoteLog:    {Ty: reflect.TypeOf(VoteInfo{}), Name: NameCloseVoteLog},
		TyUpdateMemberLog: {Ty: reflect.TypeOf(MemberInfo{}), Name: NameUpdateMemberLog},
	}
	tlog = log.New("module", "vote.types")
)

// init defines a register function
func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(VoteX))
	/ 
	types.RegFork(VoteX, InitFork)
	types.RegExec(VoteX, InitExecutor)
}

// InitFork defines register fork
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(VoteX, "Enable", 0)
}

// InitExecutor defines register executor
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(VoteX, NewType(cfg))
}

type voteType struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *voteType {
	c := &voteType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetPayload actio 
func (v *voteType) GetPayload() types.Message {
	return &VoteAction{}
}

// GeTypeMap actio i nam 
func (v *voteType) GetTypeMap() map[string]int32 {
	return actionMap
}

// GetLogMap lo 
func (v *voteType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
