#Deployment manual on chain33
##Step 1: Create 3 transactions to deploy router contracts offline
```
Transaction 1: Deployment contract: weth9
Transaction 2: deploy contract: factory
Transaction 3: deploy contract: router

./boss offline chain33 router -f 1 -k 0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944 -n "deploy router to chain33" -a 14KEKbYtKKQm4wMthSK9J4La4nAiidGozt --chainID 33

-f, --fee float: transaction fee setting, because there are only a few transactions, and the deployment transaction consumes more gas, just set 1 token directly
-k, --key string: The private key of the deployer, used to sign the transaction
-n, --note string: note information
-a, --feeToSetter: Set the transaction fee charging address (Note: This address is used to specify the address of the address that charges the transaction fee, not the address used to charge the transaction fee)
--chainID chainID of the parachain
Generate transaction files: farm.txt: router.txt

```

##Step 2: Create 5 offline farm contract transactions
```
Transaction 1: Deploy the contract: cakeToken
Transaction 2: Deploy the contract: SyrupBar
Transaction 3: Deployment contract: masterChef
Transaction 4: Transfer ownership, transfer ownership of cake token to masterchef
Transaction 5: Transfer ownership: transfer the ownership of SyrupBar to masterchef


./boss offline chain33 farm masterChef -f 1 -k 0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944 -n "deploy farm to chain33" -d 14KEKbYtKKQm4wMthSK9J4La4nAi -paraNameoz.
Generate transaction file: farm.txt
```

##Step 3: Create multiple transactions to increase lp token offline
```
./boss offline chain33 farm addPool -f 1 -k 0xcc38546e9e659d15e6b4893f0ab32a06d103931a8230b0bde71459d2b27d6944 -p 1000 -l 1HEp4BiA54iaKx5LrgN9iihkgmd3YxC2xM -m 13Ywpeqparagraph
```

##Step 4: Send the transaction in the transaction file serially
```
./boss offline chain33 send -f xxx.txt
```