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
data:  autonomy change
index: status, addr
*/

var changeOpt = &table.Option{
	Prefix:  "LODB-autonomy",
	Name:    "change",
	Primary: "heightindex",
	Index:   []string{"addr", "status", "addr_status"},
}

func NewChangeTable(kvdb db.KV) *table.Table {
	rowmeta := NewChangeRow()
	table, err := table.NewTable(rowmeta, kvdb, changeOpt)
	if err != nil {
		panic(err)
	}
	return table
}

type ChangeRow struct {
	*auty.AutonomyProposalChange
}

func NewChangeRow() *ChangeRow {
	return &ChangeRow{AutonomyProposalChange: &auty.AutonomyProposalChange{}}
}

func (r *ChangeRow) CreateRow() *table.Row {
	return &table.Row{Data: &auty.AutonomyProposalChange{}}
}

func (r *ChangeRow) SetPayload(data types.Message) error {
	if d, ok := data.(*auty.AutonomyProposalChange); ok {
		r.AutonomyProposalChange = d
		return nil
	}
	return types.ErrTypeAsset
}

func (r *ChangeRow) Get(key string) ([]byte, error) {
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
