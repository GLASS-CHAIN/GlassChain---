# paracross Participate in multi-node consensus and send consensus messages to the main chain

## Parachain transaction
 1. Filter the parachain transactions in the main chain that match the parachain title
 1. If a cross-chain contract is involved, if there are more than two parachain transactions that are judged to have failed, the transaction group will be unsuccessful in execution. (In this case, the main chain transaction will definitely be executed unsuccessfully)
 1. If no cross-chain contract is involved, then there is no regulation for the transaction group, which can be 20 ratios and 10 chains. If the main chain transaction fails, the parachain will not be executed
 1. If there is an ExecOk in the transaction group, all transactions on the main chain are ok and all can be packaged
 1. If all are ExecPack, there are two situations. One is that all transactions in the transaction group are parachain transactions, and the other is that the main chain has failed transactions and packaged transactions. LogErr needs to be checked. If there is an error, all transactions are not packaged.
 1. The transaction execution result bitmap sent by the parachain to the main chain only contains the packaged parachain tx. If it is a cross-chain transaction and the parachain transaction that has not been packaged into the block due to the failure of the main chain execution, it is not included in the bitmap Inside.

## Initial startup
 1. The consensus tick (16s) will periodically obtain the current consensus height through grpc
    * If each node is created and started, it will return -1, then enter the sync link and initiate a consensus message
    * If this node is restarted or a brand new node, it will actively synchronize the data of other nodes. During the synchronization process, no consensus data will be obtained or consensus messages will be sent. After the synchronization is over, the current consensus height will be obtained.
      Blocks before the consensus height, do not send consensus, send consensus messages from the current consensus node, enter the sync state, and participate in the consensus

## New nodes are added or restarted (including empty blocks)
   1. The node restarts and starts checking the current consensus height. If the height sent for the first time is higher than the consensus height, it will all be sent from the next height of the consensus height.
   1. New node. After the new node is started, the main chain data will be synchronized. Blocks below the consensus height will not send consensus messages until they are greater than the consensus height.

## Forking, node rollback
 1. If the height of delete is currently being sent, cancel the current sending. If you don’t cancel it, the reason for failure may be sent all the time.
 1. If the rollback height is less than the finish height during the fork, the finish height needs to be reset to the minimum value, and the finish height will be reset after the main chain consensus message comes.
 1. Stop the consensus response at the time of the fork, and release it when the new height is increased after the fork, which can ensure consistency and reduce unnecessary transaction sending and waste handling fees.

## Ordinary execution
 1. If you receive the main chain block, check whether the current transaction is in the block and the execution is successful. If the execution fails or packs, it is not counted as being on the chain and needs to be retransmitted.

## sign
 1. Export the private key from the wallet according to the configured address, and use the private key to sign on the parachain consensus. If the wallet is locked, the wallet side needs to set an error code to remind the user, and the parachain side will continue to send queries every 2s.
    Until the wallet is unlocked, the query is successful and the error code is cleared.

## Failure scenario
 1. The grpc link fails, and it will retransmit after 1s timeout. If the retransmission is not found in the two blocks of tx during the period, the retransmitted tx will also be updated to the new tx. In order to prevent mempool from being treated as a repeated transaction, tx nonce will change
 1. The transaction fee is not enough and the transaction fails
 1. The parachain has sent commit msg, the main chain is rolled back, the main chain cannot find the block corresponding to commit msg, and the parachain repeats sending until the parachain rolls back and cancels sending
 1. If the main chain of the parachain is forked, the main chain will fail to execute transactions sent from other parachains, and your own will succeed, and the main chain will be restored after the fork is rolled back.
 1. The main chain is normal, and the parachain has not reached a consensus since its inception, and debugging is required
 1. The main chain is normal, and a parachain has problems with its own calculations and cannot reach consensus with others. The transaction submitted by this parachain will fail, but the transaction will still be filtered to generate blocks without affecting the consensus. If the success is less than 2/3 of the nodes, consensus
    Will stop, each parachain will still generate blocks on its own, and the parachain itself needs to be debugged
 1. All parachains collapse at a certain height, and the consensus height lags behind the block height. After the node restarts, the consensus may have holes and need to be avoided. That is, the consensus height is the starting point, and the consensus height less than the consensus height does not need to be sent.
    Consensus greater than the consensus height, less than the height being sent, need to be obtained from the database and sent again
 1. For some reason, such as more than 2/3 node crashes or inconsistent data, the system does not generate a consensus at a certain height, and the consensus system will record the received transactions, even if the record has reached the consensus but because of the consensus height
    It is not continuous, or because the consensus is hollow, the following consensus is only a record and will not trigger done. Only the consensus commit that is continuous with the database consensus can trigger done, so once a void is generated, it needs
    Continuously send subsequent transactions from the beginning of the consensus, instead of just sending empty consensus data
 1. The main chain is in the cloud scenario, and the parachains are connected to a forked main chain node, and the parachains can reach consensus. The main chain does not have a forked node. After the main chain forked nodes later fall back and synchronize the main branch, they will be parallel. The chain nodes need to be resynchronized,
    Especially when the parachain initially sent a 20tx transaction group, the result was forked at the 10th height, the main chain consensus height is -1, and the parachain consensus is normal. When the forked master node recovers from the 10th node, the parachain node need
    Re-publish the consensus message from 0, because the current consensus height is -1

## Send failure strategy
 1. The current strategy is to either send a single transaction or a transaction group to send consensus messages, or all succeed or all fail. If it fails, that is, the transaction cannot be found in the new block, and the current block will be re-sent if it exceeds 2 blocks.
    For transactions in sending, new consensus messages will always be waiting. If the current sending tx has not entered the main block, the subsequent high-level consensus messages will never be sent. The scenario where the message fails, except for the link
    In addition to the failure, it is basically caused by the fork. The current strategy currently fails to see no problem.
 1. Another possible strategy is to send new transactions together with the current ones. It is better to have one transaction for each height, and not the transaction group. Check the transaction entry status separately. If the transaction that does not enter the chain is retransmitted, This strategy scenario
    It’s a bit complicated, and if the consensus transaction with the higher consensus succeeds and the lower fails, it is of little significance, so the first sending strategy is currently adopted.
 
## testing scenarios
 1. The main node and the parachain node are started in a docker, and the parachain node starts 120s later than the main node, which is basically when the main node is 8 heights.
 1. 6 nodes, 4 parachain nodes, the interval between two empty blocks is 4, and the other two are 3, no consensus can be reached
 1. 6 nodes, 4 parachain nodes, three empty blocks with an interval of 4, and one with 3, a consensus can be reached
 1. 6 nodes, 4 parachain nodes, 2 start first and cannot reach consensus, and the other or two start later to complete consensus
 1. 6 nodes, 4 parachain nodes, 2 start first, cannot reach consensus, another or two start later, can complete the consensus, then stop the first two, cannot reach consensus, and then start one or two of the first To complete the consensus
 1. 6 nodes, 4 parachain nodes, three are started first, and the fourth starts after 10 minutes. After starting, it will synchronize other node data and start sending from the current consensus node
 1. 6 nodes, 4 parachain nodes, all restarted at a certain height, this height was executed successfully but no consensus was sent, after restarting, check whether the unsent consensus can be resent
 1. 6 nodes, 4 parachain nodes, three or three groups, of which there are three parachain nodes in group a, and only one in group b. For the bifurcation test, stop group b first, then stop group a and start group b, and then start group a Group mining together, when group b is mining alone,
    Parachain cannot reach consensus and stays at the current height. After group a is started, the fork node of group b rolls back to reach a consensus again, and the consensus of group b parachain succeeds.

## Hard fork scenario
 1. The new version adds mining tx. If all nodes have no consensus, that is, the consensus is -1, you can delete all parachain node databases and upgrade the code to execute again without affecting the consensus
 1. If the node has a consensus and the consensus height is not N, commit msg needs to be set to> N to join the mining transaction, and all parachain databases need to be deleted, and the main chain database does not move
 1. If the node has a consensus and the height is N, the parachain does not delete the existing database. To update the version, the main chain needs to set a height that has not been reached as the consensus bifurcation point. The parachain side does not need to be set, and the previous consensus height is not sent .
 
## Test fork and fallback scenarios
 1. Enable the wallet of the CLI4 node in docker-compose.sh and transfer
    1. miner CLI4
    1. transfer CLI4
 1. Open the server chain30:8802 in nginx.conf, you can test the pause chain30 scenario
 1. If you test the parachain self-consensus, you need to set MainParaSelfConsensusForkHeight in the paracross testcase and the local setting of fork.go
    ForkParacrossCommitTx of pracross has the same height, or MainParaSelfConsensusForkHeight is greater than the fork height