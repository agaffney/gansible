BINARY=gansible

SHELL=/bin/bash

TMPDIR=$(CURDIR)/tmp

export GOPATH=$(TMPDIR)/gopath
export PATH=$(GOPATH)/bin:$(TMPDIR)/bin:$(shell echo $$PATH)

VENV_PATH=.virtualenv

PROTOC=$(TMPDIR)/bin/protoc
PROTOC_GEN_GO=$(GOPATH)/bin/protoc-gen-go

GO_FILES=$(shell find $(CURDIR) -name '*.go' | grep -v '^$(TMPDIR)')

PROTOBUF_FILES=$(shell find python/grpc -name '*.proto')
PROTOBUF_GO_FILES=$(patsubst %.proto,%.pb.go,$(PROTOBUF_FILES))
PROTOBUF_PY_FILES=$(patsubst %.proto,%_pb2.py,$(PROTOBUF_FILES))

# Disable implicit rules
.SUFFIXES:

all: build

.PHONY: all build clean

build: $(BINARY)

$(BINARY): $(GO_FILES) $(PROTOBUF_GO_FILES)
	go build -o $(BINARY)

clean:
	rm -rf $(BINARY) $(PROTOC) $(PROTOC_GEN_GO) $(PROTOBUF_GO_FILES) $(TMPDIR) $(VENV_PATH)

run: build
	./$(BINARY)

.PHONY: grpc-protos grpc-protos-go grpc-protos-py

grpc-protos: grpc-protos-go grpc-protos-py

grpc-protos-go: $(PROTOBUF_GO_FILES)

# Generate Go protobuf files
%.pb.go: %.proto $(PROTOC) $(PROTOC_GEN_GO)
	$(PROTOC) --go_out=plugins=grpc:. $(patsubst %.pb.go,%.proto,$@)

grpc-protos-py: $(PROTOBUF_PY_FILES)

%_pb2.py: %.proto $(VENV_PATH)
	source $(VENV_PATH)/bin/activate ; \
		python -m grpc_tools.protoc -I $(dir $@) --python_out=$(dir $@) --grpc_python_out=$(dir $@) $(patsubst %_pb2.py,%.proto,$@)

$(VENV_PATH):
	virtualenv $(VENV_PATH)
	touch $(VENV_PATH)/requirements.txt
	echo 'grpcio-tools' >> $(VENV_PATH)/requirements.txt
	# Fix breakage on py2
	echo 'setuptools<45 ; python_version < "3.0"' >> $(VENV_PATH)/requirements.txt
	source $(VENV_PATH)/bin/activate ; \
		pip install -r $(VENV_PATH)/requirements.txt

# Download/unzip protoc binary
$(PROTOC):
	mkdir -p $(TMPDIR)
	curl -L -o $(TMPDIR)/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.11.2/protoc-3.11.2-linux-x86_64.zip
	unzip -d $(TMPDIR) $(TMPDIR)/protoc.zip bin/protoc

# Download/build protoc-gen-go binary
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go
