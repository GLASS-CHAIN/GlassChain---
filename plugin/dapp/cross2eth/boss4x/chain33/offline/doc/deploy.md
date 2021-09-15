### Offline deployment chain33 cross-chain contract and operations
***

#### Basic steps
* Create transaction offline and sign `./boss4x chain33 offline create ...`
* Send the signed file online `./boss4x chain33 offline send -f XXX.txt`
***

#### Offline deployment chain33 cross-chain contract
* Create transactions offline
```
Transaction 1: Deployment contract: Valset
Transaction 2: Deployment contract: chain33Bridge
Transaction 3: deployment contract: Oracle
Transaction 4: Deployment contract: BridgeBank
Transaction 5: Set the BridgeBank contract address in the contract chain33Bridge
Transaction 6: Set the Oracle contract address in the contract chain33Bridge
Transaction 7: deployment contract: BridgeRegistry
Transaction 7: deployment contract: MulSign

Order:
./boss4x chain33 offline create -f 1 -k 0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6709c262ae -n 'deploy crossx to chain33' -r '1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ, [1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ, 155ooMPBTF8QQsGAknkK7ei5D78rwDEFe6, 13zBdQwuyDh7cKN79oT2odkxYuDbgQiXFv, 113ZzVamKfAtGt9dq45fX1mNsEoDiN95HG], [25, 25, 25, 25]' --chainID 33

Parameter Description:
  -f, --fee float transaction fee setting, because there are only a few transactions, and the deployment transaction consumes more gas, just set 1 token directly
  -k, --key string The private key of the deployer, used to sign the transaction
  -n, --note string note information
  -r, --valset string Constructor parameters, input'addr, [addr, addr, addr, addr], [25, 25, 25, 25]' in strict accordance with the format, where the first address is the private key of the deployer Corresponding addresses, the last 4 addresses are the addresses of different validators, and the 4 numbers are the weights of different validators

  --rpc_laddr string chain33 url address (default "https://localhost:8801")
  --chainID int32 The chainID of the parachain, default: 0 (representing the main chain)

After execution, the transaction will be written to the file:
  deployCrossX2Chain33.txt
```

* Send the signed document
```
./boss4x chain33 offline send -f deployCrossX2Chain33.txt
```
***
#### Offline deployment of ERC20 cross-chain contract
* Create transactions offline
```
Order:
./boss4x chain33 offline create_erc20 -s YCC -k 0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6709c262ae -o 1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ --chain

Parameter Description:
  -a, --amount float Minting amount, default 3300*1e8
  -k, --key string The private key of the deployer, used to sign the transaction
  -o, --owner string owner address
  -s, --symbol string token identification

After execution, the transaction will be written to the file:
  deployErc20XXXChain33.txt where XXX is the token identifier
```

#### approve_erc20
* Create transactions offline
```
Order:
./boss4x chain33 offline approve_erc20 -a 330000000000 -s 1JmWVu1GEdQYSN1opxS9C39aS4NvG57yTr -c 1998HqVnt4JUirhC9KL5V71xYU8cFRn82c -k 0xbcae7c5d33e7e7c5d33e7e7c5d33709

Parameter Description:  
  -a, --amount float Approval amount
  -s, --approve string Approve address, chain33 BridgeBank contract address
  -c, --contract string Erc20 contract address
  -k, --key string The private key of the deployer, used to sign the transaction

After execution, the transaction will be written to the file:
  approve_erc20.txt
```

#### create_add_lock_list
* Create transactions offline
```
Order:
./boss4x chain33 offline create_add_lock_list -c 1JmWVu1GEdQYSN1opxS9C39aS4NvG57yTr -k 0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6ircVCC8 -cFR82YC82 -tFR5JVCC8 -t5

Parameter Description:
  -c, --contract string bridgebank contract address
  -k, --key string The private key of the deployer, used to sign the transaction
  -s, --symbol string token identification
  -t, --token string Erc20 contract address


After execution, the transaction will be written to the file:
  create_add_lock_list.txt
```

#### create_bridge_token
* Create transactions offline
```
Order:
./boss4x chain33 offline create_bridge_token -c 1JmWVu1GEdQYSN1opxS9C39aS4NvG57yTr -k 0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6709c262ae 33s YCC --s YCC --s

Parameter Description:
  -c, --contract string bridgebank contract address
  -k, --key string The private key of the deployer, used to sign the transaction
  -s, --symbol string token identification

After execution, the transaction will be written to the file:
  create_bridge_token.txt
```
* Get the bridge_token address
```
Order:
./chain33-cli evm abi call -a 1JmWVu1GEdQYSN1opxS9C39aS4NvG57yTr -c 1N6HstkyLFS8QCeVfdvYxx1xoryXoJtvvZ -b'getToken2address(YCC)'

Parameter Description:
  -a, --address string evm contract address, here is chain33 BridgeBank contract address
  -c, --caller string the caller address, here is the deployer address
  -b, --input string call params (abi format) like foobar(param1,param2), send function
  -t, --path string abi path(optional), default to .(current directory) (default "./"), abi file address, default local "./"

Output:
  15XsGjTbV6SxQtDE1SC5oaHx8HbseQ4Lf9 - bridge_token address
```



***
#### Set offline multi-signature address information
* Create transactions offline
```
Order:
./boss4x chain33 offline multisign_setup -m 1GrhufvPtnBCtfxDrFGcCoihmYMHJafuPn -o 168Sn1DXnLrZHTcAM9stD6t2P49fNuJfJ9,13KTf57aCkVVJYNJBXBBveiA5V811SrLcT, 1JQwQWsShTHC4zxHzbUfYQK4kRBriUQdEe, 1NHuKqoKe3hyv52PF8XBAyaTmJWAqA2Jbb -k 0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6709c262ae --chainID 33

Parameter Description:
  -k, --key string Deployer's private key
  -m, --multisign string Offline multisign contract address
  -o, --owner string Multi-signature address, separated by','

After execution, the transaction will be written to the file:
  multisign_setup.txt
```

***
#### Set offline multi-signature address
* Create transactions offline
```
Order:
./boss4x chain33 offline set_offline_addr -a 16skyHQA4YPPnhrDSSpZnexDzasS8BNx1R -c 1QD5pHMKZ9QWiNb9AsH3G1aG3Hashye83o -k 0x027ca96466c71c7e7c5d65b7e1f4333e1f43

Parameter Description:
  -a, --address string Offline multi-signature address
  -c, --contract string bridgebank contract address
  -f, --fee float transaction fee
  -k, --key string Deployer private key
  -n, --note string note

After execution, the transaction will be written to the file:
  chain33_set_offline_addr.txt
```

***
#### Offline multi-signature settings
* Create transactions offline
```
Order:
./boss4x chain33 offline set_offline_token -c 1MaP3rrwiLV1wrxPhDwAfHggtei1ByaKrP -s BTY -m 100000000000 -p 50 -k 0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31chain

Parameter Description:
  -c, --contract string bridgebank contract address
  -f, --fee float transaction fee
  -k, --key string Deployer private key
  -n, --note string note
  -p, --percents uint8 percentage (default 50), after reaching the threshold, the default transfer 50% to the offline multi-signature address
  -s, --symbol string token identification
  -m, --threshold string threshold
  -t, --token string token address


After execution, the transaction will be written to the file:
  chain33_set_offline_token.txt
```

***
#### Offline multi-signature transfer
* Create a transfer transaction-online operation, you need to re-obtain nonce and other information
```
Order:
./boss4x chain33 offline create_multisign_transfer -a 10 -r 168Sn1DXnLrZHTcAM9stD6t2P49fNuJfJ9 -m 1NFDfEwne4kjuxAZrtYEh4kfSrnGSE7ap

Parameter Description:  
  -m, --address string Offline multi-signature contract address
  -a, --amount float transfer amount
  -r, --receiver string recipient address
  -t, --token string erc20 address, if empty, the default transfer BTY


After execution, the transaction will be written to the file:
  create_multisign_transfer.txt
```

* Offline multi-signature address signature transaction-offline operation
```
Order:
./boss4x chain33 offline multisign_transfer -k 0x027ca96466c71c7e7c5d73b7e1f43cb889b3bd65ebd2413eefd31c6709c262a -s 0xcd284cd17456b73619fa609bb9e3105e8eff5d059c5e0b6eb1effbebd4d64144,0xe892212221b3b58211b90194365f4662764b6d5474ef2961ef77c909e31eeed3,0x9d19a2e9a440187010634f4f08ce36e2bc7b521581436a99f05568be94dc66ea, 0x45d4ce009e25e6d5e00d8d3a50565944b2e3604aa473680a656b242d9acbff35 --chainID 33

Parameter Description:
  -f, --fee float handling fee
  -t, --file string Signature transaction file, default: create_multisign_transfer.txt
  -k, --key string Deployer private key
  -s, --keys string Multiple private keys for offline multi-signature, separated by','
  -n, --note string note


After execution, the transaction will be written to the file:
  multisign_transfer.txt
```