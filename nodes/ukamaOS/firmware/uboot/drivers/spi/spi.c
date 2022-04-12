// SPDX-License-Identifier: GPL-2.0+
/*
 * Copyright (c) 2011 The Chromium OS Authors.
 */

#include <common.h>
#include <errno.h>
#include <fdtdec.h>
#include <malloc.h>
#include <spi.h>

int spi_set_wordlen(struct spi_slave *slave, unsigned int wordlen)
{
	if (wordlen == 0 || wordlen > 32) {
		printf("spi: invalid wordlen %u\n", wordlen);
		return -1;
	}

	slave->wordlen = wordlen;

	return 0;
}

void *spi_do_alloc_slave(int offset, int size, unsigned int bus,
			 unsigned int cs)
{
	u8 *ptr;

	ptr = malloc(size);
	if (ptr) {
		struct spi_slave *slave;

		memset(ptr, '\0', size);
		slave = (struct spi_slave *)(ptr + offset);
		slave->bus = bus;
		slave->cs = cs;
		slave->wordlen = SPI_DEFAULT_WORDLEN;
	}

	return ptr;
}

__weak bool spi_is_flash_command_supported(struct spi_slave *slave,
					   const struct spi_flash_command *cmd)
{
	return false;
}

__weak int spi_exec_flash_command(struct spi_slave *slave,
				  const struct spi_flash_command *cmd)
{
	return -EINVAL;
}
