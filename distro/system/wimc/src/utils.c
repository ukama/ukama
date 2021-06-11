/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

/*
 * Misc utility functions
 */

#include <stdio.h>

#define TRUE 1
#define FALSE 0

/*
 * is_valid_url -- A valid URL is of format http://host:port/
 */

int is_valid_url(char *url) {

  if (url == NULL) {
    return FALSE;
  }

  /* XXX */
  
  return TRUE;
} 
