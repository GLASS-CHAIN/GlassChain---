#!/usr/bin/env bash
# shellcheck disable=SC2178
set +e

PWD=$(cd "$(dirname "$0")" && pwd)
export PATH="$PWD:$PATH"

NODE3="${1}_chain33_1"
CLI="docker exec ${NODE3} /root/chain33-cli"

NODE2="${1}_chain32_1"
CLI2="docker exec ${NODE2} /root/chain33-cli"

NODE1="${1}_chain31_1"
CLI3="docker exec ${NODE1} /root/chain33-cli"

NODE4="${1}_chain30_1"
CLI4="docker exec ${NODE4} /root/chain33-cli"

NODE5="${1}_chain29_1"
CLI5="docker exec ${NODE5} /root/chain33-cli"

NODE6="${1}_chain28_1"
CLI6="docker exec ${NODE6} /root/chain33-cli"

containers=("${NODE1}" "${NODE2}" "${NODE3}" "${NODE4}" "${NODE5}" "${NODE6}")
forkContainers=("${CLI3}" "${CLI2}" "${CLI}" "${CLI4}" "${CLI5}" "${CLI6}")

export COMPOSE_PROJECT_NAME="$1"

sedfix=""
if [ "$(uname)" == "Darwin" ]; then
    sedfix=".bak"
fi

DAPP=""
if [ -n "${2}" ]; then
    DAPP=$2
fi

DAPP_TEST_FILE=""

if [ -n "${DAPP}" ]; then
    testfile="fork-test.sh"
    if [ -e "$testfile" ]; then
        # shellcheck source=/dev/null
        source "${testfile}"
        DAPP_TEST_FILE="$testfile"
    fi

    DAPP_COMPOSE_FILE="docker-compose-${DAPP}.yml"
    if [ -e "$DAPP_COMPOSE_FILE" ]; then
        export COMPOSE_FILE="docker-compose.yml:${DAPP_COMPOSE_FILE}"

    fi

fi

system_coins_file="system/coins/fork-test.sh"
# shellcheck source=/dev/null
source "${system_coins_file}"

echo "=========== # env setting ============="
echo "DAPP=$DAPP"
echo "DAPP_TEST_FILE=$DAPP_TEST_FILE"
echo "COMPOSE_FILE=$COMPOSE_FILE"
echo "COMPOSE_PROJECT_NAME=$COMPOSE_PROJECT_NAME"
echo "CLI=$CLI"

function base_init() {
    # update test environment
    sed -i $sedfix 's/^Title.*/Title="local"/g' chain33.toml
    sed -i $sedfix 's/^TestNet=.*/TestNet=true/g' chain33.toml

    # p2p
    sed -i $sedfix 's/^seeds=.*/seeds=["chain33:13802","chain32:13802","chain31:13802","chain30:13802","chain29:13802","chain28:13802"]/g' chain33.toml
    sed -i $sedfix 's/^enable=.*/enable=true/g' chain33.toml
    sed -i $sedfix 's/^isSeed=.*/isSeed=true/g' chain33.toml
    sed -i $sedfix 's/^innerSeedEnable=.*/innerSeedEnable=false/g' chain33.toml
    sed -i $sedfix 's/^useGithub=.*/useGithub=false/g' chain33.toml

    # rpc
    sed -i $sedfix 's/^jrpcBindAddr=.*/jrpcBindAddr="0.0.0.0:8801"/g' chain33.toml
    sed -i $sedfix 's/^grpcBindAddr=.*/grpcBindAddr="0.0.0.0:8802"/g' chain33.toml
    sed -i $sedfix 's/^whitelist=.*/whitelist=["localhost","127.0.0.1","0.0.0.0"]/g' chain33.toml

    # wallet
    sed -i $sedfix 's/^minerdisable=.*/minerdisable=false/g' chain33.toml

}

function start() {
    # docker-compose ps
    docker-compose ps

    # remove exsit container
    docker-compose down

    # create and run docker-compose container
    #    docker-compose -f docker-compose.yml -f docker-compose-para.yml up --build -d
    docker-compose up --build -d

    local SLEEP=30
    echo "=========== sleep ${SLEEP}s ============="
    sleep ${SLEEP}

    docker-compose ps

    # query node run status

    ${CLI} block last_header
    ${CLI} net info

    ${CLI} net peer
    peersCount=$(${CLI} net peer | jq '.[] | length')
    echo "${peersCount}"
    if [ "${peersCount}" -lt 2 ]; then
        sleep 20
        peersCount=$(${CLI} net peer | jq '.[] | length')
        echo "${peersCount}"
        if [ "${peersCount}" -lt 2 ]; then
            echo "peers error"
            exit 1
        fi
    fi

    #echo "=========== # create seed for wallet ============="
    #seed=$(${CLI} seed generate -l 0 | jq ".seed")
    #if [ -z "${seed}" ]; then
    #    echo "create seed error"
    #    exit 1
    #fi

    echo "=========== # save seed to wallet ============="
    result=$(${CLI} seed save -p 1314fuzamei -s "tortoise main civil member grace happy century convince father cage beach hip maid merry rib" | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "save seed to wallet error seed, result: ${result}"
        exit 1
    fi

    sleep 1

    echo "=========== # unlock wallet ============="
    result=$(${CLI} wallet unlock -p 1314fuzamei -t 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        exit 1
    fi

    sleep 1

    echo "=========== # import private key returnAddr ============="
    result=$(${CLI} account import_key -k CC38546E9E659D15E6B4893F0AB32A06D103931A8230B0BDE71459D2B27D6944 -l returnAddr | jq ".label")
    echo "${result}"
    if [ -z "${result}" ]; then
        exit 1
    fi

    sleep 1

    echo "=========== # import private key mining ============="
    result=$(${CLI} account import_key -k 4257D8692EF7FE13C68B65D6A52F03933DB2FA5CE8FAF210B5B8B80C721CED01 -l minerAddr | jq ".label")
    echo "${result}"
    if [ -z "${result}" ]; then
        exit 1
    fi

    sleep 1
    echo "=========== # close auto mining ============="
    result=$(${CLI} wallet auto_mine -f 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        exit 1
    fi

    ## 2nd mining
    echo "=========== # save seed to wallet ============="
    result=$(${CLI4} seed save -p 1314fuzamei -s "tortoise main civil member grace happy century convince father cage beach hip maid merry rib" | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "save seed to wallet error seed, result: ${result}"
        exit 1
    fi

    sleep 1

    echo "=========== # unlock wallet ============="
    result=$(${CLI4} wallet unlock -p 1314fuzamei -t 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        exit 1
    fi

    sleep 1

    echo "=========== # import private key returnAddr ============="
    result=$(${CLI4} account import_key -k 2AFF1981291355322C7A6308D46A9C9BA311AA21D94F36B43FC6A6021A1334CF -l returnAddr | jq ".label")
    echo "${result}"
    if [ -z "${result}" ]; then
        exit 1
    fi

    sleep 1

    echo "=========== # import private key mining ============="
    result=$(${CLI4} account import_key -k 2116459C0EC8ED01AA0EEAE35CAC5C96F94473F7816F114873291217303F6989 -l minerAddr | jq ".label")
    echo "${result}"
    if [ -z "${result}" ]; then
        exit 1
    fi

    sleep 1
    echo "=========== # close auto mining ============="
    result=$(${CLI4} wallet auto_mine -f 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        exit 1
    fi

    block_wait "${CLI}" 1

    echo "=========== check genesis hash ========== "
    ${CLI} block hash -t 0
    res=$(${CLI} block hash -t 0 | jq ".hash")
    count=$(echo "$res" | grep -c "0x67c58d6ba9175313f0468ae4e0ddec946549af7748037c2fdd5d54298afd20b6")
    if [ "${count}" != 1 ]; then
        echo "genesis hash error!"
        exit 1
    fi

    echo "=========== query height ========== "
    ${CLI} block last_header
    result=$(${CLI} block last_header | jq ".height")
    if [ "${result}" -lt 1 ]; then
        exit 1
    fi

    ${CLI} wallet status
    ${CLI} account list
    ${CLI} mempool list
}

function dapp_run() {
    if [ -e "$DAPP_TEST_FILE" ]; then
        ${DAPP} "${CLI}" "${1}"
    fi

}

function optDockerfun() {
    #############################################
    #1 The first bifurcation structure: first the
    # two chains conduct joint mining, and then split
    # Don't mine, that is, the transactions when
    # the two chains are forked are different.
    #############################################
    forkType1
    #############################################
     #2 The second fork structure: including the
     # first group of dockers, the second group
     # of dockers,And the public node docker, 
     # first mine together, and then stop the
     # second group docker, back up the public node
     # docker database, in the public node docker
     # Create transaction on, sign transaction,
     # record signature, send, then close the first group
     # docker, then restore the public node 
     # docker database to the backup state,
     # Then start the second group of docker, 
     # and then send the transaction that just recorded the signature.
     # Finally start all nodes to mine together
    #############################################
    forkType2
}

function forkType1() {
    echo "=========== Start type 1 fork test ========== "
    base_init
    dapp_run forkInit

    start
    optDockerPart1
    #############################################
    #Here to join according to specific needs; such as transferring from the wallet to a specific contract account
    #1 Initialize the transaction balance
    dapp_run forkConfig

    #############################################
    optDockerPart2
    #############################################
    #Here according to specific needs to join in a test chain to send test data
    #2 Construct the first chain transaction
    dapp_run forkAGroupRun

    #############################################
    optDockerPart3
    #############################################
    #Here, according to specific needs, join to send test data in the second test chain
    #3 Construct the second chain transaction
    dapp_run forkBGroupRun

    #############################################
    optDockerPart4
    loopCount=30 #The number of cycles, the sleep time of each cycle is 100s
    checkBlockHashfun $loopCount

    #############################################
    #Add the result check according to specific needs here
    #4 Check the transaction results
    dapp_run forkCheckRst

    #############################################
    echo "=========== Type 1 fork test ends ========== "
}

function forkType2() {
    echo "=========== Start of type 2 bifurcation testing ========== "
    base_init
    dapp_run fork2Init

    start

    optDockerPart1
    #############################################
    #Here to join according to specific needs; such as transferring from the wallet to a specific contract account
    #1 Initialize the transaction balance
    initCoinsAccount
    dapp_run fork2Config

    #############################################
    type2_optDockerPart2
    #############################################
    #Here according to specific needs to join in a test chain to send test data
    #2 Construct the first chain transaction
    genFirstChainCoinstx
    dapp_run fork2AGroupRun
    #############################################
    type2_optDockerPart3
    #############################################
    #Here, according to specific needs, join to send test data in the second test chain
    #3 Construct the second chain transaction
    genSecondChainCoinstx
    dapp_run fork2BGroupRun

    #############################################
    type2_optDockerPart4
    loopCount=30 #The number of cycles, the sleep time of each cycle is 100s
    checkBlockHashfun $loopCount
    #############################################
    #Add the result check according to specific needs here
    #4 Check the transaction results
    checkCoinsResult
    dapp_run fork2CheckRst
    #############################################
    echo "=========== Type 2 fork test ends ========== "
}

function optDockerPart1() {
    echo "====== Block generation ======"
    #sleep 100
    block_wait_timeout "${CLI}" 10 100

    loopCount=20
    for ((i = 0; i < loopCount; i++)); do
        name="${CLI}"
        time=2
        needCount=6
        peersCount "${name}" $time $needCount
        peerStatus=$?
        if [ $peerStatus -eq 1 ]; then
            name="${CLI4}"
            peersCount "${name}" $time $needCount
            peerStatus=$?
            if [ $peerStatus -eq 0 ]; then
                break
            fi
        else
            break
        fi
        #Check whether the maximum number of detections has been exceeded
        if [ $i -ge $((loopCount - 1)) ]; then
            echo "====== peers not enough ======"
            exit 1
        fi
    done

    return 1

}

function optDockerPart2() {
    checkMineHeight
    status=$?
    if [ $status -eq 0 ]; then
        echo "====== All peers is the same height ======"
    else
        echo "====== All peers is the different height, syn blockchain fail======"
        exit 1
    fi

    echo "==================================="
    echo "====== Step 1: The first group of docker mining======"
    echo "==================================="

    echo "======Stop the second group of docker ======"
    docker pause "${NODE4}" "${NODE5}" "${NODE6}"

    echo "======Start the first set of docker node mining======"
    sleep 3
    result=$($CLI wallet auto_mine -f 1 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "start wallet2 mine fail"
        exit 1
    fi

    name=$CLI
    time=60
    needCount=3

    peersCount "${name}" $time $needCount
    peerStatus=$?
    if [ $peerStatus -eq 1 ]; then
        echo "====== peers not enough ======"
        exit 1
    fi

}

function optDockerPart3() {
    echo "======The first group of docker nodes is mining======"
    block_wait_timeout "${CLI}" 5 100
    echo "======Stop the first set of docker node mining======"
    result=$($CLI wallet auto_mine -f 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "stop wallet2 mine fail"
        exit 1
    fi

    echo "====== The first group of internal synchronization ======"
    names[0]="${NODE3}"
    names[1]="${NODE2}"
    names[2]="${NODE1}"
    syn_block_timeout "${CLI}" 3 50 "${names[@]}"

    echo "======================================="
    echo "======== Step 2: The second group of docker mining ======="
    echo "======================================="

    echo "======Stop the first group of docker======"
    docker pause "${NODE1}" "${NODE2}" "${NODE3}"

    echo "======sleep 5s======"
    sleep 5

    echo "======Start the second group of docker======"
    docker unpause "${NODE4}" "${NODE5}" "${NODE6}"

    echo "======sleep 20s======"
    sleep 5
    result=$($CLI4 wallet unlock -p 1314fuzamei -t 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "wallet1 unlock fail"
        exit 1
    fi

    name="${CLI4}"
    time=60
    needCount=3

    peersCount "${name}" $time $needCount
    peerStatus=$?
    if [ $peerStatus -eq 1 ]; then
        echo "====== peers not enough ======"
        exit 1
    fi

    echo "======Start the second set of docker node mining======"
    sleep 1
    result=$($CLI4 wallet auto_mine -f 1 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "start wallet1 mine fail"
        exit 1
    fi

    names[0]="${NODE4}"
    names[1]="${NODE5}"
    names[2]="${NODE6}"
    syn_block_timeout "${CLI4}" 2 100 "${names[@]}"

}

function optDockerPart4() {
    echo "======The second group of docker nodes is mining======"
    block_wait_timeout "${CLI4}" 3 50
    echo "====== The second group of internal synchronization ======"
    names[0]="${NODE4}"
    names[1]="${NODE5}"
    names[2]="${NODE6}"
    syn_block_timeout "${CLI4}" 3 50 "${names[@]}"

    echo "======================================="
    echo "====== The third step: two groups of docker mining together ======="
    echo "======================================="

    echo "======Start the first group of docker======"
    docker unpause "${NODE1}" "${NODE2}" "${NODE3}"

    echo "======sleep 20s======"
    sleep 5
    result=$($CLI wallet unlock -p 1314fuzamei -t 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "wallet2 unlock fail"
        exit 1
    fi
    echo "======Start the first set of docker node mining======"
    sleep 1
    result=$($CLI wallet auto_mine -f 1 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "start wallet2 mine fail"
        exit 1
    fi

    echo "======Two groups of docker nodes are mining together======"
    block_wait_timeout "${CLI}" 5 100
}

function copyData() {
    name="${NODE3}"
    sleep 1
    docker exec "${name}" mkdir beifen
    sleep 1
    docker exec "${name}" cp -r datadir beifen
    sleep 1
}

function restoreData() {
    name="${NODE3}"
    sleep 1
    docker exec "${name}" rm -rf datadir
    sleep 1
    docker exec "${name}" cp -r beifen/datadir ./
    sleep 1
}

function type2_optDockerPart2() {
    checkMineHeight
    status=$?
    if [ $status -eq 0 ]; then
        echo "====== All peers is the same height ======"
    else
        echo "====== All peers is the different height, syn blockchain fail======"
        exit 1
    fi

    echo "=============== Backup public node data =============="
    copyData

    echo "==================================="
    echo "====== Step 1: The first group of docker mining======"
    echo "==================================="

    echo "======Stop the second group of docker======"
    docker pause "${NODE4}" "${NODE5}" "${NODE6}"

    echo "======Start the first set of docker node mining======"
    sleep 3
    result=$($CLI wallet auto_mine -f 1 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "start wallet2 mine fail"
        exit 1
    fi

    name=$CLI
    time=60
    needCount=3

    peersCount "${name}" $time $needCount
    peerStatus=$?
    if [ $peerStatus -eq 1 ]; then
        echo "====== peers not enough ======"
        exit 1
    fi

}

function type2_optDockerPart3() {
    echo "======The first group of docker nodes is mining======"
    block_wait_timeout "${CLI}" 5 100
    echo "======Stop the first set of docker node mining======"
    result=$($CLI wallet auto_mine -f 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "stop wallet2 mine fail"
        exit 1
    fi

    echo "====== The first group of internal synchronization ======"
    names[0]="${NODE3}"
    names[1]="${NODE2}"
    names[2]="${NODE1}"
    syn_block_timeout "${CLI}" 3 50 "${names[@]}"

    echo "======================================="
    echo "======== Step 2: The second group of docker mining ======="
    echo "======================================="

    echo "======Stop docker except public nodes in the first group======"
    docker pause "${NODE1}" "${NODE2}"

    echo "=============== Restore public node data=============="
    restoreData
    docker pause "${NODE3}"

    echo "======sleep 5s======"
    sleep 5

    echo "======Start the second group of docker======"
    docker unpause "${NODE3}" "${NODE4}" "${NODE5}" "${NODE6}"

    name="${CLI}"
    time=60
    needCount=4

    peersCount "${name}" $time $needCount
    peerStatus=$?
    if [ $peerStatus -eq 1 ]; then
        echo "====== peers not enough ======"
        exit 1
    fi

    echo "======sleep 20s======"
    sleep 20
    echo "======Start the second set of docker node mining======"

    result=$($CLI wallet unlock -p 1314fuzamei -t 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "wallet1 unlock fail"
        exit 1
    fi

    sleep 1
    result=$($CLI wallet auto_mine -f 1 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "start wallet1 mine fail"
        exit 1
    fi

    sleep 1
    result=$($CLI4 wallet unlock -p 1314fuzamei -t 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "wallet2 unlock fail"
        exit 1
    fi

    sleep 1
    result=$($CLI4 wallet auto_mine -f 1 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "start wallet2 mine fail"
        exit 1
    fi

    names[0]="${NODE3}"
    names[1]="${NODE4}"
    names[2]="${NODE5}"
    names[3]="${NODE6}"
    syn_block_timeout "${CLI}" 2 100 "${names[@]}"

}

function type2_optDockerPart4() {
    echo "======The second group of docker nodes is mining======"
    block_wait_timeout "${CLI}" 3 50
    echo "====== The second group of internal synchronization ======"
    names[0]="${NODE4}"
    names[1]="${NODE5}"
    names[2]="${NODE6}"
    names[3]="${NODE3}"
    syn_block_timeout "${CLI}" 3 50 "${names[@]}"

    echo "======================================="
    echo "====== The third step: two groups of docker mining together ======="
    echo "======================================="

    echo "======Start the first group of docker======"
    docker unpause "${NODE1}" "${NODE2}"

    echo "======Two groups of docker nodes are mining together======"
    block_wait_timeout "${CLI}" 5 100
}

function checkMineHeight() {
    result=$($CLI4 wallet auto_mine -f 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "stop wallet1 mine fail"
        return 1
    fi
    sleep 1
    result=$($CLI wallet auto_mine -f 0 | jq ".isok")
    if [ "${result}" = "false" ]; then
        echo "stop wallet2 mine fail"
        return 1
    fi

    echo "====== stop all wallet mine success ======"

    echo "====== syn blockchain ======"
    syn_block_timeout "${CLI}" 5 50 "${containers[@]}"

    height=0
    height1=$($CLI4 block last_header | jq ".height")
    sleep 1
    height2=$($CLI block last_header | jq ".height")
    if [ "${height2}" -ge "${height1}" ]; then
        height=$height2
        printf 'Current maximum height %s \n' "${height}"
    else
        height=$height1
        printf 'Current maximum height %s \n' "${height}"
    fi

    if [ "${height}" -eq 0 ]; then
        echo "Failed to get the current maximum height"
        return 1
    fi
    loopCount=20
    for ((k = 0; k < ${#forkContainers[*]}; k++)); do
        for ((j = 0; j < loopCount; j++)); do
            height1=$(${forkContainers[$k]} block last_header | jq ".height")
            if [ "${height1}" -gt "${height}" ]; then #If it is larger than it means that the block has not been completely generated, replace the expected height
                height=$height1
                printf 'Inquire %s The highest height of the current block is %s \n' "${containers[$k]}" "${height}"
            elif [ "${height1}" -eq "${height}" ]; then #Find the target height
                break
            else
                printf 'Inquire %s NS %d Times, current height %d, Need height%d, synchronizing, sleep 60s Post query\n' "${containers[$k]}" $j "${height1}" "${height}"
                sleep 60
            fi
            #Check whether the maximum number of detections has been exceeded
            if [ $j -ge $((loopCount - 1)) ]; then
                echo "====== syn blockchain fail======"
                return 1
            fi
        done
    done

    return 0
}

function peersCount() {
    name=$1
    time=$2
    needCount=$3

    for ((i = 0; i < time; i++)); do
        peersCount=$($name net peer | jq '.[] | length')
        printf 'Query node %s, the required number of nodes %d, the current number of nodes %s \n' "${name}" "${needCount}" "${peersCount}"
        if [ "${peersCount}" = "$needCount" ]; then
            echo "============= Meet the requirements of the number of nodes ============="
            return 0
        else
            echo "============= Sleep for 30s to continue query ============="
            sleep 30
        fi
    done

    return 1
}

function checkBlockHashfun() {
    echo "====== syn blockchain ======"
    syn_block_timeout "${CLI}" 10 50 "${containers[@]}"

    height=0
    hash=""
    height1=$($CLI block last_header | jq ".height")
    sleep 1
    height2=$($CLI4 block last_header | jq ".height")
    if [ "${height2}" -ge "${height1}" ]; then
        height=$height2
        printf "The main chain is $CLI Current maximum height %d \\n" "${height}"
        sleep 1
        hash=$($CLI block hash -t "${height}" | jq ".hash")
    else
        height=$height1
        printf "The main chain is $CLI4 Current maximum height %d \\n" "${height}"
        sleep 1
        hash=$($CLI4 block hash -t "${height}" | jq ".hash")
    fi

    for ((j = 0; j < $1; j++)); do
        for ((k = 0; k < ${#forkContainers[*]}; k++)); do
            sleep 1
            height0[$k]=$(${forkContainers[$k]} block last_header | jq ".height")
            if [ "${height0[$k]}" -ge "${height}" ]; then
                sleep 1
                hash0[$k]=$(${forkContainers[$k]} block hash -t "${height}" | jq ".hash")
            else
                hash0[$k]="${forkContainers[$k]}"
            fi
        done

        if [ "${hash0[0]}" = "${hash}" ] && [ "${hash0[1]}" = "${hash}" ] && [ "${hash0[2]}" = "${hash}" ] && [ "${hash0[3]}" = "${hash}" ] && [ "${hash0[4]}" = "${hash}" ] && [ "${hash0[5]}" = "${hash}" ]; then
            echo "syn blockchain success break"
            break
        else
            if [ "${hash0[1]}" = "${hash0[0]}" ] && [ "${hash0[2]}" = "${hash0[0]}" ] && [ "${hash0[3]}" = "${hash0[0]}" ] && [ "${hash0[4]}" = "${hash0[0]}" ] && [ "${hash0[5]}" = "${hash0[0]}" ]; then
                echo "syn blockchain success break"
                break
            fi
        fi
        peersCount=0
        peersCount=$(${forkContainers[0]} net peer | jq '.[] | length')
        printf 'NS %d times, The network synchronization has not been queried, the current number of nodes is %d, and it will be queried after 100s \n' $j "${peersCount}"
        sleep 100
        #Check whether the maximum number of detections has been exceeded
        var=$(($1 - 1))
        if [ $j -ge "${var}" ]; then
            echo "====== syn blockchain fail======"
            exit 1
        fi
    done
    echo "====== syn blockchain success======"
}

# $1 name
# $2 txHash
function txQuery() {
    name=$1
    txHash=$2
    result=$($name tx query -s "${txHash}" | jq -r ".receipt.tyname")
    if [ "${result}" = "ExecOk" ]; then
        return 0
    fi
    return 1
}

function block_wait_timeout() {
    if [ "$#" -lt 3 ]; then
        echo "wrong block_wait params"
        exit 1
    fi
    cur_height=$(${1} block last_header | jq ".height")
    expect=$((cur_height + ${2}))
    count=0
    while true; do
        new_height=$(${1} block last_header | jq ".height")
        if [ "${new_height}" -ge "${expect}" ]; then
            break
        fi
        count=$((count + 1))
        sleep 1
        if [ $count -ge "${3}" ]; then
            echo "====== block wait timeout ======"
            break
        fi
    done
    echo "wait new block $count s"
}

function syn_block_timeout() {
    #${1} name
    #${2} minHeight
    #${3} timeout
    #${4} names

    names=${4}

    if [ "$#" -lt 3 ]; then
        echo "wrong block_wait params"
        exit 1
    fi
    cur_height=$(${1} block last_header | jq ".height")
    expect=$((cur_height + ${2}))
    count=0
    while true; do
        new_height=$(${1} block last_header | jq ".height")
        if [ "${new_height}" -lt "${expect}" ]; then
            count=$((count + 1))
            sleep 1
        else
            isSyn="true"
            for ((k = 0; k < ${#names[@]}; k++)); do
                sync_status=$(docker exec "${names[$k]}" /root/chain33-cli net is_sync)
                if [ "${sync_status}" = "false" ]; then
                    isSyn="false"
                    break
                fi
                count=$((count + 1))
                sleep 1
            done
            if [ "${isSyn}" = "true" ]; then
                break
            fi
        fi

        if [ $count -ge $(($3 + 1)) ]; then
            echo "====== syn block wait timeout ======"
            break
        fi

    done
    echo "wait block $count s"
}

function block_wait() {
    if [ "$#" -lt 2 ]; then
        echo "wrong block_wait params"
        exit 1
    fi
    cur_height=$(${1} block last_header | jq ".height")
    expect=$((cur_height + ${2}))
    count=0
    while true; do
        new_height=$(${1} block last_header | jq ".height")
        if [ "${new_height}" -ge "${expect}" ]; then
            break
        fi
        count=$((count + 1))
        sleep 1
    done
    echo "wait new block $count s"
}

optDockerfun
