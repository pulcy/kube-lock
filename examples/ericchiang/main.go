package main

import (
	"flag"
	"log"
	"time"

	kc "github.com/ericchiang/k8s"
	lock "github.com/pulcy/kube-lock"
	k8s "github.com/pulcy/kube-lock/k8s/ericchiang"
)

var (
	args struct {
		namespace      string
		daemonSetName  string
		deploymentName string
		replicaSetName string
		serviceName    string
	}
)

func init() {
	flag.StringVar(&args.namespace, "namespace", "", "Kubernetes namespace of resource")
	flag.StringVar(&args.daemonSetName, "daemonSet", "", "Kubernetes namespace of DaemonSet to store lock data in")
	flag.StringVar(&args.deploymentName, "deployment", "", "Kubernetes namespace of Deployment to store lock data in")
	flag.StringVar(&args.replicaSetName, "replicaSet", "", "Kubernetes namespace of ReplicaSet to store lock data in")
	flag.StringVar(&args.serviceName, "service", "", "Kubernetes namespace of Service to store lock data in")
}

func main() {
	flag.Parse()

	if args.namespace == "" {
		log.Fatalln("-namespace not set")
	}

	c, err := kc.NewInClusterClient()
	if err != nil {
		log.Fatalf("Cannot create k8s client: %#v\n", err)
	}

	ttl := time.Second * 30
	var l lock.KubeLock
	if args.daemonSetName != "" {
		l, err = k8s.NewDaemonSetLock(args.namespace, args.daemonSetName, c, "", "", ttl)
	} else if args.deploymentName != "" {
		l, err = k8s.NewDeploymentLock(args.namespace, args.deploymentName, c, "", "", ttl)
	} else if args.serviceName != "" {
		l, err = k8s.NewServiceLock(args.namespace, args.serviceName, c, "", "", ttl)
	} else if args.replicaSetName != "" {
		l, err = k8s.NewReplicaSetLock(args.namespace, args.replicaSetName, c, "", "", ttl)
	} else {
		l, err = k8s.NewNamespaceLock(args.namespace, c, "", "", ttl)
	}

	for {
		if err := l.Acquire(); err == nil {
			log.Println("Lock acquired")
			for i := 0; i < 3; i++ {
				time.Sleep(ttl / 2)
				// Renew lock
				if err := l.Acquire(); err != nil {
					log.Printf("Lock renewal failed: %#v\n", err)
					i = 1000
				} else {
					log.Println("Lock renewed")
				}
			}
			// Release lock
			if err := l.Release(); err != nil {
				log.Printf("Failed to release lock: %#v\n", err)
			} else {
				log.Println("Relesed lock")
			}
			time.Sleep(time.Second * 10)
		} else {
			log.Printf("Cannot acquire lock: %v\n", err)
			time.Sleep(time.Second * 2)
		}
	}
}
