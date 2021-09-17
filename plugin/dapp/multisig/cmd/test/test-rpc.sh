#!/usr/bin/env bash
# shellcheck disable=SC2128
set +e
set -o pipefail

MAIN_HTTP=""

# shellcheck source=/dev/null
source ../dapp-test-common.sh

Symbol="BTY"
Asset="coins"
AddrA="1C5xK2ytuoFqxmVGMcyz9XFKFWcDA8T3rK"
AddrB="1LDGrokrZjo1HtSmSnw8ef3oy5Vm1nctbj"
AddrE="1KHwX7ZadNeQDjBGpnweb4k2dqj2CWtAYo"

GenAddr="15wcitPEu1X1TBfrGfwN8GTkNTJoCmGc75"
PrivKeyGen="0x295710fa409bd0b0bf928efa0994645edfe80a247d89c1e1637f90dc5e303f5e"

multisigExecAddr=""
multisigAccAddr=""
execName=""

function init() {
    ispara=$(echo '"'"${MAIN_HTTP}"'"' | jq '.|contains("8901")')
    echo "ipara=$ispara"

    if [ "$ispara" == true ]; then
        execName="user.p.para.multisig"
        Symbol="para"
        multisigExecAddr=$(curl -ksd '{"method":"Chain33.ConvertExectoAddr","params":[{"execname":"user.p.para.multisig"}]}' ${MAIN_HTTP} | jq -r ".result")
    else
        execName="multisig"
        Symbol="BTY"
        multisigExecAddr=$(curl -ksd '{"method":"Chain33.ConvertExectoAddr","params":[{"execname":"multisig"}]}' ${MAIN_HTTP} | jq -r ".result")
    fi

    local main_ip=${MAIN_HTTP//8901/8801}

    if [ "$ispara" == false ]; then
        chain33_applyCoins "$GenAddr" 12000000000 "${main_ip}"
        chain33_QueryBalance "${GenAddr}" "$main_ip"

    else
        # tx fee
        chain33_applyCoins "$GenAddr" 1000000000 "${main_ip}"
        chain33_QueryBalance "${GenAddr}" "$main_ip"

        local para_ip="${MAIN_HTTP}"
        #para chain import pri key
        chain33_ImportPrivkey "0x295710fa409bd0b0bf928efa0994645edfe80a247d89c1e1637f90dc5e303f5e" "15wcitPEu1X1TBfrGfwN8GTkNTJoCmGc75" "gen" "$para_ip"

        chain33_applyCoins "$GenAddr" 12000000000 "${para_ip}"
        chain33_QueryBalance "${GenAddr}" "$para_ip"
    fi

    echo "multisigExecAddr=$multisigExecAddr"
}
function multisig_AccCreateTx() {
    echo "========== # multisig_AccCreateTx begin  =========="
    txHex=$(curl -ksd '{"method":"multisig.MultiSigAccCreateTx","params":[{"owners":[{"ownerAddr":"'$AddrA'","weight":20},{"ownerAddr":"'$AddrB'","weight":10},{"ownerAddr":"'$GenAddr'","weight":30}],"requiredWeight":15,"dailyLimit":{"symbol":"'$Symbol'","execer":"'$Asset'","dailyLimit":1000000000}}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    data=$(curl -ksd '{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccCount","payload":{}}]}' ${MAIN_HTTP} | jq -r ".result.data")
    echo "$data"

    multisigAccAddr=$(curl -ksd '{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccounts","payload":{"start":"0","end":"0"}}]}' ${MAIN_HTTP} | jq -r ".result.address[0]")

    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccountInfo","payload":{"multiSigAccAddr":"'"$multisigAccAddr"'"}}]}'
    resok='(.result.createAddr == "'$GenAddr'")'
    chain33_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"

    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccAllAddress","payload":{"multiSigAccAddr":"'$GenAddr'"}}]}'
    resok='(.result.address[0] == "'"$multisigAccAddr"'")'
    chain33_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"
    echo "========== # multisig_AccCreateTx ok  =========="
}

function multisig_TransferInTx() {
    echo "========== # multisig_TransferInTx begin =========="

    txHex=$(curl -ksd '{"method":"Chain33.CreateRawTransaction","params":[{"to":"'"$multisigExecAddr"'","amount":5000000000,"fee":1,"note":"12312","execName":"'"$execName"'"}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}


    txHex=$(curl -ksd '{"method":"multisig.MultiSigAccTransferInTx","params":[{"symbol":"'$Symbol'","amount":4000000000,"note":"test ","execname":"'$Asset'","to":"'"$multisigAccAddr"'"}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccAssets","payload":{"multiSigAddr":"'"$multisigAccAddr"'","assets":{"execer":"'$Asset'","symbol":"'$Symbol'"},"isAll":false}}]}'
    resok='(.result.accAssets[0].assets.execer == "'$Asset'") and (.result.accAssets[0].assets.symbol == "'$Symbol'") and (.result.accAssets[0].account.frozen == "4000000000")'
    chain33_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"
    echo "========== # multisig_TransferInTx end =========="
}

function multisig_TransferOutTx() {
    echo "========== # multisig_TransferOutTx begin =========="
    txHex=$(curl -ksd '{"method":"multisig.MultiSigAccTransferOutTx","params":[{"symbol":"'$Symbol'","amount":2000000000,"note":"test ","execname":"coins","to":"'$AddrB'","from":"'"$multisigAccAddr"'"}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccAssets","payload":{"multiSigAddr":"1LDGrokrZjo1HtSmSnw8ef3oy5Vm1nctbj","assets":{"execer":"coins","symbol":"'$Symbol'"},"isAll":false}}]}'
    resok='(.result.accAssets[0].assets.execer == "'$Asset'") and (.result.accAssets[0].assets.symbol == "'$Symbol'") and (.result.accAssets[0].account.balance == "2000000000")'
    chain33_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"

    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccAssets","payload":{"multiSigAddr":"'"$multisigAccAddr"'","assets":{"execer":"'$Asset'","symbol":"'$Symbol'"},"isAll":false}}]}'
    resok='(.result.accAssets[0].assets.execer == "'$Asset'") and (.result.accAssets[0].assets.symbol == "'$Symbol'") and (.result.accAssets[0].account.frozen == "2000000000")'
    chain33_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"
    echo "========== # multisig_TransferOutTx end =========="
}

function multisig_OwnerOperateTx() {
    echo "========== # multisig_OwnerOperateTx begin =========="

    txHex=$(curl -ksd '{"method":"multisig.MultiSigOwnerOperateTx","params":[{"multiSigAccAddr":"'"$multisigAccAddr"'","newOwner":"'$AddrE'","newWeight":8,"operateFlag":1}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}


    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccountInfo","payload":{"multiSigAccAddr":"'"$multisigAccAddr"'"}}]}'
    resok='(.result.owners[3].ownerAddr == "'$AddrE'")'
    chain33_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"

    txHex=$(curl -ksd '{"method":"multisig.MultiSigOwnerOperateTx","params":[{"multiSigAccAddr":"'"$multisigAccAddr"'","oldOwner":"'$AddrE'","operateFlag":2}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    txHex=$(curl -ksd '{"method":"multisig.MultiSigOwnerOperateTx","params":[{"multiSigAccAddr":"'"$multisigAccAddr"'","oldOwner":"'$AddrA'","newWeight":30,"operateFlag":3}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    txHex=$(curl -ksd '{"method":"multisig.MultiSigOwnerOperateTx","params":[{"multiSigAccAddr":"'"$multisigAccAddr"'","oldOwner":"'$AddrA'","newOwner":"'$AddrE'","operateFlag":4}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccountInfo","payload":{"multiSigAccAddr":"'"$multisigAccAddr"'"}}]}'
    resok='(.result.owners[0].ownerAddr == "'$AddrE'") and (.result.owners[0].weight == "30")'
    chain33_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"
    echo "========== # multisig_OwnerOperateTx end =========="
}

function multisig_AccOperateTx() {
    echo "========== # multisig_AccOperateTx begin =========="

    txHex=$(curl -ksd '{"method":"multisig.MultiSigAccOperateTx","params":[{"multiSigAccAddr":"'"$multisigAccAddr"'","dailyLimit":{"symbol":"'$Symbol'","execer":"'$Asset'","dailyLimit":1200000000}}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    txHex=$(curl -ksd '{"method":"multisig.MultiSigAccOperateTx","params":[{"multiSigAccAddr":"'"$multisigAccAddr"'","dailyLimit":{"symbol":"HYB","execer":"token","dailyLimit":1000000000}}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    txHex=$(curl -ksd '{"method":"multisig.MultiSigAccOperateTx","params":[{"multiSigAccAddr":"'"$multisigAccAddr"'","newRequiredWeight":16,"operateFlag":true}]}' ${MAIN_HTTP} | jq -r ".result")
    chain33_SignAndSendTx "$txHex" "$PrivKeyGen" ${MAIN_HTTP}

    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigAccTxCount","payload":{"multiSigAccAddr":"'"$multisigAccAddr"'"}}]}'
    chain33_Http "$req" ${MAIN_HTTP} '(.result.data != null)' "$FUNCNAME"

    req='{"method":"Chain33.Query","params":[{"execer":"multisig","funcName":"MultiSigTxInfo","payload":{"multiSigAddr":"'"$multisigAccAddr"'","txId":"7"}}]}'
    resok='(.result.txid == "7") and (.result.executed == true) and (.result.multiSigAddr == "'"$multisigAccAddr"'")'
    chain33_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"
    echo "========== # multisig_AccOperateTx end =========="
}

function run_test() {
    multisig_AccCreateTx
    multisig_TransferInTx
    multisig_TransferOutTx
    multisig_OwnerOperateTx
    multisig_AccOperateTx
}

function main() {
    chain33_RpcTestBegin multisi
    MAIN_HTTP="$1"
    echo "main_ip=$MAIN_HTTP"

    init
    run_test
    chain33_RpcTestRst multisi "$CASE_ERR"
}

chain33_debug_function main "$1"
