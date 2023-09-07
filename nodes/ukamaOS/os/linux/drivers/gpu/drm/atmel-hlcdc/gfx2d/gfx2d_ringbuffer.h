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
#ifndef __GFX2D_RINGBUFFER_H__
#define __GFX2D_RINGBUFFER_H__

#include <linux/list.h>
#include <linux/types.h>

struct gfx2d_gpu;

struct gfx2d_ringbuffer {
	struct gfx2d_gpu *gpu;
	uint32_t size;
	uint32_t wsize;
	uint32_t *start, *end, *cur, *tail;
	dma_addr_t paddr;
};

struct gfx2d_ringbuffer *gfx2d_ringbuffer_new(struct gfx2d_gpu *gpu);
void gfx2d_ringbuffer_destroy(struct gfx2d_ringbuffer *ring);

#endif /* __GFX2D_RINGBUFFER_H__ */
