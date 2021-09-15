# Parachain bls aggregate signature
>Multiple consensus nodes in the parachain form an internal local area network through P2P, and the consensus transactions sent by each node to the main chain are first broadcasted internally. The leader node is responsible for aggregating multiple consensus transactions into one consensus transaction and sending it to the main chain.


#1. Subscribe to P2P topic
1. Subscribe in P2P with PARA-BLS-SIGN-TOPIC as the topic, and broadcast synchronization messages through P2P between the internal nodes of the parachain. For example, here bls signature transactions and leader synchronization messages
   
#2. Negotiation leader
1. Taking into account that the leader rotates to send consensus transactions, every certain consensus height, such as 100, will rotate the next node to send transactions for the leader. After the current consensus height/100, the remaining base value of the nodegroup address is the current leader address
1. Considering that some leader nodes may be zombie nodes, each node monitors the heartbeat message of the leader node every 15s. If it is not received for more than 1 minute, it is considered that the leader node is not working and needs to skip to the next one. Here is maintenance
An offset, if it is not received over time, offset++, and the previous base are added together to determine the next leader node. If the current node finds itself as the leader node, it will start sending sync messages, and the offset will stop growing
1. As the consensus height grows, the base grows, and the offset remains the same, confirm the new leader together
1. In some special scenarios, when multiple leaders are synchronized, converge to the leader with the largest index value

#3. Send aggregate consensus transaction
1. The consensus transaction is P2P broadcast to all subscribed nodes, and the leader node is responsible for the chain after aggregation. If the collected signature transaction does not exceed 2/3 nodes, the chain transaction will not be sent, and the aggregation transaction will finally reach a consensus on the main chain
1. After the node broadcasts the consensus transaction, the consensus height does not increase for a certain period of time, and the consensus transaction is resent


#4. BLS aggregate signature algorithm
1. The private key required for BLS signature is twice as small as SECP256. The SECP256 private key is used to continuously fetch the hash until the BLS range is met as the BLS private key, and then the BLS public key is determined
1. The BLS public key is registered in the nodegroup of the main chain, and the BLS signature is verified together with the aggregate signature, while preventing the BLS leader node from cheating
1. For the same height, the consensus message signed by each node is the same, only one copy is needed, the signatures are aggregated into one, and the public key information is compressed into a bitmap and sent as a transaction
1. There are two paring curves for BLS signatures, G1 and G2. G1 generates a shorter msg, and G2 generates a longer. Generally, the public key is put on G2, and the signed message is put on G1. ETH uses the public key to put G1, the signature to G2, and the public key. The key is shorter, the message is longer,
Because the public key is statically configured in the database, the main chain is verified, and the signature is sent after the message, it takes up space. It is better than ETH. The static library can be compiled to support reversal, but it is still consistent with ETH 2.0.