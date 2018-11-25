package common

import "errors"

var (
	LockAlreadyRequiredErr = errors.New("lock has been already required")
	NetworkLocalIpNotFound = errors.New("network local ip not found")
)
