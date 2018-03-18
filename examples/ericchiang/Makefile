ROOTDIR := $(shell cd ../.. && pwd)
IMAGE := pulcy/kube-lock-example-ericchiang

all:
	docker run \
		--rm \
		-v $(ROOTDIR):/usr/code \
		-e GOPATH=/usr/code/.gobuild \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e CGO_ENABLED=0 \
		-w /usr/code/ \
		golang:1.10.0-alpine \
		go build -a -installsuffix netgo -o /usr/code/examples/ericchiang/example github.com/pulcy/kube-lock/examples/ericchiang
	docker build -t $(IMAGE) .
