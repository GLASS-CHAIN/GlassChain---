##

SRC_CONTRACT := safe-contracts/contracts
GO_OUT0 := bsctest/gnosisSafe/

PACKAGE0 := gnosisSafe

proj := "build"
.PHONY: default build clean

default: sol build

sol:
	@abigen --sol $(SRC_CONTRACT)/GnosisSafe.sol --pkg $(PACKAGE0) --out $(GO_OUT0)/gnosisSafe.go

build:
	cd bsctest && go build && cd ..


clean:
	@rm -fr $(GO_OUT)/*.go
