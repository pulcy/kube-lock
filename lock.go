package lock

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/juju/errgo"
)

// KubeLock is used to provide a distributed lock using Kubernetes annotation data.
// It works by writing data into a specific annotation key.
// Other instance trying to write into the same annotation key will be refused because a resource version is used.
type KubeLock interface {
	// Acquire tries to acquire the lock.
	// If the lock is already held by us, the lock will be updated.
	// If successfull it returns nil, otherwise it returns an error.
	// Note that Acquire will not renew the lock. To do that, call Acquire every ttl/2.
	Acquire() error

	// Release tries to release the lock.
	// If the lock is already held by us, the lock will be released.
	// If successfull it returns nil, otherwise it returns an error.
	Release() error

	// CurrentOwner fetches the current owner ID of the lock.
	// If the lock is not owner, "" is returned.
	CurrentOwner() (string, error)
}

// NewKubeLock creates a new KubeLock.
// The lock will not be aquired.
func NewKubeLock(annotationKey, ownerID string, ttl time.Duration, metaGet MetaGetter, metaUpdate MetaUpdater) (KubeLock, error) {
	if annotationKey == "" {
		annotationKey = defaultAnnotationKey
	}
	if ownerID == "" {
		id := make([]byte, 16)
		if _, err := rand.Read(id); err != nil {
			return nil, maskAny(err)
		}
		ownerID = base64.StdEncoding.EncodeToString(id)
	}
	if ttl == 0 {
		ttl = defaultTTL
	}
	if metaGet == nil {
		return nil, maskAny(fmt.Errorf("metaGet cannot be nil"))
	}
	if metaUpdate == nil {
		return nil, maskAny(fmt.Errorf("metaUpdate cannot be nil"))
	}
	return &kubeLock{
		annotationKey: annotationKey,
		ownerID:       ownerID,
		ttl:           ttl,
		getMeta:       metaGet,
		updateMeta:    metaUpdate,
	}, nil
}

const (
	defaultAnnotationKey = "pulcy.com/kube-lock"
	defaultTTL           = time.Minute
)

type kubeLock struct {
	annotationKey string
	ownerID       string
	ttl           time.Duration
	getMeta       MetaGetter
	updateMeta    MetaUpdater
}

type LockData struct {
	Owner     string    `json:"owner"`
	ExpiresAt time.Time `json:"expires_at"`
}

type MetaGetter func() (annotations map[string]string, resourceVersion string, extra interface{}, err error)
type MetaUpdater func(annotations map[string]string, resourceVersion string, extra interface{}) error

// Acquire tries to acquire the lock.
// If the lock is already held by us, the lock will be updated.
// If successfull it returns nil, otherwise it returns an error.
func (l *kubeLock) Acquire() error {
	// Get current state
	ann, rv, extra, err := l.getMeta()
	if err != nil {
		return maskAny(err)
	}

	// Get lock data
	if lockDataRaw, ok := ann[l.annotationKey]; ok && lockDataRaw != "" {
		var lockData LockData
		if err := json.Unmarshal([]byte(lockDataRaw), &lockData); err != nil {
			return maskAny(err)
		}
		if lockData.Owner != l.ownerID {
			// Lock is owned by someone else
			if time.Now().Before(lockData.ExpiresAt) {
				// Lock is held and not expired
				return maskAny(errgo.WithCausef(nil, AlreadyLockedError, "locked by %s", lockData.Owner))
			}
		}
	}

	// Try to lock it now
	expiredAt := time.Now().Add(l.ttl)
	lockDataRaw, err := json.Marshal(LockData{Owner: l.ownerID, ExpiresAt: expiredAt})
	if err != nil {
		return maskAny(err)
	}
	ann[l.annotationKey] = string(lockDataRaw)
	if err := l.updateMeta(ann, rv, extra); err != nil {
		return maskAny(err)
	}

	// Update successfull, we've acquired the lock
	return nil
}

// Release tries to release the lock.
// If the lock is already held by us, the lock will be released.
// If successfull it returns nil, otherwise it returns an error.
func (l *kubeLock) Release() error {
	// Get current state
	ann, rv, extra, err := l.getMeta()
	if err != nil {
		return maskAny(err)
	}

	// Get lock data
	if lockDataRaw, ok := ann[l.annotationKey]; ok && lockDataRaw != "" {
		var lockData LockData
		if err := json.Unmarshal([]byte(lockDataRaw), &lockData); err != nil {
			return maskAny(err)
		}
		if lockData.Owner != l.ownerID {
			// Lock is owned by someone else
			return maskAny(errgo.WithCausef(nil, NotLockedByMeError, "locked by %s", lockData.Owner))
		}
	} else if ok && lockDataRaw == "" {
		// Lock is not locked, we consider that a successfull release also.
		return nil
	}

	// Try to release lock it now
	ann[l.annotationKey] = ""
	if err := l.updateMeta(ann, rv, extra); err != nil {
		return maskAny(err)
	}

	// Update successfull, we've released the lock
	return nil
}

// CurrentOwner fetches the current owner ID of the lock.
// If the lock is not owner, "" is returned.
func (l *kubeLock) CurrentOwner() (string, error) {
	// Get current state
	ann, _, _, err := l.getMeta()
	if err != nil {
		return "", maskAny(err)
	}

	// Get lock data
	if lockDataRaw, ok := ann[l.annotationKey]; ok && lockDataRaw != "" {
		var lockData LockData
		if err := json.Unmarshal([]byte(lockDataRaw), &lockData); err != nil {
			return "", maskAny(err)
		}
		return lockData.Owner, nil
	}

	// No owner found
	return "", nil
}
