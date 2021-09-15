package executor

import (
	"github.com/33cn/chain33/types"
	et "github.com/33cn/plugin/plugin/dapp/accountmanager/types"
)

func (a *Accountmanager) Query_QueryAccountByID(in *et.QueryAccountByID) (types.Message, error) {
	return findAccountByID(a.GetLocalDB(), in.AccountID)
}

func (a *Accountmanager) Query_QueryAccountByAddr(in *et.QueryAccountByAddr) (types.Message, error) {
	return findAccountByAddr(a.GetLocalDB(), in.Addr)
}

func (a *Accountmanager) Query_QueryAccountsByStatus(in *et.QueryAccountsByStatus) (types.Message, error) {
	return findAccountListByStatus(a.GetLocalDB(), in.Status, in.Direction, in.PrimaryKey)
}

func (a *Accountmanager) Query_QueryExpiredAccounts(in *et.QueryExpiredAccounts) (types.Message, error) {
	return findAccountListByIndex(a.GetLocalDB(), in.ExpiredTime, in.PrimaryKey)
}

func (a *Accountmanager) Query_QueryBalanceByID(in *et.QueryBalanceByID) (types.Message, error) {
	return queryBalanceByID(a.GetStateDB(), a.GetLocalDB(), a.GetAPI().GetConfig(), a.GetName(), in)
}
