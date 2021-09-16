// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mm

import (
	"fmt"

	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/params"
)

type (
	StackValidationFunc func(*Stack) error
)

func MakeStackFunc(pop, push int) StackValidationFunc {
	return func(stack *Stack) error {
		if err := stack.Require(pop); err != nil {
			return err
		}

		if stack.Len()+push-pop > int(params.StackLimit) {
			return fmt.Errorf("stack limit reached %d (%d)", stack.Len(), params.StackLimit)
		}
		return nil
	}
}

func MakeDupStackFunc(n int) StackValidationFunc {
	return MakeStackFunc(n, n+1)
}

func MakeSwapStackFunc(n int) StackValidationFunc {
	return MakeStackFunc(n, n)
}

func minSwapStack(n int) int {
	return minStack(n, n)
}
func maxSwapStack(n int) int {
	return maxStack(n, n)
}

func minDupStack(n int) int {
	return minStack(n, n+1)
}
func maxDupStack(n int) int {
	return maxStack(n, n+1)
}

func maxStack(pop, push int) int {
	return int(params.StackLimit) + pop - push
}
func minStack(pops, push int) int {
	return pops
}
