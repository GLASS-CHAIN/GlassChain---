
## 
  


### 
  

#### (CreateGroup)

-  
- 
-  1
- I 

##### 
```proto

/ 
message CreateGroup {
    string   name                    = 1; / 
    repeated string admins           = 2; / , 
    repeated GroupMember members     = 3; / 
    string               description = 4; / 
}

message GroupMember {
    string addr       = 1; / 
    uint32 voteWeight = 2; / ， 1
    string nickName   = 3; / 
}

```

##### 
```proto

// 
message GroupInfo {

    string   ID                      = 1; / ID
    string   name                    = 2; / 
    uint32   memberNum               = 3; / 
    string   creator                 = 4; / 
    repeated string admins           = 5; / 
    repeated GroupMember members     = 6; / 
    string               description = 7; / 
    uint32 voteNum = 8; / 
}

```
##### 

- json rp ，Chain33.CreateTransaction
- actionName: CreateGroup

```bash
curl -kd  '{"method":"Chain33.CreateTransaction","params":[{"execer":"vote","actionName":"CreateGroup","payload":{"name":"group30","admins":[],"members":[{"addr":"1BQXS6TxaYYG5mADaWij4AxhZZUTpw95a5","voteWeight":0}],"description":""}}],"id":0}' http://localhost:8801
```

#### (UpdateGroup)
 ID  

##### 
```proto
message UpdateGroup {
    string   groupID                = 1; / ID
    repeated GroupMember addMembers = 2; / 
    repeated string removeMembers   = 3; / 
    repeated string addAdmins       = 4; / 
    repeated string removeAdmins    = 5; / 
}

message GroupMember {
    string addr       = 1; / 
    uint32 voteWeight = 2; / ， 1
    string nickName   = 3; / 
}
```
##### 
```proto

// 
message GroupInfo {

    string   ID                      = 1; / ID
    string   name                    = 2; / 
    uint32   memberNum               = 3; / 
    string   creator                 = 4; / 
    repeated string admins           = 5; / 
    repeated GroupMember members     = 6; / 
    string               description = 7; / 
    uint32 voteNum = 8; / 
}

```

##### 

- json rp ，Chain33.CreateTransaction
- actionName: UpdateGroup

```bash
curl -kd  '{"method":"Chain33.CreateTransaction","params":[{"execer":"vote","actionName":"UpdateGroup","payload":{"groupID":"g000000000000700000","addMembers":[{"addr":"member1","voteWeight":0},{"addr":"member2","voteWeight":0}],"removeMembers":["member3"],"addAdmins":["admin1"],"removeAdmins":["admin2"]}}],"id":0}' http://localhost:8801
```

#### (CreateVote)
- 
-   ID
-  
- I 

##### 
```proto

//  
message CreateVote {
    string   name                  = 1; / 
    string   groupID               = 2; / 
    repeated string voteOptions    = 3; / 
    int64           beginTimestamp = 4; / 
    int64           endTimestamp   = 5; / 
    string          description    = 6; / 
}

```

##### 
```proto

/ 
message VoteInfo {

    string   ID                        = 1;  / ID
    string   name                      = 2;  / 
    string   creator                   = 3;  / 
    string   groupID                   = 4;  / 
    repeated VoteOption voteOptions    = 5;  / 
    int64               beginTimestamp = 6;  / 
    int64               endTimestamp   = 7;  / 
    repeated CommitInfo commitInfos    = 8;  / 
    string              description    = 9;  / 
    uint32              status         = 10; / ， ， ， ， 
}

/ 
message VoteOption {
    string option = 1; / 
    uint32 score  = 2; / 
}
```

##### 

- json rp ，Chain33.CreateTransaction
- actionName: CreateVote

```bash
curl -kd  '{"method":"Chain33.CreateTransaction","params":[{"execer":"vote","actionName":"CreateVote","payload":{"name":"vote1","groupID":"g000000000000600000","voteOptions":["A","B","C"],"beginTimestamp":"1611562096","endTimestamp":"1611648496","description":""}}],"id":0}' http://localhost:8801
```

#### (CommitVote)
- 
- ID 
-  

##### 
```proto

//  
message CommitVote {
    string voteID      = 1; / ID
    uint32 optionIndex = 2; /  
}

```

##### 
```proto

/ 
message CommitInfo {
    string addr       = 1; / 
    string txHash     = 2; / 
    uint32 voteWeight = 3; / 
}
```

##### 

- json rp ，Chain33.CreateTransaction
- actionName: CommitVote

```bash
curl -kd  '{"method":"Chain33.CreateTransaction","params":[{"execer":"vote","actionName":"CommitVote","payload":{"voteID":"v000000000001300000","optionIndex":0}}],"id":0}' http://localhost:8801

```

#### (CloseVote)
-  

##### 
```proto

message CloseVote {
    string voteID = 1; // ID
}

```

##### 
```proto

message VoteInfo {

    string   ID                        = 1;  / ID
    string   name                      = 2;  / 
    string   creator                   = 3;  / 
    string   groupID                   = 4;  / 
    repeated VoteOption voteOptions    = 5;  / 
    int64               beginTimestamp = 6;  / 
    int64               endTimestamp   = 7;  / 
    repeated CommitInfo commitInfos    = 8;  / 
    string              description    = 9;  / 
    uint32              status         = 10; / ， ， ， ， 
    string              groupName      = 11; / 
}
```

##### 

- json rp ，Chain33.CreateTransaction
- actionName: CloseVote

```bash
curl -kd  '{"method":"Chain33.CreateTransaction","params":[{"execer":"vote","actionName":"CloseVote","payload":{"voteID":"v000000000001300000"}}],"id":0}' http://localhost:8801
```

#### (UpdateMember)
- 

##### 
```proto

message UpdateMember {
    string name = 1; / 
}

```

##### 
```proto

message MemberInfo {
    string   addr            = 1; / 
    string   name            = 2; / 
    repeated string groupIDs = 3; / I 
}
```

##### 

- json rp ，Chain33.CreateTransaction
- actionName: UpdateMember

```bash
 curl -kd  '{"method":"Chain33.CreateTransaction","params":[{"execer":"vote","actionName":"UpdateMember","payload":{"name":"name1"}}],"id":0}' http://localhost:8801
```


### 
  

#### (GetGroups)
 I  

##### 
```proto
message ReqStrings {
    repeated string items = 1; / I 
}
```

##### 

```proto
message GroupInfos {
    repeated GroupInfo groupList = 1; / 
}
```
##### 

- json rp ，Chain33.Query
- funcName: GetGroups

```bash
curl -ksd '{"method":"Chain33.Query","params":[{"execer":"vote","funcName":"GetGroups","payload":{"items":["g000000000001700000","g000000000001800000"]}}],"id":0}' http://localhost:8801
```

#### (GetVotes)
 I  

##### 
```proto
message ReqStrings {
    repeated string items = 1; / I 
}
```

##### 
```proto
message ReplyVoteList {
    repeated VoteInfo voteList         = 1; / 
    int64             currentTimestamp = 2; / 
}
```

##### 

- json rp ，Chain33.Query
- funcName: GetVotes

```bash
curl -kd  '{"method":"Chain33.Query","params":[{"execer":"vote","funcName":"GetVotes","payload":{"items":["v000000000001300000","v000000000001400000"]}}],"id":0}' http://localhost:8801
```

#### (GetMembers)
   

##### 
```proto
message ReqStrings {
    repeated string items = 1; / 
}
```

##### 
```proto

message MemberInfos {
    repeated MemberInfo memberList = 1; / 
}

message MemberInfo {
    string   addr            = 1; / 
    string   name            = 2; / 
    repeated string groupIDs = 3; / I 
}
```

##### 

- json rp ，Chain33.Query
- funcName: GetMembers
-
```bash
curl -kd  '{"method":"Chain33.Query","params":[{"execer":"vote","funcName":"GetMembers","payload":{"items":["1BQXS6TxaYYG5mADaWij4AxhZZUTpw95a5"]}}],"id":0}' http://localhost:8801
```

#### (ListGroup)
 

##### 
```proto
/ 
message ReqListItem {
    string startItemID = 1; / ID groupID 
    int32  count       = 2; / ,  
    int32  direction   = 3; //  I ，  I 
}
```

##### 
```proto
message GroupInfos {
    repeated GroupInfo groupList = 1; / 
}
```

##### 

- json rp ，Chain33.Query
- funcName: ListGroup

```bash
curl -kd  '{"method":"Chain33.Query","params":[{"execer":"vote","funcName":"ListGroup","payload":{"startItemID":"","count":2,"direction":0}}],"id":0}' http://localhost:8801
```

#### (ListVote)
- 
- groupID 
- , status  

##### 
```proto
/ 
message ReqListVote {
    string      groupID = 1; / ID
    ReqListItem listReq = 2; / 
    uint32      status  = 3; / ,  ， ， ， 
}

message ReqListItem {
    string startItemID = 1; / ID groupID 
    int32  count       = 2; / ,  
    int32  direction   = 3; //  I ，  I 
}
```


##### 
```proto
message ReplyVoteList {
    repeated VoteInfo voteList         = 1; / 
    int64             currentTimestamp = 2; / 
}
```

##### 

- json rp ，Chain33.Query
- funcName: ListVote

```bash
curl -kd  '{"method":"Chain33.Query","params":[{"execer":"vote","funcName":"ListVote","payload":{"groupID":"","listReq":{"startItemID":"","count":2,"direction":0}}}],"id":0}' http://localhost:8801
```

#### (ListMember)
 

##### 

```proto
/ 
message ReqListItem {
    string startItemID = 1; / ID groupID 
    int32  count       = 2; / ,  
    int32  direction   = 3; //  I ，  I 
}
```

##### 
```proto
message MemberInfos {
    repeated MemberInfo memberList = 1; / 
}
```

##### 

- json rp ，Chain33.Query
- funcName: ListMember

```bash
curl -kd  '{"method":"Chain33.Query","params":[{"execer":"vote","funcName":"ListMember","payload":{"startItemID":"","count":1,"direction":1}}],"id":0}' http://localhost:8801
```

#### 
 

   
|---|---|
errEmptyName|       
errInvalidMemberWeights | 
errDuplicateMember      | 
errDuplicateAdmin       | 
errInvalidVoteTime      | 
errInvalidVoteOption    | 
errVoteNotExist         | 
errGroupNotExist        | 
errStateDBGet           | 
errInvalidVoteID        | ID
errInvalidGroupID       | ID
errInvalidOptionIndex   | 
errAddrAlreadyVoted     | 
errVoteAlreadyFinished  | 
errVoteNotStarted       | 
errVoteAlreadyClosed    | 
errAddrPermissionDenied | 




#### 

 prot ](proto/vote.proto)
