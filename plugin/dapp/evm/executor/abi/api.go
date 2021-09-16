package abi

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"

	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/golang-collections/collections/stack"
)

func Pack(param, abiData string, readOnly bool) (methodName string, packData []byte, err error) {
	methodName, params, err := procFuncCall(param)

	if err != nil {
		return methodName, packData, err
	}

	abi, err := JSON(strings.NewReader(abiData))
	if err != nil {
		return methodName, packData, err
	}

	var method Method
	var ok bool
	if method, ok = abi.Methods[methodName]; !ok {
		err = fmt.Errorf("function %v not exists", methodName)
		return methodName, packData, err
	}

	if readOnly && !method.IsConstant() {
		return methodName, packData, errors.New("method is not readonly")
	}
	if len(params) != method.Inputs.LengthNonIndexed() {
		err = fmt.Errorf("function params error:%v", params)
		return methodName, packData, err
	}
	paramVals := []interface{}{}
	if len(params) != 0 {
		if method.Inputs.LengthNonIndexed() != len(params) {
			err = fmt.Errorf("function Params count error: %v", param)
			return methodName, packData, err
		}

		for i, v := range method.Inputs.NonIndexed() {
			paramVal, err := str2GoValue(v.Type, params[i])
			if err != nil {
				return methodName, packData, err
			}
			paramVals = append(paramVals, paramVal)
		}
	}

	packData, err = abi.Pack(methodName, paramVals...)
	return methodName, packData, err
}

func PackContructorPara(param, abiStr string) (packData []byte, err error) {
	_, params, err := procFuncCall(param)
	if err != nil {
		return nil, err
	}

	parsedAbi, err := JSON(strings.NewReader(abiStr))
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse evm code error", err)
		return
	}

	method := parsedAbi.Constructor

	paramVals := []interface{}{}
	if len(params) != 0 {
		if method.Inputs.LengthNonIndexed() != len(params) {
			err = fmt.Errorf("function Params count error: %v", param)
			return nil, err
		}

		for i, v := range method.Inputs.NonIndexed() {
			paramVal, err := str2GoValue(v.Type, params[i])
			if err != nil {
				return nil, err
			}
			paramVals = append(paramVals, paramVal)
		}
	}
	packData, err = parsedAbi.Constructor.Inputs.Pack(paramVals...)
	if err != nil {
		return nil, err
	}
	return packData, nil

}

func Unpack(data []byte, methodName, abiData string) (output []*Param, err error) {
	if len(data) == 0 {
		log.Info("Unpack", "Data len", 0, "methodName", methodName)
		return output, err
	}
	abi, err := JSON(strings.NewReader(abiData))
	if err != nil {
		return output, err
	}

	var method Method
	var ok bool
	if method, ok = abi.Methods[methodName]; !ok {
		return output, fmt.Errorf("function %v not exists", methodName)
	}

	if method.Outputs.LengthNonIndexed() == 0 {
		return output, err
	}

	values, err := method.Outputs.UnpackValues(data)
	if err != nil {
		return output, err
	}

	output = []*Param{}

	for i, v := range values {
		arg := method.Outputs[i]
		pval := &Param{Name: arg.Name, Type: arg.Type.String(), Value: v}
		if arg.Type.String() == "address" {
			pval.Value = v.(common.Hash160Address).ToAddress().String()
			log.Info("Unpack address", "address", pval.Value)
		}

		output = append(output, pval)
	}

	return
}

type Param struct {
	Name string `json:"name"`

	Type string `json:"type"`

	Value interface{} `json:"value"`
}

func convertUint(val uint64, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Uint:
		return uint(val)
	case reflect.Uint8:
		return uint8(val)
	case reflect.Uint16:
		return uint16(val)
	case reflect.Uint32:
		return uint32(val)
	case reflect.Uint64:
		return val
	}
	return val
}

func convertInt(val int64, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Int:
		return int(val)
	case reflect.Int8:
		return int8(val)
	case reflect.Int16:
		return int16(val)
	case reflect.Int32:
		return int32(val)
	case reflect.Int64:
		return val
	}
	return val
}

func str2GoValue(typ Type, val string) (res interface{}, err error) {
	switch typ.T {
	case IntTy:
		if typ.Size < 256 {
			x, err := strconv.ParseInt(val, 10, typ.Size)
			if err != nil {
				return res, err
			}
			return convertInt(x, typ.GetType().Kind()), nil
		}
		b := new(big.Int)
		b.SetString(val, 10)
		return b, err
	case UintTy:
		if typ.Size < 256 {
			x, err := strconv.ParseUint(val, 10, typ.Size)
			if err != nil {
				return res, err
			}
			return convertUint(x, typ.GetType().Kind()), nil
		}
		b := new(big.Int)
		b.SetString(val, 10)
		return b, err
	case BoolTy:
		x, err := strconv.ParseBool(val)
		if err != nil {
			return res, err
		}
		return x, nil
	case StringTy:
		return val, nil
	case SliceTy:
		subs, err := procArrayItem(val)
		if err != nil {
			return res, err
		}
		rval := reflect.MakeSlice(typ.GetType(), len(subs), len(subs))
		for idx, sub := range subs {
			subVal, er := str2GoValue(*typ.Elem, sub)
			if er != nil {
				//return res, er
				subparams, err := procArrayItem(sub) 
				if err != nil {
					return res, er
				}
				fmt.Println("subparams len", len(subparams), "subparams", subparams)

				abityps := strings.Split(typ.Elem.stringKind[1:len(typ.Elem.stringKind)-1], ",") 
				fmt.Println("abityps", abityps)


				for i := 0; i < len(subparams); i++ {
					tp, err := NewType(abityps[i], "", nil)
					if err != nil {
						fmt.Println("NewType", err.Error())
						continue
					}
					fmt.Println("tp", tp.stringKind, "types", tp.T)
					subVal, err = str2GoValue(tp, subparams[i]) 
					if err != nil {
						fmt.Println("str2GoValue err ", err.Error())
						continue
					}

					fmt.Println("subVal", subVal)
					rval.Index(idx).Field(i).Set(reflect.ValueOf(subVal))

				}

				return rval.Interface(), nil

			}
			rval.Index(idx).Set(reflect.ValueOf(subVal))
		}
		return rval.Interface(), nil
	case ArrayTy:
		rval := reflect.New(typ.GetType()).Elem()
		subs, err := procArrayItem(val)
		if err != nil {
			return res, err
		}
		for idx, sub := range subs {
			subVal, er := str2GoValue(*typ.Elem, sub)
			if er != nil {
				return res, er
			}
			rval.Index(idx).Set(reflect.ValueOf(subVal))
		}
		return rval.Interface(), nil
	case AddressTy:
		addr := common.StringToAddress(val)
		if addr == nil {
			return res, fmt.Errorf("invalid  address: %v", val)
		}
		return addr.ToHash160(), nil
	case FixedBytesTy:
		x, err := common.HexToBytes(val)
		if err != nil {
			return res, err
		}
		rval := reflect.New(typ.GetType()).Elem()
		for i, b := range x {
			rval.Index(i).Set(reflect.ValueOf(b))
		}
		return rval.Interface(), nil
	case BytesTy:
		x, err := common.HexToBytes(val)
		if err != nil {
			return res, err
		}
		return x, nil
	case HashTy:
		x, err := common.HexToBytes(val)
		if err != nil {
			return res, err
		}
		return common.BytesToHash(x), nil
	default:
		return res, fmt.Errorf("not support type: %v", typ.stringKind)
	}
}


func procArrayItem(val string) (res []string, err error) {
	ss := stack.New()
	data := []rune{}
	for _, b := range val {
		switch b {
		case ' ':
			if ss.Len() > 0 && peekRune(ss) == '"' {
				data = append(data, b)
			}
		case ',':
			if ss.Len() == 1 && peekRune(ss) == '[' {

				res = append(res, string(data))
				data = []rune{}

			} else {
				data = append(data, b)
			}
		case '"':
			if ss.Peek() == b {
				ss.Pop()
			} else {
				ss.Push(b)
			}
			//data = append(data, b)
		case '[':
			if ss.Len() == 0 {
				data = []rune{}
			} else {
				data = append(data, b)
			}
			ss.Push(b)
		case ']':
			if ss.Len() == 1 && peekRune(ss) == '[' {

				res = append(res, string(data))
			} else {
				data = append(data, b)
			}
			ss.Pop()
		default:
			data = append(data, b)
		}
	}

	if ss.Len() != 0 {
		return nil, fmt.Errorf("invalid array format:%v", val)
	}
	return res, err
}

func peekRune(ss *stack.Stack) rune {
	return ss.Peek().(rune)
}


func procFuncCall(param string) (funcName string, res []string, err error) {
	lidx := strings.Index(param, "(")
	ridx := strings.LastIndex(param, ")")

	if lidx == -1 || ridx == -1 {
		return funcName, res, fmt.Errorf("invalid function signature:%v", param)
	}

	funcName = strings.TrimSpace(param[:lidx])
	params := strings.TrimSpace(param[lidx+1 : ridx])

	if len(params) > 0 {
		res, err = procArrayItem(fmt.Sprintf("[%v]", params))
	}

	return funcName, res, err
}
