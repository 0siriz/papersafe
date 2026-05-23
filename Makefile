BINDIR := bin

CMDS := $(filter-out %_test.go,$(notdir $(wildcard cmd/*)))
BINS := $(addprefix $(BINDIR)/,$(CMDS))

.PHONY: all
all: $(BINS)

$(BINDIR)/%: cmd/%/
	go build -buildmode=pie -trimpath -o $@ ./$<

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
