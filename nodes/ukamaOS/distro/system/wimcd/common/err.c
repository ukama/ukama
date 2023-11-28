/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

/*
 * Error related utility functions.
 */

#include <stdio.h>

#include "err.h"

/*
 * error_to_str -- return string representation of the error code.
 *
 */

const char *error_to_str(int error) {

  switch (error) {
    
  case WIMC_OK:
    return WIMC_OK_STR;
    
  case WIMC_ERROR_EXIST:
    return WIMC_ERROR_EXIST_STR;

  case WIMC_ERROR_BAD_NAME:
    return WIMC_ERROR_BAD_NAME_STR;
    
  case WIMC_ERROR_BAD_ACTION:
    return WIMC_ERROR_BAD_ACTION_STR;

  case WIMC_ERROR_BAD_TYPE:
    return WIMC_ERROR_BAD_TYPE_STR;

  case WIMC_ERROR_BAD_METHOD:
    return WIMC_ERROR_BAD_METHOD_STR;
    
  case WIMC_ERROR_BAD_URL:
    return WIMC_ERROR_BAD_URL_STR;
    
  case WIMC_ERROR_BAD_ID:
    return WIMC_ERROR_BAD_ID_STR;
    
  case WIMC_ERROR_BAD_INTERVAL:
    return WIMC_ERROR_BAD_INTERVAL_STR;
    
  case WIMC_ERROR_MEMORY:
    return WIMC_ERROR_MEMORY_STR;

  case WIMC_ERROR_MISSING_CONTENT:
    return WIMC_ERROR_MISSING_CONTENT_STR;

  case WIMC_ERROR_MAX_AGENTS:
    return WIMC_ERROR_MAX_AGENTS_STR;
    
  default:
    return "";
  }

  return "";
}
