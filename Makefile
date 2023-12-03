.PHONY: build
build:
	go build code-gen.go

# usage cmd={http|services|storage} make gen-node-command
.PHONY: gen-node-command
gen-node:
	mkdir -p test && cd test && go run ../code-gen.go ../req.yaml node $(cmd) && cd ..

.PHONY: test
test:
	go test -v ./...
	