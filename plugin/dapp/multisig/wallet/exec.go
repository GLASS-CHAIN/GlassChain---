// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wallet

import (
	"github.com/33cn/chain33/types"
)

//On_MultiSigAddresList owne 
func (policy *multisigPolicy) On_MultiSigAddresList(req *types.ReqString) (types.Message, error) {
	/ 
	if req.Data == "" {
		reply, err := policy.store.listOwnerAttrs()
		if err != nil {
			bizlog.Error("On_MultiSigAddresList  listOwnerAttrs", "err", err)
		}
		return reply, err
	}
	/ owne 
	reply, err := policy.store.listOwnerAttrsByAddr(req.Data)
	if err != nil {
		bizlog.Error("On_MultiSigAddresList listOwnerAttrsByAddr", "owneraddr", req.Data, "err", err)
	}
	return reply, err
}
