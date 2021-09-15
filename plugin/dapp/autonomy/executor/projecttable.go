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
data:  autonomy project
index: status, addr
*/

var projectOpt = &table.Option{
	Prefix:  "LODB-autonomy",
	Name:    "project",
	Primary: "heightindex",
	Index:   []string{"addr", "status", "addr_status"},
}

func NewProjectTable(kvdb db.KV) *table.Table {
	rowmeta := NewProjectRow()
	table, err := table.NewTable(rowmeta, kvdb, projectOpt)
	if err != nil {
		panic(err)
	}
	return table
}

type ProjectRow struct {
	*auty.AutonomyProposalProject
}
 
func NewProjectRow() *ProjectRow {
	return &ProjectRow{AutonomyProposalProject: &auty.AutonomyProposalProject{}}
}

func (r *ProjectRow) CreateRow() *table.Row {
	return &table.Row{Data: &auty.AutonomyProposalProject{}}
}

func (r *ProjectRow) SetPayload(data types.Message) error {
	if d, ok := data.(*auty.AutonomyProposalProject); ok {
		r.AutonomyProposalProject = d
		return nil
	}
	return types.ErrTypeAsset
}

func (r *ProjectRow) Get(key string) ([]byte, error) {
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
