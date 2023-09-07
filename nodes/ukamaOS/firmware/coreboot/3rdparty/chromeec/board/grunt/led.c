/* Copyright 2018 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

#include "common.h"
#include "ec_commands.h"
#include "led_pwm.h"
#include "util.h"

const enum ec_led_id supported_led_ids[] = {
	EC_LED_ID_POWER_LED,
};
const int supported_led_ids_count = ARRAY_SIZE(supported_led_ids);

/*
 * We only have a blue and an amber LED, so setting any other colour results in
 * both LEDs being off.
 */
struct pwm_led led_color_map[EC_LED_COLOR_COUNT] = {
				/* Amber, Blue */
	[EC_LED_COLOR_RED]    = {   0,   0 },
	[EC_LED_COLOR_GREEN]  = {   0,   0 },
	[EC_LED_COLOR_BLUE]   = {   0, 100 },
	[EC_LED_COLOR_YELLOW] = {   0,   0 },
	[EC_LED_COLOR_WHITE]  = {   0,   0 },
	[EC_LED_COLOR_AMBER]  = {  100,  0 },
};

/* One logical LED with amber and blue channels. */
struct pwm_led pwm_leds[CONFIG_LED_PWM_COUNT] = {
	{
		PWM_CH_LED1_AMBER,
		PWM_CH_LED2_BLUE,
		PWM_LED_NO_CHANNEL,
	},
};

void led_get_brightness_range(enum ec_led_id led_id, uint8_t *brightness_range)
{
	memset(brightness_range, '\0',
	       sizeof(*brightness_range) * EC_LED_COLOR_COUNT);
	brightness_range[EC_LED_COLOR_AMBER] = 100;
	brightness_range[EC_LED_COLOR_BLUE] = 100;
}

int led_set_brightness(enum ec_led_id led_id, const uint8_t *brightness)
{
	enum pwm_led_id pwm_id;

	/* Convert ec_led_id to pwm_led_id. */
	if (led_id == EC_LED_ID_POWER_LED)
		pwm_id = PWM_LED0;
	else
		return EC_ERROR_UNKNOWN;

	if (brightness[EC_LED_COLOR_BLUE])
		set_pwm_led_color(pwm_id, EC_LED_COLOR_BLUE);
	else if (brightness[EC_LED_COLOR_AMBER])
		set_pwm_led_color(pwm_id, EC_LED_COLOR_AMBER);
	else
		/* Otherwise, the "color" is "off". */
		set_pwm_led_color(pwm_id, -1);

	return EC_SUCCESS;
}
