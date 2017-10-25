.PHONY: all

all: gohost

docker: Dockerfile
	docker build -t="lanrat/gohost" .

gohost: main.go
	CGO_ENABLED=0 go build -a -installsuffix cgo -o $@ $^

