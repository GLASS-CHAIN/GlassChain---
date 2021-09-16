# Parachain selectively downloads the main chain block scheme containing parachain transactions

## illustrate
 1. The new parachain node needs two pieces of information to synchronize the main chain block from the beginning, one is the header information of all main chain blocks, and the other is the main chain block containing the parachain transactions.
 1. Obtain the headers of all main chain blocks. One is to generate empty blocks, and the other is to verify. The merkel root calculation of the transaction block of the parachain transaction needs to be consistent with the header tx root

    
## Implementation of segmented download scheme
1. Synchronously align the local block and the main chain block to determine the download start height, and then determine the download end height according to the current height of the main chain minus 1w height. Before the end height, it is considered that there is no rollback problem
1. Download a list of all main chain block heights containing parachain transactions, which will take 3s
1. Divide the height list into multiple segments according to the height of 1000, and obtain a block of the main chain containing parachain transactions each time. At the same time, determine the interval of the main chain head according to the height of the head and tail of this segment, and then get 1000 each time in the interval. Block header
　　 Then combined with the previously downloaded 1,000 parallel chain transaction blocks to generate a quasi-block, the synchronization layer can be executed immediately
1. For example, the height of 1~2w, the blocks containing parachain transactions are 100,105,...1500, 1610...2700,...15300,17200...19100, the first segment is 100~1500, the acquisition area The block header interval is 1~1500,
   The next block header interval is 1501~2700, and the last block header interval is 15301~20000
  
##Comparison of several schemes:
1. Request 1000 blocks each time, download them in turn, without distinguishing whether there is a parachain transaction block. In this scenario, the server needs to read the header or the block body of the parachain transaction to read the database, and the cache is not well used.
   It takes 33 minutes to actually test 500w main chain blocks in the current test environment
1. Download from multiple main chain service nodes. The download speed can be greatly improved when the main chain node ip is determined. However, due to the current cloud node proxy method, the test results have no obvious theoretical effect.
1. Download all the main chain heads first, and download the transaction block containing the parachain. Due to serialization, it takes more than 40 minutes
1. Download the parachain transaction first, and then download the relevant main chain head section in turn, and execute in parallel, which can be increased to 30 minutes. Since there may be hundreds of thousands of main chain head sections, synchronous execution will wait and you can continue to optimize
1. First download the parachain transaction in sections, and then download the relevant main chain head section in sections, each time 1000 blocks, can reduce the waiting time of the synchronization layer, the overall download process time is basically the time to download the main chain head
    1. It takes about 1 minute to download all parachain transactions, which is very small compared to all heads
    1. The test takes 21 minutes
    1. The last plan also used in this plan