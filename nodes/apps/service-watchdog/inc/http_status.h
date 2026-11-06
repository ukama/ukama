/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef HTTP_STATUS_H_
#define HTTP_STATUS_H_

enum HttpStatusCode {
    HttpStatus_OK                  = 200,
    HttpStatus_Created             = 201,
    HttpStatus_Accepted            = 202,
    HttpStatus_NoContent           = 204,
    HttpStatus_BadRequest          = 400,
    HttpStatus_Forbidden           = 403,
    HttpStatus_NotFound            = 404,
    HttpStatus_MethodNotAllowed    = 405,
    HttpStatus_Conflict            = 409,
    HttpStatus_Gone                = 410,
    HttpStatus_InternalServerError = 500,
    HttpStatus_ServiceUnavailable  = 503,
};

#endif /* HTTP_STATUS_H_ */
