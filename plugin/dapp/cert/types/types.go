// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package types

import "github.com/33cn/chain33/types"

//cert
const (
	CertActionNew    = 1
	CertActionUpdate = 2
	CertActionNormal = 3

	AuthECDSA = 257
	AuthSM2   = 258
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, ExecerCert)
	// init executor type
	types.RegFork(CertX, InitFork)
	types.RegExec(CertX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(CertX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(CertX, NewType(cfg))
}


type CertType struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *CertType {
	c := &CertType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

func (b *CertType) GetPayload() types.Message {
	return &CertAction{}
}

func (b *CertType) GetName() string {
	return CertX
}

func (b *CertType) GetLogMap() map[int64]*types.LogInfo {
	return nil
}

func (b *CertType) GetTypeMap() map[string]int32 {
	return actionName
}
