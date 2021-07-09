/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef UTILS_H
#define UTILS_H

//char *convert_action_to_str(ActionType action);
char *convert_method_to_str(MethodType method);
MethodType convert_str_to_method(char *str);
char *convert_state_to_str(TransferState state);
char *convert_type_to_str(WReqType type);
AgentState convert_str_to_state(char *str);
ReqType convert_str_to_type(char *str);
WReqType convert_str_to_wType(char *str);
//ActionType convert_str_to_action(char *str);

#endif /* UTILS_H */
