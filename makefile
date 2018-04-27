# targets with no actual files
.PHONY : install test

# names
BINARY := shamirsplit
PREFIX := /usr/local

# find protobuf definitions that need compilation
PROTOS := $(shell find -type f -name '*.proto' | sed 's/\.proto$$/.pb.go/')

# compile binary
$(BINARY) : src/sharding/protobuf.pb.go
	go build -o $@ src/cmd/shamirsplit/main.go

# install in $PATH
install : $(BINARY)
	install -m 755 -o root -g root $(BINARY) $(PREFIX)/bin/

# compile protobuf files
%.pb.go :
	@echo $(PROTOS)
	protoc --go_out=. "$*.proto"

# run go tests TODO
test :
	go test
