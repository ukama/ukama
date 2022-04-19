/* ----------------------------------------------------------------------------
 *         ATMEL Microcontroller Software Support
 * ----------------------------------------------------------------------------
 * Copyright (c) 2014, Atmel Corporation
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
#ifndef __ACT8865_H__
#define __ACT8865_H__

#include "div.h"

/*
 * ACT8865 Register Map
 */
#define SYS_0		0x00
#define SYS_1		0x01

/* DC/DC */
#define REG1_0		0x20
#define REG1_1		0x21
#define REG1_2		0x22

#define REG2_0		0x30
#define REG2_1		0x31
#define REG2_2		0x32

#define REG3_0		0x40
#define REG3_1		0x41
#define REG3_2		0x42

#define REG4_0		0x50
#define REG4_1		0x51

#define REG5_0		0x54
#define REG5_1		0x55

#define REG6_0		0x60
#define REG6_1		0x61

#define REG7_0		0x64
#define REG7_1		0x65

/*
 * Field Definitions
 */
#define REG_MODE_BIT		(0x01 << 5)
#define	REG_MODE_POWER_SAVING		(0x00 << 5)
#define	REG_MODE_FIX_FREQ		(0x01 << 5)

#define REG_ENABLE_BIT		(0x01 << 7)

/*
 * Voltage Code
 */
#define ACT8865_1V2	0x18
#define ACT8865_1V25	0x19
#define ACT8865_1V3	0x1a
#define ACT8865_1V8	0x24
#define ACT8865_2V5	0x31
#define ACT8865_3V3	0x39

#define ACT8865_INVALID_LDO_VOLTAGE(mV) ((mV) < 600 && (mV) > 3900)

/* Get voltage code for voltage in range of 0.6-3.9V: */
#define ACT8865_mV_to_vc(mV) (unsigned int)			\
	((mV) < 1200 ? div(((mV) - 600), 25) :			\
	 ((mV) < 2400 ? div((mV), 50) :				\
	  div(((mV) + 2400), 100)))

/* Get voltage (in millivolts) for voltage code VC): */
#define ACT8865_vc_to_mV(vc)					\
	(((vc) < 0x18) ? (600 + ((vc) * 25)) :			\
	 (((vc) < 0x30) ? ((vc) * 50) :				\
	  (((vc) - 0x18) * 100)))

/*
 * Definitions
 */
#define	ACT8865_MODE_FIX_FREQ		0x01
#define	ACT8865_MODE_POWER_SAVING	0x02

/*
 * Function Declarations
 */
extern int act8865_set_reg_voltage(unsigned char volt_reg, unsigned char value);
extern int act8865_check_i2c_disabled(void);

extern int act8865_set_power_saving_mode(void);
extern void act8865_workaround(void);
extern int act8945a_suspend_charger(void);

#endif
