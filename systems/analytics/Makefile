# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2026-present, Ukama Inc.

CURMAKE := $(abspath $(lastword $(MAKEFILE_LIST)))
CURPATH := $(dir $(CURMAKE))
SYSTEMDIR := $(CURPATH)

GOPATH ?= $(shell go env GOPATH)

SYSSUBDIRS = $(shell find ./*/ -maxdepth 1 -type f -name 'Makefile' -printf '%h' | awk -F './' '{ for (i=1; i<=NF; i++) printf $$i" " }')

PB_SUBDIRS = business customer network collector

.PHONY: subdirs $(SYSSUBDIRS) gen $(PB_SUBDIRS:%=gen-%) clean info

subdirs: $(SYSSUBDIRS)

$(SYSSUBDIRS):
	$(MAKE) -C $@

gen: check-gen-env $(PB_SUBDIRS:%=gen-%)

check-gen-env:
	@test -n "$(GOPATH)" || (echo "GOPATH is not set and go env GOPATH returned empty"; exit 1)
	@mkdir -p "$(GOPATH)/src"
	@command -v protoc >/dev/null 2>&1 || (echo "missing protoc"; exit 1)
	@command -v protoc-gen-go >/dev/null 2>&1 || (echo "missing protoc-gen-go"; exit 1)
	@command -v protoc-gen-go-grpc >/dev/null 2>&1 || (echo "missing protoc-gen-go-grpc"; exit 1)
	@command -v protoc-gen-govalidators >/dev/null 2>&1 || (echo "missing protoc-gen-govalidators"; exit 1)
	@command -v mockery >/dev/null 2>&1 || (echo "missing mockery"; exit 1)

$(PB_SUBDIRS:%=gen-%):
	$(MAKE) -C $(@:gen-%=%) gen GOPATH="$(GOPATH)"

clean:
	@for dir in $(SYSSUBDIRS); do \
		$(MAKE) -C $$dir clean; \
	done

info:
	@echo $(SYSSUBDIRS)
