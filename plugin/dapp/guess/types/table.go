package types

import (
	"fmt"

	"github.com/33cn/chain33/common/db"

	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
)

/*
table  struct
data:  guess
index: addr,status,addr_status,admin,admin_status,category_status
*/

var opt_guess_user = &table.Option{
	Prefix:  "LODB-guess",
	Name:    "user",
	Primary: "index",
	Index:   []string{"addr", "startindex"},
}

func NewGuessUserTable(kvdb db.KV) *table.Table {
	rowmeta := NewGuessUserRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_guess_user)
	if err != nil {
		panic(err)
	}
	return table
}

type GuessUserRow struct {
	*UserBet
}

func NewGuessUserRow() *GuessUserRow {
	return &GuessUserRow{UserBet: &UserBet{}}
}

func (tx *GuessUserRow) CreateRow() *table.Row {
	return &table.Row{Data: &UserBet{}}
}

func (tx *GuessUserRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*UserBet); ok {
		tx.UserBet = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *GuessUserRow) Get(key string) ([]byte, error) {
	if key == "index" {
		return []byte(fmt.Sprintf("%018d", tx.Index)), nil
	} else if key == "addr" {
		return []byte(tx.Addr), nil
	} else if key == "startindex" {
		return []byte(fmt.Sprintf("%018d", tx.StartIndex)), nil
	}

	return nil, types.ErrNotFound
}

var opt_guess_game = &table.Option{
	Prefix:  "LODB-guess",
	Name:    "game",
	Primary: "startindex",
	Index:   []string{"gameid", "status", "admin", "admin_status", "category_status"},
}

func NewGuessGameTable(kvdb db.KV) *table.Table {
	rowmeta := NewGuessGameRow()
	table, err := table.NewTable(rowmeta, kvdb, opt_guess_game)
	if err != nil {
		panic(err)
	}
	return table
}

type GuessGameRow struct {
	*GuessGame
}

func NewGuessGameRow() *GuessGameRow {
	return &GuessGameRow{GuessGame: &GuessGame{}}
}

func (tx *GuessGameRow) CreateRow() *table.Row {
	return &table.Row{Data: &GuessGame{}}
}

func (tx *GuessGameRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*GuessGame); ok {
		tx.GuessGame = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *GuessGameRow) Get(key string) ([]byte, error) {
	if key == "startindex" {
		return []byte(fmt.Sprintf("%018d", tx.StartIndex)), nil
	} else if key == "gameid" {
		return []byte(tx.GameID), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%2d", tx.Status)), nil
	} else if key == "admin" {
		return []byte(tx.AdminAddr), nil
	} else if key == "admin_status" {
		return []byte(fmt.Sprintf("%s:%2d", tx.AdminAddr, tx.Status)), nil
	} else if key == "category_status" {
		return []byte(fmt.Sprintf("%s:%2d", tx.Category, tx.Status)), nil
	}

	return nil, types.ErrNotFound
}
