/* Copyright 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Register map for NPCX processor
 */

#ifndef __CROS_EC_REGISTERS_H
#define __CROS_EC_REGISTERS_H

#include "common.h"
#include "compile_time_macros.h"
#include "clock_chip.h"

/******************************************************************************/
/*
 * Macro Functions
 */
/* Bit functions */
#define SET_BIT(reg, bit)           ((reg) |= (0x1 << (bit)))
#define CLEAR_BIT(reg, bit)         ((reg) &= (~(0x1 << (bit))))
#define IS_BIT_SET(reg, bit)        ((reg >> bit) & (0x1))
#define UPDATE_BIT(reg, bit, cond)  {	if (cond) \
						SET_BIT(reg, bit); \
					else \
						CLEAR_BIT(reg, bit); }
/* Field functions */
#define GET_POS_FIELD(pos, size)    pos
#define GET_SIZE_FIELD(pos, size)   size
#define FIELD_POS(field)            GET_POS_##field
#define FIELD_SIZE(field)           GET_SIZE_##field
/* Read field functions */
#define GET_FIELD(reg, field) \
	_GET_FIELD_(reg, FIELD_POS(field), FIELD_SIZE(field))
#define _GET_FIELD_(reg, f_pos, f_size) (((reg)>>(f_pos)) & ((1<<(f_size))-1))
/* Write field functions */
#define SET_FIELD(reg, field, value) \
	_SET_FIELD_(reg, FIELD_POS(field), FIELD_SIZE(field), value)
#define _SET_FIELD_(reg, f_pos, f_size, value) \
	((reg) = ((reg) & (~(((1 << (f_size))-1) << (f_pos)))) \
			| ((value) << (f_pos)))

/******************************************************************************/
/*
 * NPCX (Nuvoton M4 EC) Register Definitions
 */

/* Global Definition */
#define I2C_7BITS_ADDR                   0
/* Switcher of features */
#define SUPPORT_LCT                      1
#define SUPPORT_WDG                      1
#define SUPPORT_P80_SEG                  0 /* Note: it uses KSO10 & KSO11 */
/* Switcher of debugging */
#define DEBUG_GPIO                       0
#define DEBUG_I2C                        0
#define DEBUG_TMR                        0
#define DEBUG_WDG                        0
#define DEBUG_FAN                        0
#define DEBUG_PWM                        0
#define DEBUG_SPI                        0
#define DEBUG_FLH                        0
#define DEBUG_PECI                       0
#define DEBUG_SHI                        0
#define DEBUG_CLK                        0
#define DEBUG_LPC                        0
#define DEBUG_ESPI                       0
#define DEBUG_WOV                        0
#define DEBUG_CEC                        0
#define DEBUG_SIB                        0

/* Modules Map */
#define NPCX_ESPI_BASE_ADDR              0x4000A000
#define NPCX_MDC_BASE_ADDR               0x4000C000
#define NPCX_PMC_BASE_ADDR               0x4000D000
#define NPCX_SIB_BASE_ADDR               0x4000E000
#define NPCX_SHI_BASE_ADDR               0x4000F000
#define NPCX_SHM_BASE_ADDR               0x40010000
#define NPCX_GDMA_BASE_ADDR              0x40011000
#define NPCX_FIU_BASE_ADDR               0x40020000
#define NPCX_KBSCAN_REGS_BASE            0x400A3000
#define NPCX_WOV_BASE_ADDR               0x400A4000
#define NPCX_APM_BASE_ADDR               0x400A4800
#define NPCX_GLUE_REGS_BASE              0x400A5000
#define NPCX_BBRAM_BASE_ADDR             0x400AF000
#define NPCX_HFCG_BASE_ADDR              0x400B5000
#define NPCX_LFCG_BASE_ADDR              0x400B5100
#define NPCX_FMUL2_BASE_ADDR             0x400B5200
#define NPCX_MTC_BASE_ADDR               0x400B7000
#define NPCX_MSWC_BASE_ADDR              0x400C1000
#define NPCX_SCFG_BASE_ADDR              0x400C3000
#define NPCX_KBC_BASE_ADDR               0x400C7000
#define NPCX_ADC_BASE_ADDR               0x400D1000
#define NPCX_SPI_BASE_ADDR               0x400D2000
#define NPCX_PECI_BASE_ADDR              0x400D4000
#define NPCX_TWD_BASE_ADDR               0x400D8000

/* Multi-Modules Map */
#define NPCX_PWM_BASE_ADDR(mdl)          (0x40080000 + ((mdl) * 0x2000L))
#define NPCX_GPIO_BASE_ADDR(mdl)         (0x40081000 + ((mdl) * 0x2000L))
#define NPCX_ITIM16_BASE_ADDR(mdl)       (0x400B0000 + ((mdl) * 0x2000L))
#define NPCX_ITIM32_BASE_ADDR            0x400BC000
#define NPCX_MIWU_BASE_ADDR(mdl)         (0x400BB000 + ((mdl) * 0x2000L))
#define NPCX_MFT_BASE_ADDR(mdl)          (0x400E1000 + ((mdl) * 0x2000L))
#define NPCX_CR_UART_BASE_ADDR(mdl)      (0x400C4000 + ((mdl) * 0x2000L))
#define NPCX_PM_CH_BASE_ADDR(mdl)        (0x400C9000 + ((mdl) * 0x2000L))
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_SMB_BASE_ADDR(mdl)        (((mdl) < 2) ? \
					(0x40009000 +  ((mdl) * 0x2000L)) : \
					((mdl) < 4) ? \
				      (0x400C0000 + (((mdl) - 2) * 0x2000L)) : \
					((mdl) == 4) ? \
					(0x40008000) : \
					(0x40017000 + (((mdl) - 5) * 0x1000L)))
#else
#define NPCX_SMB_BASE_ADDR(mdl)        (((mdl) < 2) ? \
					(0x40009000 +  ((mdl) * 0x2000L)) : \
					(0x400C0000 + (((mdl) - 2) * 0x2000L)))
#endif

/*
 * NPCX-IRQ numbers
 */
#define NPCX_IRQ_0                       0
#define NPCX_IRQ_1                       1
#define NPCX_IRQ_2                       2
#define NPCX_IRQ_3                       3
#define NPCX_IRQ_4                       4
#define NPCX_IRQ_5                       5
#define NPCX_IRQ_6                       6
#define NPCX_IRQ_7                       7
#define NPCX_IRQ_8                       8
#define NPCX_IRQ_9                       9
#define NPCX_IRQ_10                      10
#define NPCX_IRQ_11                      11
#define NPCX_IRQ_12                      12
#define NPCX_IRQ_13                      13
#define NPCX_IRQ_14                      14
#define NPCX_IRQ_15                      15
#define NPCX_IRQ_16                      16
#define NPCX_IRQ_17                      17
#define NPCX_IRQ_18                      18
#define NPCX_IRQ_19                      19
#define NPCX_IRQ_20                      20
#define NPCX_IRQ_21                      21
#define NPCX_IRQ_22                      22
#define NPCX_IRQ_23                      23
#define NPCX_IRQ_24                      24
#define NPCX_IRQ_25                      25
#define NPCX_IRQ_26                      26
#define NPCX_IRQ_27                      27
#define NPCX_IRQ_28                      28
#define NPCX_IRQ_29                      29
#define NPCX_IRQ_30                      30
#define NPCX_IRQ_31                      31
#define NPCX_IRQ_32                      32
#define NPCX_IRQ_33                      33
#define NPCX_IRQ_34                      34
#define NPCX_IRQ_35                      35
#define NPCX_IRQ_36                      36
#define NPCX_IRQ_37                      37
#define NPCX_IRQ_38                      38
#define NPCX_IRQ_39                      39
#define NPCX_IRQ_40                      40
#define NPCX_IRQ_41                      41
#define NPCX_IRQ_42                      42
#define NPCX_IRQ_43                      43
#define NPCX_IRQ_44                      44
#define NPCX_IRQ_45                      45
#define NPCX_IRQ_46                      46
#define NPCX_IRQ_47                      47
#define NPCX_IRQ_48                      48
#define NPCX_IRQ_49                      49
#define NPCX_IRQ_50                      50
#define NPCX_IRQ_51                      51
#define NPCX_IRQ_52                      52
#define NPCX_IRQ_53                      53
#define NPCX_IRQ_54                      54
#define NPCX_IRQ_55                      55
#define NPCX_IRQ_56                      56
#define NPCX_IRQ_57                      57
#define NPCX_IRQ_58                      58
#define NPCX_IRQ_59                      59
#define NPCX_IRQ_60                      60
#define NPCX_IRQ_61                      61
#define NPCX_IRQ_62                      62
#define NPCX_IRQ_63                      63

#define NPCX_IRQ0_NOUSED                 NPCX_IRQ_0
#define NPCX_IRQ1_NOUSED                 NPCX_IRQ_1
#define NPCX_IRQ_KBSCAN                  NPCX_IRQ_2
#define NPCX_IRQ_PM_CHAN_OBE             NPCX_IRQ_3
#define NPCX_IRQ_PECI                    NPCX_IRQ_4
#define NPCX_IRQ5_NOUSED                 NPCX_IRQ_5
#define NPCX_IRQ_PORT80                  NPCX_IRQ_6
#define NPCX_IRQ_MTC_WKINTAD_0           NPCX_IRQ_7
#define NPCX_IRQ_SMB8                    NPCX_IRQ_8
#define NPCX_IRQ_MFT_1                   NPCX_IRQ_9
#define NPCX_IRQ_ADC                     NPCX_IRQ_10
#define NPCX_IRQ_WKINTEFGH_0             NPCX_IRQ_11
#define NPCX_IRQ_CDMA                    NPCX_IRQ_12
#define NPCX_IRQ_SMB1                    NPCX_IRQ_13
#define NPCX_IRQ_SMB2                    NPCX_IRQ_14
#define NPCX_IRQ_WKINTC_0                NPCX_IRQ_15
#define NPCX_IRQ_SMB7                    NPCX_IRQ_16
#define NPCX_IRQ_ITIM16_3                NPCX_IRQ_17
#define NPCX_IRQ_SHI                     NPCX_IRQ_18
#define NPCX_IRQ_ESPI                    NPCX_IRQ_18
#define NPCX_IRQ_SMB5                    NPCX_IRQ_19
#define NPCX_IRQ_SMB6                    NPCX_IRQ_20
#define NPCX_IRQ_PS2                     NPCX_IRQ_21
#define NPCX_IRQ_WOV                     NPCX_IRQ_22
#define NPCX_IRQ_MFT_2                   NPCX_IRQ_23
#define NPCX_IRQ_SHM                     NPCX_IRQ_24
#define NPCX_IRQ_KBC_IBF                 NPCX_IRQ_25
#define NPCX_IRQ_PM_CHAN_IBF             NPCX_IRQ_26
#define NPCX_IRQ_ITIM16_2                NPCX_IRQ_27
#define NPCX_IRQ_ITIM16_1                NPCX_IRQ_28
#define NPCX_IRQ29_NOUSED                NPCX_IRQ_29
#define NPCX_IRQ30_NOUSED                NPCX_IRQ_30
#define NPCX_IRQ_TWD_WKINTB_0            NPCX_IRQ_31
#define NPCX_IRQ_UART2                   NPCX_IRQ_32
#define NPCX_IRQ_UART                    NPCX_IRQ_33
#define NPCX_IRQ34_NOUSED                NPCX_IRQ_34
#define NPCX_IRQ35_NOUSED                NPCX_IRQ_35
#define NPCX_IRQ_SMB3                    NPCX_IRQ_36
#define NPCX_IRQ_SMB4                    NPCX_IRQ_37
#define NPCX_IRQ38_NOUSED                NPCX_IRQ_38
#define NPCX_IRQ39_NOUSED                NPCX_IRQ_39
#define NPCX_IRQ40_NOUSED                NPCX_IRQ_40
#define NPCX_IRQ_MFT_3                   NPCX_IRQ_41
#define NPCX_IRQ42_NOUSED                NPCX_IRQ_42
#define NPCX_IRQ_ITIM16_4                NPCX_IRQ_43
#define NPCX_IRQ_ITIM16_5                NPCX_IRQ_44
#define NPCX_IRQ_ITIM16_6                NPCX_IRQ_45
#define NPCX_IRQ_ITIM32                  NPCX_IRQ_46
#define NPCX_IRQ_WKINTA_1                NPCX_IRQ_47
#define NPCX_IRQ_WKINTB_1                NPCX_IRQ_48
#define NPCX_IRQ_KSI_WKINTC_1            NPCX_IRQ_49
#define NPCX_IRQ_WKINTD_1                NPCX_IRQ_50
#define NPCX_IRQ_WKINTE_1                NPCX_IRQ_51
#define NPCX_IRQ_WKINTF_1                NPCX_IRQ_52
#define NPCX_IRQ_WKINTG_1                NPCX_IRQ_53
#define NPCX_IRQ_WKINTH_1                NPCX_IRQ_54
#define NPCX_IRQ55_NOUSED                NPCX_IRQ_55
#define NPCX_IRQ_KBC_OBE                 NPCX_IRQ_56
#define NPCX_IRQ_SPI                     NPCX_IRQ_57
#define NPCX_IRQ58_NOUSED                NPCX_IRQ_58
#define NPCX_IRQ_WKINTFG_2               NPCX_IRQ_59
#define NPCX_IRQ_WKINTA_2                NPCX_IRQ_60
#define NPCX_IRQ_WKINTB_2                NPCX_IRQ_61
#define NPCX_IRQ_WKINTC_2                NPCX_IRQ_62
#define NPCX_IRQ_WKINTD_2                NPCX_IRQ_63

#define NPCX_IRQ_COUNT                   64

/******************************************************************************/
/* Miscellaneous Device Control (MDC) registers */
#define NPCX_FWCTRL                       REG8(NPCX_MDC_BASE_ADDR + 0x007)

/* MDC register fields */
#define NPCX_FWCTRL_RO_REGION            0
#define NPCX_FWCTRL_FW_SLOT              1

/******************************************************************************/
/* High Frequency Clock Generator (HFCG) registers */
#define NPCX_HFCGCTRL                     REG8(NPCX_HFCG_BASE_ADDR + 0x000)
#define NPCX_HFCGML                       REG8(NPCX_HFCG_BASE_ADDR + 0x002)
#define NPCX_HFCGMH                       REG8(NPCX_HFCG_BASE_ADDR + 0x004)
#define NPCX_HFCGN                        REG8(NPCX_HFCG_BASE_ADDR + 0x006)
#define NPCX_HFCGP                        REG8(NPCX_HFCG_BASE_ADDR + 0x008)
#define NPCX_HFCBCD                       REG8(NPCX_HFCG_BASE_ADDR + 0x010)
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_HFCBCD1                      REG8(NPCX_HFCG_BASE_ADDR + 0x012)
#define NPCX_HFCBCD2                      REG8(NPCX_HFCG_BASE_ADDR + 0x014)
#endif

/* HFCG register fields */
#define NPCX_HFCGCTRL_LOAD               0
#define NPCX_HFCGCTRL_LOCK               2
#define NPCX_HFCGCTRL_CLK_CHNG           7

/******************************************************************************/
/* Low Frequency Clock Generator (LFCG) registers */
#define NPCX_LFCGCTL                      REG8(NPCX_LFCG_BASE_ADDR + 0x000)
#define NPCX_HFRDI                       REG16(NPCX_LFCG_BASE_ADDR + 0x002)
#define NPCX_HFRDF                       REG16(NPCX_LFCG_BASE_ADDR + 0x004)
#define NPCX_FRCDIV                      REG16(NPCX_LFCG_BASE_ADDR + 0x006)
#define NPCX_DIVCOR1                     REG16(NPCX_LFCG_BASE_ADDR + 0x008)
#define NPCX_DIVCOR2                     REG16(NPCX_LFCG_BASE_ADDR + 0x00A)
#define NPCX_LFCGCTL2                     REG8(NPCX_LFCG_BASE_ADDR + 0x014)

/* LFCG register fields */
#define NPCX_LFCGCTL_XTCLK_VAL            7
#define NPCX_LFCGCTL2_XT_OSC_SL_EN        6

/******************************************************************************/
/* CR UART Register */
#define NPCX_UTBUF(n)                    REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x000)
#define NPCX_URBUF(n)                    REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x002)
#define NPCX_UICTRL(n)                   REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x004)
#define NPCX_USTAT(n)                    REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x006)
#define NPCX_UFRS(n)                     REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x008)
#define NPCX_UMDSL(n)                    REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x00A)
#define NPCX_UBAUD(n)                    REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x00C)
#define NPCX_UPSR(n)                     REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x00E)
#ifdef NPCX_UART_FIFO_SUPPORT
 /* UART registers only used for FIFO mode */
#define NPCX_UFTSTS(n)                   REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x020)
#define NPCX_UFRSTS(n)                   REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x022)
#define NPCX_UFTCTL(n)                   REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x024)
#define NPCX_UFRCTL(n)                   REG8(NPCX_CR_UART_BASE_ADDR(n) + 0x026)

/* UART FIFO register fields */
#define NPCX_UMDSL_FIFO_MD                0

#define NPCX_UFTSTS_TEMPTY_LVL            FIELD(0, 5)
#define NPCX_UFTSTS_TEMPTY_LVL_STS        5
#define NPCX_UFTSTS_TFIFO_EMPTY_STS       6
#define NPCX_UFTSTS_NXMIP                 7

#define NPCX_UFRSTS_RFULL_LVL_STS         5
#define NPCX_UFRSTS_RFIFO_NEMPTY_STS      6
#define NPCX_UFRSTS_ERR                   7

#define NPCX_UFTCTL_TEMPTY_LVL_SEL        FIELD(0, 5)
#define NPCX_UFTCTL_TEMPTY_LVL_EN         5
#define NPCX_UFTCTL_TEMPTY_EN             6
#define NPCX_UFTCTL_NXMIPEN               7

#define NPCX_UFRCTL_RFULL_LVL_SEL         FIELD(0, 5)
#define NPCX_UFRCTL_RFULL_LVL_EN          5
#define NPCX_UFRCTL_RNEMPTY_EN            6
#define NPCX_UFRCTL_ERR_EN                7

#endif

#if defined(CHIP_FAMILY_NPCX5)
enum {
	NPCX_UART_PORT0 = 0, /* UART port 0 */
	NPCX_UART_COUNT
};
#elif defined(CHIP_FAMILY_NPCX7)
enum {
	NPCX_UART_PORT0 = 0, /* UART port 0 */
#ifdef NPCX_SECOND_UART
	NPCX_UART_PORT1 = 1, /* UART port 1 */
#endif
	NPCX_UART_COUNT
};
#else
#error "Unsupported chip family for uart ports."
#endif

/******************************************************************************/
/* KBSCAN registers */
#define NPCX_KBSIN                        REG8(NPCX_KBSCAN_REGS_BASE + 0x04)
#define NPCX_KBSINPU                      REG8(NPCX_KBSCAN_REGS_BASE + 0x05)
#define NPCX_KBSOUT0                     REG16(NPCX_KBSCAN_REGS_BASE + 0x06)
#define NPCX_KBSOUT1                     REG16(NPCX_KBSCAN_REGS_BASE + 0x08)
#define NPCX_KBS_BUF_INDX                 REG8(NPCX_KBSCAN_REGS_BASE + 0x0A)
#define NPCX_KBS_BUF_DATA                 REG8(NPCX_KBSCAN_REGS_BASE + 0x0B)
#define NPCX_KBSEVT                       REG8(NPCX_KBSCAN_REGS_BASE + 0x0C)
#define NPCX_KBSCTL                       REG8(NPCX_KBSCAN_REGS_BASE + 0x0D)
#define NPCX_KBS_CFG_INDX                 REG8(NPCX_KBSCAN_REGS_BASE + 0x0E)
#define NPCX_KBS_CFG_DATA                 REG8(NPCX_KBSCAN_REGS_BASE + 0x0F)

/* KBSCAN register fields */
#define NPCX_KBSBUFINDX                  0
#define NPCX_KBSDONE                     0
#define NPCX_KBSERR                      1
#define NPCX_KBSSTART                    0
#define NPCX_KBSMODE                     1
#define NPCX_KBSIEN                      2
#define NPCX_KBSINC                      3
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_KBHDRV_FIELD                FIELD(6, 2)
#endif
#define NPCX_KBSCFGINDX                  0

/* KBSCAN definitions */
#define KB_ROW_NUM  8  /* Rows numbers of keyboard matrix */
#define KB_COL_NUM  18 /* Columns numbers of keyboard matrix */
#define KB_ROW_MASK ((1<<KB_ROW_NUM) - 1) /* Mask of rows of keyboard matrix */

/******************************************************************************/
/* GLUE registers */
#define NPCX_GLUE_SDPD0                   REG8(NPCX_GLUE_REGS_BASE + 0x010)
#define NPCX_GLUE_SDPD1                   REG8(NPCX_GLUE_REGS_BASE + 0x012)
#define NPCX_GLUE_SDP_CTS                 REG8(NPCX_GLUE_REGS_BASE + 0x014)
#define NPCX_GLUE_SMBSEL                  REG8(NPCX_GLUE_REGS_BASE + 0x021)
#if defined(NPCX_PSL_MODE_SUPPORT)
#define NPCX_GLUE_PSL_CTS                 REG8(NPCX_GLUE_REGS_BASE + 0x027)
#endif
/******************************************************************************/
/* MIWU registers */
#define NPCX_WKEDG_ADDR(port, n)         (NPCX_MIWU_BASE_ADDR(port) + 0x00 + \
					 ((n) * 2L) + ((n) < 5 ? 0 : 0x1E))
#define NPCX_WKAEDG_ADDR(port, n)        (NPCX_MIWU_BASE_ADDR(port) + 0x01 + \
					 ((n) * 2L) + ((n) < 5 ? 0 : 0x1E))
#define NPCX_WKPND_ADDR(port, n)         (NPCX_MIWU_BASE_ADDR(port) + 0x0A + \
					 ((n) * 4L) + ((n) < 5 ? 0 : 0x10))
#define NPCX_WKPCL_ADDR(port, n)         (NPCX_MIWU_BASE_ADDR(port) + 0x0C + \
					 ((n) * 4L) + ((n) < 5 ? 0 : 0x10))
#define NPCX_WKEN_ADDR(port, n)          (NPCX_MIWU_BASE_ADDR(port) + 0x1E + \
					 ((n) * 2L) + ((n) < 5 ? 0 : 0x12))
#define NPCX_WKINEN_ADDR(port, n)        (NPCX_MIWU_BASE_ADDR(port) + 0x1F + \
					 ((n) * 2L) + ((n) < 5 ? 0 : 0x12))
#define NPCX_WKMOD_ADDR(port, n)        (NPCX_MIWU_BASE_ADDR(port) + 0x70 + (n))

#define NPCX_WKEDG(port, n)               REG8(NPCX_WKEDG_ADDR(port, n))
#define NPCX_WKAEDG(port, n)              REG8(NPCX_WKAEDG_ADDR(port, n))
#define NPCX_WKPND(port, n)               REG8(NPCX_WKPND_ADDR(port, n))
#define NPCX_WKPCL(port, n)               REG8(NPCX_WKPCL_ADDR(port, n))
#define NPCX_WKEN(port, n)                REG8(NPCX_WKEN_ADDR(port, n))
#define NPCX_WKINEN(port, n)              REG8(NPCX_WKINEN_ADDR(port, n))
#define NPCX_WKMOD(port, n)               REG8(NPCX_WKMOD_ADDR(port, n))

/* MIWU enumeration */
enum {
	MIWU_TABLE_0,
	MIWU_TABLE_1,
	MIWU_TABLE_2,
	MIWU_TABLE_COUNT
};

enum {
	MIWU_GROUP_1,
	MIWU_GROUP_2,
	MIWU_GROUP_3,
	MIWU_GROUP_4,
	MIWU_GROUP_5,
	MIWU_GROUP_6,
	MIWU_GROUP_7,
	MIWU_GROUP_8,
	MIWU_GROUP_COUNT
};

enum {
	MIWU_EDGE_RISING,
	MIWU_EDGE_FALLING,
	MIWU_EDGE_ANYING,
};

/* MIWU utilities */
#define MIWU_TABLE_WKKEY MIWU_TABLE_1
#define MIWU_GROUP_WKKEY MIWU_GROUP_3

/******************************************************************************/
/* GPIO registers */
#define NPCX_PDOUT(n)                     REG8(NPCX_GPIO_BASE_ADDR(n) + 0x000)
#define NPCX_PDIN(n)                      REG8(NPCX_GPIO_BASE_ADDR(n) + 0x001)
#define NPCX_PDIR(n)                      REG8(NPCX_GPIO_BASE_ADDR(n) + 0x002)
#define NPCX_PPULL(n)                     REG8(NPCX_GPIO_BASE_ADDR(n) + 0x003)
#define NPCX_PPUD(n)                      REG8(NPCX_GPIO_BASE_ADDR(n) + 0x004)
#define NPCX_PENVDD(n)                    REG8(NPCX_GPIO_BASE_ADDR(n) + 0x005)
#define NPCX_PTYPE(n)                     REG8(NPCX_GPIO_BASE_ADDR(n) + 0x006)
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_PLOCK_CTL(n)                 REG8(NPCX_GPIO_BASE_ADDR(n) + 0x007)
#endif

/* GPIO enumeration */
enum {
	GPIO_PORT_0,
	GPIO_PORT_1,
	GPIO_PORT_2,
	GPIO_PORT_3,
	GPIO_PORT_4,
	GPIO_PORT_5,
	GPIO_PORT_6,
	GPIO_PORT_7,
	GPIO_PORT_8,
	GPIO_PORT_9,
	GPIO_PORT_A,
	GPIO_PORT_B,
	GPIO_PORT_C,
	GPIO_PORT_D,
	GPIO_PORT_E,
	GPIO_PORT_F,
	GPIO_PORT_COUNT
};

enum {
	MASK_PIN0 = BIT(0),
	MASK_PIN1 = BIT(1),
	MASK_PIN2 = BIT(2),
	MASK_PIN3 = BIT(3),
	MASK_PIN4 = BIT(4),
	MASK_PIN5 = BIT(5),
	MASK_PIN6 = BIT(6),
	MASK_PIN7 = BIT(7),
};

/* Chip-independent aliases for port base group */
#define GPIO_0 GPIO_PORT_0
#define GPIO_1 GPIO_PORT_1
#define GPIO_2 GPIO_PORT_2
#define GPIO_3 GPIO_PORT_3
#define GPIO_4 GPIO_PORT_4
#define GPIO_5 GPIO_PORT_5
#define GPIO_6 GPIO_PORT_6
#define GPIO_7 GPIO_PORT_7
#define GPIO_8 GPIO_PORT_8
#define GPIO_9 GPIO_PORT_9
#define GPIO_A GPIO_PORT_A
#define GPIO_B GPIO_PORT_B
#define GPIO_C GPIO_PORT_C
#define GPIO_D GPIO_PORT_D
#define GPIO_E GPIO_PORT_E
#define GPIO_F GPIO_PORT_F
#define DUMMY_GPIO_BANK GPIO_PORT_0

/******************************************************************************/
/* MSWC Registers */
#define NPCX_MSWCTL1                      REG8(NPCX_MSWC_BASE_ADDR + 0x000)
#define NPCX_MSWCTL2                      REG8(NPCX_MSWC_BASE_ADDR + 0x002)
#define NPCX_HCBAL                        REG8(NPCX_MSWC_BASE_ADDR + 0x008)
#define NPCX_HCBAH                        REG8(NPCX_MSWC_BASE_ADDR + 0x00A)
#define NPCX_SRID_CR                      REG8(NPCX_MSWC_BASE_ADDR + 0x01C)
#define NPCX_SID_CR                       REG8(NPCX_MSWC_BASE_ADDR + 0x020)
#define NPCX_DEVICE_ID_CR                 REG8(NPCX_MSWC_BASE_ADDR + 0x022)

/* MSWC register fields */
#define NPCX_MSWCTL1_HRSTOB              0
#define NPCS_MSWCTL1_HWPRON              1
#define NPCX_MSWCTL1_PLTRST_ACT          2
#define NPCX_MSWCTL1_VHCFGA              3
#define NPCX_MSWCTL1_HCFGLK              4
#define NPCX_MSWCTL1_PWROFFB             6
#define NPCX_MSWCTL1_A20MB               7

/******************************************************************************/
/* System Configuration (SCFG) Registers */
#define NPCX_DEVCNT                       REG8(NPCX_SCFG_BASE_ADDR + 0x000)
#define NPCX_STRPST                       REG8(NPCX_SCFG_BASE_ADDR + 0x001)
#define NPCX_RSTCTL                       REG8(NPCX_SCFG_BASE_ADDR + 0x002)
#define NPCX_DEV_CTL4                     REG8(NPCX_SCFG_BASE_ADDR + 0x006)
#define NPCX_DEVALT(n)                   REG8(NPCX_SCFG_BASE_ADDR + 0x010 + (n))
#define NPCX_LFCGCALCNT                   REG8(NPCX_SCFG_BASE_ADDR + 0x021)
#define NPCX_PUPD_EN0                     REG8(NPCX_SCFG_BASE_ADDR + 0x028)
#define NPCX_PUPD_EN1                     REG8(NPCX_SCFG_BASE_ADDR + 0x029)
#if defined(CHIP_FAMILY_NPCX5)
#define NPCX_LV_GPIO_CTL(n)              REG8(NPCX_SCFG_BASE_ADDR + 0x02A + (n))
#elif defined(CHIP_FAMILY_NPCX7)
#define NPCX_LV_GPIO_CTL_ADDR(n)          (((n) < 5) ? \
					  (NPCX_SCFG_BASE_ADDR + 0x02A + (n)) :\
					   (NPCX_SCFG_BASE_ADDR + 0x026))
#define NPCX_LV_GPIO_CTL(n)               REG8(NPCX_LV_GPIO_CTL_ADDR(n))
#endif
#define NPCX_SCFG_VER                     REG8(NPCX_SCFG_BASE_ADDR + 0x02F)

#define TEST_BKSL                         REG8(NPCX_SCFG_BASE_ADDR + 0x037)
#define TEST0                             REG8(NPCX_SCFG_BASE_ADDR + 0x038)
#define BLKSEL                           0

/* SCFG enumeration */
enum {
	ALT_GROUP_0,
	ALT_GROUP_1,
	ALT_GROUP_2,
	ALT_GROUP_3,
	ALT_GROUP_4,
	ALT_GROUP_5,
	ALT_GROUP_6,
	ALT_GROUP_7,
	ALT_GROUP_8,
	ALT_GROUP_9,
	ALT_GROUP_A,
	ALT_GROUP_B,
	ALT_GROUP_C,
	ALT_GROUP_D,
	ALT_GROUP_E,
	ALT_GROUP_F,
	ALT_GROUP_COUNT
};

/* SCFG register fields */
#define NPCX_DEVCNT_F_SPI_TRIS           6
#define NPCX_DEVCNT_HIF_TYP_SEL_FIELD    FIELD(2, 2)
#define NPCX_DEVCNT_JEN1_HEN             5
#define NPCX_DEVCNT_JEN0_HEN             4
#define NPCX_STRPST_TRIST                1
#define NPCX_STRPST_TEST                 2
#define NPCX_STRPST_JEN1                 4
#define NPCX_STRPST_JEN0                 5
#define NPCX_STRPST_SPI_COMP             7
#define NPCX_RSTCTL_VCC1_RST_STS         0
#define NPCX_RSTCTL_DBGRST_STS           1
#define NPCX_RSTCTL_VCC1_RST_SCRATCH     3
#define NPCX_RSTCTL_LRESET_PLTRST_MODE   5
#define NPCX_RSTCTL_HIPRST_MODE          6
#define NPCX_DEV_CTL4_F_SPI_SLLK         2
#define NPCX_DEV_CTL4_SPI_SP_SEL         4
#define NPCX_DEV_CTL4_WP_IF              5
#define NPCX_DEV_CTL4_VCC1_RST_LK        6
#define NPCX_DEVPU0_I2C0_0_PUE           0
#define NPCX_DEVPU0_I2C0_1_PUE           1
#define NPCX_DEVPU0_I2C1_0_PUE           2
#define NPCX_DEVPU0_I2C2_0_PUE           4
#define NPCX_DEVPU0_I2C3_0_PUE           6
#define NPCX_DEVPU1_F_SPI_PUD_EN         7

/* DEVALT */
/* pin-mux for SPI/FIU */
#define NPCX_DEVALT0_SPIP_SL             0
#define NPCX_DEVALT0_GPIO_NO_SPIP        3
#define NPCX_DEVALT0_F_SPI_CS1_2         4
#define NPCX_DEVALT0_F_SPI_CS1_1         5
#define NPCX_DEVALT0_F_SPI_QUAD          6
#define NPCX_DEVALT0_NO_F_SPI            7

/* pin-mux for LPC/eSPI */
#define NPCX_DEVALT1_KBRST_SL            0
#define NPCX_DEVALT1_A20M_SL             1
#define NPCX_DEVALT1_SMI_SL              2
#define NPCX_DEVALT1_EC_SCI_SL           3
#define NPCX_DEVALT1_NO_PWRGD            4
#define NPCX_DEVALT1_RST_OUT_SL          5
#define NPCX_DEVALT1_CLKRN_SL            6
#define NPCX_DEVALT1_NO_LPC_ESPI         7

/* pin-mux for I2C */
#if defined(CHIP_FAMILY_NPCX5)
#define NPCX_DEVALT2_I2C0_0_SL           0
#define NPCX_DEVALT2_I2C0_1_SL           1
#define NPCX_DEVALT2_I2C1_0_SL           2
#define NPCX_DEVALT2_I2C2_0_SL           4
#define NPCX_DEVALT2_I2C3_0_SL           6
#elif defined(CHIP_FAMILY_NPCX7)
#define NPCX_DEVALT2_I2C0_0_SL           0
#define NPCX_DEVALT2_I2C7_0_SL           1
#define NPCX_DEVALT2_I2C1_0_SL           2
#define NPCX_DEVALT2_I2C6_0_SL           3
#define NPCX_DEVALT2_I2C2_0_SL           4
#define NPCX_DEVALT2_I2C5_0_SL           5
#define NPCX_DEVALT2_I2C3_0_SL           6
#define NPCX_DEVALT2_I2C4_0_SL           7
#define NPCX_DEVALT6_I2C6_1_SL           5
#define NPCX_DEVALT6_I2C5_1_SL           6
#define NPCX_DEVALT6_I2C4_1_SL           7
#endif

/* pin-mux for PS2 */
#define NPCX_DEVALT3_PS2_0_SL            0
#define NPCX_DEVALT3_PS2_1_SL            1
#define NPCX_DEVALT3_PS2_2_SL            2
#define NPCX_DEVALT3_PS2_3_SL            3
#define NPCX_DEVALTC_PS2_3_SL2           3

/* pin-mux for Tacho */
#define NPCX_DEVALT3_TA1_SL1             4
#define NPCX_DEVALT3_TB1_SL1             5
#define NPCX_DEVALT3_TA2_SL1             6
#define NPCX_DEVALT3_TB2_SL1             7
#define NPCX_DEVALTC_TA1_SL2             4
#define NPCX_DEVALTC_TB1_SL2             5
#define NPCX_DEVALTC_TA2_SL2             6
#define NPCX_DEVALTC_TB2_SL2             7

/* pin-mux for PWM */
#define NPCX_DEVALT4_PWM0_SL             0
#define NPCX_DEVALT4_PWM1_SL             1
#define NPCX_DEVALT4_PWM2_SL             2
#define NPCX_DEVALT4_PWM3_SL             3
#define NPCX_DEVALT4_PWM4_SL             4
#define NPCX_DEVALT4_PWM5_SL             5
#define NPCX_DEVALT4_PWM6_SL             6
#define NPCX_DEVALT4_PWM7_SL             7

/* pin-mux for JTAG */
#define NPCX_DEVALT5_TRACE_EN            0
#define NPCX_DEVALT5_NJEN1_EN            1
#define NPCX_DEVALT5_NJEN0_EN            2

/* pin-mux for ADC */
#define NPCX_DEVALT6_ADC0_SL             0
#define NPCX_DEVALT6_ADC1_SL             1
#define NPCX_DEVALT6_ADC2_SL             2
#define NPCX_DEVALT6_ADC3_SL             3
#define NPCX_DEVALT6_ADC4_SL             4
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_DEVALTF_ADC5_SL             0
#define NPCX_DEVALTF_ADC6_SL             1
#define NPCX_DEVALTF_ADC7_SL             2
#define NPCX_DEVALTF_ADC8_SL             3
#define NPCX_DEVALTF_ADC9_SL             4
#endif

/* pin-mux for Keyboard */
#define NPCX_DEVALT7_NO_KSI0_SL          0
#define NPCX_DEVALT7_NO_KSI1_SL          1
#define NPCX_DEVALT7_NO_KSI2_SL          2
#define NPCX_DEVALT7_NO_KSI3_SL          3
#define NPCX_DEVALT7_NO_KSI4_SL          4
#define NPCX_DEVALT7_NO_KSI5_SL          5
#define NPCX_DEVALT7_NO_KSI6_SL          6
#define NPCX_DEVALT7_NO_KSI7_SL          7
#define NPCX_DEVALT8_NO_KSO00_SL         0
#define NPCX_DEVALT8_NO_KSO01_SL         1
#define NPCX_DEVALT8_NO_KSO02_SL         2
#define NPCX_DEVALT8_NO_KSO03_SL         3
#define NPCX_DEVALT8_NO_KSO04_SL         4
#define NPCX_DEVALT8_NO_KSO05_SL         5
#define NPCX_DEVALT8_NO_KSO06_SL         6
#define NPCX_DEVALT8_NO_KSO07_SL         7
#define NPCX_DEVALT9_NO_KSO08_SL         0
#define NPCX_DEVALT9_NO_KSO09_SL         1
#define NPCX_DEVALT9_NO_KSO10_SL         2
#define NPCX_DEVALT9_NO_KSO11_SL         3
#define NPCX_DEVALT9_NO_KSO12_SL         4
#define NPCX_DEVALT9_NO_KSO13_SL         5
#define NPCX_DEVALT9_NO_KSO14_SL         6
#define NPCX_DEVALT9_NO_KSO15_SL         7
#define NPCX_DEVALTA_NO_KSO16_SL         0
#define NPCX_DEVALTA_NO_KSO17_SL         1

/* pin-mux for PSL */
#if defined(NPCX_PSL_MODE_SUPPORT)
#define NPCX_DEVALTD_PSL_IN1_AHI         0
#define NPCX_DEVALTD_NPSL_IN1_SL         1
#define NPCX_DEVALTD_PSL_IN2_AHI         2
#define NPCX_DEVALTD_NPSL_IN2_SL         3
#define NPCX_DEVALTD_PSL_IN3_AHI         4
#define NPCX_DEVALTD_PSL_IN3_SL          5
#define NPCX_DEVALTD_PSL_IN4_AHI         6
#define NPCX_DEVALTD_PSL_IN4_SL          7
#endif

/* pin-mux for Others */
#define NPCX_DEVALTA_32K_OUT_SL          2
#if !defined(NPCX_EXT32K_OSC_SUPPORT)
#define NPCX_DEVALTA_32KCLKIN_SL         3
#endif
#define NPCX_DEVALTA_NO_VCC1_RST         4
#ifdef NPCX_SECOND_UART
#define NPCX_DEVALTA_UART2_SL            5
#endif
#define NPCX_DEVALTA_NO_PECI_EN          6
#define NPCX_DEVALTA_UART_SL1            7
#define NPCX_DEVALTC_UART_SL2            0
#define NPCX_DEVALTC_SHI_SL              1

#if defined(CHIP_FAMILY_NPCX7)
/* SHI module version 2 enable bit */
#define NPCX_DEVALTF_SHI_NEW             7
#endif

/* pin-mux for WoV */
#ifdef NPCX_WOV_SUPPORT
#define NPCX_DEVALTE_WOV_SL              0
#define NPCX_DEVALTE_I2S_SL              1
#define NPCX_DEVALTE_DMCLK_FAST          2
#endif

/* Others bit definitions */
#define NPCX_LFCGCALCNT_LPREG_CTL_EN     1

/******************************************************************************/
/* Development and Debug Support (DBG) Registers */
#define NPCX_DBGCTRL                      REG8(NPCX_SCFG_BASE_ADDR + 0x074)
#define NPCX_DBGFRZEN1                    REG8(NPCX_SCFG_BASE_ADDR + 0x076)
#define NPCX_DBGFRZEN2                    REG8(NPCX_SCFG_BASE_ADDR + 0x077)
#define NPCX_DBGFRZEN3                    REG8(NPCX_SCFG_BASE_ADDR + 0x078)
/* DBG register fields */
#define NPCX_DBGFRZEN3_GLBL_FRZ_DIS      7

/******************************************************************************/
/* SMBus Registers */
#define NPCX_SMBSDA(n)                    REG8(NPCX_SMB_BASE_ADDR(n) + 0x000)
#define NPCX_SMBST(n)                     REG8(NPCX_SMB_BASE_ADDR(n) + 0x002)
#define NPCX_SMBCST(n)                    REG8(NPCX_SMB_BASE_ADDR(n) + 0x004)
#define NPCX_SMBCTL1(n)                   REG8(NPCX_SMB_BASE_ADDR(n) + 0x006)
#define NPCX_SMBADDR1(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x008)
#define NPCX_SMBTMR_ST(n)                 REG8(NPCX_SMB_BASE_ADDR(n) + 0x009)
#define NPCX_SMBCTL2(n)                   REG8(NPCX_SMB_BASE_ADDR(n) + 0x00A)
#define NPCX_SMBTMR_EN(n)                 REG8(NPCX_SMB_BASE_ADDR(n) + 0x00B)
#define NPCX_SMBADDR2(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x00C)
#define NPCX_SMBCTL3(n)                   REG8(NPCX_SMB_BASE_ADDR(n) + 0x00E)
#define NPCX_SMBADDR3(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x010)
#define NPCX_SMBADDR7(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x011)
#define NPCX_SMBADDR4(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x012)
#define NPCX_SMBADDR8(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x013)
#define NPCX_SMBADDR5(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x014)
#define NPCX_SMBADDR6(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x016)
#define NPCX_SMBCST2(n)                   REG8(NPCX_SMB_BASE_ADDR(n) + 0x018)
#define NPCX_SMBCST3(n)                   REG8(NPCX_SMB_BASE_ADDR(n) + 0x019)
#define NPCX_SMBCTL4(n)                   REG8(NPCX_SMB_BASE_ADDR(n) + 0x01A)
#define NPCX_SMBSCLLT(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x01C)
#define NPCX_SMBSCLHT(n)                  REG8(NPCX_SMB_BASE_ADDR(n) + 0x01E)

/* SMBus register fields */
#define NPCX_SMBST_XMIT                  0
#define NPCX_SMBST_MASTER                1
#define NPCX_SMBST_NMATCH                2
#define NPCX_SMBST_STASTR                3
#define NPCX_SMBST_NEGACK                4
#define NPCX_SMBST_BER                   5
#define NPCX_SMBST_SDAST                 6
#define NPCX_SMBST_SLVSTP                7
#define NPCX_SMBCST_BUSY                 0
#define NPCX_SMBCST_BB                   1
#define NPCX_SMBCST_MATCH                2
#define NPCX_SMBCST_GCMATCH              3
#define NPCX_SMBCST_TSDA                 4
#define NPCX_SMBCST_TGSCL                5
#define NPCX_SMBCST_MATCHAF              6
#define NPCX_SMBCST_ARPMATCH             7
#define NPCX_SMBCST2_MATCHA1F            0
#define NPCX_SMBCST2_MATCHA2F            1
#define NPCX_SMBCST2_MATCHA3F            2
#define NPCX_SMBCST2_MATCHA4F            3
#define NPCX_SMBCST2_MATCHA5F            4
#define NPCX_SMBCST2_MATCHA6F            5
#define NPCX_SMBCST2_MATCHA7F            6
#define NPCX_SMBCST2_INTSTS              7
#define NPCX_SMBCST3_MATCHA8F            0
#define NPCX_SMBCST3_MATCHA9F            1
#define NPCX_SMBCST3_MATCHA10F           2
#define NPCX_SMBCTL1_START               0
#define NPCX_SMBCTL1_STOP                1
#define NPCX_SMBCTL1_INTEN               2
#define NPCX_SMBCTL1_ACK                 4
#define NPCX_SMBCTL1_GCMEN               5
#define NPCX_SMBCTL1_NMINTE              6
#define NPCX_SMBCTL1_STASTRE             7
#define NPCX_SMBCTL2_ENABLE              0
#define NPCX_SMBCTL2_SCLFRQ7_FIELD       FIELD(1, 7)
#define NPCX_SMBCTL3_ARPMEN              2
#define NPCX_SMBCTL3_SCLFRQ2_FIELD       FIELD(0, 2)
#define NPCX_SMBCTL3_IDL_START           3
#define NPCX_SMBCTL3_400K                4
#define NPCX_SMBCTL3_SDA_LVL             6
#define NPCX_SMBCTL3_SCL_LVL             7
#define NPCX_SMBCTL4_HLDT_FIELD          FIELD(0, 6)
#define NPCX_SMBADDR1_SAEN               7
#define NPCX_SMBADDR2_SAEN               7
#define NPCX_SMBADDR3_SAEN               7
#define NPCX_SMBADDR4_SAEN               7
#define NPCX_SMBADDR5_SAEN               7
#define NPCX_SMBADDR6_SAEN               7
#define NPCX_SMBADDR7_SAEN               7
#define NPCX_SMBADDR8_SAEN               7
#if defined(CHIP_FAMILY_NPCX5)
#define NPCX_SMBSEL_SMB0SEL              0
#elif defined(CHIP_FAMILY_NPCX7)
#define NPCX_SMBSEL_SMB4SEL              4
#define NPCX_SMBSEL_SMB5SEL              5
#define NPCX_SMBSEL_SMB6SEL              6
#endif
/*
 * SMB enumeration
 * I2C port definitions.
 */
#if defined(CHIP_FAMILY_NPCX5)
enum {
	NPCX_I2C_PORT0_0  = 0, /* I2C port 0, bus 0 */
	NPCX_I2C_PORT0_1,      /* I2C port 0, bus 1 */
	NPCX_I2C_PORT1,        /* I2C port 1 */
	NPCX_I2C_PORT2,        /* I2C port 2 */
	NPCX_I2C_PORT3,        /* I2C port 3 */
	NPCX_I2C_COUNT,
};
#elif defined(CHIP_FAMILY_NPCX7)
enum {
	NPCX_I2C_PORT0_0  = 0, /* I2C port 0, bus 0 */
	NPCX_I2C_PORT1_0,      /* I2C port 1, bus 0 */
	NPCX_I2C_PORT2_0,      /* I2C port 2, bus 0 */
	NPCX_I2C_PORT3_0,      /* I2C port 3, bus 0 */
#if !defined(NPCX_PSL_MODE_SUPPORT)
	NPCX_I2C_PORT4_0,      /* I2C port 4, bus 0 */
#endif
	NPCX_I2C_PORT4_1,      /* I2C port 4, bus 1 */
	NPCX_I2C_PORT5_0,      /* I2C port 5, bus 0 */
	NPCX_I2C_PORT5_1,      /* I2C port 5, bus 1 */
	NPCX_I2C_PORT6_0,      /* I2C port 6, bus 0 */
	NPCX_I2C_PORT6_1,      /* I2C port 6, bus 1 */
	NPCX_I2C_PORT7_0,      /* I2C port 7, bus 0 */
	NPCX_I2C_COUNT,
};
#else
#error "Unsupported chip family for i2c ports."
#endif

/******************************************************************************/
/* Power Management Controller (PMC) Registers */
#define NPCX_PMCSR                     REG8(NPCX_PMC_BASE_ADDR + 0x000)
#define NPCX_ENIDL_CTL                 REG8(NPCX_PMC_BASE_ADDR + 0x003)
#define NPCX_DISIDL_CTL                REG8(NPCX_PMC_BASE_ADDR + 0x004)
#define NPCX_DISIDL_CTL1               REG8(NPCX_PMC_BASE_ADDR + 0x005)
#define NPCX_PWDWN_CTL_ADDR(offset)    (((offset) < 6) ? \
			(NPCX_PMC_BASE_ADDR + 0x008 + (offset)) : \
			(NPCX_PMC_BASE_ADDR + 0x024))
#define NPCX_PWDWN_CTL(offset)         REG8(NPCX_PWDWN_CTL_ADDR(offset))
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_FMUL_WIN_DLY              REG8(NPCX_PMC_BASE_ADDR + 0x010)
#define NPCX_RAM_PD(offset)            REG8(NPCX_PMC_BASE_ADDR + 0x020 + offset)
#endif

/* PMC register fields */
#define NPCX_PMCSR_DI_INSTW              0
#define NPCX_PMCSR_DHF                   1
#define NPCX_PMCSR_IDLE                  2
#define NPCX_PMCSR_NWBI                  3
#define NPCX_PMCSR_OHFC                  6
#define NPCX_PMCSR_OLFC                  7
#define NPCX_DISIDL_CTL_RAM_DID          5
#define NPCX_ENIDL_CTL_ADC_LFSL          7
#define NPCX_ENIDL_CTL_LP_WK_CTL         6
#define NPCX_ENIDL_CTL_PECI_ENI          2
#define NPCX_ENIDL_CTL_ADC_ACC_DIS       1
#define NPCX_PWDWN_CTL1_KBS_PD           0
#define NPCX_PWDWN_CTL1_SDP_PD           1
#define NPCX_PWDWN_CTL1_FIU_PD           2
#define NPCX_PWDWN_CTL1_PS2_PD           3
#define NPCX_PWDWN_CTL1_UART_PD          4
#define NPCX_PWDWN_CTL1_MFT1_PD          5
#define NPCX_PWDWN_CTL1_MFT2_PD          6
#define NPCX_PWDWN_CTL1_MFT3_PD          7
#define NPCX_PWDWN_CTL2_PWM0_PD          0
#define NPCX_PWDWN_CTL2_PWM1_PD          1
#define NPCX_PWDWN_CTL2_PWM2_PD          2
#define NPCX_PWDWN_CTL2_PWM3_PD          3
#define NPCX_PWDWN_CTL2_PWM4_PD          4
#define NPCX_PWDWN_CTL2_PWM5_PD          5
#define NPCX_PWDWN_CTL2_PWM6_PD          6
#define NPCX_PWDWN_CTL2_PWM7_PD          7
#define NPCX_PWDWN_CTL3_SMB0_PD          0
#define NPCX_PWDWN_CTL3_SMB1_PD          1
#define NPCX_PWDWN_CTL3_SMB2_PD          2
#define NPCX_PWDWN_CTL3_SMB3_PD          3
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_PWDWN_CTL3_SMB4_PD          4
#endif
#define NPCX_PWDWN_CTL3_GMDA_PD          7
#define NPCX_PWDWN_CTL4_ITIM1_PD         0
#define NPCX_PWDWN_CTL4_ITIM2_PD         1
#define NPCX_PWDWN_CTL4_ITIM3_PD         2
#define NPCX_PWDWN_CTL4_ADC_PD           4
#define NPCX_PWDWN_CTL4_PECI_PD          5
#define NPCX_PWDWN_CTL4_PWM6_PD          6
#define NPCX_PWDWN_CTL4_SPIP_PD          7
#define NPCX_PWDWN_CTL5_SHI_PD           1
#define NPCX_PWDWN_CTL5_MRFSH_DIS        2
#define NPCX_PWDWN_CTL5_C2HACC_PD        3
#define NPCX_PWDWN_CTL5_SHM_REG_PD       4
#define NPCX_PWDWN_CTL5_SHM_PD           5
#define NPCX_PWDWN_CTL5_DP80_PD          6
#define NPCX_PWDWN_CTL5_MSWC_PD          7
#define NPCX_PWDWN_CTL6_ITIM4_PD         0
#define NPCX_PWDWN_CTL6_ITIM5_PD         1
#define NPCX_PWDWN_CTL6_ITIM6_PD         2
#define NPCX_PWDWN_CTL6_ESPI_PD          7
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_PWDWN_CTL7_SMB5_PD          0
#define NPCX_PWDWN_CTL7_SMB6_PD          1
#define NPCX_PWDWN_CTL7_SMB7_PD          2
#if defined(CHIP_VARIANT_NPCX7M6FB) || defined(CHIP_VARIANT_NPCX7M6FC) || \
	defined(CHIP_VARIANT_NPCX7M7WB) || defined(CHIP_VARIANT_NPCX7M7WC)
#define NPCX_PWDWN_CTL7_ITIM64_PD        5
#define NPCX_PWDWN_CTL7_UART2_PD         6
#endif
#if defined(CHIP_VARIANT_NPCX7M7WB) || defined(CHIP_VARIANT_NPCX7M7WC)
#define NPCX_PWDWN_CTL7_WOV_PD           7
#endif
#endif

/*
 * PMC enumeration
 * Offsets from CGC_BASE registers for each peripheral.
 */
enum {
	CGC_OFFSET_KBS    = 0,
	CGC_OFFSET_UART   = 0,
	CGC_OFFSET_FAN    = 0,
	CGC_OFFSET_FIU    = 0,
	CGC_OFFSET_PWM    = 1,
	CGC_OFFSET_I2C    = 2,
	CGC_OFFSET_ADC    = 3,
	CGC_OFFSET_PECI   = 3,
	CGC_OFFSET_SPI    = 3,
	CGC_OFFSET_TIMER  = 3,
	CGC_OFFSET_LPC    = 4,
	CGC_OFFSET_ESPI   = 5,
#if defined(CHIP_FAMILY_NPCX7)
	CGC_OFFSET_I2C2   = 6,
#ifdef NPCX_SECOND_UART
	CGC_OFFSET_UART2  = 6,
#endif
#ifdef NPCX_WOV_SUPPORT
	CGC_OFFSET_WOV    = 6,
#endif
#endif
};

enum NPCX_PMC_PWDWN_CTL_T {
	NPCX_PMC_PWDWN_1    = 0,
	NPCX_PMC_PWDWN_2    = 1,
	NPCX_PMC_PWDWN_3    = 2,
	NPCX_PMC_PWDWN_4    = 3,
	NPCX_PMC_PWDWN_5    = 4,
	NPCX_PMC_PWDWN_6    = 5,
#if defined(CHIP_FAMILY_NPCX7)
	NPCX_PMC_PWDWN_7    = 6,
#endif
	NPCX_PMC_PWDWN_CNT,
};

/* TODO: set PD masks based upon actual peripheral usage */
#define CGC_KBS_MASK     BIT(NPCX_PWDWN_CTL1_KBS_PD)
#define CGC_UART_MASK    BIT(NPCX_PWDWN_CTL1_UART_PD)
#define CGC_FAN_MASK     (BIT(NPCX_PWDWN_CTL1_MFT1_PD) | \
			 BIT(NPCX_PWDWN_CTL1_MFT2_PD))
#define CGC_FIU_MASK     BIT(NPCX_PWDWN_CTL1_FIU_PD)
#if defined(CHIP_FAMILY_NPCX5)
#define CGC_I2C_MASK     (BIT(NPCX_PWDWN_CTL3_SMB0_PD) | \
			 BIT(NPCX_PWDWN_CTL3_SMB1_PD) | \
			 BIT(NPCX_PWDWN_CTL3_SMB2_PD) | \
			 BIT(NPCX_PWDWN_CTL3_SMB3_PD))
#elif defined(CHIP_FAMILY_NPCX7)
#define CGC_I2C_MASK     (BIT(NPCX_PWDWN_CTL3_SMB0_PD) | \
			 BIT(NPCX_PWDWN_CTL3_SMB1_PD) | \
			 BIT(NPCX_PWDWN_CTL3_SMB2_PD) | \
			 BIT(NPCX_PWDWN_CTL3_SMB3_PD) | \
			 BIT(NPCX_PWDWN_CTL3_SMB4_PD))
#define CGC_I2C_MASK2    (BIT(NPCX_PWDWN_CTL7_SMB5_PD) | \
			 BIT(NPCX_PWDWN_CTL7_SMB6_PD) | \
			 BIT(NPCX_PWDWN_CTL7_SMB7_PD))
#ifdef NPCX_SECOND_UART
#define CGC_UART2_MASK   BIT(NPCX_PWDWN_CTL7_UART2_PD)
#endif
#ifdef NPCX_WOV_SUPPORT
#define CGC_WOV_MASK     BIT(NPCX_PWDWN_CTL7_WOV_PD)
#endif
#endif
#define CGC_ADC_MASK     BIT(NPCX_PWDWN_CTL4_ADC_PD)
#define CGC_PECI_MASK    BIT(NPCX_PWDWN_CTL4_PECI_PD)
#define CGC_SPI_MASK     BIT(NPCX_PWDWN_CTL4_SPIP_PD)
#define CGC_TIMER_MASK   (BIT(NPCX_PWDWN_CTL4_ITIM1_PD) | \
			 BIT(NPCX_PWDWN_CTL4_ITIM2_PD) | \
			 BIT(NPCX_PWDWN_CTL4_ITIM3_PD))
#define CGC_LPC_MASK     (BIT(NPCX_PWDWN_CTL5_C2HACC_PD) | \
			 BIT(NPCX_PWDWN_CTL5_SHM_REG_PD) | \
			 BIT(NPCX_PWDWN_CTL5_SHM_PD) | \
			 BIT(NPCX_PWDWN_CTL5_DP80_PD) | \
			 BIT(NPCX_PWDWN_CTL5_MSWC_PD))
#define CGC_ESPI_MASK    BIT(NPCX_PWDWN_CTL6_ESPI_PD)

/******************************************************************************/
/* Flash Interface Unit (FIU) Registers */
#define NPCX_FIU_CFG                      REG8(NPCX_FIU_BASE_ADDR + 0x000)
#define NPCX_BURST_CFG                    REG8(NPCX_FIU_BASE_ADDR + 0x001)
#define NPCX_RESP_CFG                     REG8(NPCX_FIU_BASE_ADDR + 0x002)
#define NPCX_SPI_FL_CFG                   REG8(NPCX_FIU_BASE_ADDR + 0x014)
#define NPCX_UMA_CODE                     REG8(NPCX_FIU_BASE_ADDR + 0x016)
#define NPCX_UMA_AB0                      REG8(NPCX_FIU_BASE_ADDR + 0x017)
#define NPCX_UMA_AB1                      REG8(NPCX_FIU_BASE_ADDR + 0x018)
#define NPCX_UMA_AB2                      REG8(NPCX_FIU_BASE_ADDR + 0x019)
#define NPCX_UMA_DB0                      REG8(NPCX_FIU_BASE_ADDR + 0x01A)
#define NPCX_UMA_DB1                      REG8(NPCX_FIU_BASE_ADDR + 0x01B)
#define NPCX_UMA_DB2                      REG8(NPCX_FIU_BASE_ADDR + 0x01C)
#define NPCX_UMA_DB3                      REG8(NPCX_FIU_BASE_ADDR + 0x01D)
#define NPCX_UMA_CTS                      REG8(NPCX_FIU_BASE_ADDR + 0x01E)
#define NPCX_UMA_ECTS                     REG8(NPCX_FIU_BASE_ADDR + 0x01F)
#define NPCX_UMA_DB0_3                   REG32(NPCX_FIU_BASE_ADDR + 0x020)
#define NPCX_FIU_RD_CMD                   REG8(NPCX_FIU_BASE_ADDR + 0x030)
#define NPCX_FIU_DMM_CYC                  REG8(NPCX_FIU_BASE_ADDR + 0x032)
#define NPCX_FIU_EXT_CFG                  REG8(NPCX_FIU_BASE_ADDR + 0x033)
#define NPCX_FIU_UMA_AB0_3               REG32(NPCX_FIU_BASE_ADDR + 0x034)

/* FIU register fields */
#define NPCX_RESP_CFG_IAD_EN             0
#define NPCX_RESP_CFG_DEV_SIZE_EX        2
#define NPCX_UMA_CTS_A_SIZE              3
#define NPCX_UMA_CTS_C_SIZE              4
#define NPCX_UMA_CTS_RD_WR               5
#define NPCX_UMA_CTS_DEV_NUM             6
#define NPCX_UMA_CTS_EXEC_DONE           7
#define NPCX_UMA_ECTS_SW_CS0             0
#define NPCX_UMA_ECTS_SW_CS1             1
#define NPCX_UMA_ECTS_SEC_CS             2
#define NPCX_UMA_ECTS_UMA_LOCK           3

/******************************************************************************/
/* Shared Memory (SHM) Registers */
#define NPCX_SMC_STS                REG8(NPCX_SHM_BASE_ADDR + 0x000)
#define NPCX_SMC_CTL                REG8(NPCX_SHM_BASE_ADDR + 0x001)
#define NPCX_SHM_CTL                REG8(NPCX_SHM_BASE_ADDR + 0x002)
#define NPCX_IMA_WIN_SIZE           REG8(NPCX_SHM_BASE_ADDR + 0x005)
#define NPCX_WIN_SIZE               REG8(NPCX_SHM_BASE_ADDR + 0x007)
#define NPCX_SHAW_SEM(win)          REG8(NPCX_SHM_BASE_ADDR + 0x008 + (win))
#define NPCX_IMA_SEM                REG8(NPCX_SHM_BASE_ADDR + 0x00B)
#define NPCX_SHCFG                  REG8(NPCX_SHM_BASE_ADDR + 0x00E)
#define NPCX_WIN_WR_PROT(win)       REG8(NPCX_SHM_BASE_ADDR + 0x010 + (win*2L))
#define NPCX_WIN_RD_PROT(win)       REG8(NPCX_SHM_BASE_ADDR + 0x011 + (win*2L))
#define NPCX_IMA_WR_PROT            REG8(NPCX_SHM_BASE_ADDR + 0x016)
#define NPCX_IMA_RD_PROT            REG8(NPCX_SHM_BASE_ADDR + 0x017)
#define NPCX_WIN_BASE(win)         REG32(NPCX_SHM_BASE_ADDR + 0x020 + (win*4L))

#define NPCX_PWIN_BASEI(win)       REG16(NPCX_SHM_BASE_ADDR + 0x020 + (win*4L))
#define NPCX_PWIN_SIZEI(win)       REG16(NPCX_SHM_BASE_ADDR + 0x022 + (win*4L))

#define NPCX_IMA_BASE              REG32(NPCX_SHM_BASE_ADDR + 0x02C)
#define NPCX_RST_CFG                REG8(NPCX_SHM_BASE_ADDR + 0x03A)
#define NPCX_DP80BUF               REG16(NPCX_SHM_BASE_ADDR + 0x040)
#define NPCX_DP80STS                REG8(NPCX_SHM_BASE_ADDR + 0x042)
#define NPCX_DP80CTL                REG8(NPCX_SHM_BASE_ADDR + 0x044)
#define NPCX_HOFS_STS               REG8(NPCX_SHM_BASE_ADDR + 0x048)
#define NPCX_HOFS_CTL               REG8(NPCX_SHM_BASE_ADDR + 0x049)
#define NPCX_COFS2                 REG16(NPCX_SHM_BASE_ADDR + 0x04A)
#define NPCX_COFS1                 REG16(NPCX_SHM_BASE_ADDR + 0x04C)
#define NPCX_IHOFS2                REG16(NPCX_SHM_BASE_ADDR + 0x050)
#define NPCX_IHOFS1                REG16(NPCX_SHM_BASE_ADDR + 0x052)
#define NPCX_SHM_VER                REG8(NPCX_SHM_BASE_ADDR + 0x07F)

/* SHM register fields */
#define NPCX_SMC_STS_HRERR               0
#define NPCX_SMC_STS_HWERR               1
#define NPCX_SMC_STS_HSEM1W              4
#define NPCX_SMC_STS_HSEM2W              5
#define NPCX_SMC_STS_SHM_ACC             6
#define NPCX_SMC_CTL_HERR_IE             2
#define NPCX_SMC_CTL_HSEM1_IE            3
#define NPCX_SMC_CTL_HSEM2_IE            4
#define NPCX_SMC_CTL_ACC_IE              5
#define NPCX_SMC_CTL_PREF_EN             6
#define NPCX_SMC_CTL_HOSTWAIT            7
#define NPCX_FLASH_SIZE_STALL_HOST       6
#define NPCX_FLASH_SIZE_RD_BURST         7
#define NPCX_WIN_PROT_RW1L_RP            0
#define NPCX_WIN_PROT_RW1L_WP            1
#define NPCX_WIN_PROT_RW1H_RP            2
#define NPCX_WIN_PROT_RW1H_WP            3
#define NPCX_WIN_PROT_RW2L_RP            4
#define NPCX_WIN_PROT_RW2L_WP            5
#define NPCX_WIN_PROT_RW2H_RP            6
#define NPCX_WIN_PROT_RW2H_WP            7
#define NPCX_PWIN_SIZEI_RPROT            13
#define NPCX_PWIN_SIZEI_WPROT            14
#define NPCX_CSEM2                       6
#define NPCX_CSEM3                       7
#define NPCX_DP80STS_FWR                 5
#define NPCX_DP80STS_FNE                 6
#define NPCX_DP80STS_FOR                 7
#define NPCX_DP80CTL_DP80EN              0
#define NPCX_DP80CTL_SYNCEN              1
#define NPCX_DP80CTL_RFIFO               4
#define NPCX_DP80CTL_CIEN                5

/******************************************************************************/
/* KBC Registers */
#define NPCX_HICTRL                       REG8(NPCX_KBC_BASE_ADDR + 0x000)
#define NPCX_HIIRQC                       REG8(NPCX_KBC_BASE_ADDR + 0x002)
#define NPCX_HIKMST                       REG8(NPCX_KBC_BASE_ADDR + 0x004)
#define NPCX_HIKDO                        REG8(NPCX_KBC_BASE_ADDR + 0x006)
#define NPCX_HIMDO                        REG8(NPCX_KBC_BASE_ADDR + 0x008)
#define NPCX_KBCVER                       REG8(NPCX_KBC_BASE_ADDR + 0x009)
#define NPCX_HIKMDI                       REG8(NPCX_KBC_BASE_ADDR + 0x00A)
#define NPCX_SHIKMDI                      REG8(NPCX_KBC_BASE_ADDR + 0x00B)

/* KBC register field */
#define NPCX_HICTRL_OBFKIE               0 /* Automatic Serial IRQ1 for KBC */
#define NPCX_HICTRL_OBFMIE               1 /* Automatic Serial IRQ12 for Mouse*/
#define NPCX_HICTRL_OBECIE               2 /* KBC OBE interrupt enable */
#define NPCX_HICTRL_IBFCIE               3 /* KBC IBF interrupt enable */
#define NPCX_HICTRL_PMIHIE               4 /* Automatic Serial IRQ11 for PMC1 */
#define NPCX_HICTRL_PMIOCIE              5 /* PMC1 OBE interrupt enable */
#define NPCX_HICTRL_PMICIE               6 /* PMC1 IBF interrupt enable */
#define NPCX_HICTRL_FW_OBF               7 /* Firmware control over OBF */

#define NPCX_HIKMST_OBF                  0 /* KB output buffer is full */
/******************************************************************************/
/* PM Channel Registers */
#define NPCX_HIPMST(n)                    REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x000)
#define NPCX_HIPMDO(n)                    REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x002)
#define NPCX_HIPMDI(n)                    REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x004)
#define NPCX_SHIPMDI(n)                   REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x005)
#define NPCX_HIPMDOC(n)                   REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x006)
#define NPCX_HIPMDOM(n)                   REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x008)
#define NPCX_HIPMDIC(n)                   REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x00A)
#define NPCX_HIPMCTL(n)                   REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x00C)
#define NPCX_HIPMCTL2(n)                  REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x00D)
#define NPCX_HIPMIC(n)                    REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x00E)
#define NPCX_HIPMIE(n)                    REG8(NPCX_PM_CH_BASE_ADDR(n) + 0x010)

/* PM Channel register field */

/* NPCX_HIPMIE */
#define NPCX_HIPMIE_SCIE                 1
#define NPCX_HIPMIE_SMIE                 2

/* NPCX_HIPMCTL */
#define NPCX_HIPMCTL_IBFIE               0
#define NPCX_HIPMCTL_SCIPOL              6

/* NPCX_HIPMST */
#define NPCX_HIPMST_F0                   2 /* EC_LPC_CMDR_BUSY */
#define NPCX_HIPMST_ST0                  4 /* EC_LPC_CMDR_ACPI_BRST */
#define NPCX_HIPMST_ST1                  5 /* EC_LPC_CMDR_SCI */
#define NPCX_HIPMST_ST2                  6 /* EC_LPC_CMDR_SMI */

/* NPCX_HIPMIC */
#define NPCX_HIPMIC_SMIB                 1
#define NPCX_HIPMIC_SCIB                 2
#define NPCX_HIPMIC_SMIPOL               6

/*
 * PM Channel enumeration
 */
enum PM_CHANNEL_T {
	PM_CHAN_1,
	PM_CHAN_2,
	PM_CHAN_3,
	PM_CHAN_4
};

/******************************************************************************/
/* SuperI/O Internal Bus (SIB) Registers */
#define NPCX_IHIOA                       REG16(NPCX_SIB_BASE_ADDR + 0x000)
#define NPCX_IHD                          REG8(NPCX_SIB_BASE_ADDR + 0x002)
#define NPCX_LKSIOHA                     REG16(NPCX_SIB_BASE_ADDR + 0x004)
#define NPCX_SIOLV                       REG16(NPCX_SIB_BASE_ADDR + 0x006)
#define NPCX_CRSMAE                      REG16(NPCX_SIB_BASE_ADDR + 0x008)
#define NPCX_SIBCTRL                      REG8(NPCX_SIB_BASE_ADDR + 0x00A)
#define NPCX_C2H_VER                      REG8(NPCX_SIB_BASE_ADDR + 0x00E)
/* SIB register fields  */
#define NPCX_SIBCTRL_CSAE                0
#define NPCX_SIBCTRL_CSRD                1
#define NPCX_SIBCTRL_CSWR                2
#define NPCX_LKSIOHA_LKCFG               0
#define NPCX_LKSIOHA_LKHIKBD            11
#define NPCX_CRSMAE_CFGAE                0
#define NPCX_CRSMAE_HIKBDAE             11

/******************************************************************************/
/* Battery-Backed RAM (BBRAM) Registers */
#define NPCX_BKUP_STS                REG8(NPCX_BBRAM_BASE_ADDR + 0x100)
#define NPCX_BBRAM(offset)           REG8(NPCX_BBRAM_BASE_ADDR + offset)

/* BBRAM register fields */
#define NPCX_BKUP_STS_IBBR               7
#if defined(CHIP_VARIANT_NPCX7M6FB) || defined(CHIP_VARIANT_NPCX7M6FC) || \
	defined(CHIP_VARIANT_NPCX7M7WB) || defined(CHIP_VARIANT_NPCX7M7WC)
#define NPCX_BKUP_STS_VSBY_STS           1
#define NPCX_BKUP_STS_VCC1_STS           0
#define NPCX_BKUP_STS_ALL_MASK \
	(BIT(NPCX_BKUP_STS_IBBR) | BIT(NPCX_BKUP_STS_VSBY_STS) | \
	BIT(NPCX_BKUP_STS_VCC1_STS))
#define NPCX_BBRAM_SIZE                 128  /* Size of BBRAM */
#else
#define NPCX_BKUP_STS_ALL_MASK BIT(NPCX_BKUP_STS_IBBR)
#define NPCX_BBRAM_SIZE                  64  /* Size of BBRAM */
#endif

/******************************************************************************/
/* Timer Watch Dog (TWD) Registers */
#define NPCX_TWCFG                        REG8(NPCX_TWD_BASE_ADDR + 0x000)
#define NPCX_TWCP                         REG8(NPCX_TWD_BASE_ADDR + 0x002)
#define NPCX_TWDT0                       REG16(NPCX_TWD_BASE_ADDR + 0x004)
#define NPCX_T0CSR                        REG8(NPCX_TWD_BASE_ADDR + 0x006)
#define NPCX_WDCNT                        REG8(NPCX_TWD_BASE_ADDR + 0x008)
#define NPCX_WDSDM                        REG8(NPCX_TWD_BASE_ADDR + 0x00A)
#define NPCX_TWMT0                       REG16(NPCX_TWD_BASE_ADDR + 0x00C)
#define NPCX_TWMWD                        REG8(NPCX_TWD_BASE_ADDR + 0x00E)
#define NPCX_WDCP                         REG8(NPCX_TWD_BASE_ADDR + 0x010)

/* TWD register fields */
#define NPCX_TWCFG_LTWCFG                0
#define NPCX_TWCFG_LTWCP                 1
#define NPCX_TWCFG_LTWDT0                2
#define NPCX_TWCFG_LWDCNT                3
#define NPCX_TWCFG_WDCT0I                4
#define NPCX_TWCFG_WDSDME                5
#define NPCX_TWCFG_WDRST_MODE            6
#define NPCX_TWCFG_WDC2POR               7
#define NPCX_T0CSR_RST                   0
#define NPCX_T0CSR_TC                    1
#define NPCX_T0CSR_WDLTD                 3
#define NPCX_T0CSR_WDRST_STS             4
#define NPCX_T0CSR_WD_RUN                5
#define NPCX_T0CSR_TESDIS                7

/******************************************************************************/
/* ADC Registers */
#define NPCX_ADCSTS                 REG16(NPCX_ADC_BASE_ADDR + 0x000)
#define NPCX_ADCCNF                 REG16(NPCX_ADC_BASE_ADDR + 0x002)
#define NPCX_ATCTL                  REG16(NPCX_ADC_BASE_ADDR + 0x004)
#define NPCX_ASCADD                 REG16(NPCX_ADC_BASE_ADDR + 0x006)
#define NPCX_ADCCS                  REG16(NPCX_ADC_BASE_ADDR + 0x008)
#define NPCX_CHNDAT(n)              REG16(NPCX_ADC_BASE_ADDR + 0x040 + (2L*(n)))
#define NPCX_ADCCNF2                REG16(NPCX_ADC_BASE_ADDR + 0x020)
#define NPCX_GENDLY                 REG16(NPCX_ADC_BASE_ADDR + 0x022)
#define NPCX_MEAST                  REG16(NPCX_ADC_BASE_ADDR + 0x026)

/* ADC register fields */
#define NPCX_ATCTL_SCLKDIV_FIELD         FIELD(0, 6)
#define NPCX_ATCTL_DLY_FIELD             FIELD(8, 3)
#define NPCX_ASCADD_SADDR_FIELD          FIELD(0, 5)
#define NPCX_ADCSTS_EOCEV                0
#define NPCX_ADCCNF_ADCMD_FIELD          FIELD(1, 2)
#define NPCX_ADCCNF_ADCRPTC              3
#define NPCX_ADCCNF_INTECEN              6
#define NPCX_ADCCNF_START                4
#define NPCX_ADCCNF_ADCEN                0
#define NPCX_ADCCNF_STOP                 11
#define NPCX_CHNDAT_CHDAT_FIELD          FIELD(0, 10)
#define NPCX_CHNDAT_NEW                  15
/******************************************************************************/
/* SPI Register */
#define NPCX_SPI_DATA                    REG16(NPCX_SPI_BASE_ADDR + 0x00)
#define NPCX_SPI_CTL1                    REG16(NPCX_SPI_BASE_ADDR + 0x02)
#define NPCX_SPI_STAT                     REG8(NPCX_SPI_BASE_ADDR + 0x04)

/* SPI register fields */
#define NPCX_SPI_CTL1_SPIEN              0
#define NPCX_SPI_CTL1_SNM                1
#define NPCX_SPI_CTL1_MOD                2
#define NPCX_SPI_CTL1_EIR                5
#define NPCX_SPI_CTL1_EIW                6
#define NPCX_SPI_CTL1_SCM                7
#define NPCX_SPI_CTL1_SCIDL              8
#define NPCX_SPI_CTL1_SCDV               9
#define NPCX_SPI_STAT_BSY                0
#define NPCX_SPI_STAT_RBF                1

/******************************************************************************/
/* PECI Registers */

#define NPCX_PECI_CTL_STS                REG8(NPCX_PECI_BASE_ADDR + 0x000)
#define NPCX_PECI_RD_LENGTH              REG8(NPCX_PECI_BASE_ADDR + 0x001)
#define NPCX_PECI_ADDR                   REG8(NPCX_PECI_BASE_ADDR + 0x002)
#define NPCX_PECI_CMD                    REG8(NPCX_PECI_BASE_ADDR + 0x003)
#define NPCX_PECI_CTL2                   REG8(NPCX_PECI_BASE_ADDR + 0x004)
#define NPCX_PECI_INDEX                  REG8(NPCX_PECI_BASE_ADDR + 0x005)
#define NPCX_PECI_IDATA                  REG8(NPCX_PECI_BASE_ADDR + 0x006)
#define NPCX_PECI_WR_LENGTH              REG8(NPCX_PECI_BASE_ADDR + 0x007)
#define NPCX_PECI_CFG                    REG8(NPCX_PECI_BASE_ADDR + 0x009)
#define NPCX_PECI_RATE                   REG8(NPCX_PECI_BASE_ADDR + 0x00F)
#define NPCX_PECI_DATA_IN(i)             REG8(NPCX_PECI_BASE_ADDR + 0x010 + (i))
#define NPCX_PECI_DATA_OUT(i)            REG8(NPCX_PECI_BASE_ADDR + 0x010 + (i))

/* PECI register fields */
#define NPCX_PECI_CTL_STS_START_BUSY     0
#define NPCX_PECI_CTL_STS_DONE           1
#define NPCX_PECI_CTL_STS_AVL_ERR        2
#define NPCX_PECI_CTL_STS_CRC_ERR        3
#define NPCX_PECI_CTL_STS_ABRT_ERR       4
#define NPCX_PECI_CTL_STS_AWFCS_EN       5
#define NPCX_PECI_CTL_STS_DONE_EN        6
#define NPCX_ESTRPST_PECIST              0
#define SFT_STRP_CFG_CK50                5

/******************************************************************************/
/* PWM Registers */
#define NPCX_PRSC(n)                     REG16(NPCX_PWM_BASE_ADDR(n) + 0x000)
#define NPCX_CTR(n)                      REG16(NPCX_PWM_BASE_ADDR(n) + 0x002)
#define NPCX_PWMCTL(n)                    REG8(NPCX_PWM_BASE_ADDR(n) + 0x004)
#define NPCX_DCR(n)                      REG16(NPCX_PWM_BASE_ADDR(n) + 0x006)
#define NPCX_PWMCTLEX(n)                  REG8(NPCX_PWM_BASE_ADDR(n) + 0x00C)

/* PWM register fields */
#define NPCX_PWMCTL_INVP                 0
#define NPCX_PWMCTL_CKSEL                1
#define NPCX_PWMCTL_HB_DC_CTL_FIELD      FIELD(2, 2)
#define NPCX_PWMCTL_PWR                  7
#define NPCX_PWMCTLEX_FCK_SEL_FIELD      FIELD(4, 2)
#define NPCX_PWMCTLEX_OD_OUT             7
/******************************************************************************/
/* MFT Registers */
#define NPCX_TCNT1(n)                    REG16(NPCX_MFT_BASE_ADDR(n) + 0x000)
#define NPCX_TCRA(n)                     REG16(NPCX_MFT_BASE_ADDR(n) + 0x002)
#define NPCX_TCRB(n)                     REG16(NPCX_MFT_BASE_ADDR(n) + 0x004)
#define NPCX_TCNT2(n)                    REG16(NPCX_MFT_BASE_ADDR(n) + 0x006)
#define NPCX_TPRSC(n)                     REG8(NPCX_MFT_BASE_ADDR(n) + 0x008)
#define NPCX_TCKC(n)                      REG8(NPCX_MFT_BASE_ADDR(n) + 0x00A)
#define NPCX_TMCTRL(n)                    REG8(NPCX_MFT_BASE_ADDR(n) + 0x00C)
#define NPCX_TECTRL(n)                    REG8(NPCX_MFT_BASE_ADDR(n) + 0x00E)
#define NPCX_TECLR(n)                     REG8(NPCX_MFT_BASE_ADDR(n) + 0x010)
#define NPCX_TIEN(n)                      REG8(NPCX_MFT_BASE_ADDR(n) + 0x012)
#define NPCX_TWUEN(n)                     REG8(NPCX_MFT_BASE_ADDR(n) + 0x01A)
#define NPCX_TCFG(n)                      REG8(NPCX_MFT_BASE_ADDR(n) + 0x01C)

/* MFT register fields */
#define NPCX_TMCTRL_MDSEL_FIELD          FIELD(0, 3)
#define NPCX_TCKC_LOW_PWR                7
#define NPCX_TCKC_PLS_ACC_CLK            6
#define NPCX_TCKC_C1CSEL_FIELD           FIELD(0, 3)
#define NPCX_TCKC_C2CSEL_FIELD           FIELD(3, 3)
#define NPCX_TMCTRL_TAEN                 5
#define NPCX_TMCTRL_TBEN                 6
#define NPCX_TMCTRL_TAEDG                3
#define NPCX_TMCTRL_TBEDG                4
#define NPCX_TCFG_TADBEN                 6
#define NPCX_TCFG_TBDBEN                 7
#define NPCX_TECTRL_TAPND                0
#define NPCX_TECTRL_TBPND                1
#define NPCX_TECTRL_TCPND                2
#define NPCX_TECTRL_TDPND                3
#define NPCX_TECLR_TACLR                 0
#define NPCX_TECLR_TBCLR                 1
#define NPCX_TECLR_TCCLR                 2
#define NPCX_TECLR_TDCLR                 3
#define NPCX_TIEN_TAIEN                  0
#define NPCX_TIEN_TBIEN                  1
#define NPCX_TIEN_TCIEN                  2
#define NPCX_TIEN_TDIEN                  3
#define NPCX_TWUEN_TAWEN                 0
#define NPCX_TWUEN_TBWEN                 1
#define NPCX_TWUEN_TCWEN                 2
#define NPCX_TWUEN_TDWEN                 3
/******************************************************************************/
/* ITIM16/32 Define */
#define ITIM16_INT(module)               CONCAT2(NPCX_IRQ_, module)

/* ITIM16 registers */
#define NPCX_ITCNT(n)                     REG8(NPCX_ITIM16_BASE_ADDR(n) + 0x000)
#define NPCX_ITPRE(n)                     REG8(NPCX_ITIM16_BASE_ADDR(n) + 0x001)
#define NPCX_ITCNT16(n)                  REG16(NPCX_ITIM16_BASE_ADDR(n) + 0x002)
#define NPCX_ITCTS(n)                     REG8(NPCX_ITIM16_BASE_ADDR(n) + 0x004)

/* ITIM32 registers */
#define NPCX_ITCNT32                    REG32(NPCX_ITIM32_BASE_ADDR + 0x008)

/* ITIM16 register fields */
#define NPCX_ITCTS_TO_STS                0
#define NPCX_ITCTS_TO_IE                 2
#define NPCX_ITCTS_TO_WUE                3
#define NPCX_ITCTS_CKSEL                 4
#define NPCX_ITCTS_ITEN                  7

/* ITIM16 enumeration*/
enum ITIM16_MODULE_T {
	ITIM16_1,
	ITIM16_2,
	ITIM16_3,
	ITIM16_4,
	ITIM16_5,
	ITIM16_6,
	ITIM32,
	ITIM_MODULE_COUNT,
};

/******************************************************************************/
/* Serial Host Interface (SHI) Registers */
#define NPCX_SHICFG1                      REG8(NPCX_SHI_BASE_ADDR + 0x001)
#define NPCX_SHICFG2                      REG8(NPCX_SHI_BASE_ADDR + 0x002)
#define NPCX_I2CADDR1                     REG8(NPCX_SHI_BASE_ADDR + 0x003)
#define NPCX_I2CADDR2                     REG8(NPCX_SHI_BASE_ADDR + 0x004)
#define NPCX_EVENABLE                     REG8(NPCX_SHI_BASE_ADDR + 0x005)
#define NPCX_EVSTAT                       REG8(NPCX_SHI_BASE_ADDR + 0x006)
#define NPCX_SHI_CAPABILITY               REG8(NPCX_SHI_BASE_ADDR + 0x007)
#define NPCX_STATUS                       REG8(NPCX_SHI_BASE_ADDR + 0x008)
#define NPCX_IBUFSTAT                     REG8(NPCX_SHI_BASE_ADDR + 0x00A)
#define NPCX_OBUFSTAT                     REG8(NPCX_SHI_BASE_ADDR + 0x00B)
#if defined(CHIP_FAMILY_NPCX5)
#define NPCX_OBUF(n)                      REG8(NPCX_SHI_BASE_ADDR + 0x020 + (n))
#define NPCX_IBUF(n)                      REG8(NPCX_SHI_BASE_ADDR + 0x060 + (n))
#elif defined(CHIP_FAMILY_NPCX7)
/* Serial Host Interface (SHI) Registers - only available on SHI Version 2 */
#define NPCX_SHICFG3                      REG8(NPCX_SHI_BASE_ADDR + 0x00C)
#define NPCX_SHICFG4                      REG8(NPCX_SHI_BASE_ADDR + 0x00D)
#define NPCX_SHICFG5                      REG8(NPCX_SHI_BASE_ADDR + 0x00E)
#define NPCX_EVSTAT2                      REG8(NPCX_SHI_BASE_ADDR + 0x00F)
#define NPCX_EVENABLE2                    REG8(NPCX_SHI_BASE_ADDR + 0x010)
#define NPCX_OBUF(n)                      REG8(NPCX_SHI_BASE_ADDR + 0x020 + (n))
#define NPCX_IBUF(n)                      REG8(NPCX_SHI_BASE_ADDR + 0x0A0 + (n))
#endif

/* SHI register fields */
#define NPCX_SHICFG1_EN                  0
#define NPCX_SHICFG1_MODE                1
#define NPCX_SHICFG1_WEN                 2
#define NPCX_SHICFG1_AUTIBF              3
#define NPCX_SHICFG1_AUTOBE              4
#define NPCX_SHICFG1_DAS                 5
#define NPCX_SHICFG1_CPOL                6
#define NPCX_SHICFG1_IWRAP               7
#define NPCX_SHICFG2_SIMUL               0
#define NPCX_SHICFG2_BUSY                1
#define NPCX_SHICFG2_ONESHOT             2
#define NPCX_SHICFG2_SLWU                3
#define NPCX_SHICFG2_REEN                4
#define NPCX_SHICFG2_RESTART             5
#define NPCX_SHICFG2_REEVEN              6
#define NPCX_EVENABLE_OBEEN              0
#define NPCX_EVENABLE_OBHEEN             1
#define NPCX_EVENABLE_IBFEN              2
#define NPCX_EVENABLE_IBHFEN             3
#define NPCX_EVENABLE_EOREN              4
#define NPCX_EVENABLE_EOWEN              5
#define NPCX_EVENABLE_STSREN             6
#define NPCX_EVENABLE_IBOREN             7
#define NPCX_EVSTAT_OBE                  0
#define NPCX_EVSTAT_OBHE                 1
#define NPCX_EVSTAT_IBF                  2
#define NPCX_EVSTAT_IBHF                 3
#define NPCX_EVSTAT_EOR                  4
#define NPCX_EVSTAT_EOW                  5
#define NPCX_EVSTAT_STSR                 6
#define NPCX_EVSTAT_IBOR                 7
#define NPCX_STATUS_OBES                 6
#define NPCX_STATUS_IBFS                 7
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_SHICFG3_OBUFLVLDIS          7
#define NPCX_SHICFG4_IBUFLVLDIS          7
#define NPCX_SHICFG5_IBUFLVL2            FIELD(0, 6)
#define NPCX_SHICFG5_IBUFLVL2DIS         7
#define NPCX_EVSTAT2_IBHF2               0
#define NPCX_EVSTAT2_CSNRE               1
#define NPCX_EVSTAT2_CSNFE               2
#define NPCX_EVENABLE2_IBHF2EN           0
#define NPCX_EVENABLE2_CSNREEN           1
#define NPCX_EVENABLE2_CSNFEEN           2
#endif

/******************************************************************************/
/* Monotonic Counter (MTC) Registers */
#define NPCX_TTC                         REG32(NPCX_MTC_BASE_ADDR + 0x000)
#define NPCX_WTC                         REG32(NPCX_MTC_BASE_ADDR + 0x004)
#define NPCX_MTCTST                       REG8(NPCX_MTC_BASE_ADDR + 0x008)
#define NPCX_MTCVER                       REG8(NPCX_MTC_BASE_ADDR + 0x00C)

/* MTC register fields */
#define NPCX_WTC_PTO                     30
#define NPCX_WTC_WIE                     31

/******************************************************************************/
/* Low Power RAM definitions */
#define NPCX_LPRAM_CTRL                  REG32(0x40001044)

/******************************************************************************/
/*  eSPI Registers */
#define NPCX_ESPIID                 REG32(NPCX_ESPI_BASE_ADDR + 0X00)
#define NPCX_ESPICFG                REG32(NPCX_ESPI_BASE_ADDR + 0X04)
#define NPCX_ESPISTS                REG32(NPCX_ESPI_BASE_ADDR + 0X08)
#define NPCX_ESPIIE                 REG32(NPCX_ESPI_BASE_ADDR + 0X0C)
#define NPCX_ESPIWE                 REG32(NPCX_ESPI_BASE_ADDR + 0X10)
#define NPCX_VWREGIDX               REG32(NPCX_ESPI_BASE_ADDR + 0X14)
#define NPCX_VWREGDATA              REG32(NPCX_ESPI_BASE_ADDR + 0X18)
#define NPCX_OOBCTL                 REG32(NPCX_ESPI_BASE_ADDR + 0X24)
#define NPCX_FLASHRXRDHEAD          REG32(NPCX_ESPI_BASE_ADDR + 0X28)
#define NPCX_FLASHTXWRHEAD          REG32(NPCX_ESPI_BASE_ADDR + 0X2C)
#define NPCX_FLASHCFG               REG32(NPCX_ESPI_BASE_ADDR + 0X34)
#define NPCX_FLASHCTL               REG32(NPCX_ESPI_BASE_ADDR + 0X38)
#define NPCX_ESPIIERR               REG32(NPCX_ESPI_BASE_ADDR + 0X3C)

/* eSPI Virtual Wire channel registers */
#define NPCX_VWEVSM(n)              REG32(NPCX_ESPI_BASE_ADDR + 0x100 + (4*(n)))
#define NPCX_VWEVMS(n)              REG32(NPCX_ESPI_BASE_ADDR + 0x140 + (4*(n)))
#define NPCX_VWCTL                  REG32(NPCX_ESPI_BASE_ADDR + 0x2FC)

/* eSPI register fields */
#define NPCX_ESPICFG_PCHANEN             0
#define NPCX_ESPICFG_VWCHANEN            1
#define NPCX_ESPICFG_OOBCHANEN           2
#define NPCX_ESPICFG_FLASHCHANEN         3
#define NPCX_ESPICFG_IOMODE_FILED        FIELD(8, 9)
#define NPCX_ESPICFG_MAXFREQ_FILED       FIELD(10, 12)
#define NPCX_ESPICFG_PCCHN_SUPP          24
#define NPCX_ESPICFG_VWCHN_SUPP          25
#define NPCX_ESPICFG_OOBCHN_SUPP         26
#define NPCX_ESPICFG_FLASHCHN_SUPP       27
#define NPCX_ESPIIE_IBRSTIE              0
#define NPCX_ESPIIE_CFGUPDIE             1
#define NPCX_ESPIIE_BERRIE               2
#define NPCX_ESPIIE_OOBRXIE              3
#define NPCX_ESPIIE_FLASHRXIE            4
#define NPCX_ESPIIE_SFLASHRDIE           5
#define NPCX_ESPIIE_PERACCIE             6
#define NPCX_ESPIIE_DFRDIE               7
#define NPCX_ESPIIE_VWUPDIE              8
#define NPCX_ESPIIE_ESPIRSTIE            9
#define NPCX_ESPIIE_PLTRSTIE             10
#define NPCX_ESPIIE_AMERRIE              15
#define NPCX_ESPIIE_AMDONEIE             16
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_ESPIIE_BMTXDONEIE           19
#define NPCX_ESPIIE_PBMRXIE              20
#define NPCX_ESPIIE_PMSGRXIE             21
#define NPCX_ESPIIE_BMBURSTERRIE         22
#define NPCX_ESPIIE_BMBURSTDONEIE        23
#endif
#define NPCX_ESPIWE_IBRSTWE              0
#define NPCX_ESPIWE_CFGUPDWE             1
#define NPCX_ESPIWE_BERRWE               2
#define NPCX_ESPIWE_OOBRXWE              3
#define NPCX_ESPIWE_FLASHRXWE            4
#define NPCX_ESPIWE_PERACCWE             6
#define NPCX_ESPIWE_DFRDWE               7
#define NPCX_ESPIWE_VWUPDWE              8
#define NPCX_ESPIWE_ESPIRSTWE            9
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_ESPIWE_PBMRXWE              20
#define NPCX_ESPIWE_PMSGRXWE             21
#endif
#define NPCX_ESPISTS_IBRST               0
#define NPCX_ESPISTS_CFGUPD              1
#define NPCX_ESPISTS_BERR                2
#define NPCX_ESPISTS_OOBRX               3
#define NPCX_ESPISTS_FLASHRX             4
#define NPCX_ESPISTS_SFLASHRD            5
#define NPCX_ESPISTS_PERACC              6
#define NPCX_ESPISTS_DFRD                7
#define NPCX_ESPISTS_VWUPD               8
#define NPCX_ESPISTS_ESPIRST             9
#define NPCX_ESPISTS_PLTRST              10
#define NPCX_ESPISTS_AMERR               15
#define NPCX_ESPISTS_AMDONE              16
#if defined(CHIP_FAMILY_NPCX7)
#define NPCX_ESPISTS_VWUPDW              17
#define NPCX_ESPISTS_BMTXDONE            19
#define NPCX_ESPISTS_PBMRX               20
#define NPCX_ESPISTS_PMSGRX              21
#define NPCX_ESPISTS_BMBURSTERR          22
#define NPCX_ESPISTS_BMBURSTDONE         23
#define NPCX_ESPISTS_ESPIRST_LVL         24
#endif
/* eSPI Virtual Wire channel register fields */
#define NPCX_VWEVSM_WIRE                 FIELD(0, 4)
#define NPCX_VWEVMS_WIRE                 FIELD(0, 4)
#define NPCX_VWEVSM_VALID                FIELD(4, 4)
#define NPCX_VWEVMS_VALID                FIELD(4, 4)

/* Macro functions for eSPI CFG & IE */
#define IS_SLAVE_CHAN_ENABLE(ch)         IS_BIT_SET(NPCX_ESPICFG, ch)
#define IS_HOST_CHAN_EN(ch)              IS_BIT_SET(NPCX_ESPICFG, (ch+4))
#define ENABLE_ESPI_CHAN(ch)             SET_BIT(NPCX_ESPICFG, ch)
#define DISABLE_ESPI_CHAN(ch)            CLEAR_BIT(NPCX_ESPICFG, ch)
/* ESPI Slave Channel Support Definitions */
#define ESPI_SUPP_CH_PC                  BIT(NPCX_ESPICFG_PCCHN_SUPP)
#define ESPI_SUPP_CH_VM                  BIT(NPCX_ESPICFG_VWCHN_SUPP)
#define ESPI_SUPP_CH_OOB                 BIT(NPCX_ESPICFG_OOBCHN_SUPP)
#define ESPI_SUPP_CH_FLASH               BIT(NPCX_ESPICFG_FLASHCHN_SUPP)
#define ESPI_SUPP_CH_ALL                 (ESPI_SUPP_CH_PC | ESPI_SUPP_CH_VM | \
					  ESPI_SUPP_CH_OOB | ESPI_SUPP_CH_FLASH)
/* ESPI Interrupts Enable Definitions */
#define ESPIIE_IBRST                     BIT(NPCX_ESPIIE_IBRSTIE)
#define ESPIIE_CFGUPD                    BIT(NPCX_ESPIIE_CFGUPDIE)
#define ESPIIE_BERR                      BIT(NPCX_ESPIIE_BERRIE)
#define ESPIIE_OOBRX                     BIT(NPCX_ESPIIE_OOBRXIE)
#define ESPIIE_FLASHRX                   BIT(NPCX_ESPIIE_FLASHRXIE)
#define ESPIIE_SFLASHRD                  BIT(NPCX_ESPIIE_SFLASHRDIE)
#define ESPIIE_PERACC                    BIT(NPCX_ESPIIE_PERACCIE)
#define ESPIIE_DFRD                      BIT(NPCX_ESPIIE_DFRDIE)
#define ESPIIE_VWUPD                     BIT(NPCX_ESPIIE_VWUPDIE)
#define ESPIIE_ESPIRST                   BIT(NPCX_ESPIIE_ESPIRSTIE)
#define ESPIIE_PLTRST                    BIT(NPCX_ESPIIE_PLTRSTIE)
#define ESPIIE_AMERR                     BIT(NPCX_ESPIIE_AMERRIE)
#define ESPIIE_AMDONE                    BIT(NPCX_ESPIIE_AMDONEIE)
#if defined(CHIP_FAMILY_NPCX7)
#define ESPIIE_BMTXDONE                  BIT(NPCX_ESPIIE_BMTXDONEIE)
#define ESPIIE_PBMRX                     BIT(NPCX_ESPIIE_PBMRXIE)
#define ESPIIE_PMSGRX                    BIT(NPCX_ESPIIE_PMSGRXIE)
#define ESPIIE_BMBURSTERR                BIT(NPCX_ESPIIE_BMBURSTERRIE)
#define ESPIIE_BMBURSTDONE               BIT(NPCX_ESPIIE_BMBURSTDONEIE)
#endif
/* eSPI Interrupts for VW */
#define ESPIIE_VW                        (ESPIIE_VWUPD | ESPIIE_PLTRST)
/* eSPI Interrupts for Generic */
#define ESPIIE_GENERIC                   (ESPIIE_IBRST | ESPIIE_CFGUPD | \
					  ESPIIE_BERR | ESPIIE_ESPIRST)
/* ESPI Wake-up Enable Definitions */
#define ESPIWE_IBRST                     BIT(NPCX_ESPIWE_IBRSTWE)
#define ESPIWE_CFGUPD                    BIT(NPCX_ESPIWE_CFGUPDWE)
#define ESPIWE_BERR                      BIT(NPCX_ESPIWE_BERRWE)
#define ESPIWE_OOBRX                     BIT(NPCX_ESPIWE_OOBRXWE)
#define ESPIWE_FLASHRX                   BIT(NPCX_ESPIWE_FLASHRXWE)
#define ESPIWE_PERACC                    BIT(NPCX_ESPIWE_PERACCWE)
#define ESPIWE_DFRD                      BIT(NPCX_ESPIWE_DFRDWE)
#define ESPIWE_VWUPD                     BIT(NPCX_ESPIWE_VWUPDWE)
#define ESPIWE_ESPIRST                   BIT(NPCX_ESPIWE_ESPIRSTWE)
#if defined(CHIP_FAMILY_NPCX7)
#define ESPIWE_PBMRX                     BIT(NPCX_ESPIWE_PBMRXWE)
#define ESPIWE_PMSGRX                    BIT(NPCX_ESPIWE_PMSGRXWE)
#endif
/* eSPI  Wake-up enable for VW */
#define ESPIWE_VW                        ESPIWE_VWUPD
/* eSPI  Wake-up enable for Generic */
#define ESPIWE_GENERIC                   (ESPIWE_IBRST | ESPIWE_CFGUPD | \
					  ESPIWE_BERR)
/* Macro functions for eSPI VW */
#define ESPI_VWEVMS_NUM                  12
#define ESPI_VWEVSM_NUM                  10
#define ESPI_VW_IDX_WIRE_NUM             4
/* Determine Virtual Wire type */
#define VM_TYPE(i)              ((i >= 0   && i <=  1) ? ESPI_VW_TYPE_INT_EV : \
				 (i >= 2   && i <=  7) ? ESPI_VW_TYPE_SYS_EV : \
				 (i >= 64  && i <= 127) ? ESPI_VW_TYPE_PLT : \
				 (i >= 128 && i <= 255) ? ESPI_VW_TYPE_GPIO : \
							ESPI_VW_TYPE_NONE)

/* Bit filed manipulation for VWEVMS Value */
#define VWEVMS_INX(i)                ((i<<8)  & 0x00007F00)
#define VWEVMS_INX_EN(n)             ((n<<15) & 0x00008000)
#define VWEVMS_PLTRST_EN(p)          ((p<<17) & 0x00020000)
#define VWEVMS_INT_EN(e)             ((e<<18) & 0x00040000)
#define VWEVMS_ESPIRST_EN(r)         ((r<<19) & 0x00080000)
#if defined(CHIP_FAMILY_NPCX7)
#define VWEVMS_WK_EN(e)              ((e<<20) & 0x00100000)
#define VWEVMS_INTWK_EN(e)           (VWEVMS_INT_EN(e) | VWEVMS_WK_EN(e))
#elif defined(CHIP_FAMILY_NPCX5)
#define VWEVMS_INTWK_EN              VWEVMS_INT_EN
#endif
#define VWEVMS_FIELD(i, n, p, e, r)  (VWEVMS_INX(i) | VWEVMS_INX_EN(n) | \
				VWEVMS_PLTRST_EN(p) | VWEVMS_INTWK_EN(e) | \
				VWEVMS_ESPIRST_EN(r))
#define VWEVMS_IDX_GET(reg)          (((reg & 0x00007F00)>>8))

/* Bit filed manipulation for VWEVSM Value */
#define VWEVSM_VALID_N(v)            ((v<<4)  & 0x000000F0)
#define VWEVSM_INX(i)                ((i<<8)  & 0x00007F00)
#define VWEVSM_INX_EN(n)             ((n<<15) & 0x00008000)
#define VWEVSM_DIRTY(d)              ((d<<16) & 0x00010000)
#define VWEVSM_PLTRST_EN(p)          ((p<<17) & 0x00020000)
#define VWEVSM_CDRST_EN(c)           ((c<<19) & 0x00080000)
#define VWEVSM_FIELD(i, n, v, p, c)  (VWEVSM_INX(i) | VWEVSM_INX_EN(n) | \
				VWEVSM_VALID_N(v) | VWEVSM_PLTRST_EN(p) |\
				VWEVSM_CDRST_EN(c))
#define VWEVSM_IDX_GET(reg)          (((reg & 0x00007F00)>>8))

/* define macro to handle SMI/SCI Virtual Wire */
/* Read SMI VWire status from VWEVSM(offset 2) register. */
#define SMI_STATUS_MASK    ((uint8_t) (NPCX_VWEVSM(2) & 0x00000002))
/*
 * Read SCI VWire status from VWEVSM(offset 2) register.
 * Left shift 2 to meet the SCIB filed in HIPMIC register.
 */
#define SCI_STATUS_MASK    (((uint8_t) (NPCX_VWEVSM(2) & 0x00000001)) << 2)
#define SCIB_MASK(v)       (v << NPCX_HIPMIC_SCIB)
#define SMIB_MASK(v)       (v << NPCX_HIPMIC_SMIB)
#define NPCX_VW_SCI(level)  ((NPCX_HIPMIC(PM_CHAN_1) & 0xF9) | \
				SMI_STATUS_MASK | SCIB_MASK(level))
#define NPCX_VW_SMI(level)  ((NPCX_HIPMIC(PM_CHAN_1) & 0xF9) | \
				SCI_STATUS_MASK | SMIB_MASK(level))

/* eSPI enumeration */
/* eSPI channels */
enum {
	NPCX_ESPI_CH_PC = 0,
	NPCX_ESPI_CH_VW,
	NPCX_ESPI_CH_OOB,
	NPCX_ESPI_CH_FLASH,
	NPCX_ESPI_CH_COUNT,
	NPCX_ESPI_CH_GENERIC,
	NPCX_ESPI_CH_NONE = 0xFF
};

/* eSPI IO modes */
enum {
	NPCX_ESPI_IO_MODE_SINGLE = 0,
	NPCX_ESPI_IO_MODE_DUAL   = 1,
	NPCX_ESPI_IO_MODE_Quad   = 2,
	NPCX_ESPI_IO_MODE_ALL    = 3,
	NPCX_ESPI_IO_MODE_NONE   = 0xFF
};

/* eSPI max supported frequency */
enum {
	NPCX_ESPI_MAXFREQ_20   = 0,
	NPCX_ESPI_MAXFREQ_25   = 1,
	NPCX_ESPI_MAXFREQ_33   = 2,
	NPCX_ESPI_MAXFREQ_50   = 3,
#if defined(CHIP_FAMILY_NPCX5)
	NPCX_ESPI_MAXFREQ_66   = 4,
#endif
	NPCX_ESPI_MAXFREQ_NOOE = 0xFF
};

#if defined(CHIP_FAMILY_NPCX5)

#if (FMCLK <= 33000000)
#define NPCX_ESPI_MAXFREQ_MAX	NPCX_ESPI_MAXFREQ_33
#elif (FMCLK <= 48000000)
#define NPCX_ESPI_MAXFREQ_MAX	NPCX_ESPI_MAXFREQ_50
#else
#define NPCX_ESPI_MAXFREQ_MAX	NPCX_ESPI_MAXFREQ_66
#endif

#elif defined(CHIP_FAMILY_NPCX7)

#if (FMCLK <= 33000000)
#define NPCX_ESPI_MAXFREQ_MAX	NPCX_ESPI_MAXFREQ_33
#else
#define NPCX_ESPI_MAXFREQ_MAX	NPCX_ESPI_MAXFREQ_50
#endif

#else
#error "Please define NPCX_ESPI_MAXFREQ_MAX for your chip."
#endif

/* VW types */
enum {
	ESPI_VW_TYPE_INT_EV,            /* Interrupt event */
	ESPI_VW_TYPE_SYS_EV,            /* System Event */
	ESPI_VW_TYPE_PLT,               /* Platform specific */
	ESPI_VW_TYPE_GPIO,              /* General Purpose I/O Expander */
	ESPI_VW_TYPE_NUM,
	ESPI_VW_TYPE_NONE = 0xFF
};

/******************************************************************************/
/* GDMA (General DMA) Registers */
#define NPCX_GDMA_CTL                    REG32(NPCX_GDMA_BASE_ADDR + 0x000)
#define NPCX_GDMA_SRCB                   REG32(NPCX_GDMA_BASE_ADDR + 0x004)
#define NPCX_GDMA_DSTB                   REG32(NPCX_GDMA_BASE_ADDR + 0x008)
#define NPCX_GDMA_TCNT                   REG32(NPCX_GDMA_BASE_ADDR + 0x00C)
#define NPCX_GDMA_CSRC                   REG32(NPCX_GDMA_BASE_ADDR + 0x010)
#define NPCX_GDMA_CDST                   REG32(NPCX_GDMA_BASE_ADDR + 0x014)
#define NPCX_GDMA_CTCNT                  REG32(NPCX_GDMA_BASE_ADDR + 0x018)


/******************************************************************************/
/* GDMA register fields */
#define NPCX_GDMA_CTL_GDMAEN                         0
#define NPCX_GDMA_CTL_GDMAMS                         FIELD(2,   2)
#define NPCX_GDMA_CTL_DADIR                          4
#define NPCX_GDMA_CTL_SADIR                          5
#define NPCX_GDMA_CTL_SAFIX                          7
#define NPCX_GDMA_CTL_SIEN                           8
#define NPCX_GDMA_CTL_BME                            9
#define NPCX_GDMA_CTL_SBMS                           11
#define NPCX_GDMA_CTL_TWS                            FIELD(12,  2)
#define NPCX_GDMA_CTL_DM                             15
#define NPCX_GDMA_CTL_SOFTREQ                        16
#define NPCX_GDMA_CTL_TC                             18
#define NPCX_GDMA_CTL_GDMAERR                        20
#define NPCX_GDMA_CTL_BLOCK_BUG_CORRECTION_DISABLE   26

/******************************************************************************/
/* Nuvoton internal used only registers */
#define NPCX_INTERNAL_CTRL1               REG8(0x400DB000)
#define NPCX_INTERNAL_CTRL2               REG8(0x400DD000)
#define NPCX_INTERNAL_CTRL3               REG8(0x400DF000)

/******************************************************************************/
/* Optional M4 Registers */
#define CPU_DHCSR                        REG32(0xE000EDF0)
#define CPU_MPU_CTRL                     REG32(0xE000ED94)
#define CPU_MPU_RNR                      REG32(0xE000ED98)
#define CPU_MPU_RBAR                     REG32(0xE000ED9C)
#define CPU_MPU_RASR                     REG32(0xE000EDA0)


/******************************************************************************/
/* Flash Utiltiy definition */
/*
 *  Flash commands for the W25Q16CV SPI flash
 */
#define CMD_READ_ID                      0x9F
#define CMD_READ_MAN_DEV_ID              0x90
#define CMD_WRITE_EN                     0x06
#define CMD_WRITE_STATUS                 0x50
#define CMD_READ_STATUS_REG              0x05
#define CMD_READ_STATUS_REG2             0x35
#define CMD_WRITE_STATUS_REG             0x01
#define CMD_FLASH_PROGRAM                0x02
#define CMD_SECTOR_ERASE                 0x20
#define CMD_BLOCK_32K_ERASE              0x52
#define CMD_BLOCK_64K_ERASE              0xd8
#define CMD_PROGRAM_UINT_SIZE            0x08
#define CMD_PAGE_SIZE                    0x00
#define CMD_READ_ID_TYPE                 0x47
#define CMD_FAST_READ                    0x0B

/*
 * Status registers for the W25Q16CV SPI flash
 */
#define SPI_FLASH_SR2_SUS               BIT(7)
#define SPI_FLASH_SR2_CMP               BIT(6)
#define SPI_FLASH_SR2_LB3               BIT(5)
#define SPI_FLASH_SR2_LB2               BIT(4)
#define SPI_FLASH_SR2_LB1               BIT(3)
#define SPI_FLASH_SR2_QE                BIT(1)
#define SPI_FLASH_SR2_SRP1              BIT(0)
#define SPI_FLASH_SR1_SRP0              BIT(7)
#define SPI_FLASH_SR1_SEC               BIT(6)
#define SPI_FLASH_SR1_TB                BIT(5)
#define SPI_FLASH_SR1_BP2               BIT(4)
#define SPI_FLASH_SR1_BP1               BIT(3)
#define SPI_FLASH_SR1_BP0               BIT(2)
#define SPI_FLASH_SR1_WEL               BIT(1)
#define SPI_FLASH_SR1_BUSY              BIT(0)


/* 0: F_CS0 1: F_CS1_1(GPIO86) 2:F_CS1_2(GPIOA6) */
#define FIU_CHIP_SELECT		0
/* Create UMA control mask */
#define MASK(bit)       (0x1 << (bit))
#define A_SIZE          0x03	/* 0: No ADR field 1: 3-bytes ADR field */
#define C_SIZE          0x04	/* 0: 1-Byte CMD field 1:No CMD field */
#define RD_WR           0x05	/* 0: Read 1: Write */
#define DEV_NUM         0x06	/* 0: PVT is used 1: SHD is used */
#define EXEC_DONE       0x07
#define D_SIZE_1        0x01
#define D_SIZE_2        0x02
#define D_SIZE_3        0x03
#define D_SIZE_4        0x04
#define FLASH_SEL       MASK(DEV_NUM)

#define MASK_CMD_ONLY   (MASK(EXEC_DONE) | FLASH_SEL)
#define MASK_CMD_ADR    (MASK(EXEC_DONE) | FLASH_SEL | MASK(A_SIZE))
#define MASK_CMD_ADR_WR (MASK(EXEC_DONE) | FLASH_SEL | MASK(RD_WR)  \
			|MASK(A_SIZE) | D_SIZE_1)
#define MASK_RD_1BYTE   (MASK(EXEC_DONE) | FLASH_SEL | MASK(C_SIZE) | D_SIZE_1)
#define MASK_RD_2BYTE   (MASK(EXEC_DONE) | FLASH_SEL | MASK(C_SIZE) | D_SIZE_2)
#define MASK_RD_3BYTE   (MASK(EXEC_DONE) | FLASH_SEL | MASK(C_SIZE) | D_SIZE_3)
#define MASK_RD_4BYTE   (MASK(EXEC_DONE) | FLASH_SEL | MASK(C_SIZE) | D_SIZE_4)
#define MASK_CMD_RD_1BYTE       (MASK(EXEC_DONE) | FLASH_SEL | D_SIZE_1)
#define MASK_CMD_RD_2BYTE       (MASK(EXEC_DONE) | FLASH_SEL | D_SIZE_2)
#define MASK_CMD_RD_3BYTE       (MASK(EXEC_DONE) | FLASH_SEL | D_SIZE_3)
#define MASK_CMD_RD_4BYTE       (MASK(EXEC_DONE) | FLASH_SEL | D_SIZE_4)
#define MASK_CMD_WR_ONLY        (MASK(EXEC_DONE) | FLASH_SEL | MASK(RD_WR))
#define MASK_CMD_WR_1BYTE       (MASK(EXEC_DONE) | FLASH_SEL | MASK(RD_WR) \
				| MASK(C_SIZE) | D_SIZE_1)
#define MASK_CMD_WR_2BYTE       (MASK(EXEC_DONE) | FLASH_SEL | MASK(RD_WR) \
				| MASK(C_SIZE) | D_SIZE_2)
#define MASK_CMD_WR_ADR         (MASK(EXEC_DONE) | FLASH_SEL | MASK(RD_WR) \
				| MASK(A_SIZE))

/******************************************************************************/
/* APM (Audio Processing Module) Registers */
#define NPCX_APM_SR                 REG8(NPCX_APM_BASE_ADDR + 0x000)
#define NPCX_APM_SR2                REG8(NPCX_APM_BASE_ADDR + 0x004)
#define NPCX_APM_ICR                REG8(NPCX_APM_BASE_ADDR + 0x008)
#define NPCX_APM_IMR                REG8(NPCX_APM_BASE_ADDR + 0x00C)
#define NPCX_APM_IFR                REG8(NPCX_APM_BASE_ADDR + 0x010)
#define NPCX_APM_CR_APM             REG8(NPCX_APM_BASE_ADDR + 0x014)
#define NPCX_APM_CR_CK              REG8(NPCX_APM_BASE_ADDR + 0x018)
#define NPCX_APM_AICR_ADC           REG8(NPCX_APM_BASE_ADDR + 0x01C)
#define NPCX_APM_FCR_ADC            REG8(NPCX_APM_BASE_ADDR + 0x020)
#define NPCX_APM_CR_DMIC            REG8(NPCX_APM_BASE_ADDR + 0x02C)
#define NPCX_APM_CR_ADC             REG8(NPCX_APM_BASE_ADDR + 0x030)
#define NPCX_APM_CR_MIX             REG8(NPCX_APM_BASE_ADDR + 0x034)
#define NPCX_APM_DR_MIX             REG8(NPCX_APM_BASE_ADDR + 0x038)
#define NPCX_APM_GCR_ADCL           REG8(NPCX_APM_BASE_ADDR + 0x03C)
#define NPCX_APM_GCR_ADCR           REG8(NPCX_APM_BASE_ADDR + 0x040)
#define NPCX_APM_GCR_MIXADCL        REG8(NPCX_APM_BASE_ADDR + 0x044)
#define NPCX_APM_GCR_MIXADCR        REG8(NPCX_APM_BASE_ADDR + 0x048)
#define NPCX_APM_CR_ADC_AGC         REG8(NPCX_APM_BASE_ADDR + 0x04C)
#define NPCX_APM_DR_ADC_AGC         REG8(NPCX_APM_BASE_ADDR + 0x050)
#define NPCX_APM_SR_ADC_AGCDGL      REG8(NPCX_APM_BASE_ADDR + 0x054)
#define NPCX_APM_SR_ADC_AGCDGR      REG8(NPCX_APM_BASE_ADDR + 0x058)
#define NPCX_APM_CR_VAD             REG8(NPCX_APM_BASE_ADDR + 0x05C)
#define NPCX_APM_DR_VAD             REG8(NPCX_APM_BASE_ADDR + 0x060)
#define NPCX_APM_CR_VAD_CMD         REG8(NPCX_APM_BASE_ADDR + 0x064)
#define NPCX_APM_CR_TR              REG8(NPCX_APM_BASE_ADDR + 0x068)
#define NPCX_APM_DR_TR              REG8(NPCX_APM_BASE_ADDR + 0x06C)
#define NPCX_APM_SR_TR1             REG8(NPCX_APM_BASE_ADDR + 0x070)
#define NPCX_APM_SR_TR_SRCADC       REG8(NPCX_APM_BASE_ADDR + 0x074)

/******************************************************************************/
/* APM register fields */
#define NPCX_APM_SR_IRQ_PEND                    6
#define NPCX_APM_SR2_SMUTEIP                    6
#define NPCX_APM_ICR_INTR_MODE                  FIELD(6, 2)
#define NPCX_APM_IMR_VAD_DTC_MASK               6
#define NPCX_APM_IFR_VAD_DTC                    6
#define NPCX_APM_CR_APM_PD                      0
#define NPCX_APM_CR_APM_AGC_DIS                 FIELD(1, 2)
#define NPCX_APM_CR_CK_MCLK_FREQ                FIELD(0, 2)
#define NPCX_APM_AICR_ADC_ADC_AUDIOIF           FIELD(0, 2)
#define NPCX_APM_AICR_ADC_PD_AICR_ADC           4
#define NPCX_APM_AICR_ADC_ADC_ADWL              FIELD(6, 2)
#define NPCX_APM_FCR_ADC_ADC_FREQ               FIELD(0, 4)
#define NPCX_APM_FCR_ADC_ADC_WNF                FIELD(4, 2)
#define NPCX_APM_FCR_ADC_ADC_HPF                6
#define NPCX_APM_CR_DMIC_ADC_DMIC_SEL_RIGHT     FIELD(0, 2)
#define NPCX_APM_CR_DMIC_ADC_DMIC_SEL_LEFT      FIELD(2, 2)
#define NPCX_APM_CR_DMIC_ADC_DMIC_RATE          FIELD(4, 3)
#define NPCX_APM_CR_DMIC_PD_DMIC                7
#define NPCX_APM_CR_ADC_ADC_SOFT_MUTE           7
#define NPCX_APM_CR_MIX_MIX_ADD                 FIELD(0, 6)
#define NPCX_APM_CR_MIX_MIX_LOAD                6
#define NPCX_APM_DR_MIX_MIX_DATA                FIELD(0, 8)
#define NPCX_APM_MIX_2_AIADCR_SEL               FIELD(4, 2)
#define NPCX_APM_MIX_2_AIADCL_SEL               FIELD(6, 2)
#define NPCX_APM_GCR_ADCL_GIDL                  FIELD(0, 6)
#define NPCX_APM_GCR_ADCL_LRGID                 7
#define NPCX_APM_GCR_ADCR_GIDR                  FIELD(0, 6)
#define NPCX_APM_GCR_MIXADCL_GIMIXL             FIELD(0, 6)
#define NPCX_APM_GCR_MIXADCR_GIMIXR             FIELD(0, 6)
#define NPCX_APM_CR_ADC_AGC_ADC_AGC_ADD         FIELD(0, 6)
#define NPCX_APM_CR_ADC_AGC_ADC_AGC_LOAD        6
#define NPCX_APM_CR_ADC_AGC_ADC_AGC_EN          7
#define NPCX_APM_DR_ADC_AGC_ADC_AGC_DATA        FIELD(0, 8)
#define NPCX_ADC_AGC_0_AGC_TARGET               FIELD(2, 4)
#define NPCX_ADC_AGC_0_AGC_STEREO               6
#define NPCX_ADC_AGC_1_HOLD                     FIELD(0, 4)
#define NPCX_ADC_AGC_1_NG_THR                   FIELD(4, 3)
#define NPCX_ADC_AGC_1_NG_EN                    7
#define NPCX_ADC_AGC_2_DCY                      FIELD(0, 4)
#define NPCX_ADC_AGC_2_ATK                      FIELD(4, 4)
#define NPCX_ADC_AGC_3_AGC_MAX                  FIELD(0, 5)
#define NPCX_ADC_AGC_4_AGC_MIN                  FIELD(0, 5)
#define NPCX_APM_CR_VAD_VAD_ADD                 FIELD(0, 6)
#define NPCX_APM_CR_VAD_VAD_LOAD                6
#define NPCX_APM_CR_VAD_VAD_EN                  7
#define NPCX_APM_DR_VAD_VAD_DATA                FIELD(0, 8)
#define NPCX_APM_CR_VAD_CMD_VAD_RESTART         0
#define NPCX_APM_CR_TR_FAST_ON                  7
#define NPCX_VAD_0_VAD_INSEL                    FIELD(0, 2)
#define NPCX_VAD_0_VAD_DMIC_FREQ                FIELD(2, 3)
#define NPCX_VAD_0_VAD_ADC_WAKEUP               5
#define NPCX_VAD_0_ZCD_EN                       6
#define NPCX_VAD_1_VAD_POWER_SENS               FIELD(0, 5)
#define NPCX_APM_CONTROL_ADD                    FIELD(0, 6)
#define NPCX_APM_CONTROL_LOAD                   6

/******************************************************************************/
/* FMUL2 (Frequency Multiplier Module 2) Registers */
#define NPCX_FMUL2_FM2CTRL          REG8(NPCX_FMUL2_BASE_ADDR + 0x000)
#define NPCX_FMUL2_FM2ML            REG8(NPCX_FMUL2_BASE_ADDR + 0x002)
#define NPCX_FMUL2_FM2MH            REG8(NPCX_FMUL2_BASE_ADDR + 0x004)
#define NPCX_FMUL2_FM2N             REG8(NPCX_FMUL2_BASE_ADDR + 0x006)
#define NPCX_FMUL2_FM2P             REG8(NPCX_FMUL2_BASE_ADDR + 0x008)
#define NPCX_FMUL2_FM2_VER          REG8(NPCX_FMUL2_BASE_ADDR + 0x00A)

/******************************************************************************/
/* FMUL2 register fields */
#define NPCX_FMUL2_FM2CTRL_LOAD2                0
#define NPCX_FMUL2_FM2CTRL_LOCK2                2
#define NPCX_FMUL2_FM2CTRL_FMUL2_DIS            5
#define NPCX_FMUL2_FM2CTRL_TUNE_DIS             6
#define NPCX_FMUL2_FM2CTRL_CLK2_CHNG            7
#define NPCX_FMUL2_FM2N_FM2N                    FIELD(0, 6)
#define NPCX_FMUL2_FM2P_WFPRED                  FIELD(4, 4)

/******************************************************************************/
/* WOV (Wake-on-Voice) Registers */
#define NPCX_WOV_CLOCK_CNTL         REG32(NPCX_WOV_BASE_ADDR + 0x000)
#define NPCX_WOV_PLL_CNTL1          REG32(NPCX_WOV_BASE_ADDR + 0x004)
#define NPCX_WOV_PLL_CNTL2          REG32(NPCX_WOV_BASE_ADDR + 0x008)
#define NPCX_WOV_FIFO_CNT           REG32(NPCX_WOV_BASE_ADDR + 0x00C)
#define NPCX_WOV_FIFO_OUT           REG32(NPCX_WOV_BASE_ADDR + 0x010)
#define NPCX_WOV_STATUS             REG32(NPCX_WOV_BASE_ADDR + 0x014)
#define NPCX_WOV_WOV_INTEN          REG32(NPCX_WOV_BASE_ADDR + 0x018)
#define NPCX_WOV_APM_CTRL           REG32(NPCX_WOV_BASE_ADDR + 0x01C)
#define NPCX_WOV_I2S_CNTL(n)        REG32(NPCX_WOV_BASE_ADDR + 0x020 + (4*n))
#define NPCX_WOV_VERSION            REG32(NPCX_WOV_BASE_ADDR + 0x030)

/******************************************************************************/
/* WOV register fields */
#define NPCX_WOV_CLOCK_CNT_CLK_SEL              0
#define NPCX_WOV_CLOCK_CNT_DMIC_EN              3
#define NPCX_WOV_CLOCK_CNT_PLL_EDIV_SEL         7
#define NPCX_WOV_CLOCK_CNT_PLL_EDIV             FIELD(8, 7)
#define NPCX_WOV_CLOCK_CNT_PLL_EDIV_DC          FIELD(16, 7)
#define NPCX_WOV_CLOCK_CNT_DMIC_CKDIV_EN        24
#define NPCX_WOV_CLOCK_CNT_DMIC_CKDIV_SEL       25
#define NPCX_WOV_FIFO_CNT_FIFO_ITHRSH           FIELD(0, 6)
#define NPCX_WOV_FIFO_CNT_FIFO_WTHRSH           FIELD(6, 6)
#define NPCX_WOV_FIFO_CNT_I2S_FFRST             13
#define NPCX_WOV_FIFO_CNT_CORE_FFRST            14
#define NPCX_WOV_FIFO_CNT_CFIFO_ISEL            FIELD(16, 3)
#define NPCX_WOV_STATUS_CFIFO_CNT               FIELD(0, 8)
#define NPCX_WOV_STATUS_CFIFO_NE                8
#define NPCX_WOV_STATUS_CFIFO_OIT               9
#define NPCX_WOV_STATUS_CFIFO_OWT               10
#define NPCX_WOV_STATUS_CFIFO_OVRN              11
#define NPCX_WOV_STATUS_I2S_FIFO_OVRN           12
#define NPCX_WOV_STATUS_I2S_FIFO_UNDRN          13
#define NPCX_WOV_STATUS_BITS			FIELD(9, 6)
#define NPCX_WOV_INTEN_VAD_INTEN                0
#define NPCX_WOV_INTEN_VAD_WKEN                 1
#define NPCX_WOV_INTEN_CFIFO_NE_IE              8
#define NPCX_WOV_INTEN_CFIFO_OIT_IE             9
#define NPCX_WOV_INTEN_CFIFO_OWT_WE             10
#define NPCX_WOV_INTEN_CFIFO_OVRN_IE            11
#define NPCX_WOV_INTEN_I2S_FIFO_OVRN_IE         12
#define NPCX_WOV_INTEN_I2S_FIFO_UNDRN_IE        13
#define NPCX_WOV_APM_CTRL_APM_RST               0
#define NPCX_WOV_PLL_CNTL1_PLL_PWDEN            0
#define NPCX_WOV_PLL_CNTL1_PLL_OTDV1            FIELD(4, 4)
#define NPCX_WOV_PLL_CNTL1_PLL_OTDV2            FIELD(8, 4)
#define NPCX_WOV_PLL_CNTL1_PLL_LOCKI            15
#define NPCX_WOV_PLL_CNTL2_PLL_FBDV             FIELD(0, 12)
#define NPCX_WOV_PLL_CNTL2_PLL_INDV             FIELD(12, 4)
#define NPCX_WOV_I2S_CNTL_I2S_BCNT              FIELD(0, 5)
#define NPCX_WOV_I2S_CNTL_I2S_TRIG              5
#define NPCX_WOV_I2S_CNTL_I2S_LBHIZ             6
#define NPCX_WOV_I2S_CNTL_I2S_ST_DEL            FIELD(7, 9)
#define NPCX_WOV_I2S_CNTL_I2S_CHAN              FIELD(0, 16)
#define NPCX_WOV_I2S_CNTL0_I2S_HIZD             16
#define NPCX_WOV_I2S_CNTL0_I2S_HIZ              17
#define NPCX_WOV_I2S_CNTL0_I2S_SCLK_INV         18
#define NPCX_WOV_I2S_CNTL0_I2S_OPS              19
#define NPCX_WOV_I2S_CNTL0_I2S_OPE              20
#define NPCX_WOV_I2S_CNTL0_I2S_IPS              21
#define NPCX_WOV_I2S_CNTL0_I2S_IPE              22
#define NPCX_WOV_I2S_CNTL0_I2S_TST              23
#define NPCX_WOV_I2S_CNTL1_I2S_CHN1_DIS         24

/******************************************************************************/
/* UART registers and functions */

#if NPCX_UART_MODULE2

#ifdef CHIP_FAMILY_NPCX5
/*
 * To be used as 2nd parameter to NPCX_WK*() macro, table (1st parameter) is
 * always 1 == MIWU_TABLE_1.
 */
#define NPCX_UART_WK_GROUP 6
#define NPCX_UART_WK_BIT 4
#endif /* CHIP_FAMILY_NPCX5 */
#define NPCX_UART_MIWU_IRQ NPCX_IRQ_WKINTG_1
#define NPCX_UART_DEVALT NPCX_DEVALT(0x0C)
#define NPCX_UART_DEVALT_SL NPCX_DEVALTC_UART_SL2
#define NPCX_UART_ALT_DEVALT NPCX_DEVALT(0x0A)
#define NPCX_UART_ALT_DEVALT_SL NPCX_DEVALTA_UART_SL1
#else /* !NPCX_UART_MODULE2 */

#ifdef CHIP_FAMILY_NPCX5
#define NPCX_UART_WK_GROUP 1
#define NPCX_UART_WK_BIT 0
#endif /* CHIP_FAMILY_NPCX5 */
#define NPCX_UART_MIWU_IRQ NPCX_IRQ_WKINTB_1
#define NPCX_UART_DEVALT NPCX_DEVALT(0x0A)
#define NPCX_UART_DEVALT_SL NPCX_DEVALTA_UART_SL1
#define NPCX_UART_ALT_DEVALT NPCX_DEVALT(0x0C)
#define NPCX_UART_ALT_DEVALT_SL NPCX_DEVALTC_UART_SL2
#endif /* NPCX_UART_MODULE2 */

#ifdef CHIP_FAMILY_NPCX7
#define NPCX_UART_WK_GROUP     MIWU_GROUP_8
#define NPCX_UART_WK_BIT       7
#ifdef NPCX_SECOND_UART
#define NPCX_UART2_WK_GROUP    MIWU_GROUP_1
#define NPCX_UART2_WK_BIT      6
#endif

#endif

/* This routine checks pending bit of GPIO wake-up functionality */
#if defined(CHIP_FAMILY_NPCX5)
static inline int uart_is_wakeup_from_gpio(void)
{
	return IS_BIT_SET(NPCX_WKPND(1, NPCX_UART_WK_GROUP), NPCX_UART_WK_BIT);
}

/* This routine checks wake-up functionality from GPIO is enabled or not */
static inline int uart_is_enable_wakeup(void)
{
	return IS_BIT_SET(NPCX_WKEN(1, NPCX_UART_WK_GROUP), NPCX_UART_WK_BIT);
}

/* This routine clears the pending wake-up from GPIO on UART rx pin */
static inline void uart_clear_pending_wakeup(void)
{
	SET_BIT(NPCX_WKPCL(1, NPCX_UART_WK_GROUP), NPCX_UART_WK_BIT);
}

/* This routine enables wake-up functionality from GPIO on UART rx pin */
static inline void uart_enable_wakeup(int enable)
{
	UPDATE_BIT(NPCX_WKEN(1, NPCX_UART_WK_GROUP), NPCX_UART_WK_BIT,
		enable);
}

/* This routine checks functionality is UART rx or not */
static inline int npcx_is_uart(void)
{
	return IS_BIT_SET(NPCX_UART_DEVALT, NPCX_UART_DEVALT_SL);
}
#endif

/* Wake pin definitions, defined at board-level */
extern const enum gpio_signal hibernate_wake_pins[];
extern const int hibernate_wake_pins_used;

#endif /* __CROS_EC_REGISTERS_H */
