package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	exchangetypes "github.com/33cn/plugin/plugin/dapp/exchange/types"
)

var (

	elog = log.New("module", "exchange.executor")
)

var driverName = exchangetypes.ExchangeX

// Init register dapp
func Init(name string, cfg *types.Chain33Config, sub []byte) {
	drivers.Register(cfg, GetName(), NewExchange, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

// InitExecType Init Exec Type
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&exchange{}))
}

type exchange struct {
	drivers.DriverBase
}

//NewExchange ...
func NewExchange() drivers.Driver {
	t := &exchange{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetName get driver name
func GetName() string {
	return NewExchange().GetName()
}

//GetDriverName ...
func (e *exchange) GetDriverName() string {
	return driverName
}


func (e *exchange) CheckTx(tx *types.Transaction, index int) error {

	
	var exchange exchangetypes.ExchangeAction
	types.Decode(tx.GetPayload(), &exchange)
	cfg := e.GetAPI().GetConfig()
	if exchange.Ty == exchangetypes.TyLimitOrderAction {
		limitOrder := exchange.GetLimitOrder()
		left := limitOrder.GetLeftAsset()
		right := limitOrder.GetRightAsset()
		price := limitOrder.GetPrice()
		amount := limitOrder.GetAmount()
		op := limitOrder.GetOp()
		if !CheckExchangeAsset(cfg.GetCoinExec(), left, right) {
			return exchangetypes.ErrAsset
		}
		if !CheckPrice(price) {
			return exchangetypes.ErrAssetPrice
		}
		if !CheckAmount(amount, cfg.GetCoinPrecision()) {
			return exchangetypes.ErrAssetAmount
		}
		if !CheckOp(op) {
			return exchangetypes.ErrAssetOp
		}
	}
	if exchange.Ty == exchangetypes.TyMarketOrderAction {
		return types.ErrActionNotSupport
	}
	return nil
}


func (e *exchange) ExecutorOrder() int64 {
	return drivers.ExecLocalSameTime
}

// GetPayloadValue get payload value
func (e *exchange) GetPayloadValue() types.Message {
	return &exchangetypes.ExchangeAction{}
}
