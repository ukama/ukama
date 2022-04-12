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
#ifndef __GFX2D_GPU_H__
#define __GFX2D_GPU_H__

#include "gfx2d_ringbuffer.h"
#include <drm/atmel_drm.h>
#include <linux/clk.h>
#include <linux/io.h>
#include <linux/seq_file.h>

struct gfx2d_file_private {
	rwlock_t queuelock;
};

struct gfx2d_gpu {
	struct drm_device *dev;
	struct platform_device *pdev;
	void* mmio;
	int irq;
	uint32_t version;
	uint32_t mfn;
	struct gfx2d_ringbuffer *rb;
	struct clk *periph_clk;
};

#define GFX2D_IDLE_TIMEOUT msecs_to_jiffies(1000)

#define spin_until(X) ({                                   \
	int __ret = -ETIMEDOUT;                            \
	unsigned long __t = jiffies + GFX2D_IDLE_TIMEOUT;  \
	do {                                               \
		if (X) {                                   \
			__ret = 0;                         \
			break;                             \
		}                                          \
	} while (time_before(jiffies, __t));               \
	__ret;                                             \
})

struct gfx2d_gpu *gfx2d_gpu_init(struct device *dev);
uint32_t gfx2d_last_fence(struct gfx2d_gpu *gpu);
int gfx2d_submit(struct gfx2d_gpu *gpu, uint32_t* buf, uint32_t size);
int gfx2d_flush(struct gfx2d_gpu *gpu);
void gfx2d_idle(struct gfx2d_gpu *gpu);
#ifdef CONFIG_DEBUG_FS
void gfx2d_show(struct gfx2d_gpu *gpu, struct seq_file *m);
#endif
void gfx2d_dump_info(struct gfx2d_gpu *gpu);
void gfx2d_gpu_cleanup(struct gfx2d_gpu *gpu);

#define DBG(fmt, ...) DRM_DEBUG_DRIVER(fmt"\n", ##__VA_ARGS__)

#endif /* __GFX2D_GPU_H__ */
