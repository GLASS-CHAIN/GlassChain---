#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -x
set +e

source "./publicTest.sh"

chain33SenderAddr="14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
# validatorsAddr=["0x92c8b16afd6d423652559c6e266cbe1c29bfd84f", "0x0df9a824699bc5878232c9e612fe1a5346a5a368", "0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1", "0xd9dab021e74ecf475788ed7b61356056b2095830"]
ethValidatorAddrKeyA="3fa21584ae2e4fd74db9b58e2386f5481607dfa4d7ba0617aaa7858e5025dc1e"
ethValidatorAddrKeyB="a5f3063552f4483cfc20ac4f40f45b798791379862219de9e915c64722c1d400"
ethValidatorAddrKeyC="bbf5e65539e9af0eb0cfac30bad475111054b09c11d668fc0731d54ea777471e"
ethValidatorAddrKeyD="c9fa31d7984edf81b8ef3b40c761f1847f6fcd5711ab2462da97dc458f1f896b"
#  chain33   10 bt 
chain33Validator1="1GTxrmuWiXavhcvsaH5w9whgVxUrWsUMdV"
chain33Validator2="155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6"
chain33Validator3="13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv"
chain33Validator4="113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG"
chain33ValidatorKey1="0xd627968e445f2a41c92173225791bae1ba42126ae96c32f28f97ff8f226e5c68"
chain33ValidatorKey2="0x9d539bc5fd084eb7fe86ad631dba9aa086dba38418725c38d9751459f567da66"
chain33ValidatorKey3="0x0a6671f101e30a2cc2d79d77436b62cdf2664ed33eb631a9c9e3f3dd348a23be"
chain33ValidatorKey4="0x3818b257b05ee75b6e43ee0e3cfc2d8502342cf67caed533e3756966690b62a5"
ethReceiverAddr1="0xa4ea64a583f6e51c3799335b28a8f0529570a635"
ethReceiverAddrKey1="355b876d7cbcb930d5dfab767f66336ce327e082cbaa1877210c1bae89b1df71"
ethReceiverAddr2="0x0c05ba5c230fdaa503b53702af1962e08d0c60bf"
ethReceiverAddrKey2="9dc6df3a8ab139a54d8a984f54958ae0661f880229bf3bdbb886b87d58b56a08"

maturityDegree=10
tokenAddrBty=""
tokenAddr=""
ethUrl=""
Chain33Cli=""

function kill_ebrelayerC() {
    #shellcheck disable=SC2154
    kill_docker_ebrelayer "${dockerNamePrefix}_ebrelayerc_1"
}
function kill_ebrelayerD() {
    kill_docker_ebrelayer "${dockerNamePrefix}_ebrelayerd_1"
}

function start_ebrelayerA() {
    docker cp "./relayer.toml" "${dockerNamePrefix}_ebrelayera_1":/root/relayer.toml
    start_docker_ebrelayer "${dockerNamePrefix}_ebrelayera_1" "/root/ebrelayer" "./ebrelayera.log"
    sleep 5
}

function start_ebrelayerC() {
    start_docker_ebrelayer "${dockerNamePrefix}_ebrelayerc_1" "/root/ebrelayer" "./ebrelayerc.log"
    sleep 5
    ${CLIC} relayer unlock -p 123456hzj
    sleep 5
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"
    sleep 10
}
function start_ebrelayerD() {
    start_docker_ebrelayer "${dockerNamePrefix}_ebrelayerd_1" "/root/ebrelayer" "./ebrelayerd.log"
    sleep 5
    ${CLID} relayer unlock -p 123456hzj
    sleep 5
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"
    sleep 10
}

function InitAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    result=$(${CLIA} relayer set_pwd -p 123456hzj)
    cli_ret "${result}" "set_pwd"

    result=$(${CLIA} relayer unlock -p 123456hzj)
    cli_ret "${result}" "unlock"

    result=$(${CLIA} relayer ethereum deploy)
    cli_ret "${result}" "deploy"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function StartRelayerAndDeploy() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    # change EthProvider url
    dockerAddr=$(get_docker_addr "${dockerNamePrefix}_ganachetest_1")
    ethUrl="http://${dockerAddr}:8545"

    #  relayer.toml 
    updata_relayer_a_toml "${dockerAddr}" "${dockerNamePrefix}_ebrelayera_1" "./relayer.toml"
    # start ebrelayer A
    start_ebrelayerA
    # 
    InitAndDeploy

    #  BridgeRegistry 
    result=$(${CLIA} relayer ethereum bridgeRegistry)
    BridgeRegistry=$(cli_ret "${result}" "bridgeRegistry" ".addr")

    # kill ebrelayer A
    kill_docker_ebrelayer "${dockerNamePrefix}_ebrelayera_1"
    sleep 1

    #  relayer.toml 
    updata_relayer_toml "${BridgeRegistry}" ${maturityDegree} "./relayer.toml"
    # 
    start_ebrelayerA

    # start ebrelayer B C D
    for name in b c d; do
        local file="./relayer$name.toml"
        cp './relayer.toml' "${file}"

        # 
        for deleteName in "deployerPrivateKey" "operatorAddr" "validatorsAddr" "initPowers" "deployerPrivateKey" "deploy"; do
            delete_line "${file}" "${deleteName}"
        done

        sed -i 's/x2ethereum/x2ethereum'${name}'/g' "${file}"

        pushHost=$(get_docker_addr "${dockerNamePrefix}_ebrelayer${name}_1")
        line=$(delete_line_show "${file}" "pushHost")
        sed -i ''"${line}"' a pushHost="http://'"${pushHost}"':20000"' "${file}"

        line=$(delete_line_show "${file}" "pushBind")
        sed -i ''"${line}"' a pushBind="'"${pushHost}"':20000"' "${file}"

        docker cp "${file}" "${dockerNamePrefix}_ebrelayer${name}_1":/root/relayer.toml
        start_docker_ebrelayer "${dockerNamePrefix}_ebrelayer${name}_1" "/root/ebrelayer" "./ebrelayer${name}.log"
    done
    sleep 5

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function EthImportKey() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    #  ebrelayer 
    for name in a b c d; do
        # 
        # shellcheck disable=SC2154
        CLI="docker exec ${dockerNamePrefix}_ebrelayer${name}_1 /root/ebcli_A"

        result=$(${CLI} relayer set_pwd -p 123456hzj)

        result=$(${CLI} relayer unlock -p 123456hzj)
        cli_ret "${result}" "unlock"
    done

    result=$(${CLIA} relayer ethereum import_chain33privatekey -k "${chain33ValidatorKey1}")
    cli_ret "${result}" "import_chain33privatekey"
    result=$(${CLIB} relayer ethereum import_chain33privatekey -k "${chain33ValidatorKey2}")
    cli_ret "${result}" "import_chain33privatekey"
    result=$(${CLIC} relayer ethereum import_chain33privatekey -k "${chain33ValidatorKey3}")
    cli_ret "${result}" "import_chain33privatekey"
    result=$(${CLID} relayer ethereum import_chain33privatekey -k "${chain33ValidatorKey4}")
    cli_ret "${result}" "import_chain33privatekey"

    result=$(${CLIA} relayer chain33 import_privatekey -k "${ethValidatorAddrKeyA}")
    cli_ret "${result}" "A relayer chain33 import_privatekey"
    result=$(${CLIB} relayer chain33 import_privatekey -k "${ethValidatorAddrKeyB}")
    cli_ret "${result}" "B relayer chain33 import_privatekey"
    result=$(${CLIC} relayer chain33 import_privatekey -k "${ethValidatorAddrKeyC}")
    cli_ret "${result}" "C relayer chain33 import_privatekey"
    result=$(${CLID} relayer chain33 import_privatekey -k "${ethValidatorAddrKeyD}")
    cli_ret "${result}" "D relayer chain33 import_privatekey"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# chian33 
function InitChain33Vilators() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    #  chain33Validators 
    result=$(${Chain33Cli} account import_key -k ${chain33ValidatorKey1} -l validator1)
    check_addr "${result}" ${chain33Validator1}
    result=$(${Chain33Cli} account import_key -k ${chain33ValidatorKey2} -l validator2)
    check_addr "${result}" ${chain33Validator2}
    result=$(${Chain33Cli} account import_key -k ${chain33ValidatorKey3} -l validator3)
    check_addr "${result}" ${chain33Validator3}
    result=$(${Chain33Cli} account import_key -k ${chain33ValidatorKey4} -l validator4)
    check_addr "${result}" ${chain33Validator4}

    # SetConsensusThreshold
    hash=$(${Chain33Cli} send x2ethereum setconsensus -p 80 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"

    # add a validator
    hash=$(${Chain33Cli} send x2ethereum add -a ${chain33Validator1} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"
    hash=$(${Chain33Cli} send x2ethereum add -a ${chain33Validator2} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"
    hash=$(${Chain33Cli} send x2ethereum add -a ${chain33Validator3} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"
    hash=$(${Chain33Cli} send x2ethereum add -a ${chain33Validator4} -p 25 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"

    # query Validators
    totalPower=$(${Chain33Cli} send x2ethereum query totalpower -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv | jq .totalPower | sed 's/\"//g')
    check_number 100 "${totalPower}"

    # cions  x2ethereum 
    hash=$(${Chain33Cli} send coins send_exec -e x2ethereum -a 200 -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)

    check_tx "${Chain33Cli}" "${hash}"
    result=$(${Chain33Cli} account balance -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -e x2ethereum)
    balance_ret "${result}" "200.0000"

    # chain33Validator 
    hash=$(${Chain33Cli} send coins transfer -a 10 -t "${chain33Validator1}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"
    result=$(${Chain33Cli} account balance -a "${chain33Validator1}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${Chain33Cli} send coins transfer -a 10 -t "${chain33Validator2}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"
    result=$(${Chain33Cli} account balance -a "${chain33Validator2}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${Chain33Cli} send coins transfer -a 10 -t "${chain33Validator3}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"
    result=$(${Chain33Cli} account balance -a "${chain33Validator3}" -e coins)
    balance_ret "${result}" "10.0000"

    hash=$(${Chain33Cli} send coins transfer -a 10 -t "${chain33Validator4}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    check_tx "${Chain33Cli}" "${hash}"
    result=$(${Chain33Cli} account balance -a "${chain33Validator4}" -e coins)
    balance_ret "${result}" "10.0000"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestChain33ToEthAssets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    # token4chain33    bty
    result=$(${CLIA} relayer ethereum token4chain33 -s coins.bty)
    tokenAddrBty=$(cli_ret "${result}" "token4chain33" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # chain33 lock bty
    hash=$(${Chain33Cli} send x2ethereum lock -a 5 -t coins.bty -r ${ethReceiverAddr1} -q "${tokenAddrBty}" --node_addr "${ethUrl}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${Chain33Cli}" $((maturityDegree + 2))
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} account balance -a 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -e x2ethereum)
    balance_ret "${result}" "195.0000"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} relayer ethereum burn -m 5 -k "${ethReceiverAddrKey1}" -r "${chain33SenderAddr}" -t "${tokenAddrBty}")
    cli_ret "${result}" "burn"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth  10 
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${Chain33Cli} account balance -a "${chain33SenderAddr}" -e x2ethereum)
    balance_ret "${result}" "5"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to chain33
#   chain33   eth 
function TestETH2Chain33Assets() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth lock 0.1
    result=$(${CLIA} relayer ethereum lock -m 0.1 -k "${ethReceiverAddrKey1}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0.1"

    # eth  10 
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${Chain33Cli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0.1"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    balance=$(cli_ret "${result}" "balance" ".balance")

    hash=$(${Chain33Cli} send x2ethereum burn -a 0.1 -t eth -r ${ethReceiverAddr2} --node_addr "${ethUrl}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${Chain33Cli}" $((maturityDegree + 2))
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}+0.1" | bc)"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33Erc20() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    # token4erc20  chain33  token  mint
    tokenSymbol="testc"
    result=$(${CLIA} relayer ethereum token4erc20 -s "${tokenSymbol}")
    tokenAddr=$(cli_ret "${result}" "token4erc20" ".addr")

    #  1000
    result=$(${CLIA} relayer ethereum mint -m 1000 -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "mint"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "1000"

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    # lock 100
    result=$(${CLIA} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${chain33Validator1}" -t "${tokenAddr}")
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    # eth  10 
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${Chain33Cli} x2ethereum balance -s "${chain33Validator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "100"

    # chain33 burn 100
    hash=$(${Chain33Cli} send x2ethereum burn -a 100 -t "${tokenSymbol}" -r ${ethReceiverAddr2} -q "${tokenAddr}" --node_addr "${ethUrl}" -k "${chain33Validator1}")
    block_wait "${Chain33Cli}" $((maturityDegree + 2))
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} x2ethereum balance -s "${chain33Validator1}" -t "${tokenSymbol}" -a "${tokenAddr}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestChain33ToEthAssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    if [ "${tokenAddrBty}" == "" ]; then
        # token4chain33    bty
        result=$(${CLIA} relayer ethereum token4chain33 -s coins.bty)
        tokenAddrBty=$(cli_ret "${result}" "token4chain33" ".addr")
    fi

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # chain33 lock bty
    hash=$(${Chain33Cli} send x2ethereum lock -a 5 -t coins.bty -r ${ethReceiverAddr2} -q "${tokenAddrBty}" --node_addr "${ethUrl}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${Chain33Cli}" $((maturityDegree + 2))
    check_tx "${Chain33Cli}" "${hash}"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    start_ebrelayerC

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "5"

    # eth burn
    result=$(${CLIA} relayer ethereum burn -m 5 -k "${ethReceiverAddrKey2}" -r "${chain33Validator1}" -t "${tokenAddrBty}")
    cli_ret "${result}" "burn"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddrBty}")
    cli_ret "${result}" "balance" ".balance" "0"

    # eth  10 
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${Chain33Cli} account balance -a "${chain33Validator1}" -e x2ethereum)
    balance_ret "${result}" "0"

    start_ebrelayerD

    result=$(${Chain33Cli} account balance -a "${chain33Validator1}" -e x2ethereum)
    balance_ret "${result}" "5"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

# eth to chain33
#   chain33   eth 
function TestETH2Chain33AssetsKill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # eth lock 0.1
    result=$(${CLIA} relayer ethereum lock -m 0.1 -k "${ethReceiverAddrKey1}" -r 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0.1"

    # eth  10 
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${Chain33Cli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${Chain33Cli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0.1"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    balance=$(cli_ret "${result}" "balance" ".balance")

    kill_ebrelayerC
    kill_ebrelayerD

    hash=$(${Chain33Cli} send x2ethereum burn -a 0.1 -t eth -r ${ethReceiverAddr2} --node_addr "${ethUrl}" -k 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv)
    block_wait "${Chain33Cli}" $((maturityDegree + 2))
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} x2ethereum balance -s 12qyocayNF7Lv6C9qW4avxs2E7U41fKSfv -t eth | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}")
    cli_ret "${result}" "balance" ".balance" "0.1"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}")
    cli_ret "${result}" "balance" ".balance" "$(echo "${balance}+0.1" | bc)"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function TestETH2Chain33Erc20Kill() {
    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"
    ${CLIA} relayer unlock -p 123456hzj

    # token4erc20  chain33  token  mint
    tokenSymbol="testcc"
    result=$(${CLIA} relayer ethereum token4erc20 -s "${tokenSymbol}")
    tokenAddr2=$(cli_ret "${result}" "token4erc20" ".addr")

    #  1000
    result=$(${CLIA} relayer ethereum mint -m 1000 -o "${ethReceiverAddr1}" -t "${tokenAddr2}")
    cli_ret "${result}" "mint"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "1000"

    result=$(${CLIA} relayer ethereum bridgeBankAddr)
    bridgeBankAddr=$(cli_ret "${result}" "bridgeBankAddr" ".addr")

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "0"

    kill_ebrelayerC
    kill_ebrelayerD

    # lock 100
    result=$(${CLIA} relayer ethereum lock -m 100 -k "${ethReceiverAddrKey1}" -r "${chain33Validator1}" -t "${tokenAddr2}")
    cli_ret "${result}" "lock"

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr1}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "900"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "100"

    # eth  10 
    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    result=$(${Chain33Cli} x2ethereum balance -s "${chain33Validator1}" -t "${tokenSymbol}" -a "${tokenAddr2}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${Chain33Cli} x2ethereum balance -s "${chain33Validator1}" -t "${tokenSymbol}" -a "${tokenAddr2}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "100"

    kill_ebrelayerC
    kill_ebrelayerD

    # chain33 burn 100
    hash=$(${Chain33Cli} send x2ethereum burn -a 100 -t "${tokenSymbol}" -r ${ethReceiverAddr2} -q "${tokenAddr2}" --node_addr "${ethUrl}" -k "${chain33Validator1}")
    block_wait "${Chain33Cli}" $((maturityDegree + 2))
    check_tx "${Chain33Cli}" "${hash}"

    result=$(${Chain33Cli} x2ethereum balance -s "${chain33Validator1}" -t "${tokenSymbol}" -a "${tokenAddr2}" | jq ".res" | jq ".[]")
    balance_ret "${result}" "0"

    eth_block_wait $((maturityDegree + 2)) "${ethUrl}"

    start_ebrelayerC
    start_ebrelayerD

    result=$(${CLIA} relayer ethereum balance -o "${ethReceiverAddr2}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "100"

    result=$(${CLIA} relayer ethereum balance -o "${bridgeBankAddr}" -t "${tokenAddr2}")
    cli_ret "${result}" "balance" ".balance" "0"

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}

function AllRelayerMainTest() {
    set +e
    docker_chain33_ip=$(get_docker_addr "${dockerNamePrefix}_chain33_1")
    Chain33Cli="./chain33-cli --rpc_laddr http://${docker_chain33_ip}:8801"

    CLIA="docker exec ${dockerNamePrefix}_ebrelayera_1 /root/ebcli_A"
    CLIB="docker exec ${dockerNamePrefix}_ebrelayerb_1 /root/ebcli_A"
    CLIC="docker exec ${dockerNamePrefix}_ebrelayerc_1 /root/ebcli_A"
    CLID="docker exec ${dockerNamePrefix}_ebrelayerd_1 /root/ebcli_A"
    echo "${CLIA}"

    echo -e "${GRE}=========== $FUNCNAME begin ===========${NOC}"

    if [[ ${1} != "" ]]; then
        maturityDegree=${1}
        echo -e "${GRE}maturityDegree is ${maturityDegree} ${NOC}"
    fi

    # init
    StartRelayerAndDeploy
    InitChain33Vilators
    EthImportKey

    # test
    TestChain33ToEthAssets
    TestETH2Chain33Assets
    TestETH2Chain33Erc20

    # kill relayer and start relayer
    TestChain33ToEthAssetsKill
    TestETH2Chain33AssetsKill
    TestETH2Chain33Erc20Kill

    echo -e "${GRE}=========== $FUNCNAME end ===========${NOC}"
}
