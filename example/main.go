package main

import (
	"flag"
	"log"
	"time"

	http "github.com/YakLabs/k8s-client/http"
	lock "github.com/pulcy/kube-lock"
	k8s "github.com/pulcy/kube-lock/k8s"
)

var (
	args struct {
		namespace      string
		replicaSetName string
		serviceName    string
	}
)

func init() {
	flag.StringVar(&args.namespace, "namespace", "", "Kubernetes namespace of resource")
	flag.StringVar(&args.replicaSetName, "replicaSet", "", "Kubernetes namespace of ReplicaSet to store lock data in")
	flag.StringVar(&args.serviceName, "service", "", "Kubernetes namespace of Service to store lock data in")
}

func main() {
	flag.Parse()

	if args.namespace == "" {
		log.Fatalln("-namespace not set")
	}
	if args.replicaSetName == "" && args.serviceName == "" {
		log.Fatalln("-replicaSet or -service must be set")
	}

	c, err := http.NewInCluster()
	if err != nil {
		log.Fatalf("Cannot create k8s client: %#v\n", err)
	}

	ttl := time.Second * 30
	var l lock.KubeLock
	if args.serviceName != "" {
		l, err = k8s.NewServiceLock(args.namespace, args.serviceName, c, "", "", ttl)
	} else if args.replicaSetName != "" {
		l, err = k8s.NewReplicaSetLock(args.namespace, args.replicaSetName, c, "", "", ttl)
	} else {
		log.Fatalln("Unknown resource")
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
