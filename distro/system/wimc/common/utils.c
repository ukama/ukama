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

/* Some utility functions for serdes */
#if 0
char *convert_action_to_str(ActionType action) {

  char *str;

  switch (action) {
  case ACTION_FETCH:
    str = WIMC_ACTION_FETCH_STR;
    break;

  case ACTION_UPDATE:
    str = WIMC_ACTION_UPDATE_STR;
    break;

  case ACTION_CANCEL:
    str = WIMC_ACTION_CANCEL_STR;
    break;

  default:
    str = "";
  }

  return str;
}
#endif

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

char *convert_state_to_str(TransferState state) {

  char *str;

  switch(state) {

  case REQUEST:
    str = AGENT_TX_STATE_REQUEST_STR;
    break;

  case FETCH:
    str = AGENT_TX_STATE_FETCH_STR;
    break;

  case UNPACK:
    str = AGENT_TX_STATE_UNPACK_STR;
    break;

  case DONE:
    str = AGENT_TX_STATE_DONE_STR;
    break;

  case ERR:
    str = AGENT_TX_STATE_ERR_STR;
    break;

  default:
    str = "";
  }

  return str;
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

ReqType convert_str_to_type(char *str) {

  ReqType req;

  if (strcmp(str, AGENT_REQ_TYPE_REG_STR)==0) {
    req = REQ_REG;
  } else if (strcmp(str, AGENT_REQ_TYPE_UNREG_STR)==0) {
    req = REQ_UNREG;
  } else if (strcmp(str, AGENT_REQ_TYPE_UPDATE_STR)==0) {
    req = REQ_UPDATE;
  }

  return req;
}

WReqType convert_str_to_wType(char *str) {

  WReqType req;
#if 0
  if (strcmp(str, WIMC_REQ_TYPE_AGENT_STR)==0) {
    req = AGENT;
  } else if (strcmp(str, WIMC_REQ_TYPE_PROVIDER_STR)==0) {
    req = PROVIDER;
  }
#endif
  return req;
}

#if 0
ActionType convert_str_to_action(char *str) {

  ActionType action;

  if (strcmp(str, WIMC_ACTION_FETCH_STR)==0) {
    action = ACTION_FETCH;
  } else if (strcmp(str, WIMC_ACTION_UPDATE_STR)==0) {
    action = ACTION_UPDATE;
  } else if (strcmp(str, WIMC_ACTION_CANCEL_STR)==0) {
    action = ACTION_CANCEL;
  }

  return action;
}
#endif
