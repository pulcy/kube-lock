package lock

import (
	"github.com/juju/errgo"
)

var (
	maskAny            = errgo.MaskFunc(errgo.Any)
	AlreadyLockedError = errgo.New("already locked")
	NotLockedByMeError = errgo.New("not locked by me")
)

// IsAlreadyLocked returns true if the given error is caused by a AlreadyLockedError error.
func IsAlreadyLocked(err error) bool {
	return errgo.Cause(err) == AlreadyLockedError
}

// IsNotLockedByMe returns true if the given error is caused by a NotLockedByMeError error.
func IsNotLockedByMe(err error) bool {
	return errgo.Cause(err) == NotLockedByMeError
}
