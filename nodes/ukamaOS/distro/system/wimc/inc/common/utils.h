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

char *convert_state_to_str(TransferState state);
AgentState convert_str_to_state(char *str);

/* defined in url.c */
int validate_url(char *url);
int valid_url_format(char *url);

#endif /* UTILS_H */
