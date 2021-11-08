all: update
.PHONY: all

update: clean build

OUT_DIR=bin
build:
	mkdir -p "${OUT_DIR}"
	go build -o "${OUT_DIR}/cpa" main.go

clean:
	$(RM) ./bin/cpa
.PHONY: clean