package executor

import (
	"fmt"

	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
	ety "github.com/33cn/plugin/plugin/dapp/exchange/types"
)



const (
	KeyPrefixStateDB = "mavl-exchange-"
	KeyPrefixLocalDB = "LODB-exchange"
)

func calcOrderKey(orderID int64) []byte {
	key := fmt.Sprintf("%s"+"orderID:%022d", KeyPrefixStateDB, orderID)
	return []byte(key)
}

var opt_exchange_depth = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "depth",
	Primary: "price",
	Index:   nil,
}

var opt_exchange_order = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "order",
	Primary: "orderID",
	Index:   []string{"market_order", "addr_status"},
}

var opt_exchange_history = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "history",
	Primary: "index",
	Index:   []string{"name", "addr_status"},
}

func NewMarketDepthTable(kvdb db.KV) *table.Table {
	rowmeta := NewMarketDepthRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_exchange_depth)
	if err != nil {
		panic(err)
	}
	return table
}

//NewMarketOrderTable ...
func NewMarketOrderTable(kvdb db.KV) *table.Table {
	rowmeta := NewOrderRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_exchange_order)
	if err != nil {
		panic(err)
	}
	return table
}

//NewHistoryOrderTable ...
func NewHistoryOrderTable(kvdb db.KV) *table.Table {
	rowmeta := NewHistoryOrderRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_exchange_history)
	if err != nil {
		panic(err)
	}
	return table
}

type OrderRow struct {
	*ety.Order
}

func NewOrderRow() *OrderRow {
	return &OrderRow{Order: &ety.Order{}}
}

//CreateRow ...
func (r *OrderRow) CreateRow() *table.Row {
	return &table.Row{Data: &ety.Order{}}
}

func (r *OrderRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*ety.Order); ok {
		r.Order = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (r *OrderRow) Get(key string) ([]byte, error) {
	if key == "orderID" {
		return []byte(fmt.Sprintf("%022d", r.OrderID)), nil
	} else if key == "market_order" {
		return []byte(fmt.Sprintf("%s:%s:%d:%016d", r.GetLimitOrder().LeftAsset.GetSymbol(), r.GetLimitOrder().RightAsset.GetSymbol(), r.GetLimitOrder().Op, r.GetLimitOrder().Price)), nil
	} else if key == "addr_status" {
		return []byte(fmt.Sprintf("%s:%d", r.Addr, r.Status)), nil
	}
	return nil, types.ErrNotFound
}

type HistoryOrderRow struct {
	*ety.Order
}

//NewHistoryOrderRow ...
func NewHistoryOrderRow() *HistoryOrderRow {
	return &HistoryOrderRow{Order: &ety.Order{Value: &ety.Order_LimitOrder{LimitOrder: &ety.LimitOrder{}}}}
}

//CreateRow ...
func (m *HistoryOrderRow) CreateRow() *table.Row {
	return &table.Row{Data: &ety.Order{Value: &ety.Order_LimitOrder{LimitOrder: &ety.LimitOrder{}}}}
}

func (m *HistoryOrderRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*ety.Order); ok {
		m.Order = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (m *HistoryOrderRow) Get(key string) ([]byte, error) {
	if key == "index" {
		return []byte(fmt.Sprintf("%022d", m.Index)), nil
	} else if key == "name" {
		return []byte(fmt.Sprintf("%s:%s", m.GetLimitOrder().LeftAsset.GetSymbol(), m.GetLimitOrder().RightAsset.GetSymbol())), nil
	} else if key == "addr_status" {
		return []byte(fmt.Sprintf("%s:%d", m.Addr, m.Status)), nil
	}
	return nil, types.ErrNotFound
}

type MarketDepthRow struct {
	*ety.MarketDepth
}

func NewMarketDepthRow() *MarketDepthRow {
	return &MarketDepthRow{MarketDepth: &ety.MarketDepth{}}
}

func (m *MarketDepthRow) CreateRow() *table.Row {
	return &table.Row{Data: &ety.MarketDepth{}}
}

func (m *MarketDepthRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*ety.MarketDepth); ok {
		m.MarketDepth = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (m *MarketDepthRow) Get(key string) ([]byte, error) {
	if key == "price" {
		return []byte(fmt.Sprintf("%s:%s:%d:%016d", m.LeftAsset.GetSymbol(), m.RightAsset.GetSymbol(), m.Op, m.Price)), nil
	}
	return nil, types.ErrNotFound
}
