/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#pragma once

typedef enum {
    HttpStatus_OK                  = 200,
    HttpStatus_Accepted            = 202,
    HttpStatus_BadRequest          = 400,
    HttpStatus_Conflict            = 409,
    HttpStatus_InternalServerError = 500,
    HttpStatus_ServiceUnavailable  = 503
} HttpStatus;
