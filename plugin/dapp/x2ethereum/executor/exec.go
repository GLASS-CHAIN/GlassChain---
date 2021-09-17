package executor

import (
	"errors"

	"github.com/33cn/chain33/system/dapp"
	manTy "github.com/33cn/chain33/system/dapp/manage/types"

	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/types"
	x2eTy "github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
)

/*
 * 
 * （statedb （log）
 */

//---------------- Ethereum(eth/erc20) --> Chain33-------------------//

// chain3 ETH/ERC2 
func (x *x2ethereum) Exec_Eth2Chain33Lock(payload *x2eTy.Eth2Chain33, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}

	payload.ValidatorAddress = address.PubKeyToAddr(tx.Signature.Pubkey)

	return action.procEth2Chain33_lock(payload)
}

//----------------  Chain33(eth/erc20)------> Ethereum -------------------//
// chain3  eth
func (x *x2ethereum) Exec_Chain33ToEthBurn(payload *x2eTy.Chain33ToEth, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	return action.procChain33ToEth_burn(payload)
}

//---------------- Ethereum (bty) --> Chain33-------------------//
// et bt  chain33
func (x *x2ethereum) Exec_Eth2Chain33Burn(payload *x2eTy.Eth2Chain33, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}

	payload.ValidatorAddress = address.PubKeyToAddr(tx.Signature.Pubkey)

	return action.procEth2Chain33_burn(payload)
}

//---------------- Chain33 --> Ethereum (bty) -------------------//
//  ethereum  chain33 
func (x *x2ethereum) Exec_Chain33ToEthLock(payload *x2eTy.Chain33ToEth, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	return action.procChain33ToEth_lock(payload)
}

// 
func (x *x2ethereum) Exec_Transfer(payload *types.AssetsTransfer, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	return action.procMsgTransfer(payload)
}

func (x *x2ethereum) Exec_TransferToExec(payload *types.AssetsTransferToExec, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	if !x2eTy.IsExecAddrMatch(payload.ExecName, tx.GetRealToAddr()) {
		return nil, types.ErrToAddrNotSameToExecAddr
	}
	return action.procMsgTransferToExec(payload)
}

func (x *x2ethereum) Exec_WithdrawFromExec(payload *types.AssetsWithdraw, tx *types.Transaction, index int) (*types.Receipt, error) {
	action := newAction(x, tx, int32(index))
	if action == nil {
		return nil, errors.New("Create Action Error")
	}
	if dapp.IsDriverAddress(tx.GetRealToAddr(), x.GetHeight()) || x2eTy.IsExecAddrMatch(payload.ExecName, tx.GetRealToAddr()) {
		return action.procMsgWithDrawFromExec(payload)
	}
	return nil, errors.New("tx error")
}

//------------------------- -------------------------//

// AddValidato validator
func (x *x2ethereum) Exec_AddValidator(payload *x2eTy.MsgValidator, tx *types.Transaction, index int) (*types.Receipt, error) {
	confManager := types.ConfSub(x.GetAPI().GetConfig(), manTy.ManageX).GStrList("superManager")
	err := checkTxSignBySpecificAddr(tx, confManager)
	if err == nil {
		action := newAction(x, tx, int32(index))
		if action == nil {
			return nil, errors.New("Create Action Error")
		}
		return action.procAddValidator(payload)
	}
	return nil, err
}

// RemoveValidato validator
func (x *x2ethereum) Exec_RemoveValidator(payload *x2eTy.MsgValidator, tx *types.Transaction, index int) (*types.Receipt, error) {
	confManager := types.ConfSub(x.GetAPI().GetConfig(), manTy.ManageX).GStrList("superManager")
	err := checkTxSignBySpecificAddr(tx, confManager)
	if err == nil {
		action := newAction(x, tx, int32(index))
		if action == nil {
			return nil, errors.New("Create Action Error")
		}
		return action.procRemoveValidator(payload)
	}
	return nil, err
}

// ModifyPowe validato power
func (x *x2ethereum) Exec_ModifyPower(payload *x2eTy.MsgValidator, tx *types.Transaction, index int) (*types.Receipt, error) {
	confManager := types.ConfSub(x.GetAPI().GetConfig(), manTy.ManageX).GStrList("superManager")
	err := checkTxSignBySpecificAddr(tx, confManager)
	if err == nil {
		action := newAction(x, tx, int32(index))
		if action == nil {
			return nil, errors.New("Create Action Error")
		}
		return action.procModifyValidator(payload)
	}
	return nil, err
}

// SetConsensusThreshol validato clai 
func (x *x2ethereum) Exec_SetConsensusThreshold(payload *x2eTy.MsgConsensusThreshold, tx *types.Transaction, index int) (*types.Receipt, error) {
	confManager := types.ConfSub(x.GetAPI().GetConfig(), manTy.ManageX).GStrList("superManager")
	err := checkTxSignBySpecificAddr(tx, confManager)
	if err == nil {
		action := newAction(x, tx, int32(index))
		if action == nil {
			return nil, errors.New("Create Action Error")
		}
		return action.procMsgSetConsensusThreshold(payload)
	}
	return nil, err
}

func checkTxSignBySpecificAddr(tx *types.Transaction, addrs []string) error {
	signAddr := address.PubKeyToAddr(tx.Signature.Pubkey)
	var exist bool
	for _, addr := range addrs {
		if signAddr == addr {
			exist = true
			break
		}
	}

	if !exist {
		return x2eTy.ErrInvalidAdminAddress
	}

	if !tx.CheckSign(0) {
		return types.ErrSign
	}
	return nil
}
