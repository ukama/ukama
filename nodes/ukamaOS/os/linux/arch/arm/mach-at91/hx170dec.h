/*
 * Decoder device driver (kernel module headers)
 *
 * Copyright (C) 2009  Hantro Products Oy.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

#ifndef _HX170DEC_H_
#define _HX170DEC_H_

#include <linux/ioctl.h>
#include <linux/types.h>

struct core_desc
{
	__u32 id;    /* id of the core */
	__u32 *regs; /* pointer to user registers */
	__u32 size;  /* size of register space */
};

/* Use 'k' as magic number */
#define HX170DEC_IOC_MAGIC  'k'
/*
 * S means "Set" through a ptr,
 * T means "Tell" directly with the argument value
 * G means "Get": reply by setting through a pointer
 * Q means "Query": response is on the return value
 * X means "eXchange": G and S atomically
 * H means "sHift": T and Q atomically
 */

#define HX170DEC_IOCGHWOFFSET		_IOR(HX170DEC_IOC_MAGIC, 3, unsigned long *)
#define HX170DEC_IOCGHWIOSIZE		_IOR(HX170DEC_IOC_MAGIC, 4, unsigned int *)

#define HX170DEC_IOC_MC_OFFSETS		_IOR(HX170DEC_IOC_MAGIC, 7, unsigned long *)
#define HX170DEC_IOC_MC_CORES		_IOR(HX170DEC_IOC_MAGIC, 8, unsigned int *)
#define HX170DEC_IOCS_DEC_PUSH_REG	_IOW(HX170DEC_IOC_MAGIC, 9, struct core_desc *)
#define HX170DEC_IOCS_PP_PUSH_REG	_IOW(HX170DEC_IOC_MAGIC, 10, struct core_desc *)
#define HX170DEC_IOCH_DEC_RESERVE	_IO(HX170DEC_IOC_MAGIC, 11)
#define HX170DEC_IOCT_DEC_RELEASE	_IO(HX170DEC_IOC_MAGIC, 12)
#define HX170DEC_IOCQ_PP_RESERVE	_IO(HX170DEC_IOC_MAGIC, 13)
#define HX170DEC_IOCT_PP_RELEASE	_IO(HX170DEC_IOC_MAGIC, 14)
#define HX170DEC_IOCX_DEC_WAIT		_IOWR(HX170DEC_IOC_MAGIC, 15, struct core_desc *)
#define HX170DEC_IOCX_PP_WAIT		_IOWR(HX170DEC_IOC_MAGIC, 16, struct core_desc *)
#define HX170DEC_IOCS_DEC_PULL_REG	_IOWR(HX170DEC_IOC_MAGIC, 17, struct core_desc *)
#define HX170DEC_IOCS_PP_PULL_REG	_IOWR(HX170DEC_IOC_MAGIC, 18, struct core_desc *)

#define HX170DEC_IOX_ASIC_ID		_IOWR(HX170DEC_IOC_MAGIC, 20, __u32 *)

/*
 * Following are not used yet:
 *
 * #define HX170DEC_PP_INSTANCE		_IO(HX170DEC_IOC_MAGIC, 1)
 * #define HX170DEC_HW_PERFORMANCE	_IO(HX170DEC_IOC_MAGIC, 2)
 * #define HX170DEC_IOC_CLI		_IO(HX170DEC_IOC_MAGIC, 5)
 * #define HX170DEC_IOC_STI		_IO(HX170DEC_IOC_MAGIC, 6)
 * #define HX170DEC_IOCG_CORE_WAIT	_IOR(HX170DEC_IOC_MAGIC, 19, int *)
 * #define HX170DEC_DEBUG_STATUS	_IO(HX170DEC_IOC_MAGIC, 29)
 */

#define HX170DEC_IOC_MAXNR 29

#endif /* !_HX170DEC_H_ */
