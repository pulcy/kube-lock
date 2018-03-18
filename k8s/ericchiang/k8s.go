package ericchiang

import (
	"context"
	"fmt"
	"time"

	kc "github.com/ericchiang/k8s"
	"github.com/ericchiang/k8s/apis/core/v1"
	"github.com/ericchiang/k8s/apis/extensions/v1beta1"
	"github.com/juju/errgo"
	lock "github.com/pulcy/kube-lock"
)

// NewDaemonSetLock creates a lock that uses a DaemonSet to hold the lock data.
func NewDaemonSetLock(namespace, name string, c *kc.Client, annotationKey, ownerID string, ttl time.Duration) (lock.KubeLock, error) {
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

// NewDeploymentLock creates a lock that uses a Deployment to hold the lock data.
func NewDeploymentLock(namespace, name string, c *kc.Client, annotationKey, ownerID string, ttl time.Duration) (lock.KubeLock, error) {
	helper := &k8sHelper{
		name:      name,
		namespace: namespace,
		c:         c,
	}
	l, err := lock.NewKubeLock(annotationKey, ownerID, ttl, helper.deploymentGet, helper.deploymentUpdate)
	if err != nil {
		return nil, maskAny(err)
	}
	return l, nil
}

// NewReplicaSetLock creates a lock that uses a RepliceSet to hold the lock data.
func NewReplicaSetLock(namespace, name string, c *kc.Client, annotationKey, ownerID string, ttl time.Duration) (lock.KubeLock, error) {
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
func NewServiceLock(namespace, name string, c *kc.Client, annotationKey, ownerID string, ttl time.Duration) (lock.KubeLock, error) {
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

// NewNamespaceLock creates a lock that uses a Namespace to hold the lock data.
func NewNamespaceLock(namespace string, c *kc.Client, annotationKey, ownerID string, ttl time.Duration) (lock.KubeLock, error) {
	helper := &k8sHelper{
		name:      namespace,
		namespace: "",
		c:         c,
	}
	l, err := lock.NewKubeLock(annotationKey, ownerID, ttl, helper.namespaceGet, helper.namespaceUpdate)
	if err != nil {
		return nil, maskAny(err)
	}
	return l, nil
}

type k8sHelper struct {
	name      string
	namespace string
	c         *kc.Client
}

var (
	maskAny = errgo.MaskFunc(errgo.Any)
)

func (h *k8sHelper) daemonSetGet() (annotations map[string]string, resourceVersion string, extra interface{}, err error) {
	var daemonSet v1beta1.DaemonSet
	ctx := context.Background()
	if err := h.c.Get(ctx, h.namespace, h.name, &daemonSet); err != nil {
		return nil, "", nil, maskAny(err)
	}
	md := daemonSet.GetMetadata()
	return md.GetAnnotations(), md.GetResourceVersion(), &daemonSet, nil
}

func (h *k8sHelper) daemonSetUpdate(annotations map[string]string, resourceVersion string, extra interface{}) error {
	daemonSet, ok := extra.(*v1beta1.DaemonSet)
	if !ok {
		return maskAny(fmt.Errorf("extra must be *DaemonSet"))
	}
	md := daemonSet.GetMetadata()
	md.Annotations = annotations
	md.ResourceVersion = kc.String(resourceVersion)
	ctx := context.Background()
	if err := h.c.Update(ctx, daemonSet); err != nil {
		return maskAny(err)
	}
	return nil
}

func (h *k8sHelper) deploymentGet() (annotations map[string]string, resourceVersion string, extra interface{}, err error) {
	var deployment v1beta1.Deployment
	ctx := context.Background()
	if err := h.c.Get(ctx, h.namespace, h.name, &deployment); err != nil {
		return nil, "", nil, maskAny(err)
	}
	md := deployment.GetMetadata()
	return md.GetAnnotations(), md.GetResourceVersion(), &deployment, nil
}

func (h *k8sHelper) deploymentUpdate(annotations map[string]string, resourceVersion string, extra interface{}) error {
	deployment, ok := extra.(*v1beta1.Deployment)
	if !ok {
		return maskAny(fmt.Errorf("extra must be *Deployment"))
	}
	md := deployment.GetMetadata()
	md.Annotations = annotations
	md.ResourceVersion = kc.String(resourceVersion)
	ctx := context.Background()
	if err := h.c.Update(ctx, deployment); err != nil {
		return maskAny(err)
	}
	return nil
}

func (h *k8sHelper) replicaSetGet() (annotations map[string]string, resourceVersion string, extra interface{}, err error) {
	var replicaSet v1beta1.ReplicaSet
	ctx := context.Background()
	if err := h.c.Get(ctx, h.namespace, h.name, &replicaSet); err != nil {
		return nil, "", nil, maskAny(err)
	}
	md := replicaSet.GetMetadata()
	return md.GetAnnotations(), md.GetResourceVersion(), &replicaSet, nil
}

func (h *k8sHelper) replicaSetUpdate(annotations map[string]string, resourceVersion string, extra interface{}) error {
	replicaSet, ok := extra.(*v1beta1.ReplicaSet)
	if !ok {
		return maskAny(fmt.Errorf("extra must be *ReplicaSet"))
	}
	md := replicaSet.GetMetadata()
	md.Annotations = annotations
	md.ResourceVersion = kc.String(resourceVersion)
	ctx := context.Background()
	if err := h.c.Update(ctx, replicaSet); err != nil {
		return maskAny(err)
	}
	return nil
}

func (h *k8sHelper) serviceGet() (annotations map[string]string, resourceVersion string, extra interface{}, err error) {
	var service v1.Service
	ctx := context.Background()
	if err := h.c.Get(ctx, h.namespace, h.name, &service); err != nil {
		return nil, "", nil, maskAny(err)
	}
	md := service.GetMetadata()
	return md.GetAnnotations(), md.GetResourceVersion(), &service, nil
}

func (h *k8sHelper) serviceUpdate(annotations map[string]string, resourceVersion string, extra interface{}) error {
	service, ok := extra.(*v1.Service)
	if !ok {
		return maskAny(fmt.Errorf("extra must be *Service"))
	}
	md := service.GetMetadata()
	md.Annotations = annotations
	md.ResourceVersion = kc.String(resourceVersion)
	ctx := context.Background()
	if err := h.c.Update(ctx, service); err != nil {
		return maskAny(err)
	}
	return nil
}

func (h *k8sHelper) namespaceGet() (annotations map[string]string, resourceVersion string, extra interface{}, err error) {
	var namespace v1.Namespace
	ctx := context.Background()
	if err := h.c.Get(ctx, h.namespace, h.name, &namespace); err != nil {
		return nil, "", nil, maskAny(err)
	}
	md := namespace.GetMetadata()
	return md.GetAnnotations(), md.GetResourceVersion(), &namespace, nil
}

func (h *k8sHelper) namespaceUpdate(annotations map[string]string, resourceVersion string, extra interface{}) error {
	namespace, ok := extra.(*v1.Namespace)
	if !ok {
		return maskAny(fmt.Errorf("extra must be *Namespace"))
	}
	md := namespace.GetMetadata()
	md.Annotations = annotations
	md.ResourceVersion = kc.String(resourceVersion)
	ctx := context.Background()
	if err := h.c.Update(ctx, namespace); err != nil {
		return maskAny(err)
	}
	return nil
}
