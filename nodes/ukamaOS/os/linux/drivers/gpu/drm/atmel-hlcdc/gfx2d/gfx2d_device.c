/*
 * Copyright (C) 2018 Microchip
 * Joshua Henderson <joshua.henderson@microchip.com>
 *
 * This program is free software; you can redistribute it and/or modify it
 * under the terms of the GNU General Public License version 2 as published by
 * the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
 * FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
 * more details.
 *
 * You should have received a copy of the GNU General Public License along with
 * this program.  If not, see <http://www.gnu.org/licenses/>.
 */

#include "../atmel_hlcdc_dc.h"
#include "gfx2d_gpu.h"
#include <drm/atmel_drm.h>
#include <drm/drmP.h>
#include <linux/component.h>
#include <linux/kernel.h>
#include <linux/platform_device.h>

struct gfx2d_gpu *gfx2d_load_gpu(struct drm_device *dev)
{
	struct atmel_hlcdc_dc *priv = dev->dev_private;
	struct platform_device *pdev = priv->gpu_pdev;
	struct gfx2d_gpu *gpu = NULL;

	if (pdev)
		gpu = platform_get_drvdata(pdev);

	if (!gpu) {
		dev_err_once(dev->dev, "no GPU device was found\n");
		return NULL;
	}

	return gpu;
}

static void set_gpu_pdev(struct drm_device *dev,
			 struct platform_device *pdev)
{
	struct atmel_hlcdc_dc *priv = dev->dev_private;
	priv->gpu_pdev = pdev;
}

static int gfx2d_bind(struct device *dev, struct device *master, void *data)
{
	struct drm_device *drm = data;
	struct gfx2d_gpu *gpu;

	set_gpu_pdev(drm, to_platform_device(dev));

	gpu = gfx2d_gpu_init(dev);
	if (IS_ERR(gpu)) {
		dev_warn(dev, "failed to init gfx2d gpu\n");
		return PTR_ERR(gpu);
	}
	dev_set_drvdata(dev, gpu);

	return 0;
}

static void gfx2d_unbind(struct device *dev, struct device *master,
			 void *data)
{
	set_gpu_pdev(dev_get_drvdata(master), NULL);
}

static const struct component_ops gfx2d_ops = {
	.bind   = gfx2d_bind,
	.unbind = gfx2d_unbind,
};

static int gfx2d_probe(struct platform_device *pdev)
{
	return component_add(&pdev->dev, &gfx2d_ops);
}

static int gfx2d_remove(struct platform_device *pdev)
{
	component_del(&pdev->dev, &gfx2d_ops);
	return 0;
}

static const struct of_device_id dt_match[] = {
	{ .compatible = "microchip,sam9x60-gfx2d" },
	{}
};

static struct platform_driver gfx2d_driver = {
	.probe = gfx2d_probe,
	.remove = gfx2d_remove,
	.driver = {
		.name = "gfx2d",
		.of_match_table = dt_match,
	},
};

void __init gfx2d_register(void)
{
	platform_driver_register(&gfx2d_driver);
}

void __exit gfx2d_unregister(void)
{
	platform_driver_unregister(&gfx2d_driver);
}
