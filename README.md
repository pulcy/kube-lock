# Kube-lock

[![GoDoc](https://godoc.org/github.com/pulcy/kube-lock?status.svg)](http://godoc.org/github.com/pulcy/kube-lock)

Kube-lock is a simple Go library that implementation a distributed lock using annotations on a Kubernetes resource.

See [examples](./examples) folder for a simple example for both libraries.

# Details

In this folder you'll find the basic lock functionality.
This is abstracted using `get` and `update` functions.

In the [k8s/ericchiang](./k8s/ericchiang) folder you'll find a Kubernetes specific implementation using the lightweight yet comprehensive [ericchiang/k8s](https://github.com/ericchiang/k8s).
It implements `get` & `update` functions for various resources.

In the [k8s/yaklabs](./k8s/yaklabs) folder you'll find a Kubernetes specific implementation using the lightweight [YakLabs/k8s-client](https://github.com/YakLabs/k8s-client).
It implements `get` & `update` functions for various resources.