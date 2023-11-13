# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at https://mozilla.org/MPL/2.0/.
#
# Copyright (c) 2023-present, Ukama Inc.

.PHONY: all test clean build

build:
	yarn install  --frozen-lockfile && yarn build
test:
	yarn test --ci  --coverage --watchAll=false
