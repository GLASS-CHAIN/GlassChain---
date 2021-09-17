package types

import (
	"errors"
	"reflect"

	"github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/js/types/jsproto"
)

// action for executor
const (
	jsActionCreate = 0
	jsActionCall   = 1
)

const (
	TyLogJs = 10000
)

const JsCreator = "js-creator"

var (
	typeMap = map[string]int32{
		"Create": jsActionCreate,
		"Call":   jsActionCall,
	}
	logMap = map[int64]*types.LogInfo{
		TyLogJs: {Ty: reflect.TypeOf(jsproto.JsLog{}), Name: "TyLogJs"},
	}
)

var JsX = "jsvm"

var (
	ErrDupName            = errors.New("ErrDupName")
	ErrJsReturnNotObject  = errors.New("ErrJsReturnNotObject")
	ErrJsReturnKVSFormat  = errors.New("ErrJsReturnKVSFormat")
	ErrJsReturnLogsFormat = errors.New("ErrJsReturnLogsFormat")

	ErrInvalidFuncFormat = errors.New("chain33.js: invalid function name format")
	//ErrInvalidFuncPrefix not exec_ execloal_ query_
	ErrInvalidFuncPrefix = errors.New("chain33.js: invalid function prefix format")

	ErrFuncNotFound = errors.New("chain33.js: invalid function name not found")
	ErrSymbolName   = errors.New("chain33.js: ErrSymbolName")
	ErrExecerName   = errors.New("chain33.js: ErrExecerName")
	ErrDBType       = errors.New("chain33.js: ErrDBType")
	// ErrJsCreator
	ErrJsCreator = errors.New("ErrJsCreator")
)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(JsX))
	types.RegFork(JsX, InitFork)
	types.RegExec(JsX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(JsX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(JsX, NewType(cfg))
}

type JsType struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *JsType {
	c := &JsType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

func (t *JsType) GetPayload() types.Message {
	return &jsproto.JsAction{}
}

func (t *JsType) GetTypeMap() map[string]int32 {
	return typeMap
}

func (t *JsType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
