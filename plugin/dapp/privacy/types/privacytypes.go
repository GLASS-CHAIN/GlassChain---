// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "github.com/33cn/chain33/types"

// IsExpire FTX  true
func (ftxos *FTXOsSTXOsInOneTx) IsExpire(blockheight, blocktime int64) bool {
	valid := ftxos.GetExpire()
	if valid == 0 {
		// Expir 0 false
		return false
	}
	// expireBoun 
	if valid <= types.ExpireBound {
		return valid <= blockheight
	}
	return valid <= blocktime
}

// SetExpire 
func (ftxos *FTXOsSTXOsInOneTx) SetExpire(expire int64) {
	if expire > types.ExpireBound {
		// FTX  T 1 
		ftxos.Expire = expire + 12
	} else {
		// FTX  T +1
		ftxos.Expire = expire + 1
	}
}
