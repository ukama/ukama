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
#include "gfx2d_gpu.h"
#include "gfx2d_ringbuffer.h"
#include <linux/dma-mapping.h>
#include <linux/platform_device.h>
#include <linux/slab.h>

#define RB_ALLOC_UNIT 256
#define RB_LEN 16
#define GFX2D_GPU_RINGBUFFER_SZ (RB_LEN * RB_ALLOC_UNIT)

struct gfx2d_ringbuffer *gfx2d_ringbuffer_new(struct gfx2d_gpu *gpu)
{
	struct gfx2d_ringbuffer *ring = NULL;
	int ret;

	if (!gpu->pdev) {
		ret = -ENODEV;
		goto fail;
	}

	/* We assume everwhere that GFX2D_GPU_RINGBUFFER_SZ is a power of 2 */
	BUILD_BUG_ON(!is_power_of_2(GFX2D_GPU_RINGBUFFER_SZ));

	ring = kzalloc(sizeof(*ring), GFP_KERNEL);
	if (!ring) {
		ret = -ENOMEM;
		goto fail;
	}

	ring->gpu = gpu;
	ring->size = GFX2D_GPU_RINGBUFFER_SZ;
	ring->wsize = GFX2D_GPU_RINGBUFFER_SZ / sizeof(uint32_t);

	dma_set_coherent_mask(&gpu->pdev->dev, DMA_BIT_MASK(32));

	ring->start = dma_alloc_coherent(&gpu->pdev->dev, GFX2D_GPU_RINGBUFFER_SZ,
					 &ring->paddr, GFP_KERNEL);

	if (IS_ERR(ring->start)) {
		ret = PTR_ERR(ring->start);
		ring->start = 0;
		goto fail;
	}
	ring->end   = ring->start + ring->wsize;
	ring->cur   = ring->start;
	ring->tail  = ring->start;

	return ring;

fail:
	if (ring && ring->start)
		dma_free_wc(&gpu->pdev->dev, GFX2D_GPU_RINGBUFFER_SZ,
			    ring->start, ring->paddr);

	gfx2d_ringbuffer_destroy(ring);
	return ERR_PTR(ret);
}

void gfx2d_ringbuffer_destroy(struct gfx2d_ringbuffer *ring)
{
	if (IS_ERR_OR_NULL(ring))
		return;

	if (ring->start)
		dma_free_wc(&ring->gpu->pdev->dev, GFX2D_GPU_RINGBUFFER_SZ,
			    ring->start, ring->paddr);

	kfree(ring);
}
