BINDIR := bin

CMDS := $(filter-out %_test.go,$(notdir $(wildcard cmd/*)))

.PHONY: all
all: $(CMDS)

$(CMDS):
	go build -buildmode=pie -trimpath -o $(BINDIR)/$@ ./cmd/$@

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean
	rm -rf $(BINDIR)
