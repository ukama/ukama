/**
 * Copyright (c) 2021-present, Ukama Inc.
 * All rights reserved.
 *
 * This source code is licensed under the XXX-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#ifndef LIB_UBSP_MFG_COMMON_MFG_HELPER_H_
#define LIB_UBSP_MFG_COMMON_MFG_HELPER_H_

#include "noded_macros.h"

#define UUID_MAX_LENGTH			UUID_LENGTH
#define NAME_MAX_LENGTH			NAME_LENGTH

/* TYPE OF UNITS */
#define MODULE_TYPE_MAX			5

/* Max board Count */
#define MAX_BOARDS				5

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
