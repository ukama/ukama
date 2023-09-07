/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 Qualcomm Technologies.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 2 and
 * only version 2 as published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */
#include <types.h>
#include <soc/addressmap.h>
#include <spi-generic.h>

#ifndef __SOC_QUALCOMM_SDM845_QSPI_H__
#define __SOC_QUALCOMM_SDM845_QSPI_H__

struct sdm845_qspi_regs {
	u32 mstr_cfg;
	u32 ahb_mstr_cfg;
	u32 reserve_0;
	u32 mstr_int_en;
	u32 mstr_int_sts;
	u32 pio_xfer_ctrl;
	u32 pio_xfer_cfg;
	u32 pio_xfer_sts;
	u32 pio_dataout_1byte;
	u32 pio_dataout_4byte;
	u32 rd_fifo_cfg;
	u32 rd_fifo_sts;
	u32 rd_fifo_rst;
	u32 reserve_1[3];
	u32 next_dma_desc_addr;
	u32 current_dma_desc_addr;
	u32 current_mem_addr;
	u32 hw_version;
	u32 rd_fifo[16];
};

check_member(sdm845_qspi_regs, rd_fifo, 0x50);
static struct sdm845_qspi_regs * const sdm845_qspi = (void *) QSPI_BASE;

// MSTR_CONFIG register

#define TX_DATA_OE_DELAY_SHIFT 24
#define TX_DATA_OE_DELAY_MASK (0x3 << TX_DATA_OE_DELAY_SHIFT)
#define TX_CS_N_DELAY_SHIFT 22
#define TX_CS_N_DELAY_MASK (0x3 << TX_CS_N_DELAY_SHIFT)
#define TX_CLK_DELAY_SHIFT 20
#define TX_CLK_DELAY_MASK (0x3 << TX_CLK_DELAY_SHIFT)
#define TX_DATA_DELAY_SHIFT 18
#define TX_DATA_DELAY_MASK (0x3 << TX_DATA_DELAY_SHIFT)
#define LPA_BASE_SHIFT 14
#define LPA_BASE_MASK (0xF << LPA_BASE_SHIFT)
#define SBL_EN BIT(13)
#define CHIP_SELECT_NUM BIT(12)
#define SPI_MODE_SHIFT 10
#define SPI_MODE_MASK (0x3 << SPI_MODE_SHIFT)
#define BIG_ENDIAN_MODE BIT(9)
#define DMA_ENABLE BIT(8)
#define PIN_WPN BIT(7)
#define PIN_HOLDN BIT(6)
#define FB_CLK_EN BIT(4)
#define FULL_CYCLE_MODE BIT(3)

// MSTR_INT_ENABLE and MSTR_INT_STATUS register

#define DMA_CHAIN_DONE BIT(31)
#define TRANSACTION_DONE BIT(16)
#define WRITE_FIFO_OVERRUN BIT(11)
#define WRITE_FIFO_FULL BIT(10)
#define HRESP_FROM_NOC_ERR BIT(3)
#define RESP_FIFO_RDY BIT(2)
#define RESP_FIFO_NOT_EMPTY BIT(1)
#define RESP_FIFO_UNDERRUN BIT(0)

// PIO_TRANSFER_CONFIG register

#define TRANSFER_FRAGMENT BIT(8)
#define MULTI_IO_MODE_SHIFT 1
#define MULTI_IO_MODE_MASK (0x7 << MULTI_IO_MODE_SHIFT)
#define TRANSFER_DIRECTION BIT(0)

// PIO_TRANSFER_STATUS register

#define WR_FIFO_BYTES_SHIFT 16
#define WR_FIFO_BYTES_MASK (0xFFFF << WR_FIFO_BYTES_SHIFT)

// RD_FIFO_CONFIG register

#define CONTINUOUS_MODE BIT(0)

// RD_FIFO_STATUS register

#define FIFO_EMPTY BIT(11)
#define WR_CNTS_SHIFT 4
#define WR_CNTS_MASK (0x7F << WR_CNTS_SHIFT)
#define RDY_64BYTE BIT(3)
#define RDY_32BYTE BIT(2)
#define RDY_16BYTE BIT(1)
#define FIFO_RDY BIT(0)

// RD_FIFO_RESET register

#define RESET_FIFO BIT(0)

#define QSPI_MAX_PACKET_COUNT 0xFFC0

void quadspi_init(uint32_t hz);
int sdm845_claim_bus(const struct spi_slave *slave);
int sdm845_setup_bus(const struct spi_slave *slave);
void sdm845_release_bus(const struct spi_slave *slave);
int sdm845_xfer(const struct spi_slave *slave, const void *dout,
		size_t out_bytes, void *din, size_t in_bytes);
int sdm845_xfer_dual(const struct spi_slave *slave, const void *dout,
		     size_t out_bytes, void *din, size_t in_bytes);
#endif /* __SOC_QUALCOMM_SDM845_QSPI_H__ */
