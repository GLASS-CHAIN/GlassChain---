package types

import (
	"reflect"

	"github.com/33cn/chain33/types"
)



const (
	TyUnknowAction = iota + 100
	TyRegisterAction
	TyResetAction
	TyTransferAction
	TySuperviseAction
	TyApplyAction

	NameRegisterAction  = "Register"
	NameResetAction     = "ResetKey"
	NameTransferAction  = "Transfer"
	NameSuperviseAction = "Supervise"
	NameApplyAction     = "Apply"

	FuncNameQueryAccountByID      = "QueryAccountByID"
	FuncNameQueryAccountsByStatus = "QueryAccountsByStatus"
	FuncNameQueryExpiredAccounts  = "QueryExpiredAccounts"
	FuncNameQueryAccountByAddr    = "QueryAccountByAddr"
	FuncNameQueryBalanceByID      = "QueryBalanceByID"
)

const (
	TyUnknownLog = iota + 100
	TyRegisterLog
	TyResetLog
	TyTransferLog
	TySuperviseLog
	TyApplyLog
)

const (
	Normal = int32(iota)
	Frozen
	Locked
	Expired
)

//supervior op
const (
	UnknownSupervisorOp = int32(iota)
	Freeze
	UnFreeze
	AddExpire
	Authorize
)

//apply  op
const (
	UnknownApplyOp = int32(iota)
	RevokeReset
	EnforceReset
)

//list ...
const (
	ListDESC = int32(0)
	ListASC  = int32(1)
	ListSeek = int32(2)
)
const (
	Count = int32(10)
)

var (
	AccountmanagerX = "accountmanager"
	actionMap = map[string]int32{
		NameRegisterAction:  TyRegisterAction,
		NameResetAction:     TyResetAction,
		NameApplyAction:     TyApplyAction,
		NameTransferAction:  TyTransferAction,
		NameSuperviseAction: TySuperviseAction,
	}
	logMap = map[int64]*types.LogInfo{
		TyRegisterLog:  {Ty: reflect.TypeOf(AccountReceipt{}), Name: "TyRegisterLog"},
		TyResetLog:     {Ty: reflect.TypeOf(TransferReceipt{}), Name: "TyResetLog"},
		TyTransferLog:  {Ty: reflect.TypeOf(AccountReceipt{}), Name: "TyTransferLog"},
		TySuperviseLog: {Ty: reflect.TypeOf(SuperviseReceipt{}), Name: "TySuperviseLog"},
		TyApplyLog:     {Ty: reflect.TypeOf(AccountReceipt{}), Name: "TyApplyLog"},
	}
	//tlog = log.New("module", "accountmanager.types")
)

// init defines a register function
func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(AccountmanagerX))
	types.RegFork(AccountmanagerX, InitFork)
	types.RegExec(AccountmanagerX, InitExecutor)
}

// InitFork defines register fork
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(AccountmanagerX, "Enable", 0)
}

// InitExecutor defines register executor
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(AccountmanagerX, NewType(cfg))
}

//AccountmanagerType ...
type AccountmanagerType struct {
	types.ExecTypeBase
}

//NewType ...
func NewType(cfg *types.Chain33Config) *AccountmanagerType {
	c := &AccountmanagerType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

func (a *AccountmanagerType) GetPayload() types.Message {
	return &AccountmanagerAction{}
}

func (a *AccountmanagerType) GetTypeMap() map[string]int32 {
	return actionMap
}

func (a *AccountmanagerType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
