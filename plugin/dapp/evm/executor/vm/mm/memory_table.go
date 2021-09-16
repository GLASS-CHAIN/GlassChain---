// Copyright Fuzamei Corp. 2018 All Rights Reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mm


type (
	MemorySizeFunc func(*Stack) (size uint64, overflow bool)
)

func MemorySha3(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(0), stack.Back(1))
}

func MemoryCallDataCopy(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(0), stack.Back(2))
}

func MemoryReturnDataCopy(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(0), stack.Back(2))
}

func MemoryCodeCopy(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(0), stack.Back(2))
}

func MemoryExtCodeCopy(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(1), stack.Back(3))
}

func MemoryMLoad(stack *Stack) (uint64, bool) {
	return calcMemSize64WithUint(stack.Back(0), 32)
}

func MemoryMStore8(stack *Stack) (uint64, bool) {
	return calcMemSize64WithUint(stack.Back(0), 1)
}

func MemoryMStore(stack *Stack) (uint64, bool) {
	return calcMemSize64WithUint(stack.Back(0), 32)
}

func MemoryCreate(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(1), stack.Back(2))
}

//func memoryCreate2(stack *Stack) (uint64, bool) {
//	return calcMemSize64(stack.Back(1), stack.Back(2))
//}

func MemoryCall(stack *Stack) (uint64, bool) {
	x, overflow := calcMemSize64(stack.Back(5), stack.Back(6))
	if overflow {
		return 0, true
	}
	y, overflow := calcMemSize64(stack.Back(3), stack.Back(4))
	if overflow {
		return 0, true
	}
	if x > y {
		return x, false
	}
	return y, false
}

func MemoryDelegateCall(stack *Stack) (uint64, bool) {
	x, overflow := calcMemSize64(stack.Back(4), stack.Back(5))
	if overflow {
		return 0, true
	}
	y, overflow := calcMemSize64(stack.Back(2), stack.Back(3))
	if overflow {
		return 0, true
	}
	if x > y {
		return x, false
	}
	return y, false
}

func MemoryStaticCall(stack *Stack) (uint64, bool) {
	x, overflow := calcMemSize64(stack.Back(4), stack.Back(5))
	if overflow {
		return 0, true
	}
	y, overflow := calcMemSize64(stack.Back(2), stack.Back(3))
	if overflow {
		return 0, true
	}
	if x > y {
		return x, false
	}
	return y, false
}

func MemoryReturn(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(0), stack.Back(1))
}

func MemoryRevert(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(0), stack.Back(1))
}

func MemoryLog(stack *Stack) (uint64, bool) {
	return calcMemSize64(stack.Back(0), stack.Back(1))
}
