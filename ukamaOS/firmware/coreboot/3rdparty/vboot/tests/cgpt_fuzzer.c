// Copyright 2019 The Chromium OS Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#include <stddef.h>
#include <stdint.h>

#include "cgptlib.h"
#include "gpt.h"

struct MockDisk {
	size_t sector_shift;
	const uint8_t* data;
	size_t size;
};

// GPT disk parameters provided by the fuzzer test case. See GptData type
// definition for details.
struct GptDataParams {
	uint32_t sector_shift;
	uint32_t flags;
	uint64_t streaming_drive_sectors;
	uint64_t gpt_drive_sectors;
} __attribute__((packed));

static struct MockDisk mock_disk;

vb2_error_t VbExDiskRead(VbExDiskHandle_t h, uint64_t lba_start,
			 uint64_t lba_count, void *buffer)
{
	size_t lba_size = mock_disk.size >> mock_disk.sector_shift;
	if (lba_start > lba_size || lba_size - lba_start < lba_count) {
		return VB2_ERROR_UNKNOWN;
	}

	size_t start = lba_start << mock_disk.sector_shift;
	size_t size = lba_count << mock_disk.sector_shift;

	memcpy(buffer, &mock_disk.data[start], size);
	return VB2_SUCCESS;
}

int LLVMFuzzerTestOneInput(const uint8_t* data, size_t size);

int LLVMFuzzerTestOneInput(const uint8_t* data, size_t size) {
	struct GptDataParams params;
	if (size < sizeof(params)) {
		return 0;
	}
	memcpy(&params, data, sizeof(params));

	// Enforce a sane sector size. The sector size must accommodate the GPT
	// header (the code assumes this) and large values don't make sense
	// either (both in terms of actual hardware parameters and ability for
	// the fuzzer to deal with effectively).
	if (params.sector_shift < 9) {
		params.sector_shift = 9;  // 512 byte sectors min.
	}
	if (params.sector_shift > 12) {
		params.sector_shift = 12;  // 4K sectors max.
	}

	mock_disk.sector_shift = params.sector_shift;
	mock_disk.data = data + sizeof(params);
	mock_disk.size = size - sizeof(params);

	GptData gpt;
	memset(&gpt, 0, sizeof(gpt));
	gpt.sector_bytes = 1ULL << params.sector_shift;
	gpt.streaming_drive_sectors = params.streaming_drive_sectors;
	gpt.gpt_drive_sectors = params.gpt_drive_sectors;
	gpt.flags = params.flags;

	if (0 == AllocAndReadGptData(0, &gpt)) {
		int result = GptInit(&gpt);
		while (GPT_SUCCESS == result) {
			uint64_t part_start, part_size;
			result = GptNextKernelEntry(&gpt, &part_start,
						    &part_size);
		}
	}

	WriteAndFreeGptData(0, &gpt);

	return 0;
}
