// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"bytes"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/33cn/chain33/system/dapp"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/types"
	wcom "github.com/33cn/chain33/wallet/common"
	privacy "github.com/33cn/plugin/plugin/dapp/privacy/crypto"
	privacytypes "github.com/33cn/plugin/plugin/dapp/privacy/types"
)

func (policy *privacyPolicy) rescanAllTxAddToUpdateUTXOs() {
	accounts, err := policy.getWalletOperate().GetWalletAccounts()
	if err != nil {
		bizlog.Debug("rescanAllTxToUpdateUTXOs", "walletOperate.GetWalletAccounts error", err)
		return
	}
	bizlog.Debug("rescanAllTxToUpdateUTXOs begin!")
	for _, acc := range accounts {
		/ blockchai Account.Add 
		policy.rescanwg.Add(1)
		go policy.rescanReqTxDetailByAddr(acc.Addr, policy.rescanwg)
	}
	policy.rescanwg.Wait()
	bizlog.Debug("rescanAllTxToUpdateUTXOs success!")
}

/ blockchai add 
func (policy *privacyPolicy) rescanReqTxDetailByAddr(addr string, wg *sync.WaitGroup) {
	defer wg.Done()
	policy.reqTxDetailByAddr(addr)
}

/ blockchai add 
func (policy *privacyPolicy) reqTxDetailByAddr(addr string) {
	if len(addr) == 0 {
		bizlog.Error("reqTxDetailByAddr input addr is nil!")
		return
	}
	var txInfo types.ReplyTxInfo

	i := 0
	operater := policy.getWalletOperate()
	for {
		/ blockchai hash  
		var ReqAddr types.ReqAddr
		ReqAddr.Addr = addr
		ReqAddr.Flag = 0
		ReqAddr.Direction = 0
		ReqAddr.Count = int32(MaxTxHashsPerTime)
		if i == 0 {
			ReqAddr.Height = -1
			ReqAddr.Index = 0
		} else {
			ReqAddr.Height = txInfo.GetHeight()
			ReqAddr.Index = txInfo.GetIndex()
		}
		i++
		ReplyTxInfos, err := operater.GetAPI().GetTransactionByAddr(&ReqAddr)
		if err != nil {
			bizlog.Error("reqTxDetailByAddr", "GetTransactionByAddr error", err, "addr", addr)
			return
		}
		if ReplyTxInfos == nil {
			bizlog.Info("reqTxDetailByAddr ReplyTxInfos is nil")
			return
		}
		txcount := len(ReplyTxInfos.TxInfos)
		var ReqHashes types.ReqHashes
		ReqHashes.Hashes = make([][]byte, len(ReplyTxInfos.TxInfos))
		for index, ReplyTxInfo := range ReplyTxInfos.TxInfos {
			ReqHashes.Hashes[index] = ReplyTxInfo.GetHash()
			txInfo.Hash = ReplyTxInfo.GetHash()
			txInfo.Height = ReplyTxInfo.GetHeight()
			txInfo.Index = ReplyTxInfo.GetIndex()
		}
		operater.GetTxDetailByHashs(&ReqHashes)
		if txcount < int(MaxTxHashsPerTime) {
			return
		}
	}
}

func (policy *privacyPolicy) isRescanUtxosFlagScaning() (bool, error) {
	if privacytypes.UtxoFlagScaning == policy.GetRescanFlag() {
		return true, privacytypes.ErrRescanFlagScaning
	}
	return false, nil
}

func (policy *privacyPolicy) parseViewSpendPubKeyPair(in string) (viewPubKey, spendPubKey []byte, err error) {
	src, err := common.FromHex(in)
	if err != nil {
		return nil, nil, err
	}
	if 64 != len(src) {
		bizlog.Error("parseViewSpendPubKeyPair", "pair with len", len(src))
		return nil, nil, types.ErrPubKeyLen
	}
	viewPubKey = src[:32]
	spendPubKey = src[32:]
	return
}

func (policy *privacyPolicy) getPrivKeyByAddr(addr string) (crypto.PrivKey, error) {
	/ 
	Accountstor, err := policy.store.getAccountByAddr(addr)
	if err != nil {
		bizlog.Error("ProcSendToAddress", "GetAccountByAddr err:", err)
		return nil, err
	}

	/ passwor 
	prikeybyte, err := common.FromHex(Accountstor.GetPrivkey())
	if err != nil || len(prikeybyte) == 0 {
		bizlog.Error("ProcSendToAddress", "FromHex err", err)
		return nil, err
	}
	operater := policy.getWalletOperate()
	password := []byte(operater.GetPassword())
	privkey := wcom.CBCDecrypterPrivkey(password, prikeybyte)
	/ privke pubke addr
	cr, err := crypto.New(types.GetSignName("privacy", operater.GetSignType()))
	if err != nil {
		bizlog.Error("ProcSendToAddress", "err", err)
		return nil, err
	}
	priv, err := cr.PrivKeyFromBytes(privkey)
	if err != nil {
		bizlog.Error("ProcSendToAddress", "PrivKeyFromBytes err", err)
		return nil, err
	}
	return priv, nil
}

func (policy *privacyPolicy) getPrivacykeyPair(addr string) (*privacy.Privacy, error) {
	if accPrivacy, _ := policy.store.getWalletAccountPrivacy(addr); accPrivacy != nil {
		privacyInfo := &privacy.Privacy{}
		password := []byte(policy.getWalletOperate().GetPassword())
		copy(privacyInfo.ViewPubkey[:], accPrivacy.ViewPubkey)
		decrypteredView := wcom.CBCDecrypterPrivkey(password, accPrivacy.ViewPrivKey)
		copy(privacyInfo.ViewPrivKey[:], decrypteredView)
		copy(privacyInfo.SpendPubkey[:], accPrivacy.SpendPubkey)
		decrypteredSpend := wcom.CBCDecrypterPrivkey(password, accPrivacy.SpendPrivKey)
		copy(privacyInfo.SpendPrivKey[:], decrypteredSpend)

		return privacyInfo, nil
	}
	_, err := policy.getPrivKeyByAddr(addr)
	if err != nil {
		return nil, err
	}
	return nil, privacytypes.ErrPrivacyNotEnabled

}

func (policy *privacyPolicy) savePrivacykeyPair(addr string) (*privacy.Privacy, error) {
	priv, err := policy.getPrivKeyByAddr(addr)
	if err != nil {
		return nil, err
	}

	newPrivacy, err := privacy.NewPrivacyWithPrivKey((*[privacy.KeyLen32]byte)(unsafe.Pointer(&priv.Bytes()[0])))
	if err != nil {
		return nil, err
	}

	password := []byte(policy.getWalletOperate().GetPassword())
	encrypteredView := wcom.CBCEncrypterPrivkey(password, newPrivacy.ViewPrivKey.Bytes())
	encrypteredSpend := wcom.CBCEncrypterPrivkey(password, newPrivacy.SpendPrivKey.Bytes())
	walletPrivacy := &privacytypes.WalletAccountPrivacy{
		ViewPubkey:   newPrivacy.ViewPubkey[:],
		ViewPrivKey:  encrypteredView,
		SpendPubkey:  newPrivacy.SpendPubkey[:],
		SpendPrivKey: encrypteredSpend,
	}
	//save the privacy created to wallet db
	policy.store.setWalletAccountPrivacy(addr, walletPrivacy)
	return newPrivacy, nil
}

func (policy *privacyPolicy) enablePrivacy(req *privacytypes.ReqEnablePrivacy) (*privacytypes.RepEnablePrivacy, error) {
	var addrs []string
	if 0 == len(req.Addrs) {
		WalletAccStores, err := policy.store.getAccountByPrefix("Account")
		if err != nil || len(WalletAccStores) == 0 {
			bizlog.Info("enablePrivacy", "GetAccountByPrefix:err", err)
			return nil, types.ErrNotFound
		}
		for _, WalletAccStore := range WalletAccStores {
			addrs = append(addrs, WalletAccStore.Addr)
		}
	} else {
		addrs = append(addrs, req.Addrs...)
	}

	var rep privacytypes.RepEnablePrivacy
	for _, addr := range addrs {
		str := ""
		isOK := true
		_, err := policy.getPrivacykeyPair(addr)
		if err != nil {
			_, err = policy.savePrivacykeyPair(addr)
			if err != nil {
				isOK = false
				str = err.Error()
			}
		}

		priAddrResult := &privacytypes.PriAddrResult{
			Addr: addr,
			IsOK: isOK,
			Msg:  str,
		}

		rep.Results = append(rep.Results, priAddrResult)
	}
	return &rep, nil
}

func (policy *privacyPolicy) showPrivacyKeyPair(reqAddr *types.ReqString) (*privacytypes.ReplyPrivacyPkPair, error) {
	privacyInfo, err := policy.getPrivacykeyPair(reqAddr.GetData())
	if err != nil {
		bizlog.Error("showPrivacyKeyPair", "getPrivacykeyPair error ", err)
		return nil, err
	}

	//pair := privacyInfo.ViewPubkey[:]
	//pair = append(pair, privacyInfo.SpendPubkey[:]...)

	replyPrivacyPkPair := &privacytypes.ReplyPrivacyPkPair{
		ShowSuccessful: true,
		Pubkeypair:     makeViewSpendPubKeyPairToString(privacyInfo.ViewPubkey[:], privacyInfo.SpendPubkey[:]),
	}
	return replyPrivacyPkPair, nil
}

func (policy *privacyPolicy) getPrivacyAccountInfo(req *privacytypes.ReqPrivacyAccount) (*privacytypes.ReplyPrivacyAccount, error) {
	addr := strings.Trim(req.GetAddr(), " ")
	token := req.GetToken()
	reply := &privacytypes.ReplyPrivacyAccount{}
	reply.Displaymode = req.Displaymode
	if len(addr) == 0 {
		return nil, errors.New("Address is empty")
	}

	// 
	privacyDBStore, err := policy.store.listAvailableUTXOs(req.GetAssetExec(), token, addr)
	if err != nil {
		bizlog.Error("getPrivacyAccountInfo", "listAvailableUTXOs")
		return nil, err
	}
	utxos := make([]*privacytypes.UTXO, 0)
	for _, ele := range privacyDBStore {
		utxoBasic := &privacytypes.UTXOBasic{
			UtxoGlobalIndex: &privacytypes.UTXOGlobalIndex{
				Outindex: ele.OutIndex,
				Txhash:   ele.Txhash,
			},
			OnetimePubkey: ele.OnetimePublicKey,
		}
		utxo := &privacytypes.UTXO{
			Amount:    ele.Amount,
			UtxoBasic: utxoBasic,
		}
		utxos = append(utxos, utxo)
	}
	reply.Utxos = &privacytypes.UTXOs{Utxos: utxos}

	// 
	utxos = make([]*privacytypes.UTXO, 0)
	ftxoslice, err := policy.store.listFrozenUTXOs(req.GetAssetExec(), token, addr)
	if err == nil && ftxoslice != nil {
		for _, ele := range ftxoslice {
			utxos = append(utxos, ele.Utxos...)
		}
	}

	reply.Ftxos = &privacytypes.UTXOs{Utxos: utxos}

	return reply, nil
}

// UTX 
// UTX 1 UTXO
// 1 UTXO
// UTX       
func (policy *privacyPolicy) selectUTXO(assetExec, token, addr string, amount int64) ([]*txOutputInfo, error) {
	if len(token) == 0 || len(addr) == 0 || amount <= 0 {
		return nil, types.ErrInvalidParam
	}
	wutxos, err := policy.store.getPrivacyTokenUTXOs(assetExec, token, addr)
	if err != nil {
		return nil, types.ErrInsufficientBalance
	}
	operater := policy.getWalletOperate()
	curBlockHeight := operater.GetBlockHeight()
	var confirmUTXOs, unconfirmUTXOs []*walletUTXO
	var balance int64
	for _, wutxo := range wutxos.utxos {
		if curBlockHeight < wutxo.height {
			continue
		}
		if curBlockHeight-wutxo.height > privacytypes.UtxoMaturityDegree {
			balance += wutxo.outinfo.amount
			confirmUTXOs = append(confirmUTXOs, wutxo)
		} else {
			unconfirmUTXOs = append(unconfirmUTXOs, wutxo)
		}
	}
	if balance < amount && len(unconfirmUTXOs) > 0 {
		// UTX   
		// 
		sort.Slice(unconfirmUTXOs, func(i, j int) bool {
			return unconfirmUTXOs[i].height < unconfirmUTXOs[j].height
		})
		for _, wutxo := range unconfirmUTXOs {
			confirmUTXOs = append(confirmUTXOs, wutxo)
			balance += wutxo.outinfo.amount
			if balance >= amount {
				break
			}
		}
	}
	if balance < amount {
		return nil, types.ErrInsufficientBalance
	}
	balance = 0
	var selectedOuts []*txOutputInfo
	for balance < amount {
		index := operater.GetRandom().Intn(len(confirmUTXOs))
		selectedOuts = append(selectedOuts, confirmUTXOs[index].outinfo)
		balance += confirmUTXOs[index].outinfo.amount
		// remove selected utxo
		confirmUTXOs = append(confirmUTXOs[:index], confirmUTXOs[index+1:]...)
	}
	return selectedOuts, nil
}

/*
buildInput 
 
	1 UTX 
	2 (mixcout>0) UTX UTXO UTX 
	3  x=Hs(aR)+b   xG = Hs(ar)G+bG = Hs(aR)G+B 
*/
func (policy *privacyPolicy) buildInput(privacykeyParirs *privacy.Privacy, buildInfo *buildInputInfo) (*privacytypes.PrivacyInput, []*privacytypes.UTXOBasics, []*privacytypes.RealKeyInput, []*txOutputInfo, error) {
	operater := policy.getWalletOperate()
	/ utxo
	selectedUtxo, err := policy.selectUTXO(buildInfo.assetExec, buildInfo.assetSymbol, buildInfo.sender, buildInfo.amount)
	if err != nil {
		bizlog.Error("buildInput", "Failed to selectOutput for amount", buildInfo.amount,
			"Due to cause", err)
		return nil, nil, nil, nil, err
	}
	sort.Slice(selectedUtxo, func(i, j int) bool {
		return selectedUtxo[i].amount <= selectedUtxo[j].amount
	})

	reqGetGlobalIndex := privacytypes.ReqUTXOGlobalIndex{
		AssetExec:   buildInfo.assetExec,
		AssetSymbol: buildInfo.assetSymbol,
		MixCount:    0,
	}

	if buildInfo.mixcount > 0 {
		reqGetGlobalIndex.MixCount = common.MinInt32(int32(privacytypes.PrivacyMaxCount), common.MaxInt32(buildInfo.mixcount, 0))
	}
	for _, out := range selectedUtxo {
		reqGetGlobalIndex.Amount = append(reqGetGlobalIndex.Amount, out.amount)
	}
	//  blockchai 
	var resUTXOGlobalIndex *privacytypes.ResUTXOGlobalIndex
	if buildInfo.mixcount > 0 {
		query := &types.ChainExecutor{
			Driver:   "privacy",
			FuncName: "GetUTXOGlobalIndex",
			Param:    types.Encode(&reqGetGlobalIndex),
		}
		/ blockchai utx 
		data, err := operater.GetAPI().QueryChain(query)
		if err != nil {
			bizlog.Error("buildInput BlockChainQuery", "err", err)
			return nil, nil, nil, nil, err
		}
		resUTXOGlobalIndex = data.(*privacytypes.ResUTXOGlobalIndex)
		if resUTXOGlobalIndex == nil {
			bizlog.Info("buildInput EventBlockChainQuery is nil")
			return nil, nil, nil, nil, err
		}

		sort.Slice(resUTXOGlobalIndex.UtxoIndex4Amount, func(i, j int) bool {
			return resUTXOGlobalIndex.UtxoIndex4Amount[i].Amount <= resUTXOGlobalIndex.UtxoIndex4Amount[j].Amount
		})

		if len(selectedUtxo) != len(resUTXOGlobalIndex.UtxoIndex4Amount) {
			bizlog.Error("buildInput EventBlockChainQuery get not the same count for mix",
				"len(selectedUtxo)", len(selectedUtxo),
				"len(resUTXOGlobalIndex.UtxoIndex4Amount)", len(resUTXOGlobalIndex.UtxoIndex4Amount))
		}
	}

	/ PrivacyInput
	privacyInput := &privacytypes.PrivacyInput{}
	utxosInKeyInput := make([]*privacytypes.UTXOBasics, len(selectedUtxo))
	realkeyInputSlice := make([]*privacytypes.RealKeyInput, len(selectedUtxo))
	for i, utxo2pay := range selectedUtxo {
		var utxoIndex4Amount *privacytypes.UTXOIndex4Amount
		if nil != resUTXOGlobalIndex && i < len(resUTXOGlobalIndex.UtxoIndex4Amount) && utxo2pay.amount == resUTXOGlobalIndex.UtxoIndex4Amount[i].Amount {
			utxoIndex4Amount = resUTXOGlobalIndex.UtxoIndex4Amount[i]
			for j, utxo := range utxoIndex4Amount.Utxos {
				/ UTX  
				if bytes.Equal(utxo.OnetimePubkey, utxo2pay.onetimePublicKey) {
					utxoIndex4Amount.Utxos = append(utxoIndex4Amount.Utxos[:j], utxoIndex4Amount.Utxos[j+1:]...)
					break
				}
			}
		}

		if utxoIndex4Amount == nil {
			utxoIndex4Amount = &privacytypes.UTXOIndex4Amount{}
		}
		if utxoIndex4Amount.Utxos == nil {
			utxoIndex4Amount.Utxos = make([]*privacytypes.UTXOBasic, 0)
		}
		/ utx mi  utx  
		if len(utxoIndex4Amount.Utxos) > int(buildInfo.mixcount) {
			utxoIndex4Amount.Utxos = utxoIndex4Amount.Utxos[:len(utxoIndex4Amount.Utxos)-1]
		}

		utxo := &privacytypes.UTXOBasic{
			UtxoGlobalIndex: utxo2pay.utxoGlobalIndex,
			OnetimePubkey:   utxo2pay.onetimePublicKey,
		}
		/ utx 
		utxoIndex4Amount.Utxos = append(utxoIndex4Amount.Utxos, utxo)
		positions := operater.GetRandom().Perm(len(utxoIndex4Amount.Utxos))
		utxos := make([]*privacytypes.UTXOBasic, len(utxoIndex4Amount.Utxos))
		for k, position := range positions {
			utxos[position] = utxoIndex4Amount.Utxos[k]
		}
		utxosInKeyInput[i] = &privacytypes.UTXOBasics{Utxos: utxos}

		//x = Hs(aR) + b
		onetimePriv, err := privacy.RecoverOnetimePriKey(utxo2pay.txPublicKeyR, privacykeyParirs.ViewPrivKey, privacykeyParirs.SpendPrivKey, int64(utxo2pay.utxoGlobalIndex.Outindex))
		if err != nil {
			bizlog.Error("transPri2Pri", "Failed to RecoverOnetimePriKey", err)
			return nil, nil, nil, nil, err
		}

		realkeyInput := &privacytypes.RealKeyInput{
			Realinputkey:   int32(positions[len(positions)-1]),
			Onetimeprivkey: onetimePriv.Bytes(),
		}
		realkeyInputSlice[i] = realkeyInput

		keyImage, err := privacy.GenerateKeyImage(onetimePriv, utxo2pay.onetimePublicKey)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		keyInput := &privacytypes.KeyInput{
			Amount:   utxo2pay.amount,
			KeyImage: keyImage[:],
		}

		for _, utxo := range utxos {
			keyInput.UtxoGlobalIndex = append(keyInput.UtxoGlobalIndex, utxo.UtxoGlobalIndex)
		}
		/ inpu  ，keyImag ，
		/  utx keyinpu pubke 
		/ 
		privacyInput.Keyinput = append(privacyInput.Keyinput, keyInput)
	}

	return privacyInput, utxosInKeyInput, realkeyInputSlice, selectedUtxo, nil
}

func (policy *privacyPolicy) createTransaction(req *privacytypes.ReqCreatePrivacyTx) (*types.Transaction, error) {
	switch req.ActionType {
	case privacytypes.ActionPublic2Privacy:
		return policy.createPublic2PrivacyTx(req)
	case privacytypes.ActionPrivacy2Privacy:
		return policy.createPrivacy2PrivacyTx(req)
	case privacytypes.ActionPrivacy2Public:
		return policy.createPrivacy2PublicTx(req)
	}
	return nil, types.ErrInvalidParam
}

func (policy *privacyPolicy) createPublic2PrivacyTx(req *privacytypes.ReqCreatePrivacyTx) (*types.Transaction, error) {
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	viewPubSlice, spendPubSlice, err := parseViewSpendPubKeyPair(req.GetPubkeypair())
	if err != nil {
		bizlog.Error("createPublic2PrivacyTx", "parse view spend public key pair failed.  err ", err)
		return nil, err
	}
	amount := req.GetAmount()
	viewPublic := (*[32]byte)(unsafe.Pointer(&viewPubSlice[0]))
	spendPublic := (*[32]byte)(unsafe.Pointer(&spendPubSlice[0]))
	privacyOutput, err := generateOuts(viewPublic, spendPublic, nil, nil, amount, amount, 0, cfg.GetCoinPrecision())
	if err != nil {
		bizlog.Error("createPublic2PrivacyTx", "generate output failed.  err ", err)
		return nil, err
	}

	value := &privacytypes.Public2Privacy{
		Tokenname: req.Tokenname,
		Amount:    amount,
		Note:      req.GetNote(),
		Output:    privacyOutput,
		AssetExec: req.GetAssetExec(),
	}

	action := &privacytypes.PrivacyAction{
		Ty:    privacytypes.ActionPublic2Privacy,
		Value: &privacytypes.PrivacyAction_Public2Privacy{Public2Privacy: value},
	}
	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(privacytypes.PrivacyX)),
		Payload: types.Encode(action),
		Nonce:   policy.getWalletOperate().Nonce(),
		To:      address.ExecAddress(cfg.ExecName(privacytypes.PrivacyX)),
		ChainID: cfg.GetChainID(),
	}
	tx.SetExpire(cfg, time.Duration(req.Expire))
	tx.Signature = &types.Signature{
		Signature: types.Encode(&privacytypes.PrivacySignatureParam{
			ActionType: action.Ty,
		}),
	}
	tx.Fee, err = tx.GetRealFee(cfg.GetMinTxFeeRate())
	if err != nil {
		bizlog.Error("createPublic2PrivacyTx", "calc fee failed", err)
		return nil, err
	}

	return tx, nil
}

func (policy *privacyPolicy) createPrivacy2PrivacyTx(req *privacytypes.ReqCreatePrivacyTx) (*types.Transaction, error) {

	/ utxo
	var utxoBurnedAmount int64
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	isMainetCoins := !cfg.IsPara() && (req.AssetExec == cfg.GetCoinExec())
	if isMainetCoins {
		utxoBurnedAmount = privacytypes.PrivacyTxFee * cfg.GetCoinPrecision()
	}
	buildInfo := &buildInputInfo{
		assetExec:   req.GetAssetExec(),
		assetSymbol: req.GetTokenname(),
		sender:      req.GetFrom(),
		amount:      req.GetAmount() + utxoBurnedAmount,
		mixcount:    req.GetMixcount(),
	}
	privacyInfo, err := policy.getPrivacykeyPair(req.GetFrom())
	if err != nil {
		bizlog.Error("createPrivacy2PrivacyTx", "getPrivacykeyPair error", err)
		return nil, err
	}
	//step 1,buildInput
	privacyInput, utxosInKeyInput, realkeyInputSlice, selectedUtxo, err := policy.buildInput(privacyInfo, buildInfo)
	if err != nil {
		return nil, err
	}
	//step 2,generateOuts
	viewPublicSlice, spendPublicSlice, err := parseViewSpendPubKeyPair(req.GetPubkeypair())
	if err != nil {
		bizlog.Error("createPrivacy2PrivacyTx", "parseViewSpendPubKeyPair  ", err)
		return nil, err
	}

	viewPub4change, spendPub4change := privacyInfo.ViewPubkey.Bytes(), privacyInfo.SpendPubkey.Bytes()
	viewPublic := (*[32]byte)(unsafe.Pointer(&viewPublicSlice[0]))
	spendPublic := (*[32]byte)(unsafe.Pointer(&spendPublicSlice[0]))
	viewPub4chgPtr := (*[32]byte)(unsafe.Pointer(&viewPub4change[0]))
	spendPub4chgPtr := (*[32]byte)(unsafe.Pointer(&spendPub4change[0]))

	selectedAmounTotal := int64(0)
	for _, input := range privacyInput.Keyinput {
		selectedAmounTotal += input.Amount
	}
	/ UTXO
	privacyOutput, err := generateOuts(viewPublic, spendPublic, viewPub4chgPtr, spendPub4chgPtr, req.GetAmount(), selectedAmounTotal, utxoBurnedAmount, cfg.GetCoinPrecision())
	if err != nil {
		return nil, err
	}

	value := &privacytypes.Privacy2Privacy{
		Tokenname: req.GetTokenname(),
		Amount:    req.GetAmount(),
		Note:      req.GetNote(),
		Input:     privacyInput,
		Output:    privacyOutput,
		AssetExec: req.GetAssetExec(),
	}
	action := &privacytypes.PrivacyAction{
		Ty:    privacytypes.ActionPrivacy2Privacy,
		Value: &privacytypes.PrivacyAction_Privacy2Privacy{Privacy2Privacy: value},
	}

	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(privacytypes.PrivacyX)),
		Payload: types.Encode(action),
		Fee:     privacytypes.PrivacyTxFee * cfg.GetCoinPrecision(),
		Nonce:   policy.getWalletOperate().Nonce(),
		To:      address.ExecAddress(cfg.ExecName(privacytypes.PrivacyX)),
		ChainID: cfg.GetChainID(),
	}
	tx.SetExpire(cfg, time.Duration(req.Expire))
	if !isMainetCoins {
		tx.Fee, err = tx.GetRealFee(cfg.GetMinTxFeeRate())
		if err != nil {
			bizlog.Error("createPrivacy2PrivacyTx", "calc fee failed", err)
			return nil, err
		}
	}

	//  UTX  txHas 
	policy.saveFTXOInfo(tx.GetExpire(), req.GetAssetExec(), req.Tokenname, req.GetFrom(), hex.EncodeToString(tx.Hash()), selectedUtxo)
	tx.Signature = &types.Signature{
		Signature: types.Encode(&privacytypes.PrivacySignatureParam{
			ActionType:    action.Ty,
			Utxobasics:    utxosInKeyInput,
			RealKeyInputs: realkeyInputSlice,
		}),
	}
	return tx, nil
}

func (policy *privacyPolicy) createPrivacy2PublicTx(req *privacytypes.ReqCreatePrivacyTx) (*types.Transaction, error) {

	/ utxo
	/ utxo
	var utxoBurnedAmount int64
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	isMainetCoins := !cfg.IsPara() && (req.AssetExec == cfg.GetCoinExec())
	if isMainetCoins {
		utxoBurnedAmount = privacytypes.PrivacyTxFee * cfg.GetCoinPrecision()
	}
	buildInfo := &buildInputInfo{
		assetExec:   req.GetAssetExec(),
		assetSymbol: req.GetTokenname(),
		sender:      req.GetFrom(),
		amount:      req.GetAmount() + utxoBurnedAmount,
		mixcount:    req.GetMixcount(),
	}
	privacyInfo, err := policy.getPrivacykeyPair(req.GetFrom())
	if err != nil {
		bizlog.Error("createPrivacy2PublicTx failed to getPrivacykeyPair")
		return nil, err
	}
	//step 1,buildInput
	privacyInput, utxosInKeyInput, realkeyInputSlice, selectedUtxo, err := policy.buildInput(privacyInfo, buildInfo)
	if err != nil {
		bizlog.Error("createPrivacy2PublicTx failed to buildInput")
		return nil, err
	}

	viewPub4change, spendPub4change := privacyInfo.ViewPubkey.Bytes(), privacyInfo.SpendPubkey.Bytes()
	viewPub4chgPtr := (*[32]byte)(unsafe.Pointer(&viewPub4change[0]))
	spendPub4chgPtr := (*[32]byte)(unsafe.Pointer(&spendPub4change[0]))

	selectedAmounTotal := int64(0)
	for _, input := range privacyInput.Keyinput {
		if input.Amount <= 0 {
			return nil, errors.New("")
		}
		selectedAmounTotal += input.Amount
	}
	changeAmount := selectedAmounTotal - req.GetAmount()
	//step 2,generateOuts
	/ UTXO UTXO
	privacyOutput, err := generateOuts(nil, nil, viewPub4chgPtr, spendPub4chgPtr, 0, changeAmount, utxoBurnedAmount, cfg.GetCoinPrecision())
	if err != nil {
		return nil, err
	}

	value := &privacytypes.Privacy2Public{
		Tokenname: req.GetTokenname(),
		Amount:    req.GetAmount(),
		Note:      req.GetNote(),
		Input:     privacyInput,
		Output:    privacyOutput,
		To:        req.GetTo(),
		AssetExec: req.GetAssetExec(),
	}
	action := &privacytypes.PrivacyAction{
		Ty:    privacytypes.ActionPrivacy2Public,
		Value: &privacytypes.PrivacyAction_Privacy2Public{Privacy2Public: value},
	}

	tx := &types.Transaction{
		Execer:  []byte(cfg.ExecName(privacytypes.PrivacyX)),
		Payload: types.Encode(action),
		Fee:     privacytypes.PrivacyTxFee * cfg.GetCoinPrecision(),
		Nonce:   policy.getWalletOperate().Nonce(),
		To:      address.ExecAddress(cfg.ExecName(privacytypes.PrivacyX)),
		ChainID: cfg.GetChainID(),
	}
	tx.SetExpire(cfg, time.Duration(req.Expire))
	if !isMainetCoins {
		tx.Fee, err = tx.GetRealFee(cfg.GetMinTxFeeRate())
		if err != nil {
			bizlog.Error("createPrivacy2PublicTx", "calc fee failed", err)
			return nil, err
		}
	}
	//  UTX  txHas 
	policy.saveFTXOInfo(tx.GetExpire(), req.GetAssetExec(), req.Tokenname, req.GetFrom(), hex.EncodeToString(tx.Hash()), selectedUtxo)
	tx.Signature = &types.Signature{
		Signature: types.Encode(&privacytypes.PrivacySignatureParam{
			ActionType:    action.Ty,
			Utxobasics:    utxosInKeyInput,
			RealKeyInputs: realkeyInputSlice,
		}),
	}
	return tx, nil
}

func (policy *privacyPolicy) saveFTXOInfo(expire int64, assetExec, assetSymbol, sender, txhash string, selectedUtxos []*txOutputInfo) {
	/ utx  
	policy.store.moveUTXO2FTXO(expire, assetExec, assetSymbol, sender, txhash, selectedUtxos)
	//TODO  txhas  ，
	//TODO  FTX STXO，added by hezhengjun on 2018.6.5
}

func (policy *privacyPolicy) getPrivacyKeyPairs() ([]addrAndprivacy, error) {
	/ Accoun 
	WalletAccStores, err := policy.store.getAccountByPrefix("Account")
	if err != nil || len(WalletAccStores) == 0 {
		bizlog.Info("getPrivacyKeyPairs", "store getAccountByPrefix error", err)
		return nil, err
	}

	var infoPriRes []addrAndprivacy
	for _, AccStore := range WalletAccStores {
		if len(AccStore.Addr) != 0 {
			if privacyInfo, err := policy.getPrivacykeyPair(AccStore.Addr); err == nil {
				var priInfo addrAndprivacy
				priInfo.Addr = &AccStore.Addr
				priInfo.PrivacyKeyPair = privacyInfo
				infoPriRes = append(infoPriRes, priInfo)
			}
		}
	}

	if 0 == len(infoPriRes) {
		return nil, privacytypes.ErrPrivacyNotEnabled
	}

	return infoPriRes, nil

}

func (policy *privacyPolicy) rescanUTXOs(req *privacytypes.ReqRescanUtxos) (*privacytypes.RepRescanUtxos, error) {
	if req.Flag != 0 {
		return policy.store.getRescanUtxosFlag4Addr(req)
	}
	// Resca 
	var repRescanUtxos privacytypes.RepRescanUtxos
	repRescanUtxos.Flag = req.Flag

	operater := policy.getWalletOperate()
	if operater.IsWalletLocked() {
		return nil, types.ErrWalletIsLocked
	}
	if ok, err := policy.isRescanUtxosFlagScaning(); ok {
		return nil, err
	}
	_, err := policy.getPrivacyKeyPairs()
	if err != nil {
		return nil, err
	}
	policy.SetRescanFlag(privacytypes.UtxoFlagScaning)
	operater.GetWaitGroup().Add(1)
	go policy.rescanReqUtxosByAddr(req.Addrs)
	return &repRescanUtxos, nil
}

/ blockchai add 
func (policy *privacyPolicy) rescanReqUtxosByAddr(addrs []string) {
	defer policy.getWalletOperate().GetWaitGroup().Done()
	bizlog.Debug("RescanAllUTXO begin!")
	policy.reqUtxosByAddr(addrs)
	bizlog.Debug("RescanAllUTXO success!")
}

func (policy *privacyPolicy) reqUtxosByAddr(addrs []string) {
	// 
	var storeAddrs []string
	if len(addrs) == 0 {
		WalletAccStores, err := policy.store.getAccountByPrefix("Account")
		if err != nil || len(WalletAccStores) == 0 {
			bizlog.Info("reqUtxosByAddr", "getAccountByPrefix error", err)
			return
		}
		for _, WalletAccStore := range WalletAccStores {
			storeAddrs = append(storeAddrs, WalletAccStore.Addr)
		}
	} else {
		storeAddrs = append(storeAddrs, addrs...)
	}
	policy.store.saveREscanUTXOsAddresses(storeAddrs, privacytypes.UtxoFlagScaning)

	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	reqAddr := address.ExecAddress(cfg.ExecName(privacytypes.PrivacyX))
	var txInfo types.ReplyTxInfo
	i := 0
	operater := policy.getWalletOperate()
	for {
		select {
		case <-operater.GetWalletDone():
			return
		default:
		}

		/ exec UTXOs,
		// 1 
		var ReqAddr types.ReqAddr
		ReqAddr.Addr = reqAddr
		ReqAddr.Flag = 0
		ReqAddr.Direction = 0
		ReqAddr.Count = int32(MaxTxHashsPerTime)
		if i == 0 {
			ReqAddr.Height = -1
			ReqAddr.Index = 0
		} else {
			ReqAddr.Height = txInfo.GetHeight()
			ReqAddr.Index = txInfo.GetIndex()
			if !cfg.IsDappFork(ReqAddr.Height, privacytypes.PrivacyX, "ForkV21Privacy") { // 
				break
			}
		}
		i++
		/ 
		msg, err := operater.GetAPI().Query(privacytypes.PrivacyX, "GetTxsByAddr", &ReqAddr)
		if err != nil {
			bizlog.Error("reqUtxosByAddr", "GetTxsByAddr error", err, "addr", reqAddr)
			break
		}
		ReplyTxInfos := msg.(*types.ReplyTxInfos)
		if ReplyTxInfos == nil {
			bizlog.Info("privacy ReqTxInfosByAddr ReplyTxInfos is nil")
			break
		}
		txcount := len(ReplyTxInfos.TxInfos)

		var ReqHashes types.ReqHashes
		ReqHashes.Hashes = make([][]byte, len(ReplyTxInfos.TxInfos))
		for index, ReplyTxInfo := range ReplyTxInfos.TxInfos {
			ReqHashes.Hashes[index] = ReplyTxInfo.GetHash()
		}

		if txcount > 0 {
			txInfo.Hash = ReplyTxInfos.TxInfos[txcount-1].GetHash()
			txInfo.Height = ReplyTxInfos.TxInfos[txcount-1].GetHeight()
			txInfo.Index = ReplyTxInfos.TxInfos[txcount-1].GetIndex()
		}

		policy.getPrivacyTxDetailByHashs(&ReqHashes, addrs)
		if txcount < int(MaxTxHashsPerTime) {
			break
		}
	}
	// 
	policy.SetRescanFlag(privacytypes.UtxoFlagNoScan)
	// privacyInput
	policy.deleteScanPrivacyInputUtxo()
	policy.store.saveREscanUTXOsAddresses(storeAddrs, privacytypes.UtxoFlagScanEnd)
}

//TODO:inpu utxo, utxo
func (policy *privacyPolicy) deleteScanPrivacyInputUtxo() {
	maxUTXOsPerTime := 1000
	for {
		utxoGlobalIndexs := policy.store.setScanPrivacyInputUTXO(int32(maxUTXOsPerTime))
		policy.store.updateScanInputUTXOs(utxoGlobalIndexs)
		if len(utxoGlobalIndexs) < maxUTXOsPerTime {
			break
		}
	}
}

func (policy *privacyPolicy) getPrivacyTxDetailByHashs(ReqHashes *types.ReqHashes, addrs []string) {
	/ txhash txdetail
	TxDetails, err := policy.getWalletOperate().GetAPI().GetTransactionByHash(ReqHashes)
	if err != nil {
		bizlog.Error("getPrivacyTxDetailByHashs", "GetTransactionByHash error", err)
		return
	}
	var privacyInfo []addrAndprivacy
	if len(addrs) > 0 {
		for _, addr := range addrs {
			if privacy, err := policy.getPrivacykeyPair(addr); err == nil {
				priInfo := &addrAndprivacy{
					Addr:           &addr,
					PrivacyKeyPair: privacy,
				}
				privacyInfo = append(privacyInfo, *priInfo)
			}

		}
	} else {
		privacyInfo, _ = policy.getPrivacyKeyPairs()
	}
	policy.store.selectPrivacyTransactionToWallet(TxDetails, privacyInfo)
}

func (policy *privacyPolicy) showPrivacyAccountsSpend(req *privacytypes.ReqPrivBal4AddrToken) (*privacytypes.UTXOHaveTxHashs, error) {
	if ok, err := policy.isRescanUtxosFlagScaning(); ok {
		return nil, err
	}
	utxoHaveTxHashs, err := policy.store.listSpendUTXOs(req.GetAssetExec(), req.GetToken(), req.GetAddr())
	if err != nil {
		return nil, err
	}
	return utxoHaveTxHashs, nil
}

func (policy *privacyPolicy) signatureTx(tx *types.Transaction, privacyInput *privacytypes.PrivacyInput, utxosInKeyInput []*privacytypes.UTXOBasics, realkeyInputSlice []*privacytypes.RealKeyInput) (err error) {
	tx.Signature = nil
	data := types.Encode(tx)
	ringSign := &types.RingSignature{}
	ringSign.Items = make([]*types.RingSignatureItem, len(privacyInput.Keyinput))
	for i, input := range privacyInput.Keyinput {
		utxos := utxosInKeyInput[i]
		h := common.BytesToHash(data)
		item, err := privacy.GenerateRingSignature(h.Bytes(),
			utxos.Utxos,
			realkeyInputSlice[i].Onetimeprivkey,
			int(realkeyInputSlice[i].Realinputkey),
			input.KeyImage)
		if err != nil {
			return err
		}
		ringSign.Items[i] = item
	}
	cfg := policy.getWalletOperate().GetAPI().GetConfig()
	ringSignData := types.Encode(ringSign)
	tx.Signature = &types.Signature{
		Ty:        privacytypes.RingBaseonED25519,
		Signature: ringSignData,
		//  
		Pubkey: address.ExecPubKey(cfg.ExecName(privacytypes.PrivacyX)),
	}
	return nil
}

func (policy *privacyPolicy) buildAndStoreWalletTxDetail(param *buildStoreWalletTxDetailParam) {

	txInfo := param.txInfo
	heightstr := dapp.HeightIndexStr(txInfo.blockHeight, int64(txInfo.txIndex))
	bizlog.Debug("buildAndStoreWalletTxDetail", "heightstr", heightstr, "isRollBack", txInfo.isRollBack)
	if !txInfo.isRollBack {
		var txdetail types.WalletTxDetail
		key := calcTxKey(heightstr)
		txdetail.Tx = txInfo.tx
		txdetail.Height = txInfo.blockHeight
		txdetail.Index = int64(txInfo.txIndex)
		txdetail.Receipt = txInfo.blockDetail.Receipts[txInfo.txIndex]
		txdetail.Blocktime = txInfo.blockDetail.Block.BlockTime

		txdetail.ActionName = txInfo.actionName
		txdetail.Amount, _ = txInfo.tx.Amount()
		txdetail.Fromaddr = param.addr

		txdetailbyte := types.Encode(&txdetail)

		txInfo.batch.Set(key, txdetailbyte)
		/ 
		if sendTx == param.sendRecvFlag {
			txInfo.batch.Set(calcSendPrivacyTxKey(txInfo.assetExec, txInfo.assetSymbol, param.addr, heightstr), key)
		} else if recvTx == param.sendRecvFlag {
			txInfo.batch.Set(calcRecvPrivacyTxKey(txInfo.assetExec, txInfo.assetSymbol, param.addr, heightstr), key)
		}
	} else {
		txInfo.batch.Delete(calcTxKey(heightstr))
		if sendTx == param.sendRecvFlag {
			txInfo.batch.Delete(calcSendPrivacyTxKey(txInfo.assetExec, txInfo.assetSymbol, param.addr, heightstr))
		} else if recvTx == param.sendRecvFlag {
			txInfo.batch.Delete(calcRecvPrivacyTxKey(txInfo.assetExec, txInfo.assetSymbol, param.addr, heightstr))
		}
	}
}

func (policy *privacyPolicy) checkExpireFTXOOnTimer() {
	header := policy.getWalletOperate().GetLastHeader()
	if header == nil {
		return
	}
	curFTXOTxs, keys := policy.store.getFTXOlist()
	if len(curFTXOTxs) == 0 {
		return
	}
	dbbatch := policy.store.NewBatch(true)
	for i, ftxo := range curFTXOTxs {
		if !ftxo.IsExpire(header.GetHeight(), header.GetBlockTime()) {
			continue
		}
		policy.store.moveFTXO2UTXO(keys[i], dbbatch)
		bizlog.Debug("moveFTXO2UTXOWhenFTXOExpire", "moveFTXO2UTXO key", string(keys[i]), "expire ftxo", ftxo)
	}
	err := dbbatch.Write()
	if err != nil {
		bizlog.Error("checkExpireFTXOOnTimer", "db write err", err)
	}
}

func (policy *privacyPolicy) checkWalletStoreData() {
	operater := policy.getWalletOperate()
	defer operater.GetWaitGroup().Done()
	timecount := 10
	checkTicker := time.NewTicker(time.Duration(timecount) * time.Second)
	for {
		select {
		case <-checkTicker.C:
			policy.checkExpireFTXOOnTimer()

			//newbatch := wallet.walletStore.NewBatch(true)
			//err := wallet.procInvalidTxOnTimer(newbatch)
			//if err != nil && err != dbm.ErrNotFoundInDb {
			//	walletlog.Error("checkWalletStoreData", "procInvalidTxOnTimer error ", err)
			//	return
			//}
			//newbatch.Write()
		case <-operater.GetWalletDone():
			return
		}
	}
}

func (policy *privacyPolicy) addDelPrivacyTxsFromBlock(tx *types.Transaction, index int32, block *types.BlockDetail, batch db.Batch, addDelType int32) {

	privacyPairs, err := policy.getPrivacyKeyPairs()
	/  
	if len(privacyPairs) == 0 {
		bizlog.Debug("addDelPrivacyTxsFromBlock", "getPrivacyKeyPairs err", err)
		return
	}

	txhash := tx.Hash()
	txhashHex := hex.EncodeToString(txhash)
	var action privacytypes.PrivacyAction
	if err := types.Decode(tx.GetPayload(), &action); err != nil {
		bizlog.Error("addDelPrivacyTxsFromBlock", "txhash", txhashHex, "addDelType", addDelType, "index", index, "Decode tx.GetPayload() error", err)
		return
	}

	assetExec, assetSymbol := action.GetAssetExecSymbol()
	if assetExec == "" {
		assetExec = policy.getWalletOperate().GetAPI().GetConfig().GetCoinExec()
	}

	txInfo := &privacyTxInfo{
		tx:          tx,
		blockDetail: block,
		blockHeight: block.GetBlock().GetHeight(),
		actionName:  action.GetActionName(),
		actionTy:    action.GetTy(),
		input:       action.GetInput(),
		output:      action.GetOutput(),
		txIndex:     index,
		txHash:      txhash,
		txHashHex:   txhashHex,
		batch:       batch,
		assetSymbol: assetSymbol,
		assetExec:   assetExec,
		isRollBack:  addDelType != AddTx,
		isExecOk:    types.ExecOk == block.Receipts[index].Ty,
	}

	bizlog.Debug("addDelPrivacyTxsFromBlock", "txhash", txhashHex, "actionName", txInfo.actionName,
		"index", index, "isRollBack", txInfo.isRollBack, "isExecOk", txInfo.isExecOk)

	if txInfo.actionTy == privacytypes.ActionPublic2Privacy {

		//  
		txFrom := tx.From()
		for _, key := range privacyPairs {
			if *key.Addr == txFrom {
				param := &buildStoreWalletTxDetailParam{
					txInfo:       txInfo,
					addr:         txInfo.tx.From(),
					sendRecvFlag: sendTx,
				}
				policy.buildAndStoreWalletTxDetail(param)
				break
			}
		}
		policy.processOutputUtxos(txFrom, privacyPairs, txInfo)
	} else if txInfo.actionTy == privacytypes.ActionPrivacy2Privacy {

		sender := policy.processInputUtxos(txInfo)
		policy.processOutputUtxos(sender, privacyPairs, txInfo)
	} else if txInfo.actionTy == privacytypes.ActionPrivacy2Public {
		sender := policy.processInputUtxos(txInfo)
		policy.processOutputUtxos(sender, privacyPairs, txInfo)
		//  
		txTo := action.GetPrivacy2Public().GetTo()
		for _, key := range privacyPairs {
			if *key.Addr == txTo {
				param := &buildStoreWalletTxDetailParam{
					txInfo:       txInfo,
					addr:         txTo,
					sendRecvFlag: recvTx,
				}
				policy.buildAndStoreWalletTxDetail(param)
				break
			}
		}
	}
}

func (policy *privacyPolicy) processInputUtxos(txInfo *privacyTxInfo) string {

	bizlog.Debug("processInputUtxos", "txhash", txInfo.txHashHex, "actionName", txInfo.actionName, "isExecOk", txInfo.isExecOk, "isRollBack", txInfo.isRollBack)

	// utx 
	utxoSender := ""
	var err error
	// 
	if txInfo.isRollBack {

		/  STX  FTXO 
		stxosInOneTx, _, _ := policy.store.getWalletFtxoStxo(STXOs4Tx)
		for _, stxo := range stxosInOneTx {
			if stxo.Txhash != txInfo.txHashHex {
				continue
			}
			if txInfo.isExecOk {
				err = policy.store.moveSTXO2FTXO(txInfo.tx, txInfo.txHashHex, txInfo.batch)
				utxoSender = stxo.Sender
			}
			if err != nil {
				bizlog.Error("processInputUtxos", "txHash", txInfo.txHashHex, "moveSTXO2FTXO err", err)
			}
		}
	} else {

		ftxos, keys := policy.store.getFTXOlist()
		for i, ftxo := range ftxos {
			if ftxo.Txhash != txInfo.txHashHex {
				continue
			}

			if txInfo.isExecOk {
				err = policy.store.moveFTXO2STXO(keys[i], txInfo.txHashHex, txInfo.batch)
			} else {
				policy.store.moveFTXO2UTXO(keys[i], txInfo.batch)
			}
			utxoSender = ftxo.Sender
			if err != nil {
				bizlog.Error("processInputUtxos", "txHash", txInfo.txHashHex, "moveSTXO2FTXO err", err)
			}

		}
	}
	//  
	if utxoSender != "" {
		param := &buildStoreWalletTxDetailParam{
			txInfo:       txInfo,
			addr:         utxoSender,
			sendRecvFlag: sendTx,
		}
		policy.buildAndStoreWalletTxDetail(param)
	}
	return utxoSender
}

func (policy *privacyPolicy) processOutputUtxos(utxoSender string, keys []addrAndprivacy, txInfo *privacyTxInfo) {

	bizlog.Debug("processOutputUtxos", "txhash", txInfo.txHashHex, "actionName", txInfo.actionName, "isExecOK", txInfo.isExecOk, "isRollBack", txInfo.isRollBack)
	if txInfo.output == nil {
		return
	}
	var recvUtxos []*privacytypes.PrivacyDBStore
	// check utxo owner
	omitCheck := false
	owner := ""
	receivers := make(map[string]struct{})
	for index, out := range txInfo.output.GetKeyoutput() {

		for _, key := range keys {
			if omitCheck {
				break
			}
			owner = ""
			priv, err := privacy.RecoverOnetimePriKey(txInfo.output.GetRpubKeytx(), key.PrivacyKeyPair.ViewPrivKey, key.PrivacyKeyPair.SpendPrivKey, int64(index))
			if err != nil {
				bizlog.Error("addDelPrivacyTxsFromBlock", "txhash", txInfo.txHashHex, "actionName", txInfo.actionName, "RecoverOnetimePriKey error", err)
			}
			if bytes.Equal(priv.PubKey().Bytes()[:], out.GetOnetimepubkey()) {
				owner = *key.Addr
				receivers[owner] = struct{}{}
				break
			}
		}

		if len(owner) > 0 {
			info2store := &privacytypes.PrivacyDBStore{
				AssetExec:        txInfo.assetExec,
				Txhash:           txInfo.txHash,
				Tokenname:        txInfo.assetSymbol,
				Amount:           out.Amount,
				OutIndex:         int32(index),
				TxPublicKeyR:     txInfo.output.GetRpubKeytx(),
				OnetimePublicKey: out.Onetimepubkey,
				Owner:            owner,
				Height:           txInfo.blockHeight,
				Txindex:          txInfo.txIndex,
			}
			recvUtxos = append(recvUtxos, info2store)
		}

		// utx  utxo
		if txInfo.actionTy == privacytypes.ActionPublic2Privacy {
			omitCheck = true
		}
	}

	// add or del utxo in wallet
	for _, utxo := range recvUtxos {

		if !txInfo.isExecOk {
			/  
			break
		}
		var err error
		if !txInfo.isRollBack {
			err = policy.store.setUTXO(utxo, txInfo.txHashHex, txInfo.batch)
		} else {
			err = policy.store.unsetUTXO(utxo.AssetExec, utxo.Tokenname, utxo.Owner, txInfo.txHashHex, int(utxo.OutIndex), txInfo.batch)
		}

		if err != nil {
			bizlog.Error("processOutputUtxos", "txHash", txInfo.txHashHex, "actionName", txInfo.actionName, "isRollBack", txInfo.isRollBack, "setUtxoErr", err)
		}
	}
	// handle recv tx index
	for receiver := range receivers {
		// utx    
		// utx  
		if len(receivers) > 1 && receiver == utxoSender {
			continue
		}
		// utx  
		if txInfo.actionTy == privacytypes.ActionPrivacy2Public {
			continue
		}
		param := &buildStoreWalletTxDetailParam{
			txInfo:       txInfo,
			addr:         receiver,
			sendRecvFlag: recvTx,
		}
		policy.buildAndStoreWalletTxDetail(param)
	}
}
