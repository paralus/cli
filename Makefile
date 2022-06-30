BIN="./bin"

.PHONY: tidy
tidy:
	GOPRIVATE=github.com/paralus/* go mod tidy

.PHONY: vendor
vendor:
  GOPRIVATE=github.com/paralus/* go mod vendor

.PHONY: build
build:
	# Omit the symbol table and debug information to reduce the
	# size of binary.
	go build -ldflags "-s" -o $(BIN)/pctl

.PHONY: build-proto
build-proto:
  buf build

.PHONY: gen-proto
gen-proto:
	buf generate

.PHONY: test
test:
	go test ./...

.PHONY: check
check:
	go fmt ./...
	go vet ./...

	$(MAKE) tidy

.PHONY: clean
clean:
	rm -rf $(BIN)
