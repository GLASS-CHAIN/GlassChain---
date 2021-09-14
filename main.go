// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.9

/*
The functionality of each system is accomplished through plug-ins, which are divided into four categories:

Consensus Encrypted Dapp Storage

This go package provides officially available plug-ins.
*/
package main

import (
	_ "github.com/33cn/chain33/system"
	"github.com/33cn/chain33/util/cli"
	_ "github.com/33cn/plugin/plugin"
)

func main() {
	cli.RunChain33("", "")
}
