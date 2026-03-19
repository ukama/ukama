/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

#ifndef HTTP_STATUS_H
#define HTTP_STATUS_H

/* HTTP Status codes */
#define HttpStatus_OK                    200
#define HttpStatus_Created               201
#define HttpStatus_Accepted              202
#define HttpStatus_NoContent             204

#define HttpStatus_BadRequest            400
#define HttpStatus_Unauthorized          401
#define HttpStatus_Forbidden             403
#define HttpStatus_NotFound              404
#define HttpStatus_MethodNotAllowed      405
#define HttpStatus_Conflict              409
#define HttpStatus_UnprocessableEntity   422

#define HttpStatus_InternalServerError   500
#define HttpStatus_NotImplemented        501
#define HttpStatus_ServiceUnavailable    503

/* Get string description for status code */
static inline const char *HttpStatusStr(int status) {
    switch (status) {
    case HttpStatus_OK:                  return "OK";
    case HttpStatus_Created:             return "Created";
    case HttpStatus_Accepted:            return "Accepted";
    case HttpStatus_NoContent:           return "No Content";
    case HttpStatus_BadRequest:          return "Bad Request";
    case HttpStatus_Unauthorized:        return "Unauthorized";
    case HttpStatus_Forbidden:           return "Forbidden";
    case HttpStatus_NotFound:            return "Not Found";
    case HttpStatus_MethodNotAllowed:    return "Method Not Allowed";
    case HttpStatus_Conflict:            return "Conflict";
    case HttpStatus_UnprocessableEntity: return "Unprocessable Entity";
    case HttpStatus_InternalServerError: return "Internal Server Error";
    case HttpStatus_NotImplemented:      return "Not Implemented";
    case HttpStatus_ServiceUnavailable:  return "Service Unavailable";
    default:                             return "Unknown";
    }
}

#endif /* HTTP_STATUS_H */
