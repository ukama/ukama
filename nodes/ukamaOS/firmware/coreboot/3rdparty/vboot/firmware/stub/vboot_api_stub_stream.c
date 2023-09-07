/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Stub implementations of stream APIs.
 */

#include <stdint.h>

#include "2common.h"
#include "vboot_api.h"

/* The stub implementation assumes 512-byte disk sectors */
#define LBA_BYTES 512

/* Internal struct to simulate a stream for sector-based disks */
struct disk_stream {
	/* Disk handle */
	VbExDiskHandle_t handle;

	/* Next sector to read */
	uint64_t sector;

	/* Number of sectors left in partition */
	uint64_t sectors_left;
};

vb2_error_t VbExStreamOpen(VbExDiskHandle_t handle, uint64_t lba_start,
			   uint64_t lba_count, VbExStream_t *stream)
{
	struct disk_stream *s;

	if (!handle) {
		*stream = NULL;
		return VB2_ERROR_UNKNOWN;
	}

	s = malloc(sizeof(*s));
	s->handle = handle;
	s->sector = lba_start;
	s->sectors_left = lba_count;

	*stream = (void *)s;

	return VB2_SUCCESS;
}

vb2_error_t VbExStreamRead(VbExStream_t stream, uint32_t bytes, void *buffer)
{
	struct disk_stream *s = (struct disk_stream *)stream;
	uint64_t sectors;
	vb2_error_t rv;

	if (!s)
		return VB2_ERROR_UNKNOWN;

	/* For now, require reads to be a multiple of the LBA size */
	if (bytes % LBA_BYTES)
		return VB2_ERROR_UNKNOWN;

	/* Fail on overflow */
	sectors = bytes / LBA_BYTES;
	if (sectors > s->sectors_left)
		return VB2_ERROR_UNKNOWN;

	rv = VbExDiskRead(s->handle, s->sector, sectors, buffer);
	if (rv != VB2_SUCCESS)
		return rv;

	s->sector += sectors;
	s->sectors_left -= sectors;

	return VB2_SUCCESS;
}

void VbExStreamClose(VbExStream_t stream)
{
	struct disk_stream *s = (struct disk_stream *)stream;

	/* Allow freeing a null pointer */
	if (!s)
		return;

	free(s);
	return;
}
