/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

// Moved to common/middleware/expressApp so the consolidated server can share
// it; this shim keeps the gateway import stable until the gateway is retired.
export { configureExpress } from "../common/middleware/expressApp";
