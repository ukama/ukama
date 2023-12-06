# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

gen:
	protoc -I ./pb/dep -I ./pb --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src pb/ukama/*.proto --go-grpc_opt=require_unimplemented_servers=false
	protoc --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src pb/*.proto --go-grpc_opt=require_unimplemented_servers=false
	protoc -I ./pb/dep -I ./pb  --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
 			   pb/events/*.proto
	protoc -I ./pb/dep -I ./pb  --go_out=$(GOPATH)/src --go-grpc_out=$(GOPATH)/src --go-grpc_opt=require_unimplemented_servers=true --govalidators_out=$(GOPATH)/src \
 			   pb/msgclient/*.proto
			   
	mockery --all --recursive --dir ./


clean:
	rm mocks/*.go

test:
	go test -v ./...

# Go lint
lint:
	golangci-lint run

