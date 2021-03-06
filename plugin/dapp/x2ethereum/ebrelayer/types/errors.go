package types

import (
	"errors"
)

// relaye 
var (
	ErrInvalidEthContractAddress = errors.New("ErrInvalidEthContractAddress")
	ErrUnpack                    = errors.New("ErrUnpack")
	ErrEmptyAddress              = errors.New("ErrEmptyAddress")
	ErrAddress4Eth               = errors.New("symbol \"eth\" must have null address set as token address")
	ErrPublicKeyType             = errors.New("ErrPublicKeyType")
	ErrInvalidContractAddress    = errors.New("ErrInvalidTargetContractAddress")
	ErrNoValidatorConfigured     = errors.New("ErrNoValidatorConfigured")
)
