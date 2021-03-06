// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package executor

import (
	"github.com/33cn/chain33/common/db/table"

	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	pty "github.com/33cn/plugin/plugin/dapp/collateralize/types"
	issuanceE "github.com/33cn/plugin/plugin/dapp/issuance/types"
	tokenE "github.com/33cn/plugin/plugin/dapp/token/executor"
)

// List control
const (
	ListDESC     = int32(0)   
	ListASC      = int32(1)   
	DefaultCount = int32(20)  
	MaxCount     = int32(100) 
)

//setting
const (
	DefaultDebtCeiling        = 10000           
	DefaultLiquidationRatio   = 0.4 * 1e4       
	DefaultStabilityFeeRation = 0.08 * 1e4      
	DefaultPeriod             = 3600 * 24 * 365 
	DefaultTotalBalance       = 0               
	PriceWarningRate          = 1.3 * 1e4       
	ExpireWarningTime         = 3600 * 24 * 10  
)

// CollateralizeDB def
type CollateralizeDB struct {
	pty.Collateralize
}

// GetKVSet for CollateralizeDB
func (coll *CollateralizeDB) GetKVSet() (kvset []*types.KeyValue) {
	value := types.Encode(&coll.Collateralize)
	kvset = append(kvset, &types.KeyValue{Key: Key(coll.CollateralizeId), Value: value})
	return kvset
}

// Save for CollateralizeDB
func (coll *CollateralizeDB) Save(db dbm.KV) {
	set := coll.GetKVSet()
	for i := 0; i < len(set); i++ {
		db.Set(set[i].GetKey(), set[i].Value)
	}
}

// Key for Collateralize
func Key(id string) (key []byte) {
	key = append(key, []byte("mavl-"+pty.CollateralizeX+"-")...)
	key = append(key, []byte(id)...)
	return key
}

//ConfigKey Key for CollateralizeConfig
func ConfigKey() (key []byte) {
	key = append(key, []byte("mavl-"+pty.CollateralizeX+"-config")...)
	return key
}

//AddrKey Key for CollateralizeAddrConfig
func AddrKey() (key []byte) {
	key = append(key, []byte("mavl-"+issuanceE.IssuanceX+"-addr")...)
	return key
}

//PriceKey Key for IssuancePriceFeed
func PriceKey() (key []byte) {
	key = append(key, []byte("mavl-"+pty.CollateralizeX+"-price")...)
	return key
}

// Action struct
type Action struct {
	coinsAccount  *account.DB 
	tokenAccount  *account.DB 
	db            dbm.KV
	localDB       dbm.KVDB
	txhash        []byte
	fromaddr      string
	blocktime     int64
	height        int64
	execaddr      string
	difficulty    uint64
	index         int
	Collateralize *Collateralize
}

// NewCollateralizeAction generate New Action
func NewCollateralizeAction(c *Collateralize, tx *types.Transaction, index int) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	cfg := c.GetAPI().GetConfig()
	tokenDb, err := account.NewAccountDB(cfg, tokenE.GetName(), pty.CCNYTokenName, c.GetStateDB())
	if err != nil {
		clog.Error("NewCollateralizeAction", "Get Account DB error", "error", err)
		return nil
	}

	return &Action{
		coinsAccount: c.GetCoinsAccount(), tokenAccount: tokenDb, db: c.GetStateDB(), localDB: c.GetLocalDB(),
		txhash: hash, fromaddr: fromaddr, blocktime: c.GetBlockTime(), height: c.GetHeight(),
		execaddr: dapp.ExecAddress(string(tx.Execer)), difficulty: c.GetDifficulty(), index: index, Collateralize: c}
}

// GetCreateReceiptLog generate logs for Collateralize create action
func (action *Action) GetCreateReceiptLog(collateralize *pty.Collateralize) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogCollateralizeCreate

	c := &pty.ReceiptCollateralize{}
	c.CollateralizeId = collateralize.CollateralizeId
	c.Status = collateralize.Status
	c.AccountAddr = action.fromaddr

	log.Log = types.Encode(c)

	return log
}

// GetBorrowReceiptLog generate logs for Collateralize borrow action
func (action *Action) GetBorrowReceiptLog(collateralize *pty.Collateralize, record *pty.BorrowRecord) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogCollateralizeBorrow

	c := &pty.ReceiptCollateralize{}
	c.CollateralizeId = collateralize.CollateralizeId
	c.AccountAddr = action.fromaddr
	c.RecordId = record.RecordId
	c.Status = record.Status

	log.Log = types.Encode(c)

	return log
}

// GetRepayReceiptLog generate logs for Collateralize Repay action
func (action *Action) GetRepayReceiptLog(collateralize *pty.Collateralize, record *pty.BorrowRecord) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogCollateralizeRepay

	c := &pty.ReceiptCollateralize{}
	c.CollateralizeId = collateralize.CollateralizeId
	c.AccountAddr = action.fromaddr
	c.RecordId = record.RecordId
	c.Status = record.Status

	log.Log = types.Encode(c)

	return log
}

// GetAppendReceiptLog generate logs for Collateralize append action
func (action *Action) GetAppendReceiptLog(collateralize *pty.Collateralize, record *pty.BorrowRecord) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogCollateralizeAppend

	c := &pty.ReceiptCollateralize{}
	c.CollateralizeId = collateralize.CollateralizeId
	c.AccountAddr = action.fromaddr
	c.RecordId = record.RecordId
	c.Status = record.Status

	log.Log = types.Encode(c)

	return log
}

// GetFeedReceiptLog generate logs for Collateralize price feed action
func (action *Action) GetFeedReceiptLog(collateralize *pty.Collateralize, record *pty.BorrowRecord) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogCollateralizeFeed

	c := &pty.ReceiptCollateralize{}
	c.CollateralizeId = collateralize.CollateralizeId
	c.AccountAddr = record.AccountAddr
	c.RecordId = record.RecordId
	c.Status = record.Status

	log.Log = types.Encode(c)

	return log
}

// GetRetrieveReceiptLog generate logs for Collateralize close action
func (action *Action) GetRetrieveReceiptLog(collateralize *pty.Collateralize) *types.ReceiptLog {
	log := &types.ReceiptLog{}
	log.Ty = pty.TyLogCollateralizeRetrieve

	c := &pty.ReceiptCollateralize{}
	c.CollateralizeId = collateralize.CollateralizeId
	c.Status = collateralize.Status
	c.AccountAddr = action.fromaddr

	log.Log = types.Encode(c)

	return log
}

// GetIndex returns index in block
func (action *Action) GetIndex() int64 {
	return action.height*types.MaxTxsPerBlock + int64(action.index)
}

func getLatestLiquidationPrice(coll *pty.Collateralize) int64 {
	var latest int64
	for _, collRecord := range coll.BorrowRecords {
		if collRecord.LiquidationPrice > latest {
			latest = collRecord.LiquidationPrice
		}
	}

	return latest
}

func getLatestExpireTime(coll *pty.Collateralize) int64 {
	var latest int64 = 0x7fffffffffffffff

	for _, collRecord := range coll.BorrowRecords {
		if collRecord.ExpireTime < latest {
			latest = collRecord.ExpireTime
		}
	}

	return latest
}

// CollateralizeManage 
func (action *Action) CollateralizeManage(manage *pty.CollateralizeManage) (*types.Receipt, error) {
	var kv []*types.KeyValue
	var receipt *types.Receipt

	if !isRightAddr(issuanceE.ManageKey, action.fromaddr, action.db) {
		clog.Error("CollateralizeManage", "addr", action.fromaddr, "error", "Address has no permission to config")
		return nil, pty.ErrPermissionDeny
	}

	if manage.DebtCeiling < 0 || manage.LiquidationRatio < 0 || manage.LiquidationRatio >= 10000 ||
		manage.StabilityFeeRatio < 0 || manage.StabilityFeeRatio >= 10000 {
		return nil, pty.ErrRiskParam
	}

	manConfig, _ := getCollateralizeConfig(action.db)
	cfg := action.Collateralize.GetAPI().GetConfig()
	if manConfig == nil {
		manConfig = &pty.CollateralizeManage{
			DebtCeiling:       DefaultDebtCeiling * cfg.GetCoinPrecision(),
			LiquidationRatio:  DefaultLiquidationRatio,
			StabilityFeeRatio: DefaultStabilityFeeRation,
			Period:            DefaultPeriod,
			TotalBalance:      DefaultTotalBalance,
		}
	}

	collConfig := &pty.CollateralizeManage{}
	if manage.StabilityFeeRatio != 0 {
		collConfig.StabilityFeeRatio = manage.StabilityFeeRatio
	} else {
		collConfig.StabilityFeeRatio = manConfig.StabilityFeeRatio
	}

	if manage.Period != 0 {
		collConfig.Period = manage.Period
	} else {
		collConfig.Period = manConfig.Period
	}

	if manage.LiquidationRatio != 0 {
		collConfig.LiquidationRatio = manage.LiquidationRatio
	} else {
		collConfig.LiquidationRatio = manConfig.LiquidationRatio
	}

	if manage.DebtCeiling != 0 {
		collConfig.DebtCeiling = manage.DebtCeiling
	} else {
		collConfig.DebtCeiling = manConfig.DebtCeiling
	}

	if manage.TotalBalance != 0 {
		collConfig.TotalBalance = manage.TotalBalance
	} else {
		collConfig.TotalBalance = manConfig.TotalBalance
	}
	collConfig.CurrentTime = action.blocktime

	value := types.Encode(collConfig)
	action.db.Set(ConfigKey(), value)
	kv = append(kv, &types.KeyValue{Key: ConfigKey(), Value: value})

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: nil}
	return receipt, nil
}

func getCollateralizeConfig(db dbm.KV) (*pty.CollateralizeManage, error) {
	data, err := db.Get(ConfigKey())
	if err != nil {
		clog.Debug("getCollateralizeConfig", "error", err)
		return nil, err
	}

	var collCfg pty.CollateralizeManage
	err = types.Decode(data, &collCfg)
	if err != nil {
		clog.Debug("getCollateralizeConfig", "decode", err)
		return nil, err
	}
	return &collCfg, nil
}

func isSuperAddr(addr string, db dbm.KV) bool {
	data, err := db.Get(AddrKey())
	if err != nil {
		clog.Error("getSuperAddr", "error", err)
		return false
	}

	var item types.ConfigItem
	err = types.Decode(data, &item)
	if err != nil {
		clog.Error("isSuperAddr", "Decode", data)
		return false
	}

	for _, op := range item.GetArr().Value {
		if op == addr {
			return true
		}
	}

	return false
}

func getCollBalance(totalBalance int64, localdb dbm.KVDB, db dbm.KV) (int64, error) {
	collIDRecords, err := queryCollateralizeByStatus(localdb, pty.CollateralizeStatusCreated, "")
	if err != nil {
		clog.Debug("Query_CollateralizeByStatus", "get collateralize record error", err)
	}

	balance := totalBalance
	for _, id := range collIDRecords {
		coll, err := queryCollateralizeByID(db, id)
		if err != nil {
			clog.Error("Query_CollateralizeInfoByID", "id", id, "error", err)
			return 0, err
		}

		balance -= coll.TotalBalance
	}

	return balance, nil
}

func (action *Action) CollateralizeCreate(create *pty.CollateralizeCreate) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var receipt *types.Receipt

	if !isSuperAddr(action.fromaddr, action.db) {
		clog.Error("CollateralizeCreate", "error", "CollateralizeCreate need super address")
		return nil, pty.ErrPermissionDeny
	}

	if create.GetTotalBalance() <= 0 {
		clog.Error("CollateralizeCreate", "addr", action.fromaddr, "execaddr", action.execaddr, "total balance", create.GetTotalBalance(), "error", types.ErrAmount)
		return nil, types.ErrAmount
	}

	collcfg, err := getCollateralizeConfig(action.db)
	if err != nil {
		clog.Error("CollateralizeCreate.getCollateralizeConfig", "addr", action.fromaddr, "error", err)
		return nil, err
	}

	reBalance, err := getCollBalance(collcfg.TotalBalance, action.localDB, action.db)
	if err != nil {
		clog.Error("CollateralizeCreate.getCollBalance", "addr", action.fromaddr, "error", err)
		return nil, err
	}
	if reBalance < create.GetTotalBalance() {
		clog.Error("CollateralizeCreate.getCollBalance", "addr", action.fromaddr, "collBalance", reBalance, "create.balance", create.GetTotalBalance(), "error", pty.ErrCollateralizeLowBalance)
		return nil, pty.ErrCollateralizeLowBalance
	}

	if !action.CheckExecTokenAccount(action.fromaddr, create.TotalBalance, false) {
		clog.Error("CollateralizeCreate.CheckExecTokenAccount", "fromaddr", action.fromaddr, "balance", create.TotalBalance, "error", types.ErrInsufficientBalance)
		return nil, types.ErrInsufficientBalance
	}

	collateralizeIDs, err := queryCollateralizeByAddr(action.localDB, action.fromaddr, pty.CollateralizeStatusCreated, "")
	if err != nil && err != types.ErrNotFound {
		clog.Error("CollateralizeCreate.queryCollateralizeByAddr", "addr", action.fromaddr, "error", err)
		return nil, err
	}

	receipt, err = action.tokenAccount.ExecFrozen(action.fromaddr, action.execaddr, create.TotalBalance)
	if err != nil {
		clog.Error("CollateralizeCreate.Frozen", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", create.TotalBalance)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	var collateralizeID string
	coll := &CollateralizeDB{}
	if collateralizeIDs == nil {
		collateralizeID = common.ToHex(action.txhash)

		coll.CollateralizeId = collateralizeID
		coll.LiquidationRatio = collcfg.LiquidationRatio
		coll.TotalBalance = create.TotalBalance
		coll.DebtCeiling = collcfg.DebtCeiling
		coll.StabilityFeeRatio = collcfg.StabilityFeeRatio
		coll.Period = collcfg.Period
		coll.Balance = create.TotalBalance
		coll.CreateAddr = action.fromaddr
		coll.Status = pty.CollateralizeActionCreate
		coll.CollBalance = 0
	} else {
		collateralize, err := queryCollateralizeByID(action.db, collateralizeIDs[0])
		if err != nil {
			clog.Error("CollateralizeCreate.queryCollateralizeByID", "addr", action.fromaddr, "execaddr", action.execaddr, "collId", collateralizeIDs[0])
			return nil, err
		}
		coll.Collateralize = *collateralize
		coll.TotalBalance += create.TotalBalance
		coll.Balance += create.TotalBalance
		coll.PreStatus = coll.Status
	}
	clog.Debug("CollateralizeCreate created", "CollateralizeID", collateralizeID, "TotalBalance", create.TotalBalance)

	coll.Save(action.db)
	kv = append(kv, coll.GetKVSet()...)

	receiptLog := action.GetCreateReceiptLog(&coll.Collateralize)
	logs = append(logs, receiptLog)

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func getBtyNumToFrozen(value int64, price int64, ratio int64) (int64, error) {
	if price == 0 {
		clog.Error("Bty price should greate to 0")
		return 0, pty.ErrPriceInvalid
	}

	btyValue := (value * 1e4) / (price * ratio)
	return btyValue * 1e4, nil
}

func calcLiquidationPrice(value int64, colValue int64) int64 {
	return (value * pty.CollateralizePreLiquidationRatio) / colValue
}

func getLatestPrice(db dbm.KV) (int64, error) {
	data, err := db.Get(PriceKey())
	if err != nil {
		clog.Error("getLatestPrice", "get", err)
		return -1, err
	}
	var price pty.AssetPriceRecord
	//decode
	err = types.Decode(data, &price)
	if err != nil {
		clog.Error("getLatestPrice", "decode", err)
		return -1, err
	}

	return price.BtyPrice, nil
}

func (action *Action) CheckExecAccountBalance(fromAddr string, ToFrozen, ToActive int64) bool {
	acc := action.coinsAccount.LoadExecAccount(fromAddr, action.execaddr)
	if acc.GetBalance() >= ToFrozen && acc.GetFrozen() >= ToActive {
		return true
	}
	return false
}

func (action *Action) CheckExecTokenAccount(addr string, amount int64, isFrozen bool) bool {
	acc := action.tokenAccount.LoadExecAccount(addr, action.execaddr)
	if isFrozen {
		if acc.GetFrozen() >= amount {
			return true
		}
	} else {
		if acc.GetBalance() >= amount {
			return true
		}
	}
	return false
}

func (action *Action) CollateralizeBorrow(borrow *pty.CollateralizeBorrow) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	collateralize, err := queryCollateralizeByID(action.db, borrow.CollateralizeId)
	if err != nil {
		clog.Error("CollateralizeBorrow.queryCollateralizeByID", "CollateralizeId", borrow.CollateralizeId, "error", err)
		return nil, err
	}

	if collateralize.Status == pty.CollateralizeStatusClose {
		clog.Error("CollateralizeBorrow", "CollID", collateralize.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "status", collateralize.Status, "error", pty.ErrCollateralizeStatus)
		return nil, pty.ErrCollateralizeStatus
	}

	coll := &CollateralizeDB{*collateralize}

	if borrow.GetValue() <= 0 {
		clog.Error("CollateralizeBorrow", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "borrow value", borrow.GetValue(), "error", types.ErrInvalidParam)
		return nil, types.ErrAmount
	}

	userBalance, _ := queryCollateralizeUserBalance(action.db, action.localDB, action.fromaddr)
	if borrow.GetValue()+userBalance > coll.DebtCeiling {
		clog.Error("CollateralizeBorrow", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr,
			"borrow value", borrow.GetValue(), "current balance", userBalance, "error", pty.ErrCollateralizeExceedDebtCeiling)
		return nil, pty.ErrCollateralizeExceedDebtCeiling
	}

	if borrow.GetValue() > coll.Balance {
		clog.Error("CollateralizeBorrow", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "borrow value", borrow.GetValue(), "error", pty.ErrCollateralizeLowBalance)
		return nil, pty.ErrCollateralizeLowBalance
	}
	clog.Debug("CollateralizeBorrow", "value", borrow.GetValue())

	lastPrice, err := getLatestPrice(action.db)
	if err != nil {
		clog.Error("CollateralizeBorrow.getLatestPrice", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "error", err)
		return nil, err
	}

	btyFrozen, err := getBtyNumToFrozen(borrow.GetValue(), lastPrice, coll.LiquidationRatio)
	if err != nil {
		clog.Error("CollateralizeBorrow.getBtyNumToFrozen", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "error", err)
		return nil, err
	}

	if !action.CheckExecAccountBalance(action.fromaddr, btyFrozen, 0) {
		clog.Error("CollateralizeBorrow.CheckExecAccountBalance", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "balance", btyFrozen, "error", types.ErrNoBalance)
		return nil, types.ErrNoBalance
	}

	receipt, err := action.coinsAccount.ExecTransfer(action.fromaddr, coll.CreateAddr, action.execaddr, btyFrozen)
	if err != nil {
		clog.Error("CollateralizeBorrow.ExecTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", btyFrozen)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	receipt, err = action.coinsAccount.ExecFrozen(coll.CreateAddr, action.execaddr, btyFrozen)
	if err != nil {
		clog.Error("CollateralizeBorrow.Frozen", "addr", coll.CreateAddr, "execaddr", action.execaddr, "amount", btyFrozen)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	receipt, err = action.tokenAccount.ExecTransferFrozen(coll.CreateAddr, action.fromaddr, action.execaddr, borrow.GetValue())
	if err != nil {
		clog.Error("CollateralizeBorrow.ExecTokenTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", borrow.GetValue())
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	borrowRecord := &pty.BorrowRecord{}
	borrowRecord.RecordId = common.ToHex(action.txhash)
	borrowRecord.CollateralizeId = coll.CollateralizeId
	borrowRecord.AccountAddr = action.fromaddr
	borrowRecord.CollateralValue = btyFrozen
	borrowRecord.StartTime = action.blocktime
	borrowRecord.CollateralPrice = lastPrice
	borrowRecord.DebtValue = borrow.GetValue()
	borrowRecord.LiquidationPrice = (coll.LiquidationRatio * lastPrice * pty.CollateralizePreLiquidationRatio) / 1e8
	borrowRecord.Status = pty.CollateralizeUserStatusCreate
	borrowRecord.ExpireTime = action.blocktime + coll.Period

	if coll.LatestLiquidationPrice < borrowRecord.LiquidationPrice {
		coll.LatestLiquidationPrice = borrowRecord.LiquidationPrice
	}

	coll.BorrowRecords = append(coll.BorrowRecords, borrowRecord)
	coll.Status = pty.CollateralizeStatusCreated
	coll.Balance -= borrow.GetValue()
	coll.CollBalance += btyFrozen
	coll.LatestExpireTime = getLatestExpireTime(&coll.Collateralize)
	coll.Save(action.db)
	kv = append(kv, coll.GetKVSet()...)

	receiptLog := action.GetBorrowReceiptLog(&coll.Collateralize, borrowRecord)
	logs = append(logs, receiptLog)

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (action *Action) CollateralizeRepay(repay *pty.CollateralizeRepay) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var receipt *types.Receipt

	collateralize, err := queryCollateralizeByID(action.db, repay.CollateralizeId)
	if err != nil {
		clog.Error("CollateralizeRepay", "CollID", repay.CollateralizeId, "error", err)
		return nil, err
	}

	coll := &CollateralizeDB{*collateralize}

	if coll.Status != pty.CollateralizeStatusCreated {
		clog.Error("CollateralizeRepay", "CollID", repay.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "error", "status error", "Status", coll.Status)
		return nil, pty.ErrCollateralizeStatus
	}

	var borrowRecord *pty.BorrowRecord
	var index int
	for i, record := range coll.BorrowRecords {
		if record.RecordId == repay.RecordId {
			borrowRecord = record
			index = i
			break
		}
	}

	if borrowRecord == nil {
		clog.Error("CollateralizeRepay", "CollID", repay.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "error", "Can not find borrow record")
		return nil, pty.ErrRecordNotExist
	}

	fee := ((borrowRecord.DebtValue * coll.StabilityFeeRatio) / 1e8) * 1e4
	realRepay := borrowRecord.DebtValue + fee

	if !action.CheckExecTokenAccount(action.fromaddr, realRepay, false) {
		clog.Error("CollateralizeRepay.CheckExecTokenAccount", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "amount", realRepay, "error", types.ErrInsufficientBalance)
		return nil, types.ErrNoBalance
	}

	receipt, err = action.tokenAccount.ExecTransfer(action.fromaddr, coll.CreateAddr, action.execaddr, realRepay)
	if err != nil {
		clog.Error("CollateralizeRepay.ExecTokenTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", realRepay)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	receipt, err = action.tokenAccount.ExecFrozen(coll.CreateAddr, action.execaddr, borrowRecord.DebtValue)
	if err != nil {
		clog.Error("CollateralizeCreate.Frozen", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", borrowRecord.DebtValue)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	receipt, err = action.coinsAccount.ExecTransferFrozen(coll.CreateAddr, action.fromaddr, action.execaddr, borrowRecord.CollateralValue)
	if err != nil {
		clog.Error("CollateralizeRepay.ExecTransferFrozen", "addr", coll.CreateAddr, "execaddr", action.execaddr, "amount", borrowRecord.CollateralValue)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	borrowRecord.PreStatus = borrowRecord.Status
	borrowRecord.Status = pty.CollateralizeUserStatusClose

	coll.Balance += borrowRecord.DebtValue
	coll.CollBalance -= borrowRecord.CollateralValue
	coll.BorrowRecords = append(coll.BorrowRecords[:index], coll.BorrowRecords[index+1:]...)
	coll.InvalidRecords = append(coll.InvalidRecords, borrowRecord)
	coll.LatestLiquidationPrice = getLatestLiquidationPrice(&coll.Collateralize)
	coll.LatestExpireTime = getLatestExpireTime(&coll.Collateralize)
	coll.Save(action.db)
	kv = append(kv, coll.GetKVSet()...)

	receiptLog := action.GetRepayReceiptLog(&coll.Collateralize, borrowRecord)
	logs = append(logs, receiptLog)

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (action *Action) CollateralizeAppend(cAppend *pty.CollateralizeAppend) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if cAppend.GetCollateralValue() <= 0 {
		clog.Error("CollateralizeAppend", "addr", action.fromaddr, "execaddr", action.execaddr, "append value", cAppend.GetCollateralValue(), "error", types.ErrAmount)
		return nil, types.ErrAmount
	}

	collateralize, err := queryCollateralizeByID(action.db, cAppend.CollateralizeId)
	if err != nil {
		clog.Error("CollateralizeAppend", "CollateralizeId", cAppend.CollateralizeId, "error", err)
		return nil, err
	}

	coll := &CollateralizeDB{*collateralize}

	if coll.Status != pty.CollateralizeStatusCreated {
		clog.Error("CollateralizeAppend", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "status", coll.Status, "error", pty.ErrCollateralizeStatus)
		return nil, pty.ErrCollateralizeStatus
	}

	var borrowRecord *pty.BorrowRecord
	for _, record := range coll.BorrowRecords {
		if record.RecordId == cAppend.RecordId {
			borrowRecord = record
		}
	}

	if borrowRecord == nil {
		clog.Error("CollateralizeAppend", "CollID", cAppend.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "error", "Can not find borrow record")
		return nil, pty.ErrRecordNotExist
	}

	clog.Debug("CollateralizeAppend", "value", cAppend.CollateralValue)

	lastPrice, err := getLatestPrice(action.db)
	if err != nil {
		clog.Error("CollateralizeBorrow", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "error", err)
		return nil, err
	}

	if !action.CheckExecAccountBalance(action.fromaddr, cAppend.CollateralValue, 0) {
		clog.Error("CollateralizeBorrow.CheckExecAccountBalance", "CollID", coll.CollateralizeId, "addr", action.fromaddr, "execaddr", action.execaddr, "amount", cAppend.CollateralValue, "error", types.ErrNoBalance)
		return nil, types.ErrNoBalance
	}

	receipt, err := action.coinsAccount.ExecTransfer(action.fromaddr, coll.CreateAddr, action.execaddr, cAppend.CollateralValue)
	if err != nil {
		clog.Error("CollateralizeBorrow.ExecTransfer", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", cAppend.CollateralValue, "error", err)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	receipt, err = action.coinsAccount.ExecFrozen(coll.CreateAddr, action.execaddr, cAppend.CollateralValue)
	if err != nil {
		clog.Error("CollateralizeBorrow.Frozen", "addr", coll.CreateAddr, "execaddr", action.execaddr, "amount", cAppend.CollateralValue, "error", err)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	borrowRecord.CollateralValue += cAppend.CollateralValue
	borrowRecord.CollateralPrice = lastPrice
	borrowRecord.LiquidationPrice = calcLiquidationPrice(borrowRecord.DebtValue, borrowRecord.CollateralValue)
	if borrowRecord.LiquidationPrice*PriceWarningRate < lastPrice {

		if borrowRecord.Status == pty.CollateralizeUserStatusWarning {
			borrowRecord.PreStatus = borrowRecord.Status
			borrowRecord.Status = pty.CollateralizeStatusCreated
		}
	}

	coll.CollBalance += cAppend.CollateralValue
	coll.LatestLiquidationPrice = getLatestLiquidationPrice(&coll.Collateralize)
	coll.LatestExpireTime = getLatestExpireTime(&coll.Collateralize)

	coll.Save(action.db)

	kv = append(kv, coll.GetKVSet()...)

	receiptLog := action.GetAppendReceiptLog(&coll.Collateralize, borrowRecord)
	logs = append(logs, receiptLog)

	receipt = &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func getManageKey(key string, db dbm.KV) ([]byte, error) {
	manageKey := types.ManageKey(key)
	value, err := db.Get([]byte(manageKey))
	if err != nil {
		return nil, err
	}
	return value, nil
}

func isRightAddr(key string, addr string, db dbm.KV) bool {
	value, err := getManageKey(key, db)
	if err != nil {
		clog.Error("isRightAddr", "Key", key)
		return false
	}
	if value == nil {
		clog.Error("isRightAddr", "key", key, "error", "Found key nil value")
		return false
	}

	var item types.ConfigItem
	err = types.Decode(value, &item)
	if err != nil {
		clog.Error("isRightAddr", "Decode", value)
		return false
	}

	for _, op := range item.GetArr().Value {
		if op == addr {
			return true
		}
	}
	return false

}

func getGuarantorAddr(db dbm.KV) (string, error) {
	value, err := getManageKey(issuanceE.GuarantorKey, db)
	if err != nil {
		clog.Error("CollateralizePriceFeed", "getGuarantorAddr", err)
		return "", err
	}
	if value == nil {
		clog.Error("CollateralizePriceFeed guarantorKey found nil value")
		return "", err
	}

	var item types.ConfigItem
	err = types.Decode(value, &item)
	if err != nil {
		clog.Error("CollateralizePriceFeed", "getGuarantorAddr", err)
		return "", err
	}

	return item.GetArr().Value[0], nil
}

func removeLiquidateRecord(borrowRecords []*pty.BorrowRecord, remove pty.BorrowRecord) []*pty.BorrowRecord {
	for index, record := range borrowRecords {
		if record.RecordId == remove.RecordId {
			borrowRecords = append(borrowRecords[:index], borrowRecords[index+1:]...)
			break
		}
	}

	return borrowRecords
}

func (action *Action) systemLiquidation(coll *pty.Collateralize, price int64) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var removeRecord []*pty.BorrowRecord

	for _, borrowRecord := range coll.BorrowRecords {
		if (borrowRecord.LiquidationPrice*PriceWarningRate)/1e4 < price {

			if borrowRecord.Status == pty.CollateralizeUserStatusWarning {
				borrowRecord.PreStatus = borrowRecord.Status
				borrowRecord.Status = pty.CollateralizeUserStatusCreate
				log := action.GetFeedReceiptLog(coll, borrowRecord)
				logs = append(logs, log)
			}
			continue
		}

		if borrowRecord.LiquidationPrice >= price {

			clog.Debug("systemLiquidation", "coll id", borrowRecord.CollateralizeId, "record id", borrowRecord.RecordId, "account", borrowRecord.AccountAddr, "price", price)

			getGuarantorAddr, err := getGuarantorAddr(action.db)
			if err != nil {
				if err != nil {
					clog.Error("systemLiquidation", "getGuarantorAddr", err)
					continue
				}
			}

			receipt, err := action.coinsAccount.ExecTransferFrozen(coll.CreateAddr, getGuarantorAddr, action.execaddr, borrowRecord.CollateralValue)
			if err != nil {
				clog.Error("systemLiquidation", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", borrowRecord.CollateralValue, "error", err)
				continue
			}
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)

			borrowRecord.LiquidateTime = action.blocktime
			borrowRecord.PreStatus = borrowRecord.Status
			borrowRecord.Status = pty.CollateralizeUserStatusSystemLiquidate
			coll.InvalidRecords = append(coll.InvalidRecords, borrowRecord)
			removeRecord = append(removeRecord, borrowRecord)
			coll.CollBalance -= borrowRecord.CollateralValue

			log := action.GetFeedReceiptLog(coll, borrowRecord)
			logs = append(logs, log)
			continue
		}


		if borrowRecord.Status != pty.CollateralizeUserStatusWarning {
			borrowRecord.PreStatus = borrowRecord.Status
			borrowRecord.Status = pty.CollateralizeUserStatusWarning
			log := action.GetFeedReceiptLog(coll, borrowRecord)
			logs = append(logs, log)
		}
	}

	for _, record := range removeRecord {
		coll.BorrowRecords = removeLiquidateRecord(coll.BorrowRecords, *record)
	}


	coll.LatestLiquidationPrice = getLatestLiquidationPrice(coll)
	coll.LatestExpireTime = getLatestExpireTime(coll)
	collDB := &CollateralizeDB{*coll}
	collDB.Save(action.db)
	kv = append(kv, collDB.GetKVSet()...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (action *Action) expireLiquidation(coll *pty.Collateralize) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var removeRecord []*pty.BorrowRecord

	for _, borrowRecord := range coll.BorrowRecords {
		if borrowRecord.ExpireTime-ExpireWarningTime > action.blocktime {
			continue
		}

		if borrowRecord.ExpireTime <= action.blocktime {
			clog.Debug("expireLiquidation", "coll id", borrowRecord.CollateralizeId, "record id", borrowRecord.RecordId, "account", borrowRecord.AccountAddr, "time", action.blocktime)

			getGuarantorAddr, err := getGuarantorAddr(action.db)
			if err != nil {
				if err != nil {
					clog.Error("expireLiquidation", "getGuarantorAddr", err)
					continue
				}
			}

			receipt, err := action.coinsAccount.ExecTransferFrozen(coll.CreateAddr, getGuarantorAddr, action.execaddr, borrowRecord.CollateralValue)
			if err != nil {
				clog.Error("expireLiquidation", "addr", action.fromaddr, "execaddr", action.execaddr, "amount", borrowRecord.CollateralValue, "error", err)
				continue
			}
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)

			borrowRecord.LiquidateTime = action.blocktime
			borrowRecord.PreStatus = borrowRecord.Status
			borrowRecord.Status = pty.CollateralizeUserStatusExpireLiquidate
			coll.InvalidRecords = append(coll.InvalidRecords, borrowRecord)
			removeRecord = append(removeRecord, borrowRecord)
			coll.CollBalance -= borrowRecord.CollateralValue

			log := action.GetFeedReceiptLog(coll, borrowRecord)
			logs = append(logs, log)
			continue
		}

		if borrowRecord.Status != pty.CollateralizeUserStatusExpire {
			borrowRecord.PreStatus = borrowRecord.Status
			borrowRecord.Status = pty.CollateralizeUserStatusExpire
			log := action.GetFeedReceiptLog(coll, borrowRecord)
			logs = append(logs, log)
		}
	}

	for _, record := range removeRecord {
		coll.BorrowRecords = removeLiquidateRecord(coll.BorrowRecords, *record)
	}

	coll.LatestLiquidationPrice = getLatestLiquidationPrice(coll)
	coll.LatestExpireTime = getLatestExpireTime(coll)
	collDB := &CollateralizeDB{*coll}
	collDB.Save(action.db)
	kv = append(kv, collDB.GetKVSet()...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func pricePolicy(feed *pty.CollateralizeFeed) int64 {
	var totalPrice int64
	var totalVolume int64
	for _, volume := range feed.Volume {
		totalVolume += volume
	}

	if totalVolume == 0 {
		clog.Error("collateralize price feed volume empty")
		return 0
	}

	for i, price := range feed.Price {
		totalPrice += (price * feed.Volume[i]) / totalVolume
	}

	return totalPrice
}

func (action *Action) CollateralizeFeed(feed *pty.CollateralizeFeed) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if feed == nil || len(feed.Price) == 0 || len(feed.Price) != len(feed.Volume) {
		clog.Error("CollateralizePriceFeed", types.ErrInvalidParam)
		return nil, types.ErrInvalidParam
	}

	if !isRightAddr(issuanceE.PriceFeedKey, action.fromaddr, action.db) {
		clog.Error("CollateralizePriceFeed", "addr", action.fromaddr, "error", "Address has no permission to feed price")
		return nil, pty.ErrPermissionDeny
	}

	price := pricePolicy(feed)
	if price <= 0 {
		clog.Error("CollateralizePriceFeed", "price", price, "error", pty.ErrPriceInvalid)
		return nil, pty.ErrPriceInvalid
	}

	ids, err := queryCollateralizeByStatus(action.localDB, pty.CollateralizeStatusCreated, "")
	if err != nil {
		clog.Debug("CollateralizePriceFeed", "get collateralize record error", err)
	}

	for _, collID := range ids {
		coll, err := queryCollateralizeByID(action.db, collID)
		if err != nil {
			clog.Error("CollateralizePriceFeed", "Collateralize ID", coll.CollateralizeId, "get collateralize record by id error", err)
			continue
		}

		if coll.LatestExpireTime-ExpireWarningTime <= action.blocktime {
			receipt, err := action.expireLiquidation(coll)
			if err != nil {
				clog.Error("CollateralizePriceFeed", "Collateralize ID", coll.CollateralizeId, "expire liquidation error", err)
				continue
			}
			logs = append(logs, receipt.Logs...)
			kv = append(kv, receipt.KV...)
		}

		receipt, err := action.systemLiquidation(coll, price)
		if err != nil {
			clog.Error("CollateralizePriceFeed", "Collateralize ID", coll.CollateralizeId, "system liquidation error", err)
			continue
		}
		logs = append(logs, receipt.Logs...)
		kv = append(kv, receipt.KV...)
	}

	var priceRecord pty.AssetPriceRecord
	priceRecord.BtyPrice = price
	priceRecord.RecordTime = action.blocktime

	pricekv := &types.KeyValue{Key: PriceKey(), Value: types.Encode(&priceRecord)}
	action.db.Set(pricekv.Key, pricekv.Value)
	kv = append(kv, pricekv)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

func (action *Action) CollateralizeRetrieve(retrieve *pty.CollateralizeRetrieve) (*types.Receipt, error) {
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var receipt *types.Receipt

	collateralize, err := queryCollateralizeByID(action.db, retrieve.CollateralizeId)
	if err != nil {
		clog.Error("CollateralizeRetrieve", "CollateralizeId", retrieve.CollateralizeId, "error", err)
		return nil, err
	}

	if action.fromaddr != collateralize.CreateAddr {
		clog.Error("CollateralizeRetrieve", "CollateralizeId", retrieve.CollateralizeId, "error", "account error", "create", collateralize.CreateAddr, "from", action.fromaddr)
		return nil, pty.ErrPermissionDeny
	}

	if retrieve.Balance > collateralize.Balance {
		clog.Error("CollateralizeRetrieve", "CollateralizeId", retrieve.CollateralizeId, "error", "balance error", "retrieve balance", retrieve.Balance, "available balance", collateralize.Balance)
		return nil, types.ErrAmount
	}

	receipt, err = action.tokenAccount.ExecActive(action.fromaddr, action.execaddr, retrieve.Balance)
	if err != nil {
		clog.Error("IssuanceClose.ExecActive", "addr", action.fromaddr, "execaddr", action.execaddr, "balance", retrieve.Balance)
		return nil, err
	}
	logs = append(logs, receipt.Logs...)
	kv = append(kv, receipt.KV...)

	clog.Debug("CollateralizeRetrieve", "ID", retrieve.CollateralizeId, "balance", retrieve.Balance)

	coll := &CollateralizeDB{*collateralize}
	coll.TotalBalance -= retrieve.Balance
	coll.Balance -= retrieve.Balance
	coll.PreStatus = coll.Status
	if coll.TotalBalance == 0 {
		coll.Status = pty.CollateralizeStatusClose
	}
	coll.Save(action.db)
	kv = append(kv, coll.GetKVSet()...)

	receiptLog := action.GetRetrieveReceiptLog(&coll.Collateralize)
	logs = append(logs, receiptLog)

	return &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}, nil
}

func queryCollateralizeByID(db dbm.KV, collateralizeID string) (*pty.Collateralize, error) {
	data, err := db.Get(Key(collateralizeID))
	if err != nil {
		clog.Debug("queryCollateralizeByID", "error", err)
		return nil, err
	}

	var coll pty.Collateralize
	err = types.Decode(data, &coll)
	if err != nil {
		clog.Debug("queryCollateralizeByID", "decode", err)
		return nil, err
	}
	return &coll, nil
}

func queryCollateralizeByStatus(localdb dbm.KVDB, status int32, collID string) ([]string, error) {
	query := pty.NewCollateralizeTable(localdb).GetQuery(localdb)
	var primary []byte
	if len(collID) > 0 {
		primary = []byte(collID)
	}

	var data = &pty.ReceiptCollateralize{
		CollateralizeId: collID,
		Status:          status,
	}
	rows, err := query.List("status", data, primary, DefaultCount, ListDESC)
	if err != nil {
		clog.Debug("queryCollateralizeByStatus.List", "error", err)
		return nil, err
	}

	var ids []string
	for _, row := range rows {
		ids = append(ids, string(row.Primary))
	}

	return ids, nil
}

func queryCollateralizeByAddr(localdb dbm.KVDB, addr string, status int32, collID string) ([]string, error) {
	query := pty.NewCollateralizeTable(localdb).GetQuery(localdb)
	var primary []byte
	if len(collID) > 0 {
		primary = []byte(collID)
	}

	var data = &pty.ReceiptCollateralize{
		CollateralizeId: collID,
		Status:          status,
		AccountAddr:     addr,
	}
	var rows []*table.Row
	var err error
	if status == 0 {
		rows, err = query.List("addr", data, primary, DefaultCount, ListDESC)
		if err != nil {
			clog.Debug("queryCollateralizeByAddr.List", "index", "addr", "error", err)
			return nil, err
		}
	} else {
		rows, err = query.List("addr_status", data, primary, DefaultCount, ListDESC)
		if err != nil {
			clog.Debug("queryCollateralizeByAddr.List", "index", "addr_status", "error", err)
			return nil, err
		}
	}
	var ids []string
	for _, row := range rows {
		ids = append(ids, string(row.Primary))
	}

	return ids, nil
}

func queryCollateralizeRecordByID(db dbm.KV, collateralizeID string, recordID string) (*pty.BorrowRecord, error) {
	coll, err := queryCollateralizeByID(db, collateralizeID)
	if err != nil {
		clog.Debug("queryCollateralizeRecordByID", "error", err)
		return nil, err
	}

	for _, record := range coll.BorrowRecords {
		if record.RecordId == recordID {
			return record, nil
		}
	}

	for _, record := range coll.InvalidRecords {
		if record.RecordId == recordID {
			return record, nil
		}
	}

	return nil, types.ErrNotFound
}

func queryCollateralizeRecordByAddr(db dbm.KV, localdb dbm.KVDB, addr string, status int32, collID string, recordID string) ([]*pty.BorrowRecord, error) {
	query := pty.NewRecordTable(localdb).GetQuery(localdb)
	var primary []byte
	if len(recordID) > 0 {
		primary = []byte(recordID)
	}

	var data = &pty.ReceiptCollateralize{
		AccountAddr:     addr,
		Status:          status,
		CollateralizeId: collID,
	}

	var rows []*table.Row
	var err error
	if len(collID) != 0 {
		rows, err = query.List("id_addr", data, primary, DefaultCount, ListDESC)
		if err != nil {
			clog.Debug("queryCollateralizeRecordByAddr.List", "index", "id_addr", "error", err)
			return nil, err
		}
	} else if status != 0 {
		rows, err = query.List("addr_status", data, primary, DefaultCount, ListDESC)
		if err != nil {
			clog.Debug("queryCollateralizeRecordByAddr.List", "index", "addr_status", "error", err)
			return nil, err
		}
	} else {
		rows, err = query.List("addr", data, primary, DefaultCount, ListDESC)
		if err != nil {
			clog.Debug("queryCollateralizeRecordByAddr.List", "index", "addr", "error", err)
			return nil, err
		}
	}

	var records []*pty.BorrowRecord
	for _, row := range rows {
		record, err := queryCollateralizeRecordByID(db, row.Data.(*pty.ReceiptCollateralize).CollateralizeId, row.Data.(*pty.ReceiptCollateralize).RecordId)
		if err != nil {
			clog.Debug("queryCollateralizeRecordByAddr.queryCollateralizeRecordByID", "error", err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

func queryCollateralizeRecordByStatus(db dbm.KV, localdb dbm.KVDB, status int32, collID string, recordID string) ([]*pty.BorrowRecord, error) {
	query := pty.NewRecordTable(localdb).GetQuery(localdb)
	var primary []byte
	if len(recordID) > 0 {
		primary = []byte(recordID)
	}

	var data = &pty.ReceiptCollateralize{
		Status:          status,
		CollateralizeId: collID,
	}

	var rows []*table.Row
	var err error
	if len(collID) == 0 {
		rows, err = query.List("status", data, primary, DefaultCount, ListDESC)
		if err != nil {
			clog.Debug("queryCollateralizeRecordByStatus.List", "index", "status", "error", err)
			return nil, err
		}
	} else {
		rows, err = query.List("id_status", data, primary, DefaultCount, ListDESC)
		if err != nil {
			clog.Debug("queryCollateralizeRecordByStatus.List", "index", "id_status", "error", err)
			return nil, err
		}
	}

	var records []*pty.BorrowRecord
	for _, row := range rows {
		record, err := queryCollateralizeRecordByID(db, row.Data.(*pty.ReceiptCollateralize).CollateralizeId, row.Data.(*pty.ReceiptCollateralize).RecordId)
		if err != nil {
			clog.Debug("queryCollateralizeRecordByStatus.queryCollateralizeRecordByID", "error", err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

func queryCollateralizeUserBalanceStatus(db dbm.KV, localdb dbm.KVDB, addr string, status int32) (int64, error) {
	var totalBalance int64
	query := pty.NewRecordTable(localdb).GetQuery(localdb)
	var primary []byte
	var data = &pty.ReceiptCollateralize{
		AccountAddr: addr,
		Status:      status,
	}

	var rows []*table.Row
	var err error

	for {
		rows, err = query.List("addr_status", data, primary, DefaultCount, ListDESC)
		if err != nil {
			return -1, err
		}

		for _, row := range rows {
			record, err := queryCollateralizeRecordByID(db, row.Data.(*pty.ReceiptCollateralize).CollateralizeId, row.Data.(*pty.ReceiptCollateralize).RecordId)
			if err != nil {
				continue
			}
			totalBalance += record.DebtValue
		}

		if len(rows) < int(DefaultCount) {
			break
		}
		primary = []byte(rows[DefaultCount-1].Data.(*pty.ReceiptCollateralize).RecordId)
	}

	return totalBalance, nil
}

func queryCollateralizeUserBalance(db dbm.KV, localdb dbm.KVDB, addr string) (int64, error) {
	var totalBalance int64

	balance, err := queryCollateralizeUserBalanceStatus(db, localdb, addr, pty.CollateralizeUserStatusCreate)
	if err != nil {
		if err != types.ErrNotFound {
			clog.Error("queryCollateralizeUserBalance", "err", err)
		}
	} else {
		totalBalance += balance
	}

	balance, err = queryCollateralizeUserBalanceStatus(db, localdb, addr, pty.CollateralizeUserStatusWarning)
	if err != nil {
		if err != types.ErrNotFound {
			clog.Error("queryCollateralizeUserBalance", "err", err)
		}
	} else {
		totalBalance += balance
	}

	balance, err = queryCollateralizeUserBalanceStatus(db, localdb, addr, pty.CollateralizeUserStatusExpire)
	if err != nil {
		if err != types.ErrNotFound {
			clog.Error("queryCollateralizeUserBalance", "err", err)
		}
	} else {
		totalBalance += balance
	}

	return totalBalance, nil
}
