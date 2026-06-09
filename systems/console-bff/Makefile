# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

.PHONY: all test test-unit clean build

build:
	pnpm install --frozen-lockfile && pnpm build
test:
	pnpm test --ci --coverage --watchAll=false
test-unit:
	pnpm test:unit
