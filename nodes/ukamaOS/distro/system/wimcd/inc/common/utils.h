/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef UTILS_H
#define UTILS_H

char *convert_state_to_str(TransferState state);
AgentState convert_str_to_state(char *str);

/* defined in url.c */
int validate_url(char *url);
int valid_url_format(char *url);

#endif /* UTILS_H */
