.DEFAULT: all
.PHONY: all

GENDIR= aurora/gen

GOPATH := ${shell pwd}/.gopath:${shell pwd}/vendor
export GOPATH

all: gru

gru: $(GENDIR)/api/*.go main.go aurora
	go build -o $@ ./$(@D)

aurora: aurora/*.go

aurora/gen/api/*.go: aurora/api.thrift | $(GENDIR)
	thrift -out $(GENDIR) \
	--gen go:package_prefix=github.com/jc-m/gru/$(GENDIR)/ \
	aurora/api.thrift

$(GENDIR):
	mkdir -p $(GENDIR)
