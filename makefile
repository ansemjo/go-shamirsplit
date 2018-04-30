# targets with no actual files
.PHONY : install test clean

# names
BINARY := shamirsplit
PREFIX := /usr/local

# find protobuf definitions that need compilation
PROTOS := $(shell find -type f -name '*.proto' | sed 's/\.proto$$/.pb.go/')

# compile binary
$(BINARY) : sharding/protobuf.pb.go
	go build -o $@ cmd/shamirsplit/main.go

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

# clean anything not tracked by git
clean :
	git clean -fx
