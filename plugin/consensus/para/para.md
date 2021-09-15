#Support the master node in the cloud, and allow free switching, load balancing
 1. The main node of the parachain is set in the cloud as a cluster, providing a unified IP and port
 
## Scenes
 1. The master node cluster synchronizes the main network data, and each master node remains independent
 1. When a master node hangs, the cloud automatically switches to the new master node. The parachain node needs to check that the block obtained at the new node is consistent with the previous block hash, which means that the new block is add type, parentHash And before
    The hash of the main chain is the same, the block hash of del type is the same as the hash of the previous main chain
 1. If they are inconsistent, search for the seq of the master node blockHash of the parachain record on the new master node as the next seq to obtain tx
 1. If the mainBlockHash of the current parachain block cannot be found on the new node, it may be a fork scenario. You need to find the fork, delete the future parachain block, and synchronize the parachain data from the next seq at the fork.

## testing scenarios
 1. The parachain switches to the master node before blockHeight=1, and the parachain re-synchronizes data from seq=startSeq (startSeq=0 or non-zero scenario)
 1. The master node is switched, the new master node seq and blockhash are exactly the same as the old one
 1. The main node is switched, the old blockhash cannot be found on the new node, and the mainHash of the last node of the parachain can be found
 1. The main node is switched, the old blockhash can not be found on the new node, the mainHash of the last node of the parachain cannot be found, it is found back, and the node is resynchronized after deleting the fork
 1. The main node is switched, the old blockhash cannot be found on the new node, the mainHash of the last node of the parachain cannot be found, the fallback is found, the deletion of the forked node fails, and the synchronization is re-searched
 1. The master node is switched, the old blockhash cannot be found on the new node, and all the nodes of the parachain cannot find the scene, it will search in an infinite loop until the new master node is switched
 1. The system restarts and the master node switches scenes