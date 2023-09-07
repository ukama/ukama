/* Copyright 2015 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Intersil ISL-9237/8 battery charger driver.
 */

#ifndef __CROS_EC_ISL923X_H
#define __CROS_EC_ISL923X_H

#define ISL923X_ADDR_FLAGS	(0x09)

/* Registers */
#define ISL923X_REG_CHG_CURRENT      0x14
#define ISL923X_REG_ADAPTER_CURRENT1 0x3f
#define ISL923X_REG_ADAPTER_CURRENT2 0x3b
#define ISL923X_REG_SYS_VOLTAGE_MAX  0x15
#define ISL923X_REG_SYS_VOLTAGE_MIN  0x3e
#define ISL923X_REG_PROCHOT_AC       0x47
#define ISL923X_REG_PROCHOT_DC       0x48
#define ISL923X_REG_T1_T2            0x38
#define ISL923X_REG_CONTROL0         0x39
#define ISL923X_REG_CONTROL1         0x3c
#define ISL923X_REG_CONTROL2         0x3d
#define ISL9238_REG_CONTROL3         0x4c
#define ISL9238_REG_CONTROL4         0x4e
#define ISL923X_REG_INFO             0x3a
#define ISL9238_REG_INFO2            0x4d
#define ISL923X_REG_OTG_VOLTAGE      0x49
#define ISL923X_REG_OTG_CURRENT      0x4a
#define ISL9238_REG_INPUT_VOLTAGE    0x4b
#define ISL923X_REG_MANUFACTURER_ID  0xfe
#define ISL923X_REG_DEVICE_ID        0xff

/* Sense resistor default values in mOhm */
#define ISL923X_DEFAULT_SENSE_RESISTOR_AC 20
#define ISL923X_DEFAULT_SENSE_RESISTOR 10

/* Maximum charging current register value */
#define ISL923X_CURRENT_REG_MAX 0x17c0 /* bit<12:2> 10111110000 */

/* 2-level adpater current limit duration T1 & T2 in micro seconds */
#define ISL923X_T1_10000 0x00
#define ISL923X_T1_20000 0x01
#define ISL923X_T1_15000 0x02
#define ISL923X_T1_5000  0x03
#define ISL923X_T1_1000  0x04
#define ISL923X_T1_500   0x05
#define ISL923X_T1_100   0x06
#define ISL923X_T1_0     0x07
#define ISL923X_T2_10    (0x00 << 8)
#define ISL923X_T2_100   (0x01 << 8)
#define ISL923X_T2_500   (0x02 << 8)
#define ISL923X_T2_1000  (0x03 << 8)
#define ISL923X_T2_300   (0x04 << 8)
#define ISL923X_T2_750   (0x05 << 8)
#define ISL923X_T2_2000  (0x06 << 8)
#define ISL923X_T2_10000 (0x07 << 8)

#define ISL9237_SYS_VOLTAGE_REG_MAX 13824
#define ISL9238_SYS_VOLTAGE_REG_MAX 18304
#define ISL923X_SYS_VOLTAGE_REG_MIN 2048

/* PROCHOT# debounce time and duration time in micro seconds */
#define ISL923X_PROCHOT_DURATION_10000  (0 << 6)
#define ISL923X_PROCHOT_DURATION_20000  BIT(6)
#define ISL923X_PROCHOT_DURATION_15000  (2 << 6)
#define ISL923X_PROCHOT_DURATION_5000   (3 << 6)
#define ISL923X_PROCHOT_DURATION_1000   (4 << 6)
#define ISL923X_PROCHOT_DURATION_500    (5 << 6)
#define ISL923X_PROCHOT_DURATION_100000 (6 << 6)
#define ISL923X_PROCHOT_DURATION_0      (7 << 6)
#define ISL923X_PROCHOT_DURATION_MASK   (7 << 6)

#define ISL923X_PROCHOT_DEBOUNCE_10   (0 << 9)
#define ISL923X_PROCHOT_DEBOUNCE_100  BIT(9)
#define ISL923X_PROCHOT_DEBOUNCE_500  (2 << 9)
#define ISL923X_PROCHOT_DEBOUNCE_1000 (3 << 9)
#define ISL923X_PROCHOT_DEBOUNCE_MASK (3 << 9)

/* Maximum PROCHOT register value */
#define ISL923X_PROCHOT_AC_REG_MAX 6400
#define ISL923X_PROCHOT_DC_REG_MAX 12800

/* Control0: adapter voltage regulation reference */
#define ISL9237_C0_VREG_REF_3900 0
#define ISL9237_C0_VREG_REF_4200 1
#define ISL9237_C0_VREG_REF_4500 2
#define ISL9237_C0_VREG_REF_4800 3
#define ISL9237_C0_VREG_REF_MASK 0x03

/* Control0: disable adapter voltaqe regulation */
#define ISL923X_C0_DISABLE_VREG BIT(2)

/* Control0: battery DCHOT reference for RS2 == 20mOhm */
#define ISL923X_C0_DCHOT_6A   (0 << 3)
#define ISL923X_C0_DCHOT_5A   BIT(3)
#define ISL923X_C0_DCHOT_4A   (2 << 3)
#define ISL923X_C0_DCHOT_3A   (3 << 3)
#define ISL923X_C0_DCHOT_MASK (3 << 3)

/* Control1: general purpose comparator debounce time in micro second */
#define ISL923X_C1_GP_DEBOUNCE_2       (0 << 14)
#define ISL923X_C1_GP_DEBOUNCE_12      BIT(14)
#define ISL923X_C1_GP_DEBOUNCE_2000    (2 << 14)
#define ISL923X_C1_GP_DEBOUNCE_5000000 (3 << 14)
#define ISL923X_C1_GP_DEBOUNCE_MASK    (3 << 14)

/* Control1: learn mode */
#define ISL923X_C1_LEARN_MODE_AUTOEXIT BIT(13)
#define ISL923X_C1_LEARN_MODE_ENABLE   BIT(12)

/* Control1: OTG enable */
#define ISL923X_C1_OTG BIT(11)

/* Control1: audio filter */
#define ISL923X_C1_AUDIO_FILTER BIT(10)

/* Control1: switch frequency, ISL9238 defines bit 7 as unused */
#define ISL923X_C1_SWITCH_FREQ_PROG (0 << 7) /* 1000kHz or PROG */
#define ISL9237_C1_SWITCH_FREQ_913K BIT(7)
#define ISL923X_C1_SWITCH_FREQ_839K (2 << 7)
#define ISL9237_C1_SWITCH_FREQ_777K (3 << 7)
#define ISL923X_C1_SWITCH_FREQ_723K (4 << 7)
#define ISL9237_C1_SWITCH_FREQ_676K (5 << 7)
#define ISL923X_C1_SWITCH_FREQ_635K (6 << 7)
#define ISL9237_C1_SWITCH_FREQ_599K (7 << 7)
#define ISL923X_C1_SWITCH_FREQ_MASK (7 << 7)

/* Control1: turbo mode */
#define ISL923X_C1_TURBO_MODE BIT(6)

/* Control1: AMON & BMON */
#define ISL923X_C1_DISABLE_MON BIT(5)
#define ISL923X_C1_SELECT_BMON BIT(4)

/* Control1: PSYS, VSYS, VSYSLO */
#define ISL923X_C1_ENABLE_PSYS BIT(3)
#define ISL923X_C1_ENABLE_VSYS BIT(2)
#define ISL923X_C1_VSYSLO_REF_6000 0
#define ISL923X_C1_VSYSLO_REF_6300 1
#define ISL923X_C1_VSYSLO_REF_6600 2
#define ISL923X_C1_VSYSLO_REF_6900 3
#define ISL923X_C1_VSYSLO_REF_MASK 3

/* Control2: trickle charging current in mA */
#define ISL923X_C2_TRICKLE_256  (0 << 14)
#define ISL923X_C2_TRICKLE_128  BIT(14)
#define ISL923X_C2_TRICKLE_64   (2 << 14)
#define ISL923X_C2_TRICKLE_512  (3 << 14)
#define ISL923X_C2_TRICKLE_MASK (3 << 14)

/* Control2: OTGEN debounce time in ms */
#define ISL923X_C2_OTG_DEBOUNCE_1300 (0 << 13)
#define ISL923X_C2_OTG_DEBOUNCE_150  BIT(13)
#define ISL923X_C2_OTG_DEBOUNCE_MASK BIT(13)

/* Control2: 2-level adapter over current */
#define ISL923X_C2_2LVL_OVERCURRENT BIT(12)

/* Control2: adapter insertion debounce time in ms */
#define ISL923X_C2_ADAPTER_DEBOUNCE_1300 (0 << 11)
#define ISL923X_C2_ADAPTER_DEBOUNCE_150  BIT(11)
#define ISL923X_C2_ADAPTER_DEBOUNCE_MASK BIT(11)

/* Control2: PROCHOT debounce time in uS */
#define ISL9238_C2_PROCHOT_DEBOUNCE_7    (0 << 9)
#define ISL9237_C2_PROCHOT_DEBOUNCE_10   (0 << 9)
#define ISL923X_C2_PROCHOT_DEBOUNCE_100  BIT(9)
#define ISL923X_C2_PROCHOT_DEBOUNCE_500  (2 << 9)
#define ISL923X_C2_PROCHOT_DEBOUNCE_1000 (3 << 9)
#define ISL923X_C2_PROCHOT_DEBOUNCE_MASK (3 << 9)

/* Control2: min PROCHOT duration in uS */
#define ISL923X_C2_PROCHOT_DURATION_10000 (0 << 6)
#define ISL923X_C2_PROCHOT_DURATION_20000 BIT(6)
#define ISL923X_C2_PROCHOT_DURATION_15000 (2 << 6)
#define ISL923X_C2_PROCHOT_DURATION_5000  (3 << 6)
#define ISL923X_C2_PROCHOT_DURATION_1000  (4 << 6)
#define ISL923X_C2_PROCHOT_DURATION_500   (5 << 6)
#define ISL923X_C2_PROCHOT_DURATION_100   (6 << 6)
#define ISL923X_C2_PROCHOT_DURATION_0     (7 << 6)
#define ISL923X_C2_PROCHOT_DURATION_MASK  (7 << 6)

/* Control2: turn off ASGATE in OTG mode */
#define ISL923X_C2_ASGATE_OFF BIT(5)

/* Control2: CMIN, general purpose comparator reference in mV */
#define ISL923X_C2_CMIN_2000 (0 << 4)
#define ISL923X_C2_CMIN_1200 BIT(4)

/* Control2: general purpose comparator enable */
#define ISL923X_C2_COMPARATOR BIT(3)

/* Control2: invert CMOUT, general purpose comparator output, polarity */
#define ISL923X_C2_INVERT_CMOUT BIT(2)

/* Control2: disable WOC, way over current */
#define ISL923X_C2_WOC_OFF BIT(1)

/* Control2: PSYS gain in uA/W (ISL9237 only) */
#define ISL9237_C2_PSYS_GAIN BIT(0)

/*
 * Control3: Buck-Boost switching period
 * 0: x1 frequency, 1: half frequency.
 */
#define ISL9238_C3_BB_SWITCHING_PERIOD BIT(1)

/*
 * Control3: AMON/BMON direction.
 * 0: adapter/charging, 1:OTG/discharging (ISL9238 only)
 */
#define ISL9238_C3_AMON_BMON_DIRECTION BIT(3)

/*
 * Control3: Disables Autonomous Charing
 *
 * Note: This is disabled automatically when ever we set the current limit
 * manually (which we always do).
 */
#define ISL9238_C3_DISABLE_AUTO_CHARING BIT(7)

/* Control3: PSYS gain in uA/W (ISL9238 only) */
#define ISL9238_C3_PSYS_GAIN BIT(9)

/* Control3: Don't reload ACLIM on ACIN. */
#define ISL9238_C3_NO_RELOAD_ACLIM_ON_ACIN BIT(14)

/* Control3: Don't reread PROG pin. */
#define ISL9238_C3_NO_REREAD_PROG_PIN BIT(15)

/* OTG voltage limit in mV, current limit in mA */
#define ISL9237_OTG_VOLTAGE_MIN 4864
#define ISL9237_OTG_VOLTAGE_MAX 5376
#define ISL9238_OTG_VOLTAGE_MAX 27456
#define ISL923X_OTG_CURRENT_MAX 4096

#define ISL9238_OTG_VOLTAGE_STEP 12
#define ISL9238_OTG_VOLTAGE_SHIFT 3
#define ISL923X_OTG_CURRENT_STEP 128
#define ISL923X_OTG_CURRENT_SHIFT 7

/* Input voltage regulation voltage reference */
#define ISL9238_INPUT_VOLTAGE_REF_STEP 341
#define ISL9238_INPUT_VOLTAGE_REF_SHIFT 8

/* Info register fields */
#define ISL9237_INFO_PROG_RESISTOR_MASK 0xf
#define ISL923X_INFO_TRICKLE_ACTIVE_MASK BIT(4)
#define ISL9237_INFO_PSTATE_SHIFT 5
#define ISL9237_INFO_PSTATE_MASK 3

enum isl9237_power_stage {
	BUCK_MODE,
	BOOST_MODE,
	BUCK_BOOST_MODE,
	REVERSE_BUCK_MODE
};

#define ISL9237_INFO_FSM_STATE_SHIFT 7
#define ISL9237_INFO_FSM_STATE_MASK 7

enum isl9237_fsm_state {
	FSM_OFF,
	FSM_BAT,
	FSM_ADPT,
	FSM_ACOK,
	FSM_VSYS,
	FSM_CHRG,
	FSM_ENTOG,
	FSM_OTG
};

#define ISL923X_INFO_VSYSLO BIT(10)
#define ISL923X_INFO_DCHOT BIT(11)
#define ISL9237_INFO_ACHOT BIT(12)

#if defined(CONFIG_CHARGER_ISL9237)
#define CHARGER_NAME  "isl9237"
#define CHARGE_V_MAX  ISL9237_SYS_VOLTAGE_REG_MAX
#define CHARGE_V_MIN  ISL923X_SYS_VOLTAGE_REG_MIN
#define CHARGE_V_STEP 8
#elif defined(CONFIG_CHARGER_ISL9238)
#define CHARGER_NAME  "isl9238"
#define CHARGE_V_MAX  ISL9238_SYS_VOLTAGE_REG_MAX
#define CHARGE_V_MIN  ISL923X_SYS_VOLTAGE_REG_MIN
#define CHARGE_V_STEP 8
#endif

#define CHARGE_I_MAX  ISL923X_CURRENT_REG_MAX
#define CHARGE_I_MIN  4
#define CHARGE_I_OFF  0
#define CHARGE_I_STEP 4
#define INPUT_I_MAX   ISL923X_CURRENT_REG_MAX
#define INPUT_I_MIN   4
#define INPUT_I_STEP  4

#define I2C_ADDR_CHARGER_FLAGS ISL923X_ADDR_FLAGS
#endif /* __CROS_EC_ISL923X_H */

/**
 * Initialize AC & DC prochot threshold
 *
 * @param	AC Porchot threshold current in mA:
 *			multiple of 128 up to 6400 mA
 *			DC Porchot threshold current in mA:
 *			multiple of 128 up to 12800 mA
 * 		Bits below 128mA are truncated (ignored).
 * @return enum ec_error_list
 */
int isl923x_set_ac_prochot(uint16_t ma);
int isl923x_set_dc_prochot(uint16_t ma);

#define ISL923X_AC_PROCHOT_CURRENT_MAX	6400	/* mA */
#define ISL923X_DC_PROCHOT_CURRENT_MAX	12800	/* mA */
