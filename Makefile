.PHONY: build
build:
	go build code-gen.go

.PHONY: gen-node
gen-node:
	rm -rf ./test && mkdir test && cd test && go run ../code-gen.go ../req.yaml node init && cd ..

.PHONY: test
test:
	go test -v ./...
	