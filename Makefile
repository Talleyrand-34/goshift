# Go parameters
GOCMD=go
GOFILE=main.go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=goshift
BINARY_UNIX=$(BINARY_NAME)_unix
BTRFS_ROOT_1=/
BTRFS_ROOT_2=/mnt/btrfs-test-mnt

all: build test 

build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(GOFILE)  

test: setup-test testlistroot testlistmnt

setup-btrfs-testdeps:
	dd if=/dev/zero of=btrfs.mount bs=1M count=1000
	mkfs.btrfs btrfs.mount

setup-test:
	doas mkdir -p /mnt/btrfs-test-mnt
	-doas mount ./btrfs.mount /mnt/btrfs-test-mnt # || true


testlistroot:
	doas ./$(BINARY_NAME) subvolume list $(BTRFS_ROOT_1)

testlistmnt:
	doas ./$(BINARY_NAME) subvolume list $(BTRFS_ROOT_2)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOCMD) mod tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

docker-build:
	docker build -t $(BINARY_NAME):latest .

.PHONY: all build test clean run deps build-linux docker-build

