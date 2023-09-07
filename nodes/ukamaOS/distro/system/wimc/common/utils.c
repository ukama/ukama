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
#include <string.h>

#include "wimc.h"
#include "agent.h"

#include "utils.h"

char *convert_method_to_str(MethodType method) {

  char *str;

  switch(method) {

  case CHUNK:
    str = WIMC_METHOD_CHUNK_STR;
    break;

  case TEST:
    str = WIMC_METHOD_TEST_STR;
    break;

  default:
    str = "";
  }

  return str;
}

MethodType convert_str_to_method(char *str) {

  MethodType method;

  if (strcmp(str, WIMC_METHOD_CHUNK_STR)==0) {
    method = CHUNK;
  } else if (strcmp(str, WIMC_METHOD_TEST_STR)==0) {
    method = TEST;
  }

  return method;
}

char *convert_tx_state_to_str(TransferState state) {

  if (state == REQUEST) {
    return strdup(AGENT_TX_STATE_REQUEST_STR);
  } else if (state == FETCH) {
    return strdup(AGENT_TX_STATE_FETCH_STR);
  } else if (state == UNPACK) {
    return strdup(AGENT_TX_STATE_UNPACK_STR);
  } else if (state == DONE) {
    return strdup(AGENT_TX_STATE_DONE_STR);
  } else if (state == ERR) {
    return strdup(AGENT_TX_STATE_ERR_STR);
  } else {
    return strdup("");
  }
}

AgentState convert_str_to_state(char *str) {

  AgentState state;

  if (strcmp(str, AGENT_STATE_REGISTER_STR)==0) {
    state = REGISTER;
  } else if (strcmp(str, AGENT_STATE_ACTIVE_STR)==0) {
    state = ACTIVE;
  }  else if (strcmp(str, AGENT_STATE_UNREGISTER_STR)==0) {
    state = UNREGISTER;
  }  else if (strcmp(str, AGENT_STATE_INACTIVE_STR)==0) {
    state = INACTIVE;
  }

  return state;
}

char *convert_uuid_to_str(uuid_t uuid) {

  char *str;

  str = (char *)malloc(36+1); /* 36-byte string + trailing '\0' */
  uuid_unparse(uuid, str);

  return str;
}

TransferState convert_str_to_tx_state(char *state) {

  if (strcmp(state, AGENT_TX_STATE_REQUEST_STR)==0) {
    return REQUEST;
  } else if (strcmp(state, AGENT_TX_STATE_FETCH_STR)==0) {
    return FETCH;
  }  else if (strcmp(state, AGENT_TX_STATE_UNPACK_STR)==0) {
    return UNPACK;
  }  else if (strcmp(state, AGENT_TX_STATE_DONE_STR)==0) {
    return DONE;
  }  else if (strcmp(state, AGENT_TX_STATE_ERR_STR)==0) {
    return ERR;
  }
}
