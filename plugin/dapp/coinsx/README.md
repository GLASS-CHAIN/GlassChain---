# coinsx contract

## Foreword
1. In order to comply with the BSN consortium chain’s restriction rules on p2p native token transfers, in order not to affect the existing coins contract, the coinsx contract has been added.
2. Coinsx contract transfer, creation and other functions are exactly the same as coins, adding management functions for transfers


## use
1. The coinsx executor is different from the coins executor. It is configured through the toml configuration file coinExec="". The default is the coins executor
1. If the configuration is coinExec="coinsx", the native token is the coinsx contract, and the creation to the coinsx contract
1. The coins cli will modify the corresponding executor according to the configuration file, the transfer, withdraw and other commands are consistent with coins,
1. json-rpc construction transaction needs to explicitly adopt the configured coinExec, the default is coins
1. Parachain asset transfer
   1. The old interface assetTransfer/assetWithdraw still only accepts coins. If the main chain is coinsx, it will fail and the new interface needs to be used
   1. The new interface crossTransfer, through the required transaction executor parameters, ensures that the assets minted by the parachain are consistent with the main chain, rather than the configuration of the parachain.
        1. For example, the main chain coinsx, the default token of the parachain is coins, and the assets transferred from the main chain to the parachain are still coinsx.bty
        1. For example, the main chain coins, the default token of the parachain is coinsx, and the assets transferred from the parachain to the main chain are coinsx.symbol

## coinsx management function
1. Only the node super administrator can configure
1. Configure transfer enable and restricted flags
1. Configure the administrator group on the chain, add or delete administrators

## p2p transfer restriction rules
1. The system default transfer is limited
1. Node super administrators are not restricted to transfer from or to
1. If the transfer enable is configured, any user transfer is not restricted
1. If the transfer limited function is configured, transfers with the super administrator are not limited, and transfers between individuals are limited