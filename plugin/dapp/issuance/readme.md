## Issuing contract table structure

### Issuer table structure of the total issuance table
Field name|Type|Description
---|---|---
issuanceId|string|Total issuance ID, primary key
status|int32|Release status (1: issued 2: offline)

### Total issue table issuer table index
Index name|Description
---|---
status|Query the total issuance ID according to the issuance status

### Debt table structure
Field name|Type|Description
---|---|---
debtId|string|Issuance ID of major account, primary key
issuanceId|string|Total issuance ID
accountAddr|string|User address
status|int32|Issuance status (1: issued 2: price clearing alarm 3: price clearing 4: overtime clearing alarm 5: overtime clearing 6: closed)

### The index of the debt table of the large-scale release table
Index name|Description
---|---
status|Query the issuance ID of the major account according to the status of the major account issuance
addr|Query the ID issued by the major account according to the address of the major account
addr_status|Query the issuance ID of the major account based on the issuance status and the address of the major account