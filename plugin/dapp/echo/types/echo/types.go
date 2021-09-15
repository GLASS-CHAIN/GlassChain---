package echo

import (
	"reflect"

	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
)

const (
	ActionPing = iota
	ActionPang
)

const (
	TyLogPing = 100001
	TyLogPang = 100002
)

var (
	EchoX = "echo"
	actionName = map[string]int32{
		"Ping": ActionPing,
		"Pang": ActionPang,
	}
	logInfo = map[int64]*types.LogInfo{
		TyLogPing: {Ty: reflect.TypeOf(PingLog{}), Name: "PingLog"},
		TyLogPang: {Ty: reflect.TypeOf(PangLog{}), Name: "PangLog"},
	}
)
var elog = log.New("module", EchoX)

func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(EchoX))
	types.RegFork(EchoX, InitFork)
	types.RegExec(EchoX, InitExecutor)
}

//InitFork ...
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(EchoX, "Enable", 0)
}

//InitExecutor ...
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(EchoX, NewType(cfg))
}

type Type struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *Type {
	c := &Type{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

func (b *Type) GetPayload() types.Message {
	return &EchoAction{}
}

func (b *Type) GetName() string {
	return EchoX
}

func (b *Type) GetTypeMap() map[string]int32 {
	return actionName
}

func (b *Type) GetLogMap() map[int64]*types.LogInfo {
	return logInfo
}
