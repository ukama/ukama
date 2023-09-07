// SPDX-License-Identifier: GPL-2.0+
/*
 * Copyright (C) 2018, STMicroelectronics - All Rights Reserved
 */

#include <common.h>
#include <command.h>
#include <misc.h>
#include <errno.h>
#include <dm/device.h>
#include <dm/uclass.h>

#define STM32MP_OTP_BANK	0

/*
 * The 'fuse' command API
 */
int fuse_read(u32 bank, u32 word, u32 *val)
{
	int ret = 0;
	struct udevice *dev;

	switch (bank) {
	case STM32MP_OTP_BANK:
		ret = uclass_get_device_by_driver(UCLASS_MISC,
						  DM_GET_DRIVER(stm32mp_bsec),
						  &dev);
		if (ret)
			return ret;
		ret = misc_read(dev, word * 4 + STM32_BSEC_SHADOW_OFFSET,
				val, 4);
		break;

	default:
		printf("stm32mp %s: wrong value for bank %i\n", __func__, bank);
		ret = -EINVAL;
		break;
	}

	return ret;
}

int fuse_prog(u32 bank, u32 word, u32 val)
{
	struct udevice *dev;
	int ret;

	switch (bank) {
	case STM32MP_OTP_BANK:
		ret = uclass_get_device_by_driver(UCLASS_MISC,
						  DM_GET_DRIVER(stm32mp_bsec),
						  &dev);
		if (ret)
			return ret;
		ret = misc_write(dev, word * 4 + STM32_BSEC_OTP_OFFSET,
				 &val, 4);
		break;

	default:
		printf("stm32mp %s: wrong value for bank %i\n", __func__, bank);
		ret = -EINVAL;
		break;
	}

	return ret;
}

int fuse_sense(u32 bank, u32 word, u32 *val)
{
	struct udevice *dev;
	int ret;

	switch (bank) {
	case STM32MP_OTP_BANK:
		ret = uclass_get_device_by_driver(UCLASS_MISC,
						  DM_GET_DRIVER(stm32mp_bsec),
						  &dev);
		if (ret)
			return ret;
		ret = misc_read(dev, word * 4 + STM32_BSEC_OTP_OFFSET, val, 4);
		break;

	default:
		printf("stm32mp %s: wrong value for bank %i\n", __func__, bank);
		ret = -EINVAL;
		break;
	}

	return ret;
}

int fuse_override(u32 bank, u32 word, u32 val)
{
	struct udevice *dev;
	int ret;

	switch (bank) {
	case STM32MP_OTP_BANK:
		ret = uclass_get_device_by_driver(UCLASS_MISC,
						  DM_GET_DRIVER(stm32mp_bsec),
						  &dev);
		if (ret)
			return ret;
		ret = misc_write(dev, word * 4 + STM32_BSEC_SHADOW_OFFSET,
				 &val, 4);
		break;

	default:
		printf("stm32mp %s: wrong value for bank %i\n",
		       __func__, bank);
		ret = -EINVAL;
		break;
	}

	return ret;
}
