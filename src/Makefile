BINARY=gansible

.PHONY: all build clean

all: build

build: $(BINARY)

$(BINARY): $(shell find . -name '*.go')
	go build -o $(BINARY)

clean:
	rm -f $(BINARY)
