/* SPDX-License-Identifier: GPL-2.0+ */
/*
 * Copyright (C) 2018 Synopsys, Inc. All rights reserved.
 */

#ifndef _CONFIG_EMDK_H_
#define _CONFIG_EMDK_H_

#include <linux/sizes.h>

#define CONFIG_SYS_MONITOR_BASE		CONFIG_SYS_TEXT_BASE

#define CONFIG_SYS_SDRAM_BASE		0x10000000
#define CONFIG_SYS_SDRAM_SIZE		SZ_8M

#define CONFIG_SYS_INIT_SP_ADDR		(CONFIG_SYS_SDRAM_BASE + SZ_1M)

#define CONFIG_SYS_MALLOC_LEN		SZ_64K
#define CONFIG_SYS_LOAD_ADDR		CONFIG_SYS_SDRAM_BASE

/* Required by DW MMC driver */
#define CONFIG_BOUNCE_BUFFER

/*
 * Environment
 */
#define CONFIG_ENV_SIZE			SZ_4K
#define CONFIG_BOOTFILE			"app.bin"
#define CONFIG_LOADADDR			CONFIG_SYS_LOAD_ADDR

#define CONFIG_EXTRA_ENV_SETTINGS \
	"upgrade_image=u-boot.bin\0" \
	"upgrade=emdk rom unlock && " \
		"fatload mmc 0 ${loadaddr} ${upgrade_image} && " \
		"cp.b ${loadaddr} 0 ${filesize} && " \
		"dcache flush && " \
		"emdk rom lock\0"

#endif /* _CONFIG_EMDK_H_ */

