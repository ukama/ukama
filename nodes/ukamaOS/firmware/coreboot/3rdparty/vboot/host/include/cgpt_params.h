/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#ifndef VBOOT_REFERENCE_CGPT_PARAMS_H_
#define VBOOT_REFERENCE_CGPT_PARAMS_H_

#include <stdint.h>

#include "gpt.h"

#ifdef __cplusplus
extern "C" {
#endif  /* __cplusplus */

enum {
	CGPT_OK = 0,
	CGPT_FAILED,
};

typedef struct CgptCreateParams {
	const char *drive_name;
	uint64_t drive_size;
	int zap;
	uint64_t padding;
} CgptCreateParams;

typedef struct CgptAddParams {
	const char *drive_name;
	uint64_t drive_size;
	uint32_t partition;
	uint64_t begin;
	uint64_t size;
	Guid type_guid;
	Guid unique_guid;
	const char *label;
	int successful;
	int tries;
	int priority;
	int required;
	int legacy_boot;
	uint32_t raw_value;
	int set_begin;
	int set_size;
	int set_type;
	int set_unique;
	int set_successful;
	int set_tries;
	int set_priority;
	int set_required;
	int set_legacy_boot;
	int set_raw;
} CgptAddParams;

typedef struct CgptEditParams {
	const char *drive_name;
	uint64_t drive_size;
	Guid unique_guid;
	int set_unique;
} CgptEditParams;

typedef struct CgptShowParams {
	const char *drive_name;
	uint64_t drive_size;
	int numeric;
	int verbose;
	int quick;
	uint32_t partition;
	int single_item;
	int debug;
	int num_partitions;
} CgptShowParams;

typedef struct CgptRepairParams {
	const char *drive_name;
	uint64_t drive_size;
	int verbose;
} CgptRepairParams;

typedef struct CgptBootParams {
	const char *drive_name;
	uint64_t drive_size;
	uint32_t partition;
	const char *bootfile;
	int create_pmbr;
} CgptBootParams;

typedef struct CgptPrioritizeParams {
	const char *drive_name;
	uint64_t drive_size;
	uint32_t set_partition;
	int set_friends;
	int max_priority;
	int orig_priority;
} CgptPrioritizeParams;

struct CgptFindParams;
typedef void (*CgptFindShowFn)(struct CgptFindParams *params,
			       const char *filename, int partnum,
			       GptEntry *entry);
typedef struct CgptFindParams {
	const char *drive_name;
	uint64_t drive_size;
	int verbose;
	int set_unique;
	int set_type;
	int set_label;
	int oneonly;
	int numeric;
	uint8_t *matchbuf;
	uint64_t matchlen;
	uint64_t matchoffset;
	uint8_t *comparebuf;
	Guid unique_guid;
	Guid type_guid;
	const char *label;
	int hits;
	int match_partnum;           /* 1-based; 0 means no match */
	/* when working with MTD, we actually work on a temp file, but we still
	 * need to print the device name. so this parameter is here to properly
	 * show the correct device name in that special case. */
	CgptFindShowFn show_fn;
} CgptFindParams;

enum {
	CGPT_LEGACY_MODE_LEGACY = 0,
	CGPT_LEGACY_MODE_EFIPART,
	CGPT_LEGACY_MODE_IGNORE_PRIMARY,
};

typedef struct CgptLegacyParams {
	const char *drive_name;
	uint64_t drive_size;
	int mode;
} CgptLegacyParams;

#ifdef __cplusplus
}
#endif  /* __cplusplus */

#endif  /* VBOOT_REFERENCE_CGPT_PARAMS_H_ */
