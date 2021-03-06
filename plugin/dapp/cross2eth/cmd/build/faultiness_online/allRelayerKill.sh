#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"
source "./relayerPublic.sh"

ethDeployAddr="0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a"
ethDeployKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"
ethValidatorAddrKeyA="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

chain33DeployAddr="1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ"

chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"

ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
#ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"

maturityDegree=10

Chain33Cli="../../chain33-cli"
chain33BridgeBank=""
ethBridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
chain33EthTokenAddr=""
ethereumBtyTokenAddr=""
chain33YccTokenAddr=""
ethereumYccTokenAddr=""
chain33ZBCErc20Addr=""
ethBridgeToeknZBCAddr=""
chain33ID=0

CLIA="./ebcli_A"

# chain33 lock BTY, eth burn BTY
function TestChain33ToEthAssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # chain33 lock bty
    hash=$(${Chain33Cli} send evm call -f 1 -a 5 -k "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33BtyTokenAddr}, 500000000)" --chainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "5.0000"

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    kill_ebrelayerC

    # eth burn
    result=$(${CLIA} ethereum burn -m 3 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumBtyTokenAddr}") #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    result=$(${Chain33Cli} account balance -a "${chain33ReceiverAddr}" -e evm)
    balance_ret "${result}" "0.0000"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${Chain33Cli} account balance -a "${chain33ReceiverAddr}" -e evm)
    balance_ret "${result}" "3.0000"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "2.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chain33 lock ZBC, eth burn ZBC
function TestChain33ToEthZBCAssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    echo -e "${GRE}=========== chain33 lock ZBC, eth burn ZBC ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${Chain33Cli} evm query -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # chain33 lock ZBC
    hash=$(${Chain33Cli} send evm call -f 1 -k "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33ZBCErc20Addr}, 900000000)" --chainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} evm query -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "900000000"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "9"

    kill_ebrelayerC

    # eth burn
    result=$(${CLIA} ethereum burn -m 8 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethBridgeToeknZBCAddr}") #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethBridgeToeknZBCAddr}")
    cli_ret "${result}" "balance" ".balance" "1"

    sleep ${maturityDegree}

    result=$(${Chain33Cli} evm query -a "${chain33ZBCErc20Addr}" -c "${chain33ReceiverAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "0"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${Chain33Cli} evm query -a "${chain33ZBCErc20Addr}" -c "${chain33ReceiverAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    is_equal "${result}" "800000000"

    result=$(${Chain33Cli} evm query -a "${chain33ZBCErc20Addr}" -c "${chain33BridgeBank}" -b "balanceOf(${chain33BridgeBank})")
    is_equal "${result}" "100000000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33AssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    result=$(${CLIA} ethereum lock -m 11 -k "${ethValidatorAddrKeyA}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "11"

    sleep ${maturityDegree}

    result=$(${Chain33Cli} evm query -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    is_equal "${result}" "0"

    start_ebrelayerC

    result=$(${Chain33Cli} evm query -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    is_equal "${result}" "1100000000"

    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "100"

    kill_ebrelayerC

    echo '#5.burn ETH from Chain33 ETH(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33EthTokenAddr}"

    sleep ${maturityDegree}

    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "100"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}")
    cli_ret "${result}" "balance" ".balance" "105"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm query -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    is_equal "${result}" "600000000"

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "6"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33YccKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
 
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    result=$(${CLIA} ethereum lock -m 7 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "lock"

    eth_block_wait 2
    sleep ${maturityDegree}

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "7"

    result=$(${Chain33Cli} evm query -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    is_equal "${result}" "0"

    start_ebrelayerC

    result=$(${Chain33Cli} evm query -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    is_equal "${result}" "700000000"

    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC

    echo '#5.burn YCC from Chain33 YCC(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethReceiverAddr1}" -t "${chain33YccTokenAddr}"

    sleep ${maturityDegree}

    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${CLIA} ethereum balance -o "${ethReceiverAddr1}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm query -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    is_equal "${result}" "200000000"

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# shellcheck disable=SC2120
function mainTest() {
    if [[ $# -ge 1 && ${1} != "" ]]; then
        chain33ID="${1}"
    fi
    StartChain33
    start_trufflesuite
    AllRelayerStart

    TestChain33ToEthAssetsKill
    TestChain33ToEthZBCAssetsKill
    TestETH2Chain33AssetsKill
    TestETH2Chain33YccKill
}

mainTest "${1}"
