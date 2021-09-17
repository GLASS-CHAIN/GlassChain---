package testnode

//DefaultConfig default config for testnode
var DefaultConfig = `
Title="user.p.test."
CoinSymbol="bty"
ChainID=33
# TestNet=true

[crypto]
[log]
#  debug(dbug)/info/warn/error(eror)/crit
loglevel = "debug"
logConsoleLevel = "info"
#   
logFile = "logs/chain33.para.log"
#   ）
maxFileSize = 300
# 
maxBackups = 100
#   ）
maxAge = 28
#  UT ）
localTime = true
#  gz）
compress = true
# 
callerFile = false
# 
callerFunction = false


[blockchain]
defCacheSize=128
maxFetchBlockNum=128
timeoutSeconds=5
batchBlockNum=128
driver="leveldb"
dbPath="paradatadir"
dbCache=64
isStrongConsistency=true
singleMode=true
batchsync=false
isRecordBlockSequence=false
isParaChain = true
enableTxQuickIndex=false

[p2p]
enable=false
msgCacheSize=10240
driver="leveldb"
dbPath="paradatadir/addrbook"
dbCache=4
grpcLogFile="grpc33.log"


[rpc]
# 
jrpcBindAddr="localhost:8901"
grpcBindAddr="localhost:8902"
whitelist=["127.0.0.1"]
jrpcFuncWhitelist=["*"]
grpcFuncWhitelist=["*"]


[mempool]
name="timeline"
poolCacheSize=10240
minTxFeeRate=100000
maxTxNumPerAccount=10000

[mempool.sub.para]
poolCacheSize=102400

[consensus]
name="para"
genesisBlockTime=1514533390
genesis="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
minerExecs=["paracross"]

[mver.consensus]
fundKeyAddr = "1BQXS6TxaYYG5mADaWij4AxhZZUTpw95a5"
powLimitBits = "0x1f00ffff"
maxTxNumber = 1600      #160


[mver.consensus.ticket]
coinReward = 18
coinDevFund = 12
ticketPrice = 10000
retargetAdjustmentFactor = 4
futureBlockTime = 16
ticketFrozenTime = 5    #5s only for test
ticketWithdrawTime = 10 #10s only for test
ticketMinerWaitTime = 2 #2s only for test
targetTimespan = 2304
targetTimePerBlock = 16

[mver.consensus.paracross]
coinReward = 18
coinDevFund = 12
minerMode="normal"


[consensus.sub.para]
 grp ip i  “101.37.227.226:8802,39.97.20.242:8802,47.107.15.126:8802,jiedian2.33.cn”
ParaRemoteGrpcClient=""
 
startHeight=1
  
writeBlockSeconds=2
    
authAccount="1EbDHAXpoiewjPLX9uqoz38HsKqMXayZrF"
   2
waitBlocks4CommitMsg=2
  block blockhas 
searchHashMatchedBlockDepth=10000
 
genesisAmount=100000000
mainBlockHashForkHeight=1
mainForkParacrossCommitTx=1
mainLoopCheckCommitTxDoneForkHeight=11
selfConsensEnablePreContract=["0-1000"]
emptyBlockInterval=["0:2"]


[store]
name="mavl"
driver="leveldb"
dbPath="paradatadir/mavltree"
dbCache=128
enableMavlPrefix=false
enableMVCC=false
enableMavlPrune=false
pruneHeight=10000

[wallet]
minFee=100000
driver="leveldb"
dbPath="parawallet"
dbCache=16
signType="secp256k1"
minerdisable=true

[exec]
enableStat=false

[exec.sub.relay]
genesis="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"

[exec.sub.manage]
superManager=[
    "1Bsg9j6gW83sShoee1fZAt9TkUjcrCgA9S",
    "12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv",
    "1Q8hGLfoGe63efeWa8fJ4Pnukhkngt6poK"
]

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


[pprof]
listenAddr = "localhost:6062"
`
