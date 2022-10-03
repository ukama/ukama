# Dependencies

In order to work on cloud services you will need below tools:
- go 1.18
- [Protoc](https://grpc.io/docs/protoc-installation/)
- [Protoc_gen_go](https://grpc.io/docs/languages/go/quickstart/)
- [Go-proto-validators](https://github.com/mwitkow/go-proto-validators)
- [Mockery](https://github.com/vektra/mockery)

Make sure that go/bin directory is part of the PATH:
```
export PATH=${PATH}:${GOPATH}/bin
```

You can install all tools by running `make go-install` in the services directory
