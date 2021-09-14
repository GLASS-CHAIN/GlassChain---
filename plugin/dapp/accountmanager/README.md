# accountmanager contract

## Foreword
In order to adapt to the [Financial Distributed Account Technology Security Specification] (http://www.cfstc.org/bzgk/gk/view/yulan.jsp?i_id=1855) issued by the Central Bank, to meet the financial supervision in the alliance chain, account Public key reset, black and white list and other requirements, specially developed on chain33
accountmanager contract

## use
The contract is designed according to the centralized financial service, and there are administrators to supervise the accounts under the contract
Provides the following functions:
Function|Content
----|----
Account creation|Ordinary accounts, administrator accounts, and other system accounts with special permissions. In the accountmanager contract, accountID is unique and can be used as an identity identifier
Account authorization|Ordinary authority is authorized when registering, and special authority requires the administrator to authorize
Account freeze and unfreeze|Account freeze is initiated by the administrator, the frozen account cannot be traded, and the assets under the account will be frozen
Account lock and recovery|used to reset the external private key when the private key is lost, a certain lock-up period is required, and the assets under the account cannot be transferred during the lock-up period
Account cancellation | The account should have a lifespan, the default time is five years, the expired account will be cancelled, and the query interface for the cancelled account is provided
Account assets|Account assets can be transferred normally under the accountmanager contract

The contract interface, online construction transaction and query interface reuse the CreateTransaction and Query interfaces in the framework respectively. For details, please refer to
[CreateTransaction interface](https://github.com/33cn/chain33/blob/master/rpc/jrpchandler.go#L1101) and [Query interface](https://github.com/33cn/chain33/blob/master /rpc/jrpchandler.go#L838)

Query method name|function
-----|----
QueryAccountByID|Query account information according to the account ID, which can be used to check whether the account ID is registered
QueryAccountsByStatus|Query account information based on status
QueryExpiredAccounts|Query Expired Time
QueryAccountByAddr|Query account information based on user address
QueryBalanceByID|Query account asset balance based on account ID


You can refer to the relevant test cases in account_test.go to construct relevant transactions for testing

## Precautions

**Table structure description**

Table name|Primary key|Index|Purpose|Description
 ---|---|---|---|---
 account|index|accountID,addr,status|Record registered account information|index is a composite index composed of {expiretime*1e5+index (the index of the registered exchange in the block)}

**Description of relevant parameters in the table**

Parameter name|Description
----|----
Asset|Asset name
op|Operation types are divided into supervisor op 1 for freezing, 2 for unfreezing, 3 for increasing the validity period, 4 for authorization apply op 1 to revoke the account public key reset, 2 after the lock period is over, perform the reset public key operation
status|Account status, 0 is normal, 1 means frozen, 2 means locked 3, expired and cancelled
level|Account authority 0 Normal, others can be customized according to business needs
index|Time stamp of account overdue*1e5+index of registered transaction in block, occupying 15%015d