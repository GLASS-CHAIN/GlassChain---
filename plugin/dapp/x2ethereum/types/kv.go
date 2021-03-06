package types

import "strings"

/*
 * k ，ke 
 * key = keyPrefix + userKey
 *  ’- 
 */

var (
	//KeyPrefixStateDB state db ke 
	KeyPrefixStateDB = "mavl-x2ethereum-"
	//KeyPrefixLocalDB local d ke 
	KeyPrefixLocalDB = "LODB-x2ethereum-"
)

//CalProphecyPrefix ...
func CalProphecyPrefix(id string) []byte {
	return []byte(KeyPrefixStateDB + string(ProphecyKey) + id)
}

//CalEth2Chain33Prefix ...
func CalEth2Chain33Prefix() []byte {
	return []byte(KeyPrefixStateDB + string(Eth2Chain33Key))
}

//CalWithdrawEthPrefix ...
func CalWithdrawEthPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(WithdrawEthKey))
}

//CalChain33ToEthPrefix ...
func CalChain33ToEthPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(Chain33ToEthKey))
}

//CalWithdrawChain33Prefix ...
func CalWithdrawChain33Prefix() []byte {
	return []byte(KeyPrefixStateDB + string(WithdrawChain33Key))
}

//CalValidatorMapsPrefix ...
func CalValidatorMapsPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(ValidatorMapsKey))
}

//CalLastTotalPowerPrefix ...
func CalLastTotalPowerPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(LastTotalPowerKey))
}

//CalConsensusThresholdPrefix ...
func CalConsensusThresholdPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(ConsensusThresholdKey))
}

//CalTokenSymbolTotalLockOrBurnAmount ...
func CalTokenSymbolTotalLockOrBurnAmount(symbol, tokenAddress, direction, txType string) []byte {
	return []byte(KeyPrefixLocalDB + string(TokenSymbolTotalLockOrBurnAmountKey) + direction + "-" + txType + "-" + strings.ToLower(symbol) + "-" + strings.ToLower(tokenAddress))
}

//CalTokenSymbolToTokenAddress ...
func CalTokenSymbolToTokenAddress(symbol string) []byte {
	return []byte(KeyPrefixLocalDB + string(TokenSymbolToTokenAddressKey) + strings.ToLower(symbol))
}
