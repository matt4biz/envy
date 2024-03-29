version=$(shell git describe --tags --long --always --dirty 2>/dev/null)
SOURCES := $(wildcard internal/*.go cmd/*.go)

## NOTE: we can't use go install because it
## doesn't have the -o option to name the file

envy: envy.go $(SOURCES)
	go build -ldflags "-X main.version=$(version)" -o $@ ./cmd
	install -d $(GOPATH)/bin
	install $@ $(GOPATH)/bin

child: hack/main.go
	go build -o $@ ./hack

lint:
	golangci-lint run

test: envy.go $(SOURCES)
	go test -v ./... -coverprofile=c.out -covermode=count

demo: envy child
	envy add top a=b
	envy exec top ./child

clean:
	rm -rf envy child

