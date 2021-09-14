package executor

import (
	"fmt"

	"github.com/33cn/chain33/common/db"

	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	auty "github.com/33cn/plugin/plugin/dapp/autonomy/types"
)

/*
table  struct
data:  autonomy rule
index: status, addr
*/

var ruleOpt = &table.Option{
	Prefix:  "LODB-autonomy",
	Name:    "rule",
	Primary: "heightindex",
	Index:   []string{"addr", "status", "addr_status"},
}

func NewRuleTable(kvdb db.KV) *table.Table {
	rowmeta := NewRuleRow()
	table, err := table.NewTable(rowmeta, kvdb, ruleOpt)
	if err != nil {
		panic(err)
	}
	return table
}

type RuleRow struct {
	*auty.AutonomyProposalRule
}

func NewRuleRow() *RuleRow {
	return &RuleRow{AutonomyProposalRule: &auty.AutonomyProposalRule{}}
}

func (r *RuleRow) CreateRow() *table.Row {
	return &table.Row{Data: &auty.AutonomyProposalRule{}}
}

func (r *RuleRow) SetPayload(data types.Message) error {
	if d, ok := data.(*auty.AutonomyProposalRule); ok {
		r.AutonomyProposalRule = d
		return nil
	}
	return types.ErrTypeAsset
}

func (r *RuleRow) Get(key string) ([]byte, error) {
	if key == "heightindex" {
		return []byte(dapp.HeightIndexStr(r.Height, int64(r.Index))), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%2d", r.Status)), nil
	} else if key == "addr" {
		return []byte(r.Address), nil
	} else if key == "addr_status" {
		return []byte(fmt.Sprintf("%s:%2d", r.Address, r.Status)), nil
	}
	return nil, types.ErrNotFound
}
