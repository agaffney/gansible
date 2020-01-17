BINARY=gansible

# Use bash to run shell commands
SHELL=/bin/bash

PYTHON=python

# Temp dir to use during builds
TMPDIR=$(CURDIR)/tmp

export GOPATH=$(TMPDIR)/gopath
export PATH=$(GOPATH)/bin:$(TMPDIR)/bin:$(shell echo $$PATH)

VENV_PATH=$(TMPDIR)/virtualenv

PROTOC=$(TMPDIR)/bin/protoc
PROTOC_GEN_GO=$(GOPATH)/bin/protoc-gen-go

GO_FILES=$(shell find $(CURDIR) -name '*.go' | grep -v '^$(TMPDIR)')

PROTOBUF_FILES=$(shell find python/grpc -name '*.proto')
PROTOBUF_GO_FILES=$(patsubst %.proto,%.pb.go,$(PROTOBUF_FILES))
PROTOBUF_PY_FILES=$(patsubst %.proto,%_pb2.py,$(PROTOBUF_FILES)) $(patsubst %.proto,%_pb2_grpc.py,$(PROTOBUF_FILES))

# Paths for Python gRPC stub
GRPC_PY_DIR=python/grpc
GRPC_PY_FILES=$(GRPC_PY_DIR)/__main__.py
GRPC_PY_ZIP=$(GRPC_PY_DIR)/python_grpc.zip
GRPC_PY_ZIP_GO_FILE=$(GRPC_PY_DIR)/python_grpc_zip.go

# Disable implicit rules
.SUFFIXES:

all: build

.PHONY: all build clean

build: $(BINARY)

$(BINARY): $(GO_FILES) $(PROTOBUF_GO_FILES) $(GRPC_PY_ZIP_GO_FILE)
	go build -o $(BINARY)

clean:
	rm -rf $(BINARY) $(PROTOC) $(PROTOC_GEN_GO) $(PROTOBUF_GO_FILES) $(TMPDIR) $(VENV_PATH) $(GRPC_PY_ZIP)

run: build
	./$(BINARY)

# Build generated gRPC files
.PHONY: grpc-protos grpc-protos-go grpc-protos-py

grpc-protos: grpc-protos-go grpc-protos-py

grpc-protos-go: $(PROTOBUF_GO_FILES)

# Generate Go protobuf files
%.pb.go: %.proto $(PROTOC) $(PROTOC_GEN_GO)
	cd $(dir $@); \
		$(PROTOC) -I . --go_out=plugins=grpc:. $(notdir $(patsubst %.pb.go,%.proto,$@))

grpc-protos-py: $(PROTOBUF_PY_FILES)

%_pb2.py: %.proto $(VENV_PATH)
	source $(VENV_PATH)/bin/activate ; \
		python -m grpc_tools.protoc -I $(dir $@) --python_out=$(dir $@) --grpc_python_out=$(dir $@) $(patsubst %_pb2.py,%.proto,$@)

# Build virtualenv used for generating Python gRPC files
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

# Build Python gRPC stub ZIP
.PHONY: grpc-py-zip

grpc-py-zip: $(GRPC_PY_ZIP)

$(GRPC_PY_ZIP): $(PROTOBUF_PY_FILES) $(GRPC_PY_FILES)
	cd $(dir $(GRPC_PY_ZIP)) ; \
		$(PYTHON) -m zipfile -c $(notdir $(GRPC_PY_ZIP)) $(notdir $(PROTOBUF_PY_FILES) $(GRPC_PY_FILES))

# Build Go file with Python gRPC stub ZIP embedded
$(GRPC_PY_ZIP_GO_FILE): $(GRPC_PY_ZIP)
	echo -e '// This file is generated. Do not edit!\n\npackage grpc\n\nconst pyGrpcZipContent = `' > $(GRPC_PY_ZIP_GO_FILE)
	base64 $(GRPC_PY_ZIP) >> $(GRPC_PY_ZIP_GO_FILE)
	echo -e '`' >> $(GRPC_PY_ZIP_GO_FILE)
