### Offline deployment of ethereum cross-chain contracts and operations
***

#### Basic steps
* Create transaction online `./boss4x ethereum offline create ...` Need to check nonce and other information online
* Offline signature transaction `./boss4x ethereum offline sign -f xxx.txt -k 8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230`
* Send the signed file online `./boss4x ethereum offline send -f deploysigntxs.txt` The default file name after signing is deploysigntxs.txt
***

#### Offline deployment of ethereum cross-chain contract
* Create transactions online
```
Transaction 1: Deployment contract: Valset
Transaction 2: Deployment contract: EthereumBridge
Transaction 3: deployment contract: Oracle
Transaction 4: Deployment contract: BridgeBank
Transaction 5: Set the BridgeBank contract address in the EthereumBridge contract
Transaction 6: Set the Oracle contract address in the contract EthereumBridge
Transaction 7: deployment contract: BridgeRegistry
Transaction 7: deployment contract: MulSign

Order:
./boss4x ethereum offline create -p 25,25,25,25 -o 0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a -v 0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a, 0x0df9a824699bc5878232c9e612fe1a5346a5a368,0xcb074cb21cdddf3ce9c3c0a7ac4497d633c9d9f1,0xd9dab021e74ecf475788ed7b61356056b2095830

Parameter Description:
  -p, --initPowers string verifier weight, as: '25,25,25,25'
  -o, --owner string Deployer address
  -v, --validatorsAddrs string validator address, as:'addr,addr,addr,addr'

  --rpc_laddr_ethereum string ethereum url address (default "http://localhost:7545")

Output:
tx is written to file: deploytxs.txt

Write transaction information into the file
```

* Offline signature transaction
```
./boss4x ethereum offline sign -k 8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230

Parameter Description:
  -f, --file string The file to be signed, default: deploytxs.txt (default "deploytxs.txt")
  -k, --key string Deployer's private key
```

* Send the signed document
```
./boss4x ethereum offline send -f deploysigntxs.txt
```
***
#### Offline deployment of ERC20 cross-chain contract
* Create transactions online
```
Order:
./boss4x ethereum offline create_erc20 -m 33000000000000000000 -s YCC -o 0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a -d 0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a

Parameter Description:
  -m, --amount string amount
  -d, --deployAddr string Deployer address
  -o, --owner string owner address
  -s, --symbol string erc20 symbol

Output
tx is written to file: deployErc20YCC.txt
Write the transaction information into the deployErc20XXX.txt file, where XXX is the erc20 symbol
```

* Offline signature transaction
```
./boss4x ethereum offline sign -f deployErc20YCC.txt -k 8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230
```

* Send the signed document
```
./boss4x ethereum offline send -f deploysigntxs.txt
```

***
#### create_add_lock_list
* Create transactions online
```
Order:
./boss4x ethereum offline create_add_lock_list -s YCC -t 0x20a32A5680EBf55740B0C98B54cDE8e6FD5a4FB0 -c 0xC65B02a22B714b55D708518E2426a22ffB79113d -d 0x8afdad5a1087c9a1dd634

Parameter Description:
  -c, --contract string bridgebank contract address
  -d, --deployAddr string Deployer address
  -s, --symbol string token symbol
  -t, --token string token addr

Output
tx is written to file: create_add_lock_list.txt
```

***
#### Create bridge token
* Create transactions online
```
Order:
./boss4x ethereum offline create_bridge_token -s BTY -c 0xC65B02a22B714b55D708518E2426a22ffB79113d -d 0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a

Parameter Description:
  -c, --contract string bridgebank contract address
  -d, --deployAddr string Deployer address
  -s, --symbol string token symbol

Output
tx is written to file: create_bridge_token.txt
```

***
#### Set offline multi-signature address information
* Create transactions online
```
Order:
./boss4x ethereum offline multisign_setup -m 0xbf271b2B23DA4fA8Dc93Ce86D27dd09796a7Bf54 -d 0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a -o 0x4c85848a7E2985B76f06a7Ed338FCB3aF94a7DCf, 0x6F163E6daf0090D897AD7016484f10e0cE844994,0xbc333839E37bc7fAAD0137aBaE2275030555101f, 0x495953A743ef169EC5D4aC7b5F786BF2Bd56aFd5

Parameter Description:
  -d, --deployAddr string Deployer address
  -m, --multisign string Offline multisign contract address
  -o, --owner string Multi-signature address, separated by','

Output
tx is written to file: multisign_setup.txt
```

***
#### Set offline multi-signature address
* Create transactions online
```
Order:
./boss4x ethereum offline set_offline_addr -a 0xbf271b2B23DA4fA8Dc93Ce86D27dd09796a7Bf54 -c 0xC65B02a22B714b55D708518E2426a22ffB79113d -d 0x8afdadfc88a1087c9a1d6634c9a1d6

Parameter Description:
  -a, --address string Offline multi-signature address
  -c, --contract string bridgebank contract address
  -d, --deployAddr string deployment contract address

Output
tx is written to file: set_offline_addr.txt
```

***
#### Offline multi-signature settings
* Create transactions online
```
Order:
./boss4x ethereum offline set_offline_token -s ETH -m 20 -c 0xC65B02a22B714b55D708518E2426a22ffB79113d -d 0x8afdadfc88a1087c9a1d6c0f5dd04634b87f303a

Parameter Description:
  -c, --contract string bridgebank contract address
  -d, --deployAddr string deploy deployer address
  -p, --percents uint8 percentage (default 50), after reaching the threshold, the default transfer 50% to the offline multi-signature address
  -s, --symbol string token identification
  -m, --threshold float threshold
  -t, --token string token address

Output
tx is written to file: set_offline_token.txt
```

***
#### Offline multi-signature transfer
* Preparatory transaction for transfer-online operation
```
Order:
./boss4x ethereum offline multisign_transfer_prepare -a 3 -r 0xC65B02a22B714b55D708518E2426a22ffB79113d -c 0xbf271b2B23DA4fA8Dc93Ce86D27dd09796a7Bf54 -d 0x0df9a3466995bc5368a5346e612fe1878232c9

Parameter Description:
  -a, --amount float transfer amount
  -c, --contract string Offline multi-signature contract address
  -r, --receiver string recipient address
  -d, --sendAddr string The address to send this transaction, some handling fees need to be deducted
  -t, --token string erc20 address, if empty, transfer ETH by default

Output
tx is written to file: multisign_transfer_prepare.txt
```

* Offline multi-signature address signature transaction-offline operation
```
Order:
./boss4x ethereum offline sign_multisign_tx -k 0x5e8aadb91eaa0fce4df0bcc8bd1af9e703a1d6db78e7a4ebffd6cf045e053574,0x0504bcb22b21874b85b15f1bfae19ad62fc2ad89caefc5344dc669c57efa60db, 0x0c61f5a879d70807686e43eccc1f52987a15230ae0472902834af4d1933674f2,0x2809477ede1261da21270096776ba7dc68b89c9df5f029965eaa5fe7f0b80697

Parameter Description:
  -f, --file string tx file, default: multisign_transfer_prepare.txt (default "multisign_transfer_prepare.txt")
  -k, --keys string owners' private key, separated by','

Output
tx is written to file: sign_multisign_tx.txt
```

* Create a transfer transaction-online operation, you need to re-obtain nonce and other information
```
Order:
./boss4x ethereum offline create_multisign_tx

Output
tx is written to file: create_multisign_tx.txt
```

* Offline signature transaction
```
./boss4x ethereum offline sign -f create_multisign_tx.txt -k 8656d2bc732a8a816a461ba5e2d8aac7c7f85c26a813df30d5327210465eb230
```

* Send the signed document
```
./boss4x ethereum offline send -f deploysigntxs.txt
```