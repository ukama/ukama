/*
 * mfg_helper.h
 *
 *  Created on: Jun 25, 2021
 *      Author: vishal
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


int verify_board_name(char* name);

int verify_uuid(char* uuid);

#endif /* LIB_UBSP_MFG_COMMON_MFG_HELPER_H_ */
