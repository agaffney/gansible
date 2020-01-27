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

PYTHON_DIR=python
PROTOBUF_DIR=$(PYTHON_DIR)/protobuf
PYHELPER_DIR=$(PYTHON_DIR)/pyhelper
PYHELPER_GRPC_DIR=$(PYHELPER_DIR)/gansible/grpc
GO_GRPC_DIR=$(PYTHON_DIR)/grpc

PROTOBUF_FILES=$(shell find $(PROTOBUF_DIR) -name '*.proto')
PROTOBUF_GO_FILES=$(patsubst $(PROTOBUF_DIR)/%,$(GO_GRPC_DIR)/%,$(patsubst %.proto,%.pb.go,$(PROTOBUF_FILES)))
PROTOBUF_PY_FILES=$(patsubst $(PROTOBUF_DIR)/%,$(PYHELPER_GRPC_DIR)/%,$(patsubst %.proto,%_pb2.py,$(PROTOBUF_FILES))) $(patsubst $(PROTOBUF_DIR)/%,$(PYHELPER_GRPC_DIR)/%,$(patsubst %.proto,%_pb2_grpc.py,$(PROTOBUF_FILES)))

# Paths for Python gRPC stub
PYHELPER_FILES=$(shell find $(PYHELPER_DIR) -name '*.py')
PYHELPER_ZIP=$(PYHELPER_DIR)/python_grpc.zip
PYHELPER_ZIP_GO_FILE=$(GO_GRPC_DIR)/python_grpc_zip.go

# Disable implicit rules
.SUFFIXES:

all: build

.PHONY: all build clean

build: $(BINARY)

$(BINARY): $(GO_FILES) $(PROTOBUF_GO_FILES) $(PYHELPER_ZIP_GO_FILE)
	go build -o $(BINARY)

clean:
	# Fix weird permissions in gopath/
	test -d $(TMPDIR) && chmod -R u+rw $(TMPDIR) || true
	rm -rf $(BINARY) $(TMPDIR) $(shell find $(PYHELPER_GRPC_DIR) -name '*.pyc') $(PYHELPER_ZIP) $(PYHELPER_ZIP_GO_FILE) $(PROTOBUF_PY_FILES) $(PROTOBUF_GO_FILES)

run: build
	./$(BINARY)

# Build generated gRPC files
.PHONY: grpc-protos grpc-protos-go grpc-protos-py

grpc-protos: grpc-protos-go grpc-protos-py

grpc-protos-go: $(PROTOBUF_GO_FILES)

# Generate Go protobuf files
$(PROTOBUF_GO_FILES): $(PROTOBUF_FILES) $(PROTOC) $(PROTOC_GEN_GO)
	$(PROTOC) -I $(PROTOBUF_DIR) --go_out=plugins=grpc:$(GO_GRPC_DIR) $(PROTOBUF_DIR)/$(patsubst %.pb.go,%.proto,$(notdir $@))

grpc-protos-py: $(PROTOBUF_PY_FILES)

# Generate Python protobuf files
$(PROTOBUF_PY_FILES): $(PROTOBUF_FILES) $(VENV_PATH)
	source $(VENV_PATH)/bin/activate ; \
		python -m grpc_tools.protoc -I $(PROTOBUF_DIR) --python_out=$(PYHELPER_GRPC_DIR) --grpc_python_out=$(PYHELPER_GRPC_DIR) $(PROTOBUF_DIR)/$(patsubst %_pb2.py,%.proto,$(notdir $@))

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
.PHONY: pyhelper-zip

pyhelper-zip: $(PYHELPER_ZIP)

$(PYHELPER_ZIP): $(PROTOBUF_PY_FILES) $(PYHELPER_FILES)
	cd $(PYHELPER_DIR) ; \
		$(PYTHON) -m zipfile -c $(notdir $(PYHELPER_ZIP)) __main__.py gansible

# Build Go file with Python gRPC stub ZIP embedded
$(PYHELPER_ZIP_GO_FILE): $(PYHELPER_ZIP)
	echo -e '// This file is generated. Do not edit!\n\npackage grpc\n\nconst pyGrpcZipContent = `' > $(PYHELPER_ZIP_GO_FILE)
	base64 $(PYHELPER_ZIP) >> $(PYHELPER_ZIP_GO_FILE)
	echo -e '`' >> $(PYHELPER_ZIP_GO_FILE)
