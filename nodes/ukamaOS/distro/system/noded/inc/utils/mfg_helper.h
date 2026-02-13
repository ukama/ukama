/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2021-present, Ukama Inc.
 */

#ifndef LIB_UBSP_MFG_COMMON_MFG_HELPER_H_
#define LIB_UBSP_MFG_COMMON_MFG_HELPER_H_

#include "noded_macros.h"

#define UUID_MAX_LENGTH			UUID_LENGTH
#define NAME_MAX_LENGTH			NAME_LENGTH

/* TYPE OF UNITS */
#define MODULE_TYPE_MAX			6

/* Max board Count */
#define MAX_BOARDS				6

/**
 * @fn      int verify_board_name(char*)
 * @brief   Verifying board name.
 *
 * @param   name
 * @return  On Success, 0
 *          On failure, -1
 */
int verify_board_name(char* name);

/**
 * @fn      int verify_uuid(char*)
 * @brief   Verify the UUID by length
 *
 * @param   uuid
 * @return  On Success, 0
 *          On failure, -1
 */
int verify_uuid(char* uuid);

#endif /* LIB_UBSP_MFG_COMMON_MFG_HELPER_H_ */
