// SPDX-License-Identifier: GPL-2.0+
/*
 * SPI flash interface
 *
 * Copyright (C) 2008 Atmel Corporation
 * Copyright (C) 2010 Reinhard Meyer, EMK Elektronik
 */

#include <common.h>
#include <spi.h>
#include <spi_flash.h>

#include "sf_internal.h"

static void spi_flash_addr(u32 addr, u8 addr_len, u8 *cmd_buf)
{
	u8 i;

	for (i = 0; i < addr_len; i++)
		cmd_buf[i] = addr >> ((addr_len - 1 - i) * 8);
}

static u8 spi_compute_num_dummy_bytes(enum spi_flash_protocol proto,
				      u8 num_dummy_clock_cycles)
{
	int shift = fls(spi_flash_protocol_get_addr_nbits(proto)) - 1;

	if (shift < 0)
		shift = 0;
	return (num_dummy_clock_cycles << shift) >> 3;
}

static int spi_flash_exec(struct spi_flash *flash,
			  const struct spi_flash_command *cmd)
{
	struct spi_slave *spi = flash->spi;
	u8 cmd_buf[SPI_FLASH_CMD_LEN];
	size_t cmd_len, num_dummy_bytes;
	unsigned long flags = SPI_XFER_BEGIN;
	int ret;

	if (spi_is_flash_command_supported(spi, cmd))
		return spi_exec_flash_command(spi, cmd);

	if (cmd->data_len == 0)
		flags |= SPI_XFER_END;

	cmd_buf[0] = cmd->inst;
	spi_flash_addr(cmd->addr, cmd->addr_len, cmd_buf + 1);
	cmd_len = 1 + cmd->addr_len;

	num_dummy_bytes = spi_compute_num_dummy_bytes(cmd->proto,
						      cmd->num_mode_cycles +
						      cmd->num_wait_states);
	memset(cmd_buf + cmd_len, 0xff, num_dummy_bytes);
	cmd_len += num_dummy_bytes;

	ret = spi_xfer(spi, cmd_len * 8, cmd_buf, NULL, flags);
	if (ret) {
		debug("SF: Failed to send command (%zu bytes): %d\n",
		      cmd_len, ret);
	} else if (cmd->data_len != 0) {
		ret = spi_xfer(spi, cmd->data_len * 8,
			       cmd->tx_data, cmd->rx_data,
			       SPI_XFER_END);
		if (ret)
			debug("SF: Failed to transfer %zu bytes of data: %d\n",
			      cmd->data_len, ret);
	}

	return ret;
}

int spi_flash_cmd_read(struct spi_flash *flash,
		       const struct spi_flash_command *cmd)
{
	return spi_flash_exec(flash, cmd);
}

int spi_flash_cmd(struct spi_flash *flash, u8 instr, void *response, size_t len)
{
	struct spi_flash_command cmd;
	u8 flags = (response && len) ? SPI_FCMD_READ_REG : SPI_FCMD_WRITE_REG;

	spi_flash_command_init(&cmd, instr, 0, flags);
	cmd.data_len = len;
	cmd.rx_data = response;
	return spi_flash_exec(flash, &cmd);
}

int spi_flash_cmd_write(struct spi_flash *flash,
			const struct spi_flash_command *cmd)
{
	return spi_flash_exec(flash, cmd);
}
