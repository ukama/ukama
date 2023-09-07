/* Copyright 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 */

/* GPIO common functionality for Chrome EC */

#include "common.h"
#include "console.h"
#include "gpio.h"
#include "host_command.h"
#include "system.h"
#include "util.h"

static uint8_t last_val[(GPIO_COUNT + 7) / 8];

/**
 * Find a GPIO signal by name.
 *
 * @param name		Signal name to find
 *
 * @return the signal index, or GPIO_COUNT if no match.
 */
static enum gpio_signal find_signal_by_name(const char *name)
{
	int i;

	if (!name || !*name)
		return GPIO_COUNT;

	for (i = 0; i < GPIO_COUNT; i++)
		if (gpio_is_implemented(i) &&
		    !strcasecmp(name, gpio_get_name(i)))
			return i;

	return GPIO_COUNT;
}

/**
 * Update last_val
 *
 * @param i		Index of last_val[] to update
 * @param v		New value for last_val[i]
 *
 * @return 1 if last_val[i] was updated, 0 if last_val[i]==v.
 */
static int last_val_changed(int i, int v)
{
	if (v && !(last_val[i / 8] & (1 << (i % 8)))) {
		last_val[i / 8] |= 1 << (i % 8);
		return 1;
	} else if (!v && last_val[i / 8] & (1 << (i % 8))) {
		last_val[i / 8] &= ~(1 << (i % 8));
		return 1;
	} else {
		return 0;
	}
}

static enum ec_error_list set(const char *name, int value)
{
	enum gpio_signal signal = find_signal_by_name(name);

	if (signal == GPIO_COUNT)
		return EC_ERROR_INVAL;

	if (!gpio_is_implemented(signal))
		return EC_ERROR_INVAL;

	if (!(gpio_get_default_flags(signal) & GPIO_OUTPUT))
		return EC_ERROR_INVAL;

	gpio_set_level(signal, value);

	return EC_SUCCESS;
}
/*****************************************************************************/
/* Console commands */

static void print_gpio_info(int gpio)
{
	int changed, v, flags;

	if (!gpio_is_implemented(gpio))
		return;  /* Skip unsupported signals */

	v = gpio_get_level(gpio);
#ifdef CONFIG_CMD_GPIO_EXTENDED
	flags = gpio_get_flags(gpio);
#else
	flags = 0;
#endif
	changed = last_val_changed(gpio, v);

	ccprintf("  %d%c %s%s%s%s%s%s%s%s%s%s%s%s\n", v,
		 (changed ? '*' : ' '),
		 (flags & GPIO_INPUT ? "I " : ""),
		 (flags & GPIO_OUTPUT ? "O " : ""),
		 (flags & GPIO_LOW ? "L " : ""),
		 (flags & GPIO_HIGH ? "H " : ""),
		 (flags & GPIO_ANALOG ? "A " : ""),
		 (flags & GPIO_OPEN_DRAIN ? "ODR " : ""),
		 (flags & GPIO_PULL_UP ? "PU " : ""),
		 (flags & GPIO_PULL_DOWN ? "PD " : ""),
		 (flags & GPIO_ALTERNATE ? "ALT " : ""),
		 (flags & GPIO_SEL_1P8V ? "1P8 " : ""),
		 (flags & GPIO_LOCKED ? "LCK " : ""),
		 gpio_get_name(gpio));

	/* Flush console to avoid truncating output */
	cflush();
}

static int command_gpio_get(int argc, char **argv)
{
	int i;

	/* If a signal is specified, print only that one */
	if (argc == 2) {
		i = find_signal_by_name(argv[1]);
		if (i == GPIO_COUNT)
			return EC_ERROR_PARAM1;
		print_gpio_info(i);

		return EC_SUCCESS;
	}

	/* Otherwise print them all */
	for (i = 0; i < GPIO_COUNT; i++) {
		if (!gpio_is_implemented(i))
			continue;  /* Skip unsupported signals */

		print_gpio_info(i);
	}

	return EC_SUCCESS;
}
DECLARE_SAFE_CONSOLE_COMMAND(gpioget, command_gpio_get,
			     "[name]",
			     "Read GPIO value(s)");

static int command_gpio_set(int argc, char **argv)
{
#ifdef CONFIG_CMD_GPIO_EXTENDED
	int gpio;
	int flags = 0;
	int af = -1;
	char *e;

	if (argc < 3)
		return EC_ERROR_PARAM_COUNT;

	gpio = find_signal_by_name(argv[1]);
	if (gpio == GPIO_COUNT)
		return EC_ERROR_PARAM1;

	if (strcasecmp(argv[2], "IN") == 0)
		flags = GPIO_INPUT;
	else if (strcasecmp(argv[2], "1") == 0)
		flags = GPIO_OUT_HIGH;
	else if (strcasecmp(argv[2], "0") == 0)
		flags = GPIO_OUT_LOW;
	else if (strcasecmp(argv[2], "A") == 0)
		flags = GPIO_ANALOG;
	else if (strcasecmp(argv[2], "ALT") == 0) {
		if (argc >= 4) {
			af = strtoi(argv[3], &e, 0);
			if (*e || af < 0 || af > 5)
				return EC_ERROR_PARAM2;
		}
		flags = GPIO_ALTERNATE;
	} else
		return EC_ERROR_PARAM2;

	/* Update alt function if requested. */
	if (af >= 0) {
		const struct gpio_info *g = gpio_list + gpio;

		gpio_set_alternate_function(g->port, g->mask, af);
	}
	/* Update GPIO flags. */
	gpio_set_flags(gpio, flags);
#else
	char *e;
	int v;

	if (argc < 3)
		return EC_ERROR_PARAM_COUNT;

	v = strtoi(argv[2], &e, 0);
	if (*e)
		return EC_ERROR_PARAM2;

	if (set(argv[1], v) != EC_SUCCESS)
		return EC_ERROR_PARAM1;
#endif
	return EC_SUCCESS;
}
DECLARE_CONSOLE_COMMAND(gpioset, command_gpio_set,
#ifdef CONFIG_CMD_GPIO_EXTENDED
			"name <0 | 1 | IN | A | ALT [func]>",
#else
			"name <0 | 1>",
#endif
			"Set a GPIO");

/*****************************************************************************/
/* Host commands */

static enum ec_status gpio_command_get(struct host_cmd_handler_args *args)
{
	const struct ec_params_gpio_get_v1 *p_v1 = args->params;
	struct ec_response_gpio_get_v1 *r_v1 = args->response;
	int i, len;

	if (args->version == 0) {
		const struct ec_params_gpio_get *p = args->params;
		struct ec_response_gpio_get *r = args->response;

		i = find_signal_by_name(p->name);
		if (i == GPIO_COUNT)
			return EC_RES_ERROR;

		r->val = gpio_get_level(i);
		args->response_size = sizeof(struct ec_response_gpio_get);
		return EC_RES_SUCCESS;
	}

	switch (p_v1->subcmd) {
	case EC_GPIO_GET_BY_NAME:
		i = find_signal_by_name(p_v1->get_value_by_name.name);
		if (i == GPIO_COUNT)
			return EC_RES_ERROR;

		r_v1->get_value_by_name.val = gpio_get_level(i);
		args->response_size = sizeof(r_v1->get_value_by_name);
		break;
	case EC_GPIO_GET_COUNT:
		r_v1->get_count.val = GPIO_COUNT;
		args->response_size = sizeof(r_v1->get_count);
		break;
	case EC_GPIO_GET_INFO:
		if (p_v1->get_info.index >= GPIO_COUNT)
			return EC_RES_ERROR;

		i = p_v1->get_info.index;
		len = strlen(gpio_get_name(i));
		memcpy(r_v1->get_info.name, gpio_get_name(i), len+1);
		r_v1->get_info.val = gpio_get_level(i);
		r_v1->get_info.flags = gpio_get_default_flags(i);
		args->response_size = sizeof(r_v1->get_info);
		break;
	default:
		return EC_RES_INVALID_PARAM;
	}

	return EC_RES_SUCCESS;

}
DECLARE_HOST_COMMAND(EC_CMD_GPIO_GET, gpio_command_get,
		     EC_VER_MASK(0) | EC_VER_MASK(1));

static enum ec_status gpio_command_set(struct host_cmd_handler_args *args)
{
	const struct ec_params_gpio_set *p = args->params;

	if (system_is_locked())
		return EC_RES_ACCESS_DENIED;

	if (set(p->name, p->val) != EC_SUCCESS)
		return EC_RES_ERROR;

	return EC_RES_SUCCESS;
}
DECLARE_HOST_COMMAND(EC_CMD_GPIO_SET, gpio_command_set, EC_VER_MASK(0));
