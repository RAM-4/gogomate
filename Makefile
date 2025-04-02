.PHONY: build clean

build:
	go build -o gogomate ./cmd/gogomate

clean:
	rm -f gogomate

clean-all: clean
	rm -rf ~/.gogomate/data/
