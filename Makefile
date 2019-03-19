.PHONY: all docker

all: gohost

docker: Dockerfile main.go
	docker build -t="lanrat/gohost" .

gohost: main.go
	CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ $^

deps: go.mod
	go mod download

fmt:
	gofmt -s -w -l .

clean: 
	rm gohost

