/*
 * mfg_helper.h
 *
 *  Created on: Jun 25, 2021
 *      Author: vishal
 */

#ifndef LIB_UBSP_MFG_COMMON_MFG_HELPER_H_
#define LIB_UBSP_MFG_COMMON_MFG_HELPER_H_

#define UUID_MAX_LENGTH			24
#define NAME_MAX_LENGTH			24

/* TYPE OF UNITS */
#define MODULE_TYPE_MAX			5

/* Max board Count */
#define MAX_BOARDS				5

int verify_uuid(char* uuid);
int verify_boardname(char* name);


#endif /* LIB_UBSP_MFG_COMMON_MFG_HELPER_H_ */
