/*
 * This file is part of the coreboot project.
 *
 * Copyright 2018 MediaTek Inc.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#include <delay.h>
#include <soc/rtc_common.h>
#include <soc/rtc.h>
#include <soc/pmic_wrap.h>
#include <timer.h>

/* ensure rtc write success */
int rtc_busy_wait(void)
{
	struct stopwatch sw;
	u16 bbpu;

	stopwatch_init_usecs_expire(&sw, RTC_CBUSY_TIMEOUT_US);

	do {
		rtc_read(RTC_BBPU, &bbpu);
		/* Time > 1sec, time out and set recovery mode enable.*/
		if (stopwatch_expired(&sw)) {
			rtc_info("BBPU CBUSY time out !!\n");
			return 0;
		}
	} while (bbpu & RTC_BBPU_CBUSY);

	return 1;
}

int rtc_write_trigger(void)
{
	rtc_write(RTC_WRTGR, 1);
	return rtc_busy_wait();
}

/* unlock rtc write interface */
int rtc_writeif_unlock(void)
{
	rtc_write(RTC_PROT, RTC_PROT_UNLOCK1);
	if (!rtc_write_trigger())
		return 0;
	rtc_write(RTC_PROT, RTC_PROT_UNLOCK2);
	if (!rtc_write_trigger())
		return 0;

	return 1;
}

/* set rtc time */
int rtc_set(const struct rtc_time *time)
{
	return -1;
}

/* get rtc time */
int rtc_get(struct rtc_time *time)
{
	u16 value;

	rtc_read(RTC_TC_SEC, &value);
	time->sec = value;
	rtc_read(RTC_TC_MIN, &value);
	time->min = value;
	rtc_read(RTC_TC_HOU, &value);
	time->hour = value;
	rtc_read(RTC_TC_DOM, &value);
	time->mday = value;
	rtc_read(RTC_TC_MTH, &value);
	time->mon = value;
	rtc_read(RTC_TC_YEA, &value);
	time->year = (value + RTC_MIN_YEAR_OFFSET) % 100;

	return 0;
}

/* set rtc xosc setting */
int rtc_xosc_write(u16 val)
{
	u16 bbpu;

	rtc_write(RTC_OSC32CON, RTC_OSC32CON_UNLOCK1);
	if (!rtc_busy_wait())
		return 0;
	rtc_write(RTC_OSC32CON, RTC_OSC32CON_UNLOCK2);
	if (!rtc_busy_wait())
		return 0;

	rtc_write(RTC_OSC32CON, val);
	if (!rtc_busy_wait())
		return 0;

	rtc_read(RTC_BBPU, &bbpu);
	bbpu |= RTC_BBPU_KEY | RTC_BBPU_RELOAD;
	rtc_write(RTC_BBPU, bbpu);

	return rtc_write_trigger();
}

/* initialize rtc related registers */
int rtc_reg_init(void)
{
	u16 irqsta;

	rtc_write(RTC_IRQ_EN, 0);
	rtc_write(RTC_CII_EN, 0);
	rtc_write(RTC_AL_MASK, 0);
	rtc_write(RTC_AL_YEA, 1970 - RTC_MIN_YEAR);
	rtc_write(RTC_AL_MTH, 1);
	rtc_write(RTC_AL_DOM, 1);
	rtc_write(RTC_AL_DOW, 4);
	rtc_write(RTC_AL_HOU, 0);
	rtc_write(RTC_AL_MIN, 0);
	rtc_write(RTC_AL_SEC, 0);

	rtc_write(RTC_DIFF, 0);
	rtc_write(RTC_CALI, 0);
	if (!rtc_write_trigger())
		return 0;

	rtc_read(RTC_IRQ_STA, &irqsta);  /* read clear */

	/* init time counters after resetting RTC_DIFF and RTC_CALI */
	rtc_write(RTC_TC_YEA, RTC_DEFAULT_YEA - RTC_MIN_YEAR);
	rtc_write(RTC_TC_MTH, RTC_DEFAULT_MTH);
	rtc_write(RTC_TC_DOM, RTC_DEFAULT_DOM);
	rtc_write(RTC_TC_DOW, RTC_DEFAULT_DOW);
	rtc_write(RTC_TC_HOU, 0);
	rtc_write(RTC_TC_MIN, 0);
	rtc_write(RTC_TC_SEC, 0);

	return rtc_write_trigger();
}

static u8 rtc_check_state(void)
{
	u16 con;
	u16 pwrky1;
	u16 pwrky2;

	rtc_read(RTC_CON, &con);
	rtc_read(RTC_POWERKEY1, &pwrky1);
	rtc_read(RTC_POWERKEY2, &pwrky2);

	rtc_info("con=%x, pwrkey1=%x, pwrkey2=%x\n", con, pwrky1, pwrky2);

	if (con & RTC_CON_LPSTA_RAW)
		return RTC_STATE_INIT;

	if (!rtc_busy_wait())
		return RTC_STATE_RECOVER;

	if (!rtc_writeif_unlock())
		return RTC_STATE_RECOVER;

	if (pwrky1 != RTC_POWERKEY1_KEY || pwrky2 != RTC_POWERKEY2_KEY)
		return RTC_STATE_INIT;
	else
		return RTC_STATE_REBOOT;
}

void rtc_boot_common(void)
{
	u16 bbpu;
	u16 con;
	u16 irqsta;

	switch (rtc_check_state()) {
	case RTC_STATE_REBOOT:
		pwrap_write_field(RTC_BBPU, RTC_BBPU_KEY | RTC_BBPU_RELOAD,
				  0xFFFF, 0);
		rtc_write_trigger();
		rtc_osc_init();
		break;
	case RTC_STATE_RECOVER:
		rtc_init(1);
		break;
	case RTC_STATE_INIT:
	default:
		if (rtc_init(0))
			rtc_init(1);
		break;
	}

	rtc_read(RTC_IRQ_STA, &irqsta);  /* Read clear */
	rtc_read(RTC_BBPU, &bbpu);
	rtc_read(RTC_CON, &con);

	rtc_info("irqsta=%x, bbpu=%x, con=%x\n", irqsta, bbpu, con);
}
