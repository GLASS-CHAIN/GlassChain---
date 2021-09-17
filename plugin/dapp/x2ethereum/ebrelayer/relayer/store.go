package relayer

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/33cn/chain33/common/crypto"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
)

var (
	storelog = log15.New("relayer manager", "store")
)

const (
	keyEncryptionFlag     = "Encryption"
	keyEncryptionCompFlag = "EncryptionFlag" //   
	keyPasswordHash       = "PasswordHash"
)

// CalcEncryptionFlag Key
func calcEncryptionFlag() []byte {
	return []byte(keyEncryptionFlag)
}

// calckeyEncryptionCompFlag Key
func calckeyEncryptionCompFlag() []byte {
	return []byte(keyEncryptionCompFlag)
}

// CalcPasswordHash has Key
func calcPasswordHash() []byte {
	return []byte(keyPasswordHash)
}

// NewStore 
func NewStore(db db.DB) *Store {
	return &Store{db: db}
}

// Store  
type Store struct {
	db db.DB
}

// Close 
func (store *Store) Close() {
	store.db.Close()
}

// GetDB 
func (store *Store) GetDB() db.DB {
	return store.db
}

// NewBatch 
func (store *Store) NewBatch(sync bool) db.Batch {
	return store.db.NewBatch(sync)
}

// Get 
func (store *Store) Get(key []byte) ([]byte, error) {
	return store.db.Get(key)
}

// Set 
func (store *Store) Set(key []byte, value []byte) (err error) {
	return store.db.Set(key, value)
}

// NewListHelper 
func (store *Store) NewListHelper() *db.ListHelper {
	return db.NewListHelper(store.db)
}

// SetEncryptionFlag 
func (store *Store) SetEncryptionFlag(batch db.Batch) error {
	var flag int64 = 1
	data, err := json.Marshal(flag)
	if err != nil {
		storelog.Error("SetEncryptionFlag marshal flag", "err", err)
		return types.ErrMarshal
	}

	batch.Set(calcEncryptionFlag(), data)
	return nil
}

// GetEncryptionFlag 
func (store *Store) GetEncryptionFlag() int64 {
	var flag int64
	data, err := store.Get(calcEncryptionFlag())
	if data == nil || err != nil {
		data, err = store.Get(calckeyEncryptionCompFlag())
		if data == nil || err != nil {
			return 0
		}
	}
	err = json.Unmarshal(data, &flag)
	if err != nil {
		storelog.Error("GetEncryptionFlag unmarshal", "err", err)
		return 0
	}
	return flag
}

// SetPasswordHash 
func (store *Store) SetPasswordHash(password string, batch db.Batch) error {
	var WalletPwHash types.WalletPwHash
	/ 
	randstr := fmt.Sprintf("fuzamei:$@%s", crypto.CRandHex(16))
	WalletPwHash.Randstr = randstr

	/ passwor has 
	pwhashstr := fmt.Sprintf("%s:%s", password, WalletPwHash.Randstr)
	pwhash := sha256.Sum256([]byte(pwhashstr))
	WalletPwHash.PwHash = pwhash[:]

	pwhashbytes, err := json.Marshal(WalletPwHash)
	if err != nil {
		storelog.Error("SetEncryptionFlag marshal flag", "err", err)
		return types.ErrMarshal
	}
	batch.Set(calcPasswordHash(), pwhashbytes)
	return nil
}

// VerifyPasswordHash 
func (store *Store) VerifyPasswordHash(password string) bool {
	var WalletPwHash types.WalletPwHash
	pwhashbytes, err := store.Get(calcPasswordHash())
	if pwhashbytes == nil || err != nil {
		return false
	}
	err = json.Unmarshal(pwhashbytes, &WalletPwHash)
	if err != nil {
		storelog.Error("VerifyPasswordHash unmarshal", "err", err)
		return false
	}
	pwhashstr := fmt.Sprintf("%s:%s", password, WalletPwHash.Randstr)
	pwhash := sha256.Sum256([]byte(pwhashstr))
	Pwhash := pwhash[:]
	/ pwhas 
	return bytes.Equal(WalletPwHash.GetPwHash(), Pwhash)
}
