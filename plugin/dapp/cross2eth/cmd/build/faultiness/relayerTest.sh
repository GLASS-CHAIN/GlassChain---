#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"
source "./relayerPublic.sh"

chain33DeployAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
#chain33DeployKey="0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944"

ethDeployAddr="0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a"
ethDeployKey="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"

# validatorsAddr=["0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
#ethValidatorAddrA="0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a"
ethValidatorAddrKeyA="8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230"
chain33ReceiverAddr="12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv"
chain33ReceiverAddrKey="4257d8692ef7fe13c68b65d6a52f03933db2fa5ce8faf210b5b8b80c721ced01"
ethValidatorAddrB="0x0df9a824699bc5878232c9e612fe1a5346a5a368"
#ethValidatorAddrKeyB="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"

maturityDegree=10
Chain33Cli="../../chain33-cli"

CLIA="./ebcli_A"
chain33ID=0
#BridgeRegistryOnChain33=""
chain33BridgeBank=""
#BridgeRegistryOnEth=""
ethBridgeBank=""
chain33BtyTokenAddr="1111111111111111111114oLvT2"
chain33EthTokenAddr=""
ethereumBtyTokenAddr=""
chain33YccTokenAddr=""
ethereumYccTokenAddr=""

# chain33 lock BTY, eth burn BTY
function TestChain33ToEthAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm)
    #    balance=$(cli_ret "${result}" "balance" ".balance")

    # chain33 lock bty
    hash=$(${Chain33Cli} send evm call -f 1 -a 5 -k "${chain33DeployAddr}" -e "${chain33BridgeBank}" -p "lock(${ethDeployAddr}, ${chain33BtyTokenAddr}, 500000000)" --khainID "${chain33ID}")
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a "${chain33DeployAddr}" -e evm)
    #    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}-5" | bc)"
    #balance_ret "${result}" "195.0000"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "5.0000"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} ethereum burn -m 3 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumBtyTokenAddr}") #--node_addr https://ropsten.infura.io/v3/9e83f296716142ffbaeaafc05790f26c)
    cli_ret "${result}" "burn"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethDeployAddr}" -t "${ethereumBtyTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    sleep ${maturityDegree}

    result=$(${Chain33Cli} account balance -a "${chain33ReceiverAddr}" -e evm)
    balance_ret "${result}" "3.0000"

    result=$(${Chain33Cli} account balance -a "${chain33BridgeBank}" -e evm)
    balance_ret "${result}" "2.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33Assets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${CLIA} ethereum lock -m 11 -k "${ethValidatorAddrKeyA}" -r "${chain33ReceiverAddr}")
    cli_ret "${result}" "lock"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "11"

    sleep ${maturityDegree}

    result=$(${Chain33Cli} evm query -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    is_equal "${result}" "1100000000"


    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrB}")
    cli_ret "${result}" "balance" ".balance" "100"

    echo '#5.burn ETH from Chain33 ETH(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethValidatorAddrB}" -t "${chain33EthTokenAddr}"

    sleep ${maturityDegree}

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm query -a "${chain33EthTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
 
    is_equal "${result}" "600000000"

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}")
    cli_ret "${result}" "balance" ".balance" "6"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrB}")
    cli_ret "${result}" "balance" ".balance" "105"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33Ycc() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${CLIA} ethereum lock -m 7 -k "${ethDeployKey}" -r "${chain33ReceiverAddr}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "lock"

    eth_block_wait 2

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "7"

    sleep ${maturityDegree}

    result=$(${Chain33Cli} evm query -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")
    
    is_equal "${result}" "700000000"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrB}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo '#5.burn YCC from Chain33 YCC(Chain33)-----> Ethereum'
    ${CLIA} chain33 burn -m 5 -k "${chain33ReceiverAddrKey}" -r "${ethValidatorAddrB}" -t "${chain33YccTokenAddr}"

    sleep ${maturityDegree}

    echo "check the balance on chain33"
    result=$(${Chain33Cli} evm query -a "${chain33YccTokenAddr}" -c "${chain33DeployAddr}" -b "balanceOf(${chain33ReceiverAddr})")

    is_equal "${result}" "200000000"

    result=$(${CLIA} ethereum balance -o "${ethBridgeBank}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "2"

    result=$(${CLIA} ethereum balance -o "${ethValidatorAddrB}" -t "${ethereumYccTokenAddr}")
    cli_ret "${result}" "balance" ".balance" "5"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# shellcheck disable=SC2120
function mainTest() {
    if [[ $# -ge 1 && ${1} != "" ]]; then
        chain33ID="${1}"
    fi
    StartChain33
    start_trufflesuite
    StartOneRelayer

    TestChain33ToEthAssets
    TestETH2Chain33Assets
    TestETH2Chain33Ycc
}

mainTest "${1}"
