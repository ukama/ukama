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
    HttpStatus_Unauthorized        = 401,
    HttpStatus_Forbidden           = 403,
    HttpStatus_NotFound            = 404,
    HttpStatus_MethodNotAllowed    = 405,
    HttpStatus_Conflict            = 409,
    HttpStatus_UnprocessableEntity = 422,
    HttpStatus_InternalServerError = 500,
    HttpStatus_NotImplemented      = 501,
    HttpStatus_ServiceUnavailable  = 503,
};

static inline const char *HttpStatusStr(int code) {
    switch (code) {
    case HttpStatus_OK:
        return "OK";
    case HttpStatus_Created:
        return "Created";
    case HttpStatus_Accepted:
        return "Accepted";
    case HttpStatus_NoContent:
        return "No Content";
    case HttpStatus_BadRequest:
        return "Bad Request";
    case HttpStatus_Unauthorized:
        return "Unauthorized";
    case HttpStatus_Forbidden:
        return "Forbidden";
    case HttpStatus_NotFound:
        return "Not Found";
    case HttpStatus_MethodNotAllowed:
        return "Method Not Allowed";
    case HttpStatus_Conflict:
        return "Conflict";
    case HttpStatus_UnprocessableEntity:
        return "Unprocessable Entity";
    case HttpStatus_InternalServerError:
        return "Internal Server Error";
    case HttpStatus_NotImplemented:
        return "Not Implemented";
    case HttpStatus_ServiceUnavailable:
        return "Service Unavailable";
    default:
        return "";
    }
}

#endif /* HTTP_STATUS_H_ */
