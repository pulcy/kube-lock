package k8s

import (
	"fmt"
	"time"

	kc "github.com/YakLabs/k8s-client"
	"github.com/juju/errgo"
	lock "github.com/pulcy/kube-lock"
)

// NewDaemonSetLock creates a lock that uses a DaemonSet to hold the lock data.
func NewDaemonSetLock(namespace, name string, c kc.Client, annotationKey, ownerID string, ttl time.Duration) (lock.KubeLock, error) {
	helper := &k8sHelper{
		name:      name,
		namespace: namespace,
		c:         c,
	}
	l, err := lock.NewKubeLock(annotationKey, ownerID, ttl, helper.daemonSetGet, helper.daemonSetUpdate)
	if err != nil {
		return nil, maskAny(err)
	}
	return l, nil
}

// NewReplicaSetLock creates a lock that uses a RepliceSet to hold the lock data.
func NewReplicaSetLock(namespace, name string, c kc.Client, annotationKey, ownerID string, ttl time.Duration) (lock.KubeLock, error) {
	helper := &k8sHelper{
		name:      name,
		namespace: namespace,
		c:         c,
	}
	l, err := lock.NewKubeLock(annotationKey, ownerID, ttl, helper.replicaSetGet, helper.replicaSetUpdate)
	if err != nil {
		return nil, maskAny(err)
	}
	return l, nil
}

// NewServiceLock creates a lock that uses a Service to hold the lock data.
func NewServiceLock(namespace, name string, c kc.Client, annotationKey, ownerID string, ttl time.Duration) (lock.KubeLock, error) {
	helper := &k8sHelper{
		name:      name,
		namespace: namespace,
		c:         c,
	}
	l, err := lock.NewKubeLock(annotationKey, ownerID, ttl, helper.serviceGet, helper.serviceUpdate)
	if err != nil {
		return nil, maskAny(err)
	}
	return l, nil
}

type k8sHelper struct {
	name      string
	namespace string
	c         kc.Client
}

var (
	maskAny = errgo.MaskFunc(errgo.Any)
)

func (h *k8sHelper) daemonSetGet() (annotations map[string]string, resourceVersion string, extra interface{}, err error) {
	daemonSet, err := h.c.GetDaemonSet(h.namespace, h.name)
	if err != nil {
		return nil, "", nil, maskAny(err)
	}
	return daemonSet.ObjectMeta.Annotations, daemonSet.ObjectMeta.ResourceVersion, daemonSet, nil
}

func (h *k8sHelper) daemonSetUpdate(annotations map[string]string, resourceVersion string, extra interface{}) error {
	daemonSet, ok := extra.(*kc.DaemonSet)
	if !ok {
		return maskAny(fmt.Errorf("extra must be *DaemonSet"))
	}
	daemonSet.ObjectMeta.Annotations = annotations
	daemonSet.ObjectMeta.ResourceVersion = resourceVersion
	if _, err := h.c.UpdateDaemonSet(h.namespace, daemonSet); err != nil {
		return maskAny(err)
	}
	return nil
}

func (h *k8sHelper) replicaSetGet() (annotations map[string]string, resourceVersion string, extra interface{}, err error) {
	replicaSet, err := h.c.GetReplicaSet(h.namespace, h.name)
	if err != nil {
		return nil, "", nil, maskAny(err)
	}
	return replicaSet.ObjectMeta.Annotations, replicaSet.ObjectMeta.ResourceVersion, replicaSet, nil
}

func (h *k8sHelper) replicaSetUpdate(annotations map[string]string, resourceVersion string, extra interface{}) error {
	replicaSet, ok := extra.(*kc.ReplicaSet)
	if !ok {
		return maskAny(fmt.Errorf("extra must be *ReplicaSet"))
	}
	replicaSet.ObjectMeta.Annotations = annotations
	replicaSet.ObjectMeta.ResourceVersion = resourceVersion
	if _, err := h.c.UpdateReplicaSet(h.namespace, replicaSet); err != nil {
		return maskAny(err)
	}
	return nil
}

func (h *k8sHelper) serviceGet() (annotations map[string]string, resourceVersion string, extra interface{}, err error) {
	service, err := h.c.GetService(h.namespace, h.name)
	if err != nil {
		return nil, "", nil, maskAny(err)
	}
	return service.ObjectMeta.Annotations, service.ObjectMeta.ResourceVersion, service, nil
}

func (h *k8sHelper) serviceUpdate(annotations map[string]string, resourceVersion string, extra interface{}) error {
	service, ok := extra.(*kc.Service)
	if !ok {
		return maskAny(fmt.Errorf("extra must be *Service"))
	}
	service.ObjectMeta.Annotations = annotations
	service.ObjectMeta.ResourceVersion = resourceVersion
	if _, err := h.c.UpdateService(h.namespace, service); err != nil {
		return maskAny(err)
	}
	return nil
}
