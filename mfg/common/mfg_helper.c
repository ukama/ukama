/*
 * common.c
 *
 *  Created on: Jun 25, 2021
 *      Author: vishal
 */

#include "utils/mfg_helper.h"

#include "usys_log.h"
#include "usys_string.h"

/* Board Name */
static char* boardname[MODULE_TYPE_MAX] = {
    "ComV1",
    "LTE",
    "MASK",
    "RF CTRL BOARD",
    "RF BOARD"
};

/* Verify UUID */
int verify_uuid(char* uuid) {
  int ret = 0;
  unsigned int len = strlen(uuid);
  log_trace("MFG:: UUID %s Length %d \n", uuid, len);
  if ((!len) || (len > UUID_MAX_LENGTH)) {
    usys_log_error("MFG:: Error:: UUID length should be greater than 0 and less than 24 characters.");
    ret = -1;
  }
  return ret;
}

/* Verify Board name */
int verify_board_name(char* name) {
  int ret = -1;
  unsigned int len = strlen(name);
  log_trace("MFG:: Name %s Length %d \n", name, len);
  if ((!len) || (len > UUID_MAX_LENGTH)) {
    usys_log_error("MFG:: Error:: Name length should be greater than 0 and less than 24 characters.");
  } else {

    /* Make sure board name is proper one */
    for(unsigned short int idx = 0; idx < MODULE_TYPE_MAX; idx++) {
      if (!strcmp(name, boardname[idx])) {
        ret = 0;
        break;
      }
    }
    if (ret) {
      usys_log_error("MFG:: Error:: Check the module name %s.", name);
    }

  }
  return ret;
}


