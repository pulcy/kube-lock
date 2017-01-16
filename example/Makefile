ROOTDIR := $(shell cd .. && pwd)
IMAGE := pulcy/kube-lock-example

all:
	docker run \
		--rm \
		-v $(ROOTDIR):/usr/code \
		-e GOPATH=/usr/code/.gobuild \
		-e GOOS=linux \
		-e GOARCH=amd64 \
		-e CGO_ENABLED=0 \
		-w /usr/code/ \
		golang:1.7.4-alpine \
		go build -a -installsuffix netgo -o /usr/code/example/example github.com/pulcy/kube-lock/example
	docker build -t $(IMAGE) .
