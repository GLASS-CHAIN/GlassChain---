package types

import (
	"reflect"

	"github.com/33cn/chain33/types"
)


const (
	TyUnknowAction = iota + 200
	TyLimitOrderAction
	TyMarketOrderAction
	TyRevokeOrderAction

	NameLimitOrderAction  = "LimitOrder"
	NameMarketOrderAction = "MarketOrder"
	NameRevokeOrderAction = "RevokeOrder"

	FuncNameQueryMarketDepth      = "QueryMarketDepth"
	FuncNameQueryHistoryOrderList = "QueryHistoryOrderList"
	FuncNameQueryOrder            = "QueryOrder"
	FuncNameQueryOrderList        = "QueryOrderList"
)

const (
	TyUnknownLog = iota + 200
	TyLimitOrderLog
	TyMarketOrderLog
	TyRevokeOrderLog
)

// OP
const (
	OpBuy = iota + 1
	OpSell
)

//order status
const (
	Ordered = iota
	Completed
	Revoked
)

//const
const (
	ListDESC = int32(0)
	ListASC  = int32(1)
	ListSeek = int32(2)
)

const (
	Count = int32(10)
	MaxMatchCount = 100
)

var (
	ExchangeX = "exchange"
	actionMap = map[string]int32{
		NameLimitOrderAction:  TyLimitOrderAction,
		NameMarketOrderAction: TyMarketOrderAction,
		NameRevokeOrderAction: TyRevokeOrderAction,
	}
	logMap = map[int64]*types.LogInfo{
		TyLimitOrderLog:  {Ty: reflect.TypeOf(ReceiptExchange{}), Name: "TyLimitOrderLog"},
		TyMarketOrderLog: {Ty: reflect.TypeOf(ReceiptExchange{}), Name: "TyMarketOrderLog"},
		TyRevokeOrderLog: {Ty: reflect.TypeOf(ReceiptExchange{}), Name: "TyRevokeOrderLog"},
	}
	//tlog = log.New("module", "exchange.types")
)

// init defines a register function
func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(ExchangeX))

	types.RegFork(ExchangeX, InitFork)
	types.RegExec(ExchangeX, InitExecutor)
}

// InitFork defines register fork
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(ExchangeX, "Enable", 0)
}

// InitExecutor defines register executor
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(ExchangeX, NewType(cfg))
}

//ExchangeType ...
type ExchangeType struct {
	types.ExecTypeBase
}

//NewType ...
func NewType(cfg *types.Chain33Config) *ExchangeType {
	c := &ExchangeType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}


func (e *ExchangeType) GetPayload() types.Message {
	return &ExchangeAction{}
}

func (e *ExchangeType) GetTypeMap() map[string]int32 {
	return actionMap
}

func (e *ExchangeType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
