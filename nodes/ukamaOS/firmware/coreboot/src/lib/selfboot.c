/*
 * This file is part of the coreboot project.
 *
 * Copyright (C) 2003 Eric W. Biederman <ebiederm@xmission.com>
 * Copyright (C) 2009 Ron Minnich <rminnich@gmail.com>
 * Copyright (C) 2016 George Trudeau <george.trudeau@usherbrooke.ca>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <commonlib/compression.h>
#include <commonlib/endian.h>
#include <console/console.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>
#include <symbols.h>
#include <cbfs.h>
#include <lib.h>
#include <bootmem.h>
#include <program_loading.h>
#include <timestamp.h>
#include <cbmem.h>

/* The type syntax for C is essentially unparsable. -- Rob Pike */
typedef int (*checker_t)(struct cbfs_payload_segment *cbfssegs, void *args);

/* Decode a serialized cbfs payload segment
 * from memory into native endianness.
 */
static void cbfs_decode_payload_segment(struct cbfs_payload_segment *segment,
		const struct cbfs_payload_segment *src)
{
	segment->type        = read_be32(&src->type);
	segment->compression = read_be32(&src->compression);
	segment->offset      = read_be32(&src->offset);
	segment->load_addr   = read_be64(&src->load_addr);
	segment->len         = read_be32(&src->len);
	segment->mem_len     = read_be32(&src->mem_len);
}

static int segment_targets_type(void *dest, unsigned long memsz,
		enum bootmem_type dest_type)
{
	uintptr_t d = (uintptr_t) dest;
	if (bootmem_region_targets_type(d, memsz, dest_type))
		return 1;

	if (payload_arch_usable_ram_quirk(d, memsz))
		return 1;

	printk(BIOS_ERR, "SELF segment doesn't target RAM: 0x%p, %lu bytes\n", dest, memsz);
	bootmem_dump_ranges();
	return 0;
}

static int load_one_segment(uint8_t *dest,
			    uint8_t *src,
			    size_t len,
			    size_t memsz,
			    uint32_t compression,
			    int flags)
{
		unsigned char *middle, *end;
		printk(BIOS_DEBUG, "Loading Segment: addr: 0x%p memsz: 0x%016zx filesz: 0x%016zx\n",
		       dest, memsz, len);

		/* Compute the boundaries of the segment */
		end = dest + memsz;

		/* Copy data from the initial buffer */
		switch (compression) {
		case CBFS_COMPRESS_LZMA: {
			printk(BIOS_DEBUG, "using LZMA\n");
			timestamp_add_now(TS_START_ULZMA);
			len = ulzman(src, len, dest, memsz);
			timestamp_add_now(TS_END_ULZMA);
			if (!len) /* Decompression Error. */
				return 0;
			break;
		}
		case CBFS_COMPRESS_LZ4: {
			printk(BIOS_DEBUG, "using LZ4\n");
			timestamp_add_now(TS_START_ULZ4F);
			len = ulz4fn(src, len, dest, memsz);
			timestamp_add_now(TS_END_ULZ4F);
			if (!len) /* Decompression Error. */
				return 0;
			break;
		}
		case CBFS_COMPRESS_NONE: {
			printk(BIOS_DEBUG, "it's not compressed!\n");
			memcpy(dest, src, len);
			break;
		}
		default:
			printk(BIOS_INFO,  "CBFS:  Unknown compression type %d\n", compression);
			return 0;
		}
		/* Calculate middle after any changes to len. */
		middle = dest + len;
		printk(BIOS_SPEW, "[ 0x%08lx, %08lx, 0x%08lx) <- %08lx\n",
			(unsigned long)dest,
			(unsigned long)middle,
			(unsigned long)end,
			(unsigned long)src);

		/* Zero the extra bytes between middle & end */
		if (middle < end) {
			printk(BIOS_DEBUG,
				"Clearing Segment: addr: 0x%016lx memsz: 0x%016lx\n",
				(unsigned long)middle,
				(unsigned long)(end - middle));

			/* Zero the extra bytes */
			memset(middle, 0, end - middle);
		}

		/*
		 * Each architecture can perform additional operations
		 * on the loaded segment
		 */
		prog_segment_loaded((uintptr_t)dest, memsz, flags);


	return 1;
}

/* Note: this function is a bit dangerous so is not exported.
 * It assumes you're smart enough not to call it with the very
 * last segment, since it uses seg + 1 */
static int last_loadable_segment(struct cbfs_payload_segment *seg)
{
	return read_be32(&(seg + 1)->type) == PAYLOAD_SEGMENT_ENTRY;
}

static int check_payload_segments(struct cbfs_payload_segment *cbfssegs,
		void *args)
{
	uint8_t *dest;
	size_t memsz;
	struct cbfs_payload_segment *seg, segment;
	enum bootmem_type dest_type = *(enum bootmem_type *)args;

	for (seg = cbfssegs;; ++seg) {
		printk(BIOS_DEBUG, "Checking segment from ROM address 0x%p\n", seg);
		cbfs_decode_payload_segment(&segment, seg);
		dest = (uint8_t *)(uintptr_t)segment.load_addr;
		memsz = segment.mem_len;
		if (segment.type == PAYLOAD_SEGMENT_ENTRY)
			break;
		if (!segment_targets_type(dest, memsz, dest_type))
			return -1;
	}
	return 0;
}

static int load_payload_segments(struct cbfs_payload_segment *cbfssegs, uintptr_t *entry)
{
	uint8_t *dest, *src;
	size_t filesz, memsz;
	uint32_t compression;
	struct cbfs_payload_segment *first_segment, *seg, segment;
	int flags = 0;

	for (first_segment = seg = cbfssegs;; ++seg) {
		printk(BIOS_DEBUG, "Loading segment from ROM address 0x%p\n", seg);

		cbfs_decode_payload_segment(&segment, seg);
		dest = (uint8_t *)(uintptr_t)segment.load_addr;
		memsz = segment.mem_len;
		compression = segment.compression;
		filesz = segment.len;

		switch (segment.type) {
		case PAYLOAD_SEGMENT_CODE:
		case PAYLOAD_SEGMENT_DATA:
			printk(BIOS_DEBUG, "  %s (compression=%x)\n",
				segment.type == PAYLOAD_SEGMENT_CODE
				?  "code" : "data", segment.compression);
			src = ((uint8_t *)first_segment) + segment.offset;
			printk(BIOS_DEBUG,
				"  New segment dstaddr 0x%p memsize 0x%zx srcaddr 0x%p filesize 0x%zx\n",
			       dest, memsz, src, filesz);

			/* Clean up the values */
			if (filesz > memsz)  {
				filesz = memsz;
				printk(BIOS_DEBUG, "  cleaned up filesize 0x%zx\n", filesz);
			}
			break;

		case PAYLOAD_SEGMENT_BSS:
			printk(BIOS_DEBUG, "  BSS 0x%p (%d byte)\n", (void *)
				(intptr_t)segment.load_addr, segment.mem_len);
			filesz = 0;
			src = ((uint8_t *)first_segment) + segment.offset;
			compression = CBFS_COMPRESS_NONE;
			break;

		case PAYLOAD_SEGMENT_ENTRY:
			printk(BIOS_DEBUG, "  Entry Point 0x%p\n", (void *)
				(intptr_t)segment.load_addr);

			*entry = segment.load_addr;
			/* Per definition, a payload always has the entry point
			 * as last segment. Thus, we use the occurrence of the
			 * entry point as break condition for the loop.
			 */
			return 0;

		default:
			/* We found something that we don't know about. Throw
			 * hands into the sky and run away!
			 */
			printk(BIOS_EMERG, "Bad segment type %x\n", segment.type);
			return -1;
		}
		/* Note that the 'seg + 1' is safe as we only call this
		 * function on "not the last" * items, since entry
		 * is always last. */
		if (last_loadable_segment(seg))
			flags = SEG_FINAL;
		if (!load_one_segment(dest, src, filesz, memsz, compression, flags))
			return -1;
	}

	return 1;
}

__weak int payload_arch_usable_ram_quirk(uint64_t start, uint64_t size)
{
	return 0;
}

static void *selfprepare(struct prog *payload)
{
	void *data;
	data = rdev_mmap_full(prog_rdev(payload));
	return data;
}

static bool _selfload(struct prog *payload, checker_t f, void *args)
{
	uintptr_t entry = 0;
	struct cbfs_payload_segment *cbfssegs;
	void *data;

	data = selfprepare(payload);
	if (data == NULL)
		return false;

	cbfssegs = &((struct cbfs_payload *)data)->segments;

	if (f && f(cbfssegs, args))
		goto out;

	if (load_payload_segments(cbfssegs, &entry))
		goto out;

	printk(BIOS_SPEW, "Loaded segments\n");

	rdev_munmap(prog_rdev(payload), data);

	/* Pass cbtables to payload if architecture desires it. */
	prog_set_entry(payload, (void *)entry, cbmem_find(CBMEM_ID_CBTABLE));

	return true;
out:
	rdev_munmap(prog_rdev(payload), data);
	return false;
}

bool selfload_check(struct prog *payload, enum bootmem_type dest_type)
{
	return _selfload(payload, check_payload_segments, &dest_type);
}

bool selfload(struct prog *payload)
{
	return _selfload(payload, NULL, 0);
}
