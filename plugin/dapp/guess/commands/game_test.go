package commands

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/33cn/chain33/blockchain"
	"github.com/33cn/chain33/common/address"
	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/common/limits"
	"github.com/33cn/chain33/common/log"
	"github.com/33cn/chain33/executor"
	"github.com/33cn/chain33/mempool"
	"github.com/33cn/chain33/p2p"
	"github.com/33cn/chain33/queue"
	"github.com/33cn/chain33/rpc"
	"github.com/33cn/chain33/store"
	"github.com/33cn/chain33/system/consensus/solo"
	"github.com/33cn/chain33/types"
	"github.com/33cn/chain33/util"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	_ "github.com/33cn/chain33/system"
	_ "github.com/33cn/plugin/plugin/store/init"

	cty "github.com/33cn/chain33/system/dapp/coins/types"
	pty "github.com/33cn/plugin/plugin/dapp/norm/types"
)

var (
	secp crypto.Crypto

	jrpcURL = "127.0.0.1:8803"
	grpcURL = "127.0.0.1:8804"

	config = `# Title is local, Indicating that this configuration file is the configuration of a local single node. At this time, there is only one node in the chain where the local node is located, and the consensus module generally adopts the solo mode。
Title="local"
TestNet=true
FixTime=false
[crypto]
[log]
# debug(dbug)/info/warn/error(eror)/crit
loglevel = "info"
logConsoleLevel = "info"
logFile = "logs/chain33.log"
maxFileSize = 300
maxBackups = 100
maxAge = 28
localTime = true
compress = true
callerFile = false
callerFunction = false
[blockchain]
defCacheSize=128
maxFetchBlockNum=128
timeoutSeconds=5
driver="leveldb"
dbPath="datadir"
dbCache=64
singleMode=true
batchsync=false
isRecordBlockSequence=true
isParaChain=false
enableTxQuickIndex=false

[p2p]
types=["dht"]
msgCacheSize=10240
driver="leveldb"
dbPath="datadir/addrbook"
dbCache=4
grpcLogFile="grpc33.log"


[rpc]

jrpcBindAddr="localhost:8803"

grpcBindAddr="localhost:8804"

whitelist=["127.0.0.1"]

jrpcFuncWhitelist=["*"]
grpcFuncWhitelist=["*"]
enableTLS=false
certFile="cert.pem"
keyFile="key.pem"
[mempool]
name="timeline"
poolCacheSize=10240
minTxFeeRate=100000
maxTxNumPerAccount=10000
[mempool.sub.timeline]

poolCacheSize=10240
minTxFeeRate=100000
maxTxNumPerAccount=10000
[mempool.sub.score]
poolCacheSize=10240
minTxFeeRate=100000
maxTxNumPerAccount=10000
timeParam=1
priceConstant=1544
pricePower=1
[mempool.sub.price]
poolCacheSize=10240
minTxFeeRate=100000
maxTxNumPerAccount=10000
[consensus]
name="solo"
minerstart=true
genesisBlockTime=1514533394

genesis="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
[mver.consensus]

fundKeyAddr = "1BQXS6TxaYYG5mADaWij4AxhZZUTpw95a5"

coinReward = 18

coinDevFund = 12

ticketPrice = 10000

powLimitBits = "0x1f00ffff"
retargetAdjustmentFactor = 4
futureBlockTime = 16

ticketFrozenTime = 5    #5s only for test
ticketWithdrawTime = 10 #10s only for test
ticketMinerWaitTime = 2 #2s only for test

maxTxNumber = 1600      #160
targetTimespan = 2304

targetTimePerBlock = 16
[consensus.sub.solo]

genesis="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
genesisBlockTime=1514533394

waitTxMs=10
[store]
name="mavl"
driver="leveldb"
dbPath="datadir/mavltree"
dbCache=128
localdbVersion="1.0.0"
[store.sub.mavl]
enableMavlPrefix=false
enableMVCC=false
enableMavlPrune=false
pruneHeight=10000
[wallet]
minFee=100000
driver="leveldb"
dbPath="wallet"
dbCache=16
signType="secp256k1"
[wallet.sub.ticket]
minerdisable=false
minerwhitelist=["*"]
[exec]
enableStat=false
enableMVCC=false
alias=["token1:token","token2:token","token3:token"]
[exec.sub.token]

saveTokenTxList=true

tokenApprs = [
    "1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S",
    "1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK",
    "1LY8GFia5EiyoTodMLfkB5PHNNpXRqxhyB",
    "1GCzJDS6HbgTQ2emade7mEJGGWFfA15pS9",
    "1JYB8sxi4He5pZWHCd3Zi2nypQ4JMB6AxN",
    "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv",
]
[exec.sub.cert]

enable=false

cryptoPath="authdir/crypto"

signType="auth_ecdsa"
[exec.sub.relay]

genesis="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
[exec.sub.manage]

superManager=[
    "1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S", 
    "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv", 
    "1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK"
]
`
)

var (
	random *rand.Rand

	loopCount = 1
	conn      *grpc.ClientConn
	c         types.Chain33Client
	adminPriv = "CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"
	adminAddr = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"

	//userAPubkey = "03EF0E1D3112CF571743A3318125EDE2E52A4EB904BCBAA4B1F75020C2846A7EB4"
	userAAddr = "14BqNZoysh5Br8wiyJLyLu6aTzkFhbvZNm"
	userAPriv = "880D055D311827AB427031959F43238500D08F3EDF443EB903EC6D1A16A5783A"

	//userBPubkey = "027848E7FA630B759DB406940B5506B666A344B1060794BBF314EB459D40881BB3"
	userBAddr = "1PLxDdfdLxpsgR52T6DGDKQSTeJULR2Gre"
	userBPriv = "2C1B7E0B53548209E6915D8C5EB3162E9AF43D051C1316784034C9D549DD5F61"
)

const fee = 1e6

func init() {
	err := limits.SetLimits()
	if err != nil {
		panic(err)
	}
	log.SetLogLevel("info")
	random = rand.New(rand.NewSource(types.Now().UnixNano()))

	cr2, err := crypto.New(types.GetSignName("", types.SECP256K1))
	if err != nil {
		fmt.Println("crypto.New failed for types.ED25519")
		return
	}
	secp = cr2
}

func Init() {
	fmt.Println("=======Init Data1!=======")
	os.RemoveAll("datadir")
	os.RemoveAll("wallet")
	os.Remove("chain33.test.toml")

	ioutil.WriteFile("chain33.test.toml", []byte(config), 0664)
}

func clearTestData() {
	fmt.Println("=======start clear test data!=======")

	os.Remove("chain33.test.toml")
	os.RemoveAll("wallet")
	err := os.RemoveAll("datadir")
	if err != nil {
		fmt.Println("delete datadir have a err:", err.Error())
	}

	fmt.Println("test data clear successfully!")
}

func TestGuess(t *testing.T) {
	Init()
	testGuessImp(t)
	fmt.Println("=======start clear test data!=======")
	clearTestData()
}

func testGuessImp(t *testing.T) {
	fmt.Println("=======start guess test!=======")
	q, chain, s, mem, exec, cs, p2p, cmd := initEnvGuess()
	cfg := q.GetConfig()
	defer chain.Close()
	defer mem.Close()
	defer exec.Close()
	defer s.Close()
	defer q.Close()
	defer cs.Close()
	defer p2p.Close()
	err := createConn()
	for err != nil {
		err = createConn()
	}
	time.Sleep(2 * time.Second)
	fmt.Println("=======start NormPut!=======")

	for i := 0; i < loopCount; i++ {
		NormPut(cfg)
		time.Sleep(time.Second)
	}

	fmt.Println("=======start sendTransferTx sendTransferToExecTx!=======")

	sendTransferTx(cfg, adminPriv, userAAddr, 2000000000000)
	sendTransferTx(cfg, adminPriv, userBAddr, 2000000000000)

	time.Sleep(2 * time.Second)
	in := &types.ReqBalance{}
	in.Addresses = append(in.Addresses, userAAddr)
	acct, err1 := c.GetBalance(context.Background(), in)
	if err1 != nil || len(acct.Acc) == 0 {
		fmt.Println("no balance for ", userAAddr)
	} else {
		fmt.Println(userAAddr, " balance:", acct.Acc[0].Balance, "frozen:", acct.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct.Acc[0].Balance == 2000000000000)

	in2 := &types.ReqBalance{}
	in2.Addresses = append(in.Addresses, userBAddr)
	acct2, err2 := c.GetBalance(context.Background(), in2)
	if err2 != nil || len(acct2.Acc) == 0 {
		fmt.Println("no balance for ", userBAddr)
	} else {
		fmt.Println(userBAddr, " balance:", acct2.Acc[0].Balance, "frozen:", acct2.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct2.Acc[0].Balance == 2000000000000)


	sendTransferToExecTx(cfg, userAPriv, "guess", 1000000000000)
	sendTransferToExecTx(cfg, userBPriv, "guess", 1000000000000)
	time.Sleep(2 * time.Second)

	fmt.Println("=======start GetBalance!=======")

	in3 := &types.ReqBalance{}
	in3.Addresses = append(in3.Addresses, userAAddr)
	acct3, err3 := c.GetBalance(context.Background(), in3)
	if err3 != nil || len(acct3.Acc) == 0 {
		fmt.Println("no balance for ", userAAddr)
	} else {
		fmt.Println(userAAddr, " balance:", acct3.Acc[0].Balance, "frozen:", acct3.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct3.Acc[0].Balance == 1000000000000)

	in4 := &types.ReqBalance{}
	in4.Addresses = append(in4.Addresses, userBAddr)
	acct4, err4 := c.GetBalance(context.Background(), in4)
	if err4 != nil || len(acct4.Acc) == 0 {
		fmt.Println("no balance for ", userBAddr)
	} else {
		fmt.Println(userBAddr, " balance:", acct4.Acc[0].Balance, "frozen:", acct4.Acc[0].Frozen)
	}
	assert.Equal(t, true, acct4.Acc[0].Balance == 1000000000000)
	testCmd(cmd)
	time.Sleep(2 * time.Second)
}

func initEnvGuess() (queue.Queue, *blockchain.BlockChain, queue.Module, queue.Module, *executor.Executor, queue.Module, queue.Module, *cobra.Command) {
	flag.Parse()

	chain33Cfg := types.NewChain33Config(types.ReadFile("chain33.test.toml"))
	var q = queue.New("channel")
	q.SetConfig(chain33Cfg)
	cfg := chain33Cfg.GetModuleConfig()
	cfg.Log.LogFile = ""
	sub := chain33Cfg.GetSubConfig()
	chain := blockchain.New(chain33Cfg)
	chain.SetQueueClient(q.Client())

	exec := executor.New(chain33Cfg)
	exec.SetQueueClient(q.Client())
	chain33Cfg.SetMinFee(0)
	s := store.New(chain33Cfg)
	s.SetQueueClient(q.Client())

	cs := solo.New(cfg.Consensus, sub.Consensus["solo"])
	cs.SetQueueClient(q.Client())

	mem := mempool.New(chain33Cfg)
	mem.SetQueueClient(q.Client())
	network := p2p.NewP2PMgr(chain33Cfg)

	network.SetQueueClient(q.Client())

	rpc.InitCfg(cfg.RPC)
	gapi := rpc.NewGRpcServer(q.Client(), nil)
	go gapi.Listen()

	japi := rpc.NewJSONRPCServer(q.Client(), nil)
	go japi.Listen()

	cmd := GuessCmd()
	return q, chain, s, mem, exec, cs, network, cmd
}

func createConn() error {
	var err error
	url := grpcURL
	fmt.Println("grpc url:", url)
	conn, err = grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	c = types.NewChain33Client(conn)
	return nil
}

func generateKey(i, valI int) string {
	key := make([]byte, valI)
	binary.PutUvarint(key[:10], uint64(valI))
	binary.PutUvarint(key[12:24], uint64(i))
	if _, err := rand.Read(key[24:]); err != nil {
		os.Exit(1)
	}
	return string(key)
}

func generateValue(i, valI int) string {
	value := make([]byte, valI)
	binary.PutUvarint(value[:16], uint64(i))
	binary.PutUvarint(value[32:128], uint64(i))
	if _, err := rand.Read(value[128:]); err != nil {
		os.Exit(1)
	}
	return string(value)
}

func getprivkey(key string) crypto.PrivKey {
	bkey, err := hex.DecodeString(key)
	if err != nil {
		panic(err)
	}
	priv, err := secp.PrivKeyFromBytes(bkey)
	if err != nil {
		panic(err)
	}
	return priv
}

func prepareTxList(cfg *types.Chain33Config) *types.Transaction {
	var key string
	var value string
	var i int

	key = generateKey(i, 32)
	value = generateValue(i, 180)

	nput := &pty.NormAction_Nput{Nput: &pty.NormPut{Key: []byte(key), Value: []byte(value)}}
	action := &pty.NormAction{Value: nput, Ty: pty.NormActionPut}
	tx := &types.Transaction{Execer: []byte("norm"), Payload: types.Encode(action), Fee: fee}
	tx.To = address.ExecAddress("norm")
	tx.Nonce = random.Int63()
	tx.ChainID = cfg.GetChainID()
	tx.Sign(types.SECP256K1, getprivkey("CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944"))
	return tx
}

func NormPut(cfg *types.Chain33Config) {
	tx := prepareTxList(cfg)

	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		return
	}
}

func sendTransferTx(cfg *types.Chain33Config, fromKey, to string, amount int64) bool {
	signer := util.HexToPrivkey(fromKey)
	var tx *types.Transaction
	transfer := &cty.CoinsAction{}
	v := &cty.CoinsAction_Transfer{Transfer: &types.AssetsTransfer{Amount: amount, Note: []byte(""), To: to}}
	transfer.Value = v
	transfer.Ty = cty.CoinsActionTransfer
	execer := []byte(cfg.GetCoinExec())
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(transfer), To: to, Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("in sendTransferTx formatTx failed")
		return false
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendTransferTx SendTransaction failed")

		return false
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendTransferTx SendTransaction failed,reply not ok.")

		return false
	}
	fmt.Println("sendTransferTx ok")

	return true
}

func sendTransferToExecTx(cfg *types.Chain33Config, fromKey, execName string, amount int64) bool {
	signer := util.HexToPrivkey(fromKey)
	var tx *types.Transaction
	transfer := &cty.CoinsAction{}
	execAddr := address.ExecAddress(execName)
	v := &cty.CoinsAction_TransferToExec{TransferToExec: &types.AssetsTransferToExec{Amount: amount, Note: []byte(""), ExecName: execName, To: execAddr}}
	transfer.Value = v
	transfer.Ty = cty.CoinsActionTransferToExec
	execer := []byte(cfg.GetCoinExec())
	tx = &types.Transaction{Execer: execer, Payload: types.Encode(transfer), To: address.ExecAddress("guess"), Fee: fee}
	tx, err := types.FormatTx(cfg, string(execer), tx)
	if err != nil {
		fmt.Println("sendTransferToExecTx formatTx failed.")

		return false
	}

	tx.Sign(types.SECP256K1, signer)
	reply, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Println("in sendTransferToExecTx SendTransaction failed")

		return false
	}
	if !reply.IsOk {
		fmt.Fprintln(os.Stderr, errors.New(string(reply.GetMsg())))
		fmt.Println("in sendTransferToExecTx SendTransaction failed,reply not ok.")

		return false
	}

	fmt.Println("sendTransferToExecTx ok")

	return true
}

func testCmd(cmd *cobra.Command) {
	var rootCmd = &cobra.Command{
		Use:   "chain33-cli",
		Short: "chain33 client tools",
	}

	chain33Cfg := types.NewChain33Config(types.ReadFile("chain33.test.toml"))
	types.SetCliSysParam(chain33Cfg.GetTitle(), chain33Cfg)

	rootCmd.PersistentFlags().String("title", chain33Cfg.GetTitle(), "get title name")

	rootCmd.PersistentFlags().String("rpc_laddr", "http://"+grpcURL, "http url")
	rootCmd.AddCommand(cmd)

	rootCmd.SetArgs([]string{"guess", "start", "--topic", "WorldCup Final", "--options", "A:France;B:Claodia", "--category", "football", "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	strGameID := "0x20d78aabfa67701f7492261312575614f684b084efa2c13224f1505706c846e8"

	rootCmd.SetArgs([]string{"guess", "bet", "--gameId", strGameID, "--option", "A", "--betsNumber", "500000000", "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "stop", "--gameId", strGameID, "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "abort", "--gameId", strGameID, "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "publish", "--gameId", strGameID, "--result", "A", "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "query", "--type", "ids", "--gameIDs", strGameID, "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "query", "--type", "id", "--gameId", strGameID, "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "query", "--type", "addr", "--addr", userAAddr, "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "query", "--type", "status", "--status", "10", "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "query", "--type", "adminAddr", "--adminAddr", adminAddr, "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "query", "--type", "addrStatus", "--addr", userAAddr, "--status", "11", "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "query", "--type", "adminStatus", "--adminAddr", adminAddr, "--status", "11", "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()

	rootCmd.SetArgs([]string{"guess", "query", "--type", "categoryStatus", "--category", "football", "--status", "11", "--rpc_laddr", "http://" + jrpcURL})
	rootCmd.Execute()
}
