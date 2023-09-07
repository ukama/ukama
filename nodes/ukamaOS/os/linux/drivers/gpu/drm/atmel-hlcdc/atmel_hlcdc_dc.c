// SPDX-License-Identifier: GPL-2.0-only
/*
 * Copyright (C) 2014 Traphandler
 * Copyright (C) 2014 Free Electrons
 * Copyright (C) 2014 Atmel
 *
 * Author: Jean-Jacques Hiblot <jjhiblot@traphandler.com>
 * Author: Boris BREZILLON <boris.brezillon@free-electrons.com>
 */

#include <linux/clk.h>
#include <linux/irq.h>
#include <linux/irqchip.h>
#include <linux/mfd/atmel-hlcdc.h>
#include <linux/module.h>
#include <linux/pm_runtime.h>
#include <linux/platform_device.h>

#include <drm/drm_atomic.h>
#include <drm/drm_atomic_helper.h>
#include <drm/drm_debugfs.h>
#include <drm/drm_drv.h>
#include <drm/drm_fb_helper.h>
#include <drm/drm_gem_cma_helper.h>
#include <drm/drm_gem_framebuffer_helper.h>
#include <drm/drm_irq.h>
#include <drm/drm_probe_helper.h>
#include <drm/drm_vblank.h>

#include "atmel_hlcdc_dc.h"
#include "gfx2d/gfx2d_gpu.h"
#include <drm/atmel_drm.h>
#include <drm/drm_of.h>
#include <linux/component.h>

struct gfx2d_gpu *gfx2d_load_gpu(struct drm_device *dev);
void __init gfx2d_register(void);
void __exit gfx2d_unregister(void);

#define ATMEL_HLCDC_LAYER_IRQS_OFFSET		8

static const struct atmel_hlcdc_layer_desc atmel_hlcdc_at91sam9n12_layers[] = {
	{
		.name = "base",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x40,
		.id = 0,
		.type = ATMEL_HLCDC_BASE_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.xstride = { 2 },
			.default_color = 3,
			.general_config = 4,
		},
		.clut_offset = 0x400,
	},
};

static const struct atmel_hlcdc_dc_desc atmel_hlcdc_dc_at91sam9n12 = {
	.min_width = 0,
	.min_height = 0,
	.max_width = 1280,
	.max_height = 860,
	.max_spw = 0x3f,
	.max_vpw = 0x3f,
	.max_hpw = 0xff,
	.conflicting_output_formats = true,
	.nlayers = ARRAY_SIZE(atmel_hlcdc_at91sam9n12_layers),
	.layers = atmel_hlcdc_at91sam9n12_layers,
};

static const struct atmel_hlcdc_layer_desc atmel_hlcdc_at91sam9x5_layers[] = {
	{
		.name = "base",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x40,
		.id = 0,
		.type = ATMEL_HLCDC_BASE_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.xstride = { 2 },
			.default_color = 3,
			.general_config = 4,
			.disc_pos = 5,
			.disc_size = 6,
		},
		.clut_offset = 0x400,
	},
	{
		.name = "overlay1",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x100,
		.id = 1,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.pstride = { 5 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
		},
		.clut_offset = 0x800,
	},
	{
		.name = "high-end-overlay",
		.formats = &atmel_hlcdc_plane_rgb_and_yuv_formats,
		.regs_offset = 0x280,
		.id = 2,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x4c,
		.layout = {
			.pos = 2,
			.size = 3,
			.memsize = 4,
			.xstride = { 5, 7 },
			.pstride = { 6, 8 },
			.default_color = 9,
			.chroma_key = 10,
			.chroma_key_mask = 11,
			.general_config = 12,
			.scaler_config = 13,
			.csc = 14,
		},
		.clut_offset = 0x1000,
	},
	{
		.name = "cursor",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x340,
		.id = 3,
		.type = ATMEL_HLCDC_CURSOR_LAYER,
		.max_width = 128,
		.max_height = 128,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
		},
		.clut_offset = 0x1400,
	},
};

static const struct atmel_hlcdc_dc_desc atmel_hlcdc_dc_at91sam9x5 = {
	.min_width = 0,
	.min_height = 0,
	.max_width = 800,
	.max_height = 600,
	.max_spw = 0x3f,
	.max_vpw = 0x3f,
	.max_hpw = 0xff,
	.conflicting_output_formats = true,
	.nlayers = ARRAY_SIZE(atmel_hlcdc_at91sam9x5_layers),
	.layers = atmel_hlcdc_at91sam9x5_layers,
};

static const struct atmel_hlcdc_layer_desc atmel_hlcdc_sama5d3_layers[] = {
	{
		.name = "base",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x40,
		.id = 0,
		.type = ATMEL_HLCDC_BASE_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.xstride = { 2 },
			.default_color = 3,
			.general_config = 4,
			.disc_pos = 5,
			.disc_size = 6,
		},
		.clut_offset = 0x600,
	},
	{
		.name = "overlay1",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x140,
		.id = 1,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.pstride = { 5 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
		},
		.clut_offset = 0xa00,
	},
	{
		.name = "overlay2",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x240,
		.id = 2,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.pstride = { 5 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
		},
		.clut_offset = 0xe00,
	},
	{
		.name = "high-end-overlay",
		.formats = &atmel_hlcdc_plane_rgb_and_yuv_formats,
		.regs_offset = 0x340,
		.id = 3,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x4c,
		.layout = {
			.pos = 2,
			.size = 3,
			.memsize = 4,
			.xstride = { 5, 7 },
			.pstride = { 6, 8 },
			.default_color = 9,
			.chroma_key = 10,
			.chroma_key_mask = 11,
			.general_config = 12,
			.scaler_config = 13,
			.phicoeffs = {
				.x = 17,
				.y = 33,
			},
			.csc = 14,
		},
		.clut_offset = 0x1200,
	},
	{
		.name = "cursor",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x440,
		.id = 4,
		.type = ATMEL_HLCDC_CURSOR_LAYER,
		.max_width = 128,
		.max_height = 128,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.pstride = { 5 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
			.scaler_config = 13,
		},
		.clut_offset = 0x1600,
	},
};

static const struct atmel_hlcdc_dc_desc atmel_hlcdc_dc_sama5d3 = {
	.min_width = 0,
	.min_height = 0,
	.max_width = 2048,
	.max_height = 2048,
	.max_spw = 0x3f,
	.max_vpw = 0x3f,
	.max_hpw = 0x1ff,
	.conflicting_output_formats = true,
	.nlayers = ARRAY_SIZE(atmel_hlcdc_sama5d3_layers),
	.layers = atmel_hlcdc_sama5d3_layers,
};

static const struct atmel_hlcdc_layer_desc atmel_hlcdc_sama5d4_layers[] = {
	{
		.name = "base",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x40,
		.id = 0,
		.type = ATMEL_HLCDC_BASE_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.xstride = { 2 },
			.default_color = 3,
			.general_config = 4,
			.disc_pos = 5,
			.disc_size = 6,
		},
		.clut_offset = 0x600,
	},
	{
		.name = "overlay1",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x140,
		.id = 1,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.pstride = { 5 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
		},
		.clut_offset = 0xa00,
	},
	{
		.name = "overlay2",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x240,
		.id = 2,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.pstride = { 5 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
		},
		.clut_offset = 0xe00,
	},
	{
		.name = "high-end-overlay",
		.formats = &atmel_hlcdc_plane_rgb_and_yuv_formats,
		.regs_offset = 0x340,
		.id = 3,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x4c,
		.layout = {
			.pos = 2,
			.size = 3,
			.memsize = 4,
			.xstride = { 5, 7 },
			.pstride = { 6, 8 },
			.default_color = 9,
			.chroma_key = 10,
			.chroma_key_mask = 11,
			.general_config = 12,
			.scaler_config = 13,
			.phicoeffs = {
				.x = 17,
				.y = 33,
			},
			.csc = 14,
		},
		.clut_offset = 0x1200,
	},
};

static const struct atmel_hlcdc_dc_desc atmel_hlcdc_dc_sama5d4 = {
	.min_width = 0,
	.min_height = 0,
	.max_width = 2048,
	.max_height = 2048,
	.max_spw = 0xff,
	.max_vpw = 0xff,
	.max_hpw = 0x3ff,
	.nlayers = ARRAY_SIZE(atmel_hlcdc_sama5d4_layers),
	.layers = atmel_hlcdc_sama5d4_layers,
};

static const struct atmel_hlcdc_layer_desc atmel_hlcdc_sam9x60_layers[] = {
	{
		.name = "base",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x60,
		.id = 0,
		.type = ATMEL_HLCDC_BASE_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.xstride = { 2 },
			.default_color = 3,
			.general_config = 4,
			.disc_pos = 5,
			.disc_size = 6,
		},
		.clut_offset = 0x600,
	},
	{
		.name = "overlay1",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x160,
		.id = 1,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.pstride = { 5 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
		},
		.clut_offset = 0xa00,
	},
	{
		.name = "overlay2",
		.formats = &atmel_hlcdc_plane_rgb_formats,
		.regs_offset = 0x260,
		.id = 2,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x2c,
		.layout = {
			.pos = 2,
			.size = 3,
			.xstride = { 4 },
			.pstride = { 5 },
			.default_color = 6,
			.chroma_key = 7,
			.chroma_key_mask = 8,
			.general_config = 9,
		},
		.clut_offset = 0xe00,
	},
	{
		.name = "high-end-overlay",
		.formats = &atmel_hlcdc_plane_rgb_and_yuv_formats,
		.regs_offset = 0x360,
		.id = 3,
		.type = ATMEL_HLCDC_OVERLAY_LAYER,
		.cfgs_offset = 0x4c,
		.layout = {
			.pos = 2,
			.size = 3,
			.memsize = 4,
			.xstride = { 5, 7 },
			.pstride = { 6, 8 },
			.default_color = 9,
			.chroma_key = 10,
			.chroma_key_mask = 11,
			.general_config = 12,
			.scaler_config = 13,
			.phicoeffs = {
				.x = 17,
				.y = 33,
			},
			.csc = 14,
		},
		.clut_offset = 0x1200,
	},
};

static const struct atmel_hlcdc_dc_desc atmel_hlcdc_dc_sam9x60 = {
	.min_width = 0,
	.min_height = 0,
	.max_width = 2048,
	.max_height = 2048,
	.max_spw = 0xff,
	.max_vpw = 0xff,
	.max_hpw = 0x3ff,
	.fixed_clksrc = true,
	.nlayers = ARRAY_SIZE(atmel_hlcdc_sam9x60_layers),
	.layers = atmel_hlcdc_sam9x60_layers,
};

static const struct of_device_id atmel_hlcdc_of_match[] = {
	{
		.compatible = "atmel,at91sam9n12-hlcdc",
		.data = &atmel_hlcdc_dc_at91sam9n12,
	},
	{
		.compatible = "atmel,at91sam9x5-hlcdc",
		.data = &atmel_hlcdc_dc_at91sam9x5,
	},
	{
		.compatible = "atmel,sama5d2-hlcdc",
		.data = &atmel_hlcdc_dc_sama5d4,
	},
	{
		.compatible = "atmel,sama5d3-hlcdc",
		.data = &atmel_hlcdc_dc_sama5d3,
	},
	{
		.compatible = "atmel,sama5d4-hlcdc",
		.data = &atmel_hlcdc_dc_sama5d4,
	},
	{
		.compatible = "microchip,sam9x60-hlcdc",
		.data = &atmel_hlcdc_dc_sam9x60,
	},
	{ /* sentinel */ },
};
MODULE_DEVICE_TABLE(of, atmel_hlcdc_of_match);

enum drm_mode_status
atmel_hlcdc_dc_mode_valid(struct atmel_hlcdc_dc *dc,
			  const struct drm_display_mode *mode)
{
	int vfront_porch = mode->vsync_start - mode->vdisplay;
	int vback_porch = mode->vtotal - mode->vsync_end;
	int vsync_len = mode->vsync_end - mode->vsync_start;
	int hfront_porch = mode->hsync_start - mode->hdisplay;
	int hback_porch = mode->htotal - mode->hsync_end;
	int hsync_len = mode->hsync_end - mode->hsync_start;

	if (hsync_len > dc->desc->max_spw + 1 || hsync_len < 1)
		return MODE_HSYNC;

	if (vsync_len > dc->desc->max_spw + 1 || vsync_len < 1)
		return MODE_VSYNC;

	if (hfront_porch > dc->desc->max_hpw + 1 || hfront_porch < 1 ||
	    hback_porch > dc->desc->max_hpw + 1 || hback_porch < 1 ||
	    mode->hdisplay < 1)
		return MODE_H_ILLEGAL;

	if (vfront_porch > dc->desc->max_vpw + 1 || vfront_porch < 1 ||
	    vback_porch > dc->desc->max_vpw || vback_porch < 0 ||
	    mode->vdisplay < 1)
		return MODE_V_ILLEGAL;

	return MODE_OK;
}

static void atmel_hlcdc_layer_irq(struct atmel_hlcdc_layer *layer)
{
	if (!layer)
		return;

	if (layer->desc->type == ATMEL_HLCDC_BASE_LAYER ||
	    layer->desc->type == ATMEL_HLCDC_OVERLAY_LAYER ||
	    layer->desc->type == ATMEL_HLCDC_CURSOR_LAYER)
		atmel_hlcdc_plane_irq(atmel_hlcdc_layer_to_plane(layer));
}

static irqreturn_t atmel_hlcdc_dc_irq_handler(int irq, void *data)
{
	struct drm_device *dev = data;
	struct atmel_hlcdc_dc *dc = dev->dev_private;
	unsigned long status;
	unsigned int imr, isr;
	int i;

	regmap_read(dc->hlcdc->regmap, ATMEL_HLCDC_IMR, &imr);
	regmap_read(dc->hlcdc->regmap, ATMEL_HLCDC_ISR, &isr);
	status = imr & isr;
	if (!status)
		return IRQ_NONE;

	if (status & ATMEL_HLCDC_SOF)
		atmel_hlcdc_crtc_irq(dc->crtc);

	for (i = 0; i < ATMEL_HLCDC_MAX_LAYERS; i++) {
		if (ATMEL_HLCDC_LAYER_STATUS(i) & status)
			atmel_hlcdc_layer_irq(dc->layers[i]);
	}

	return IRQ_HANDLED;
}

static struct drm_framebuffer *atmel_hlcdc_fb_create(struct drm_device *dev,
		struct drm_file *file_priv, const struct drm_mode_fb_cmd2 *mode_cmd)
{
	return drm_gem_fb_create(dev, file_priv, mode_cmd);
}

struct atmel_hlcdc_dc_commit {
	struct work_struct work;
	struct drm_device *dev;
	struct drm_atomic_state *state;
};

static void
atmel_hlcdc_dc_atomic_complete(struct atmel_hlcdc_dc_commit *commit)
{
	struct drm_device *dev = commit->dev;
	struct atmel_hlcdc_dc *dc = dev->dev_private;
	struct drm_atomic_state *old_state = commit->state;

	/* Apply the atomic update. */
	drm_atomic_helper_commit_modeset_disables(dev, old_state);
	drm_atomic_helper_commit_planes(dev, old_state, 0);
	drm_atomic_helper_commit_modeset_enables(dev, old_state);

	drm_atomic_helper_wait_for_vblanks(dev, old_state);

	drm_atomic_helper_cleanup_planes(dev, old_state);

	drm_atomic_state_put(old_state);

	/* Complete the commit, wake up any waiter. */
	spin_lock(&dc->commit.wait.lock);
	dc->commit.pending = false;
	wake_up_all_locked(&dc->commit.wait);
	spin_unlock(&dc->commit.wait.lock);

	kfree(commit);
}

static void atmel_hlcdc_dc_atomic_work(struct work_struct *work)
{
	struct atmel_hlcdc_dc_commit *commit =
		container_of(work, struct atmel_hlcdc_dc_commit, work);

	atmel_hlcdc_dc_atomic_complete(commit);
}

static int atmel_hlcdc_dc_atomic_commit(struct drm_device *dev,
					struct drm_atomic_state *state,
					bool async)
{
	struct atmel_hlcdc_dc *dc = dev->dev_private;
	struct atmel_hlcdc_dc_commit *commit;
	int ret;

	ret = drm_atomic_helper_prepare_planes(dev, state);
	if (ret)
		return ret;

	/* Allocate the commit object. */
	commit = kzalloc(sizeof(*commit), GFP_KERNEL);
	if (!commit) {
		ret = -ENOMEM;
		goto error;
	}

	INIT_WORK(&commit->work, atmel_hlcdc_dc_atomic_work);
	commit->dev = dev;
	commit->state = state;

	spin_lock(&dc->commit.wait.lock);
	ret = wait_event_interruptible_locked(dc->commit.wait,
					      !dc->commit.pending);
	if (ret == 0)
		dc->commit.pending = true;
	spin_unlock(&dc->commit.wait.lock);

	if (ret)
		goto err_free;

	/* We have our own synchronization through the commit lock. */
	BUG_ON(drm_atomic_helper_swap_state(state, false) < 0);

	/* Swap state succeeded, this is the point of no return. */
	drm_atomic_state_get(state);
	if (async)
		queue_work(dc->wq, &commit->work);
	else
		atmel_hlcdc_dc_atomic_complete(commit);

	return 0;

err_free:
	kfree(commit);
error:
	drm_atomic_helper_cleanup_planes(dev, state);
	return ret;
}

static const struct drm_mode_config_funcs mode_config_funcs = {
	.fb_create = atmel_hlcdc_fb_create,
	.atomic_check = drm_atomic_helper_check,
	.atomic_commit = atmel_hlcdc_dc_atomic_commit,
};

static int atmel_hlcdc_dc_modeset_init(struct drm_device *dev)
{
	struct atmel_hlcdc_dc *dc = dev->dev_private;
	int ret;

	drm_mode_config_init(dev);

	ret = atmel_hlcdc_create_outputs(dev);
	if (ret) {
		dev_err(dev->dev, "failed to create HLCDC outputs: %d\n", ret);
		return ret;
	}

	ret = atmel_hlcdc_create_planes(dev);
	if (ret) {
		dev_err(dev->dev, "failed to create planes: %d\n", ret);
		return ret;
	}

	ret = atmel_hlcdc_crtc_create(dev);
	if (ret) {
		dev_err(dev->dev, "failed to create crtc\n");
		return ret;
	}

	dev->mode_config.min_width = dc->desc->min_width;
	dev->mode_config.min_height = dc->desc->min_height;
	dev->mode_config.max_width = dc->desc->max_width;
	dev->mode_config.max_height = dc->desc->max_height;
	dev->mode_config.funcs = &mode_config_funcs;

	return 0;
}

static int atmel_hlcdc_dc_load(struct drm_device *dev)
{
	struct platform_device *pdev = to_platform_device(dev->dev);
	const struct of_device_id *match;
	struct atmel_hlcdc_dc *dc;
	int ret;

	match = of_match_node(atmel_hlcdc_of_match, dev->dev->parent->of_node);
	if (!match) {
		dev_err(&pdev->dev, "invalid compatible string\n");
		return -ENODEV;
	}

	if (!match->data) {
		dev_err(&pdev->dev, "invalid hlcdc description\n");
		return -EINVAL;
	}

	dc = devm_kzalloc(dev->dev, sizeof(*dc), GFP_KERNEL);
	if (!dc)
		return -ENOMEM;

	dc->wq = alloc_ordered_workqueue("atmel-hlcdc-dc", 0);
	if (!dc->wq)
		return -ENOMEM;

	init_waitqueue_head(&dc->commit.wait);
	dc->desc = match->data;
	dc->hlcdc = dev_get_drvdata(dev->dev->parent);
	dev->dev_private = dc;

	if (dc->desc->fixed_clksrc) {
		ret = clk_prepare_enable(dc->hlcdc->sys_clk);
		if (ret) {
			dev_err(dev->dev, "failed to enable sys_clk\n");
			goto err_destroy_wq;
		}
	}

	ret = clk_prepare_enable(dc->hlcdc->periph_clk);
	if (ret) {
		dev_err(dev->dev, "failed to enable periph_clk\n");
		goto err_sys_clk_disable;
	}

	pm_runtime_enable(dev->dev);

	ret = drm_vblank_init(dev, 1);
	if (ret < 0) {
		dev_err(dev->dev, "failed to initialize vblank\n");
		goto err_periph_clk_disable;
	}

	ret = atmel_hlcdc_dc_modeset_init(dev);
	if (ret < 0) {
		dev_err(dev->dev, "failed to initialize mode setting\n");
		goto err_periph_clk_disable;
	}

	drm_mode_config_reset(dev);

	pm_runtime_get_sync(dev->dev);
	ret = drm_irq_install(dev, dc->hlcdc->irq);
	pm_runtime_put_sync(dev->dev);
	if (ret < 0) {
		dev_err(dev->dev, "failed to install IRQ handler\n");
		goto err_periph_clk_disable;
	}

	platform_set_drvdata(pdev, dev);

	drm_kms_helper_poll_init(dev);

	return 0;

err_periph_clk_disable:
	pm_runtime_disable(dev->dev);
	clk_disable_unprepare(dc->hlcdc->periph_clk);
err_sys_clk_disable:
	if (dc->desc->fixed_clksrc)
		clk_disable_unprepare(dc->hlcdc->sys_clk);

err_destroy_wq:
	destroy_workqueue(dc->wq);

	return ret;
}

static void atmel_hlcdc_dc_unload(struct drm_device *dev)
{
	struct atmel_hlcdc_dc *dc = dev->dev_private;

	flush_workqueue(dc->wq);
	drm_kms_helper_poll_fini(dev);
	drm_atomic_helper_shutdown(dev);
	drm_mode_config_cleanup(dev);

	pm_runtime_get_sync(dev->dev);
	drm_irq_uninstall(dev);
	pm_runtime_put_sync(dev->dev);

	dev->dev_private = NULL;

	pm_runtime_disable(dev->dev);
	clk_disable_unprepare(dc->hlcdc->periph_clk);
	if (dc->desc->fixed_clksrc)
		clk_disable_unprepare(dc->hlcdc->sys_clk);
	destroy_workqueue(dc->wq);
}

static int atmel_hlcdc_dc_irq_postinstall(struct drm_device *dev)
{
	struct atmel_hlcdc_dc *dc = dev->dev_private;
	unsigned int cfg = 0;
	int i;

	/* Enable interrupts on activated layers */
	for (i = 0; i < ATMEL_HLCDC_MAX_LAYERS; i++) {
		if (dc->layers[i])
			cfg |= ATMEL_HLCDC_LAYER_STATUS(i);
	}

	regmap_write(dc->hlcdc->regmap, ATMEL_HLCDC_IER, cfg);

	return 0;
}

static void atmel_hlcdc_dc_irq_uninstall(struct drm_device *dev)
{
	struct atmel_hlcdc_dc *dc = dev->dev_private;
	unsigned int isr;

	regmap_write(dc->hlcdc->regmap, ATMEL_HLCDC_IDR, 0xffffffff);
	regmap_read(dc->hlcdc->regmap, ATMEL_HLCDC_ISR, &isr);
}



DEFINE_DRM_GEM_CMA_FOPS(fops);

/*
ioctl to export the physical address of GEM to user space for
video decoder
*/
int atmel_drm_gem_get_ioctl(struct drm_device *drm, void *data,
				      struct drm_file *file_priv)
{
	struct drm_gem_object *gem_obj;
	struct drm_gem_cma_object *cma_obj;
	struct drm_mode_map_dumb *args = data;

	mutex_lock(&drm->struct_mutex);

	gem_obj = drm_gem_object_lookup(file_priv, args->handle);
	if (!gem_obj) {
		dev_err(drm->dev, "failed to lookup gem object\n");
		mutex_unlock(&drm->struct_mutex);
		return -EINVAL;
	}

	cma_obj = to_drm_gem_cma_obj(gem_obj);
	args->offset = (__u64)cma_obj->paddr;

	drm_gem_object_put(gem_obj);

	mutex_unlock(&drm->struct_mutex);

	return 0;
}

static int gfx2d_ioctl_submit(struct drm_device *dev, void *data,
			      struct drm_file *file)
{
	struct atmel_hlcdc_dc *priv = dev->dev_private;
	struct gfx2d_gpu *gpu = priv->gpu;
	struct drm_gfx2d_submit *args = data;

	if (!gpu)
		return -ENXIO;

	return gfx2d_submit(gpu, (uint32_t *)args->buf, args->size);
}

static int gfx2d_ioctl_flush(struct drm_device *dev, void *data,
			     struct drm_file *file)
{
	struct atmel_hlcdc_dc *priv = dev->dev_private;
	struct gfx2d_gpu *gpu = priv->gpu;
	return gfx2d_flush(gpu);
}

static int gfx2d_ioctl_gem_addr(struct drm_device *dev, void *data,
				struct drm_file *file_priv)
{
	struct drm_gfx2d_gem_addr *args = data;
	struct drm_gem_object *obj;

	if (!drm_core_check_feature(dev, DRIVER_GEM))
		return -ENODEV;

	mutex_lock(&dev->object_name_lock);
	obj = idr_find(&dev->object_name_idr, (int) args->name);
	if (!obj) {
		mutex_unlock(&dev->object_name_lock);
		return -ENOENT;
	}

	args->paddr = to_drm_gem_cma_obj(obj)->paddr;
	args->size = obj->size;

	mutex_unlock(&dev->object_name_lock);

	return 0;
}

static const struct drm_ioctl_desc atmel_ioctls[] = {
	DRM_IOCTL_DEF_DRV(ATMEL_GEM_GET, atmel_drm_gem_get_ioctl, DRM_UNLOCKED),
	DRM_IOCTL_DEF_DRV(GFX2D_SUBMIT, gfx2d_ioctl_submit, DRM_AUTH|DRM_RENDER_ALLOW),
	DRM_IOCTL_DEF_DRV(GFX2D_FLUSH, gfx2d_ioctl_flush, DRM_AUTH|DRM_RENDER_ALLOW),
	DRM_IOCTL_DEF_DRV(GFX2D_GEM_ADDR, gfx2d_ioctl_gem_addr, DRM_AUTH|DRM_RENDER_ALLOW),
};

#ifdef CONFIG_DEBUG_FS
static int atmel_hlcdc_dc_gpu_show(struct drm_device *dev, struct seq_file *m)
{
	struct atmel_hlcdc_dc *priv = dev->dev_private;
	struct gfx2d_gpu *gpu = priv->gpu;

	if (gpu) {
		gfx2d_show(gpu, m);
	}

	return 0;
}

static int show_locked(struct seq_file *m, void *arg)
{
	struct drm_info_node *node = (struct drm_info_node *) m->private;
	struct drm_device *dev = node->minor->dev;
	int (*show)(struct drm_device *dev, struct seq_file *m) =
		node->info_ent->data;
	int ret;

	ret = mutex_lock_interruptible(&dev->struct_mutex);
	if (ret)
		return ret;

	ret = show(dev, m);

	mutex_unlock(&dev->struct_mutex);

	return ret;
}

static struct drm_info_list atmel_hlcdc_dc_debugfs_list[] = {
	{"gpu", show_locked, 0, atmel_hlcdc_dc_gpu_show},
};

int atmel_hlcdc_dc_debugfs_init(struct drm_minor *minor)
{
	struct drm_device *dev = minor->dev;
	int ret;

	ret = drm_debugfs_create_files(atmel_hlcdc_dc_debugfs_list,
				       ARRAY_SIZE(atmel_hlcdc_dc_debugfs_list),
				       minor->debugfs_root, minor);

	if (ret) {
		dev_err(dev->dev, "could not install atmel_hlcdc_dc_debugfs_list\n");
		return ret;
	}

	return ret;
}
#endif

static void load_gpu(struct drm_device *dev)
{
	static DEFINE_MUTEX(init_lock);
	struct atmel_hlcdc_dc *priv = dev->dev_private;

	mutex_lock(&init_lock);

	if (!priv->gpu)
		priv->gpu = gfx2d_load_gpu(dev);

	mutex_unlock(&init_lock);
}

static struct drm_driver atmel_hlcdc_dc_driver = {
	.driver_features = DRIVER_GEM | DRIVER_MODESET | DRIVER_ATOMIC,
	.irq_handler = atmel_hlcdc_dc_irq_handler,
	.irq_preinstall = atmel_hlcdc_dc_irq_uninstall,
	.irq_postinstall = atmel_hlcdc_dc_irq_postinstall,
	.irq_uninstall = atmel_hlcdc_dc_irq_uninstall,
	.gem_free_object_unlocked = drm_gem_cma_free_object,
	.gem_vm_ops = &drm_gem_cma_vm_ops,
	.prime_handle_to_fd = drm_gem_prime_handle_to_fd,
	.prime_fd_to_handle = drm_gem_prime_fd_to_handle,
	.gem_prime_get_sg_table = drm_gem_cma_prime_get_sg_table,
	.gem_prime_import_sg_table = drm_gem_cma_prime_import_sg_table,
	.gem_prime_vmap = drm_gem_cma_prime_vmap,
	.gem_prime_vunmap = drm_gem_cma_prime_vunmap,
	.gem_prime_mmap = drm_gem_cma_prime_mmap,
	.dumb_create = drm_gem_cma_dumb_create,
#ifdef CONFIG_DEBUG_FS
	.debugfs_init       = atmel_hlcdc_dc_debugfs_init,
#endif
	.ioctls	= atmel_ioctls,
	.num_ioctls= ARRAY_SIZE(atmel_ioctls),
	.fops = &fops,
	.name = "atmel-hlcdc",
	.desc = "Atmel HLCD Controller DRM",
	.date = "20141504",
	.major = 1,
	.minor = 0,
};

static int compare_of(struct device *dev, void *data)
{
	return dev->of_node == data;
}

static int atmel_hlcdc_dc_bind(struct device *dev)
{
	struct drm_device *ddev;
	int ret;

	ddev = drm_dev_alloc(&atmel_hlcdc_dc_driver, dev);
	if (IS_ERR(ddev))
		return PTR_ERR(ddev);

	ret = atmel_hlcdc_dc_load(ddev);
	if (ret)
		goto err_put;

	dev_set_drvdata(dev, ddev);

	ret = component_bind_all(dev, ddev);
	if (ret < 0)
		goto out_bind;

	load_gpu(ddev);

	ret = drm_dev_register(ddev, 0);
	if (ret)
		goto out_register;

	drm_fbdev_generic_setup(ddev, 24);

	return 0;

out_register:
	component_unbind_all(dev, ddev);
out_bind:
	atmel_hlcdc_dc_unload(ddev);

err_put:
	drm_dev_put(ddev);

	return ret;
}

static void atmel_hlcdc_dc_unbind(struct device *dev)
{
	/* TODO */
}

static const struct component_master_ops atmel_hlcdc_dc_master_ops = {
	.bind = atmel_hlcdc_dc_bind,
	.unbind = atmel_hlcdc_dc_unbind,
};

static int atmel_hlcdc_dc_drm_probe(struct platform_device *pdev)
{
	struct component_match *match = NULL;
	struct device_node *core_node;
	struct drm_device *ddev;
	int ret;

	for_each_compatible_node(core_node, NULL, "microchip,sam9x60-gfx2d") {
		if (!of_device_is_available(core_node))
				continue;

		component_match_add(&pdev->dev, &match,
				    compare_of, core_node);
	}

	if (match) {
		return component_master_add_with_match(&pdev->dev,
						       &atmel_hlcdc_dc_master_ops,
						       match);
	}

	/* Fall through old probe routine */
	ddev = drm_dev_alloc(&atmel_hlcdc_dc_driver, &pdev->dev);
	if (IS_ERR(ddev))
		return PTR_ERR(ddev);

	ret = atmel_hlcdc_dc_load(ddev);
	if (ret)
		goto err_put;

	ret = drm_dev_register(ddev, 0);
	if (ret)
		goto err_unload;

	drm_fbdev_generic_setup(ddev, 24);

	return 0;

err_unload:
	atmel_hlcdc_dc_unload(ddev);

err_put:
	drm_dev_put(ddev);

	return ret;
}

static int atmel_hlcdc_dc_drm_remove(struct platform_device *pdev)
{
	struct drm_device *ddev = platform_get_drvdata(pdev);

	drm_dev_unregister(ddev);
	atmel_hlcdc_dc_unload(ddev);
	drm_dev_put(ddev);

	return 0;
}

#ifdef CONFIG_PM_SLEEP
static int atmel_hlcdc_dc_drm_suspend(struct device *dev)
{
	struct drm_device *drm_dev = dev_get_drvdata(dev);
	struct atmel_hlcdc_dc *dc = drm_dev->dev_private;
	struct regmap *regmap = dc->hlcdc->regmap;
	struct drm_atomic_state *state;

	state = drm_atomic_helper_suspend(drm_dev);
	if (IS_ERR(state))
		return PTR_ERR(state);

	dc->suspend.state = state;

	regmap_read(regmap, ATMEL_HLCDC_IMR, &dc->suspend.imr);
	regmap_write(regmap, ATMEL_HLCDC_IDR, dc->suspend.imr);
	clk_disable_unprepare(dc->hlcdc->periph_clk);
	if (dc->desc->fixed_clksrc)
		clk_disable_unprepare(dc->hlcdc->sys_clk);

	return 0;
}

static int atmel_hlcdc_dc_drm_resume(struct device *dev)
{
	struct drm_device *drm_dev = dev_get_drvdata(dev);
	struct atmel_hlcdc_dc *dc = drm_dev->dev_private;

	if (dc->desc->fixed_clksrc)
		clk_prepare_enable(dc->hlcdc->sys_clk);
	clk_prepare_enable(dc->hlcdc->periph_clk);
	regmap_write(dc->hlcdc->regmap, ATMEL_HLCDC_IER, dc->suspend.imr);

	return drm_atomic_helper_resume(drm_dev, dc->suspend.state);
}
#endif

static SIMPLE_DEV_PM_OPS(atmel_hlcdc_dc_drm_pm_ops,
		atmel_hlcdc_dc_drm_suspend, atmel_hlcdc_dc_drm_resume);

static const struct of_device_id atmel_hlcdc_dc_of_match[] = {
	{ .compatible = "atmel,hlcdc-display-controller" },
	{ },
};

static struct platform_driver atmel_hlcdc_dc_platform_driver = {
	.probe	= atmel_hlcdc_dc_drm_probe,
	.remove	= atmel_hlcdc_dc_drm_remove,
	.driver	= {
		.name	= "atmel-hlcdc-display-controller",
		.pm	= &atmel_hlcdc_dc_drm_pm_ops,
		.of_match_table = atmel_hlcdc_dc_of_match,
	},
};

static int __init atmel_hlcdc_dc_drm_init(void)
{
	gfx2d_register();
	return platform_driver_register(&atmel_hlcdc_dc_platform_driver);
}
module_init(atmel_hlcdc_dc_drm_init);

static void __exit atmel_hlcdc_dc_drm_exit(void)
{
	gfx2d_unregister();
	platform_driver_unregister(&atmel_hlcdc_dc_platform_driver);
}
module_exit(atmel_hlcdc_dc_drm_exit);

MODULE_AUTHOR("Jean-Jacques Hiblot <jjhiblot@traphandler.com>");
MODULE_AUTHOR("Boris Brezillon <boris.brezillon@free-electrons.com>");
MODULE_DESCRIPTION("Atmel HLCDC Display Controller DRM Driver");
MODULE_LICENSE("GPL");
MODULE_ALIAS("platform:atmel-hlcdc-dc");
