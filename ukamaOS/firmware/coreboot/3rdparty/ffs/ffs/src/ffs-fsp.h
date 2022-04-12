/* IBM_PROLOG_BEGIN_TAG                                                   */
/* This is an automatically generated prolog.                             */
/*                                                                        */
/* $Source: ffs/src/ffs-fsp.h $                                           */
/*                                                                        */
/* OpenPOWER FFS Project                                                  */
/*                                                                        */
/* Contributors Listed Below - COPYRIGHT 2014,2015                        */
/* [+] International Business Machines Corp.                              */
/*                                                                        */
/*                                                                        */
/* Licensed under the Apache License, Version 2.0 (the "License");        */
/* you may not use this file except in compliance with the License.       */
/* You may obtain a copy of the License at                                */
/*                                                                        */
/*     http://www.apache.org/licenses/LICENSE-2.0                         */
/*                                                                        */
/* Unless required by applicable law or agreed to in writing, software    */
/* distributed under the License is distributed on an "AS IS" BASIS,      */
/* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or        */
/* implied. See the License for the specific language governing           */
/* permissions and limitations under the License.                         */
/*                                                                        */
/* IBM_PROLOG_END_TAG                                                     */
#ifndef __FFS_BOOT_H__
#define __FFS_BOOT_H__

#include "ffs.h"

/*
 * Values to use in USER_DATA_VOL
 *
 * @note: Side 0/1 must be defined as even/odd values.  Code in the IPL depends
 * on this to be able to use the appropriate volume based on the boot bank.
 */
enum vol {
	FFS_VOL_IPL0       = 0,
	FFS_VOL_IPL1       = 1,
	FFS_VOL_SPL0       = 2,
	FFS_VOL_SPL1       = 3,
	FFS_VOL_BOOTENV0   = 4,
	FFS_VOL_BOOTENV1   = 5,
	FFS_VOL_KERNEL0    = 6,
	FFS_VOL_KERNEL1    = 7,
	FFS_VOL_INITRAMFS0 = 8,
	FFS_VOL_INITRAMFS1 = 9,
	FFS_VOL_DTB0       = 10,
	FFS_VOL_DTB1       = 11,
	FFS_VOL_SERIES0    = 12,
	FFS_VOL_SERIES1    = 13,
	FFS_VOL_CARD0      = 14,
	FFS_VOL_CARD1      = 15,
	FFS_VOL_DUMP0      = 16,
	FFS_VOL_DUMP1      = 17,
	FFS_VOL_DUMP2      = 18,
	FFS_VOL_DUMP3      = 19,
};

#endif /* __FFS_BOOT_H__ */
