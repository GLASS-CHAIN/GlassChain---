package executor

/*
 * k ，ke 
 * key = keyPrefix + userKey
 *  ’- 
 */

var (
	//keyPrefixStateDB state db ke 
	keyPrefixStateDB = "mavl-vote-"
	//keyPrefixLocalDB local d ke 
	keyPrefixLocalDB = "LODB-vote-"
)

// groupID or voteID
func formatStateIDKey(id string) []byte {
	return []byte(keyPrefixStateDB + id)
}
