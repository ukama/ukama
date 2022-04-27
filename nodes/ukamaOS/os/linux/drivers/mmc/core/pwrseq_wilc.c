// SPDX-License-Identifier: GPL-2.0
// pwrseq_wilc.c - power sequence support for WILC Wifi chips
//
// Copyright (C) 2019 Microchip Technology Inc.
// Copyright (C) 2019 Claudiu Beznea <claudiu.beznea@microchip.com>
//
// Based on the original work at pwrseq_sd8787.c
//  Copyright (C) 2016 Matt Ranostay <matt@ranostay.consulting>

#include <linux/delay.h>
#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/platform_device.h>
#include <linux/module.h>
#include <linux/slab.h>
#include <linux/device.h>
#include <linux/err.h>
#include <linux/gpio/consumer.h>

#include <linux/mmc/host.h>

#include "pwrseq.h"

struct mmc_pwrseq_wilc {
	struct mmc_pwrseq pwrseq;
	struct gpio_desc *reset_gpio;
	struct gpio_desc *pwrdn_gpio;
};

#define to_pwrseq_wilc(p) container_of(p, struct mmc_pwrseq_wilc, pwrseq)

static void mmc_pwrseq_wilc_pre_power_on(struct mmc_host *host)
{
	struct mmc_pwrseq_wilc *pwrseq = to_pwrseq_wilc(host->pwrseq);

	gpiod_set_value_cansleep(pwrseq->pwrdn_gpio, 1);
	usleep_range(5000, 7000);
	gpiod_set_value_cansleep(pwrseq->reset_gpio, 1);
}

static void mmc_pwrseq_wilc_power_off(struct mmc_host *host)
{
	struct mmc_pwrseq_wilc *pwrseq = to_pwrseq_wilc(host->pwrseq);

	gpiod_set_value_cansleep(pwrseq->reset_gpio, 0);
	gpiod_set_value_cansleep(pwrseq->pwrdn_gpio, 0);
}

static const struct mmc_pwrseq_ops mmc_pwrseq_wilc_ops = {
	.pre_power_on = mmc_pwrseq_wilc_pre_power_on,
	.power_off = mmc_pwrseq_wilc_power_off,
};

static const struct of_device_id mmc_pwrseq_wilc_of_match[] = {
	{ .compatible = "mmc-pwrseq-wilc",},
	{/* sentinel */},
};
MODULE_DEVICE_TABLE(of, mmc_pwrseq_wilc_of_match);

static int mmc_pwrseq_wilc_probe(struct platform_device *pdev)
{
	struct mmc_pwrseq_wilc *pwrseq;
	struct device *dev = &pdev->dev;

	pwrseq = devm_kzalloc(dev, sizeof(*pwrseq), GFP_KERNEL);
	if (!pwrseq)
		return -ENOMEM;

	pwrseq->pwrdn_gpio = devm_gpiod_get(dev, "powerdown", GPIOD_OUT_LOW);
	if (IS_ERR(pwrseq->pwrdn_gpio))
		return PTR_ERR(pwrseq->pwrdn_gpio);

	pwrseq->reset_gpio = devm_gpiod_get(dev, "reset", GPIOD_OUT_LOW);
	if (IS_ERR(pwrseq->reset_gpio))
		return PTR_ERR(pwrseq->reset_gpio);

	pwrseq->pwrseq.dev = dev;
	pwrseq->pwrseq.ops = &mmc_pwrseq_wilc_ops;
	pwrseq->pwrseq.owner = THIS_MODULE;
	platform_set_drvdata(pdev, pwrseq);

	return mmc_pwrseq_register(&pwrseq->pwrseq);
}

static int mmc_pwrseq_wilc_remove(struct platform_device *pdev)
{
	struct mmc_pwrseq_wilc *pwrseq = platform_get_drvdata(pdev);

	mmc_pwrseq_unregister(&pwrseq->pwrseq);

	return 0;
}

static struct platform_driver mmc_pwrseq_wilc_driver = {
	.probe = mmc_pwrseq_wilc_probe,
	.remove = mmc_pwrseq_wilc_remove,
	.driver = {
		.name = "pwrseq_wilc",
		.of_match_table = mmc_pwrseq_wilc_of_match,
	},
};

module_platform_driver(mmc_pwrseq_wilc_driver);
MODULE_LICENSE("GPL v2");
