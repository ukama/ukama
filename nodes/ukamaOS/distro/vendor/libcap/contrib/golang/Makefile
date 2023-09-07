GO := go

all: posix.go posix-cgo.go
	CGO_ENABLED=0 $(GO) build posix.go
	$(GO) build posix-cgo.go

posix.go: ptest.go posix.stub_go Makefile
	sed -e '/\/\/ main_here/ r posix.stub_go' ptest.go > $@

posix-cgo.go: ptest.go posix-cgo.stub_go Makefile
	sed -e '/\/\/ main_here/ r posix-cgo.stub_go' ptest.go > $@

clean:
	rm -f posix.go posix
	rm -f posix-cgo.go posix-cgo
