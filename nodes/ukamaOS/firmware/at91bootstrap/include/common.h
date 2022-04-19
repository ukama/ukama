/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2012, Atmel Corporation
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * - Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the disclaimer below.
 *
 * Atmel's name may not be used to endorse or promote products derived from
 * this software without specific prior written permission.
 *
 * DISCLAIMER: THIS SOFTWARE IS PROVIDED BY ATMEL "AS IS" AND ANY EXPRESS OR
 * IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT ARE
 * DISCLAIMED. IN NO EVENT SHALL ATMEL BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
 * LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA,
 * OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF
 * LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING
 * NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */
#ifndef __COMMON_H__
#define __COMMON_H__

#define ARRAY_SIZE(array) (sizeof(array) / sizeof((array)[0]))

#define ALIGN(size, align)	(((size) + (align) - 1) & (~((align) - 1)))
#define OF_ALIGN(size)		ALIGN(size, 4)

#define min(a, b)	(((a) < (b)) ? (a) : (b))
#define max(a, b)	(((a) > (b)) ? (a) : (b))

#ifndef NULL
#define	NULL	0
#endif

#define FILENAME_BUF_LEN	32
#define CMDLINE_BUF_LEN		256

enum {
	KERNEL_IMAGE,
	DT_BLOB,
};

/* structure definition */
struct image_info
{
#if defined(CONFIG_DATAFLASH) || defined(CONFIG_NANDFLASH) || defined(CONFIG_FLASH)
	unsigned int offset;
	unsigned int length;
#endif
#ifdef CONFIG_SDCARD
	char *filename;
#ifdef CONFIG_OVERRIDE_CMDLINE_FROM_EXT_FILE
	char *cmdline_file;
	char *cmdline_args;
#endif
#endif
	unsigned char *dest;

#ifdef CONFIG_OF_LIBFDT
#if defined(CONFIG_DATAFLASH) || defined(CONFIG_NANDFLASH) || defined(CONFIG_FLASH)
	unsigned int of_offset;
	unsigned int of_length;
#endif
#ifdef CONFIG_SDCARD
	char *of_filename;
#endif
	unsigned char *of_dest;
#endif
};

typedef int (*load_function)(struct image_info *image);

extern load_function load_image;
extern void init_load_image(struct image_info *image);
extern void load_image_done(int retval);

extern int load_kernel(struct image_info *image);

extern int kernel_size(unsigned char *addr);

static inline unsigned int swap_uint32(unsigned int data)
{
	volatile unsigned int a, b, c, d;

	a = ((data) & 0xff000000) >> 24;
	b = ((data) & 0x00ff0000) >> 8;
	c = ((data) & 0x0000ff00) << 8;
	d = ((data) & 0x000000ff) << 24;

	return a | b | c | d;
}

#endif /* #ifdef __COMMON_H__ */
