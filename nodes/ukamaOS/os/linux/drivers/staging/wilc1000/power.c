#include <linux/delay.h>
#include <linux/of.h>
#include <linux/version.h>
#include <linux/of_gpio.h>
#include <linux/gpio.h>

#include "netdev.h"

/**
 * wilc_of_parse_power_pins() - parse power sequence pins; to keep backward
 *		compatibility with old device trees that doesn't provide
 *		power sequence pins we check for default pins on proper boards
 *
 * @wilc:	wilc data structure
 *
 * Returns:	 0 on success, negative error number on failures.
 */
int wilc_of_parse_power_pins(struct wilc *wilc)
{
	static const struct wilc_power_gpios default_gpios[] = {
		{ .reset = GPIO_NUM_RESET,	.chip_en = GPIO_NUM_CHIP_EN, },
	};

	static const struct of_device_id wilc_default_pins_ids[] = {
		{
			.compatible = "atmel,sama5d4-xplained",
			.data = &default_gpios[0],
		},
		{ /* Sentinel. */ }
	};

	struct device_node *of = wilc->dt_dev->of_node;
	struct wilc_power *power = &wilc->power;
	const struct wilc_power_gpios *gpios;
	const struct of_device_id *of_id;
	struct device_node *np;
	int ret = 0;

	/*
	 * The maching here is to keep backward compatibility with old DT that
	 * doesn't provide reset-gpios and chip_en.
	 */
	np = of_find_matching_node_and_match(NULL, wilc_default_pins_ids,
					     &of_id);
	if (np)
		gpios = of_id->data;

	power->gpios.reset = of_get_named_gpio_flags(of, "reset-gpios", 0,
						     NULL);
	if (!gpio_is_valid(power->gpios.reset) && np)
		power->gpios.reset = gpios->reset;
	else
		goto put_node;

	power->gpios.chip_en = of_get_named_gpio_flags(of, "chip_en-gpios", 0,
						       NULL);
	if (!gpio_is_valid(power->gpios.chip_en) && np)
		power->gpios.chip_en = gpios->chip_en;
	else
		goto put_node;

	ret = devm_gpio_request(wilc->dev, power->gpios.chip_en, "CHIP_EN");
	if (ret)
		goto put_node;

	ret = devm_gpio_request(wilc->dev, power->gpios.reset, "RESET");
	if (ret)
		goto put_node;

	return 0;

put_node:
	of_node_put(np);
	return ret;
}

/**
 * wilc_wlan_power() - handle power on/off commands
 *
 * @wilc:	wilc data structure
 * @on:		requested power status
 *
 * Returns:	none
 */
void wilc_wlan_power(struct wilc *wilc, bool on)
{
	if (!gpio_is_valid(wilc->power.gpios.chip_en) ||
	    !gpio_is_valid(wilc->power.gpios.reset)) {
		/* In case SDIO power sequence driver is used to power this
		 * device then the powering sequence is handled by the bus
		 * via pm_runtime_* functions. */
		return;
	}

	if (on) {
		gpio_direction_output(wilc->power.gpios.chip_en, 1);
		mdelay(5);
		gpio_direction_output(wilc->power.gpios.reset, 1);
	} else {
		gpio_direction_output(wilc->power.gpios.chip_en, 0);
		gpio_direction_output(wilc->power.gpios.reset, 0);
	}
}
