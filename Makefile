BIN="./bin"

## Build

.PHONY: build
build:
	@echo 'Building pctl...'
	go build -ldflags "-s" -o $(BIN)/pctl

## Quality control

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor:
	go mod vendor

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
