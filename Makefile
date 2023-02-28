LOCAL_BIN := $(CURDIR)/bin

build:
	go build -o $(LOCAL_BIN)/fts -tags fts5

.PHONY: \
	build