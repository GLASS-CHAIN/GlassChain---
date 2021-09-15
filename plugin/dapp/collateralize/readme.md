## Loan Contract Table Structure

### Loan table coller table structure
Field name|Type|Description
---|---|---
collateralizeId|string|Lending ID, primary key
accountAddr|string|big account address
recordId|string|Loan ID
status|int32|Lending status (1: loaned 2: recovered)

### Loan table coller table index
Index name|Description
---|---
status|Query the loan ID according to the loan status
addr|Query the loan ID according to the address of the big account
addr_status|Query the loan ID according to the loan status and the address of the major account

### Borrow table structure
Field name|Type|Description
---|---|---
recordId|string|Loan ID, primary key
collateralizeId|string|Lending ID
accountAddr|string|User address
status|int32|Lending status (1: issued 2: price clearing alarm 3: price clearing 4: overtime clearing alarm 5: overtime clearing 6: cleared)

### Loan table borrow table index
Index name|Description
---|---
status|Query the loan ID according to the loan status
addr|Query the loan ID according to the address of the big account
addr_status|Query the loan ID according to the loan status and user address
id_status|Query the loan ID according to the loan ID and loan status
id_addr|Query the loan ID according to the loan ID and user address