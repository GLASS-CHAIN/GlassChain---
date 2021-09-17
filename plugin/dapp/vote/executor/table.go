package executor

import (
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
	vty "github.com/33cn/plugin/plugin/dapp/vote/types"
)

var (
	groupTablePrimary  = "groupid"
	voteTablePrimary   = "voteid"
	memberTablePrimary = "addr"
)
var groupTableOpt = &table.Option{
	Prefix:  keyPrefixLocalDB,
	Name:    "group",
	Primary: groupTablePrimary,
}

var voteTableOpt = &table.Option{
	Prefix:  keyPrefixLocalDB,
	Name:    "vote",
	Primary: voteTablePrimary,
	Index:   []string{groupTablePrimary},
}

var memberTableOpt = &table.Option{
	Prefix:  keyPrefixLocalDB,
	Name:    "member",
	Primary: memberTablePrimary,
}

/ 
func newTable(kvDB db.KV, rowMeta table.RowMeta, opt *table.Option) *table.Table {
	table, err := table.NewTable(rowMeta, kvDB, opt)
	if err != nil {
		panic(err)
	}
	return table
}

func newGroupTable(kvDB db.KV) *table.Table {
	return newTable(kvDB, &groupTableRow{}, groupTableOpt)
}

func newVoteTable(kvDB db.KV) *table.Table {
	return newTable(kvDB, &voteTableRow{}, voteTableOpt)
}

func newMemberTable(kvDB db.KV) *table.Table {
	return newTable(kvDB, &memberTableRow{}, memberTableOpt)
}

//groupTableRow table meta 
type groupTableRow struct {
	*vty.GroupInfo
}

//CreateRow 
func (r *groupTableRow) CreateRow() *table.Row {
	return &table.Row{Data: &vty.GroupInfo{}}
}

//SetPayload 
func (r *groupTableRow) SetPayload(data types.Message) error {
	if d, ok := data.(*vty.GroupInfo); ok {
		r.GroupInfo = d
		return nil
	}
	return types.ErrTypeAsset
}

//Get indexName  indexValue
func (r *groupTableRow) Get(key string) ([]byte, error) {
	if key == groupTablePrimary {
		return []byte(r.GroupInfo.GetID()), nil
	}
	return nil, types.ErrNotFound
}

//voteTableRow table meta 
type voteTableRow struct {
	*vty.VoteInfo
}

//CreateRow 
func (r *voteTableRow) CreateRow() *table.Row {
	return &table.Row{Data: &vty.VoteInfo{}}
}

//SetPayload 
func (r *voteTableRow) SetPayload(data types.Message) error {
	if d, ok := data.(*vty.VoteInfo); ok {
		r.VoteInfo = d
		return nil
	}
	return types.ErrTypeAsset
}

//Get indexName  indexValue
func (r *voteTableRow) Get(key string) ([]byte, error) {
	if key == voteTablePrimary {
		return []byte(r.VoteInfo.GetID()), nil
	} else if key == groupTablePrimary {
		return []byte(r.VoteInfo.GetGroupID()), nil
	}
	return nil, types.ErrNotFound
}

type memberTableRow struct {
	*vty.MemberInfo
}

//CreateRow 
func (r *memberTableRow) CreateRow() *table.Row {
	return &table.Row{Data: &vty.MemberInfo{}}
}

//SetPayload 
func (r *memberTableRow) SetPayload(data types.Message) error {
	if d, ok := data.(*vty.MemberInfo); ok {
		r.MemberInfo = d
		return nil
	}
	return types.ErrTypeAsset
}

//Get indexName  indexValue
func (r *memberTableRow) Get(key string) ([]byte, error) {
	if key == memberTablePrimary {
		return []byte(r.MemberInfo.GetAddr()), nil
	}
	return nil, types.ErrNotFound
}
