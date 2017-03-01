CLEAN_PATH := .

include ../../common.mk

GEN_PATH ?= .
CLEAN_PATH ?= .

.PHONY: generate
generate:
	$(GO) generate $(GEN_PATH)

.PHONY: clean
clean:
	git clean -f -d -X $(CLEAN_PATH)
