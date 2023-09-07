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
#include <drm/drmP.h>
#include <drm/drm_gem.h>
#include <linux/spinlock.h>

#define REG_GFX2D_GC 0x00
#define REG_GFX2D_GE 0x04
#define REG_GFX2D_GE_ENABLE (1<<0)
#define REG_GFX2D_GD 0x08
#define REG_GFX2D_GD_DISABLE (1<<0)
#define REG_GFX2D_GD_WFERES (1<<8)
#define REG_GFX2D_GS 0x0C
#define REG_GFX2D_GS_STATUS (1<<0)
#define REG_GFX2D_GS_BUSY (1<<4)
#define REG_GFX2D_GS_WFEIP (1<<8)
#define REG_GFX2D_IE 0x10
#define REG_GFX2D_IE_RBEMPTY (1<<0)
#define REG_GFX2D_IE_EXEND (1<<1)
#define REG_GFX2D_IE_RERR (1<<2)
#define REG_GFX2D_IE_BERR (1<<3)
#define REG_GFX2D_IE_IERR (1<<4)
#define REG_GFX2D_ID 0x14
#define REG_GFX2D_IM 0x18
#define REG_GFX2D_IS 0x1C
#define REG_GFX2D_IS_RBEMPTY (1<<0)
#define REG_GFX2D_IS_EXEND (1<<1)
#define REG_GFX2D_IS_RERR (1<<2)
#define REG_GFX2D_IS_BERR (1<<3)
#define REG_GFX2D_IS_IERR (1<<4)
#define REG_GFX2D_PC0 0x20
#define REG_GFX2D_MC0 0x24
#define REG_GFX2D_PC1 0x28
#define REG_GFX2D_MC1 0x2C
#define REG_GFX2D_BASE 0x30
#define REG_GFX2D_LEN 0x34
#define REG_GFX2D_HEAD 0x38
#define REG_GFX2D_TAIL 0x3C
#define REG_GFX2D_VERSION 0xFC

static inline void gpu_write(struct gfx2d_gpu *gpu, u32 reg, u32 data)
{
	writel(data, gpu->mmio + reg);
}

static inline u32 gpu_read(struct gfx2d_gpu *gpu, u32 reg)
{
	return readl(gpu->mmio + reg);
}

static inline uint32_t get_wptr(struct gfx2d_gpu *gpu)
{
	return gpu_read(gpu, REG_GFX2D_HEAD);
}

static inline uint32_t get_rptr(struct gfx2d_gpu *gpu)
{
	return gpu_read(gpu, REG_GFX2D_TAIL);
}

static int gfx2d_hw_init(struct gfx2d_gpu *gpu)
{
	gpu->version = gpu_read(gpu, REG_GFX2D_VERSION) & 0x7ff;
	gpu->mfn = (gpu_read(gpu, REG_GFX2D_VERSION) >> 16) & 0x7;

	dev_dbg(&gpu->pdev->dev, "GPU version %d (%d)", gpu->version, gpu->mfn);

	/* disable rb */
	gpu_write(gpu, REG_GFX2D_GD, REG_GFX2D_GD_DISABLE);

	/* initialize head/tail pointers */
	gpu_write(gpu, REG_GFX2D_HEAD, 0);
	gpu_write(gpu, REG_GFX2D_TAIL, 0);

	/* configure */
	gpu_write(gpu, REG_GFX2D_BASE, gpu->rb->paddr);
	gpu_write(gpu, REG_GFX2D_LEN, (gpu->rb->size / 256) - 1);

	/* enable rb */
	gpu_write(gpu, REG_GFX2D_GE, REG_GFX2D_GE_ENABLE);

	/* enable interrupt */
	gpu_write(gpu, REG_GFX2D_IE, REG_GFX2D_IE_EXEND | REG_GFX2D_IE_RBEMPTY | REG_GFX2D_IE_IERR);

	return 0;
}

static inline unsigned ring_freewords(struct gfx2d_ringbuffer *ring)
{
    if (ring->cur >= ring->tail) {
        return ring->end - ring->cur + ring->tail - ring->start;
    } else {
        return ring->tail - ring->cur;
    }
}

int gfx2d_submit(struct gfx2d_gpu *gpu, uint32_t* buf, uint32_t nwords)
{
	struct gfx2d_ringbuffer *ring = gpu->rb;
	int ret, size, nremainingwords, freewords;

	if (!nwords)
		return 0;

	/*
	 * This case should not happend, libm2d should submit as soon as
	 * there are more than 512 words.
	 */
	if (nwords > ring->wsize)
		return -EINVAL;

	freewords = ring_freewords(ring);
	/*
	 * One extra free word needed because is head and tail are at the same
	 * position, the hardware assumes the buffer is empty and not full.
	 */
	if (freewords < nwords + 1) {
		/* Check the hardware tail position. */
		ring->tail = ring->start + gpu_read(gpu, REG_GFX2D_TAIL);
		freewords = ring_freewords(ring);
		/*
		 * If we still have not enough place, let the hardware drain the
		 * command buffer.
		 */
		if (freewords < nwords + 1) {
			while (gpu_read(gpu, REG_GFX2D_GS) & REG_GFX2D_GS_BUSY);
			ring->tail = ring->start + gpu_read(gpu, REG_GFX2D_TAIL);
		}
	}

	/* Copy the commands to the ring buffer. */
	size = ring->end - ring->cur;
	if ((nwords > size)) {
		ret = copy_from_user(ring->cur, (const void __user *)buf, size * 4);
		if (ret) {
			ret = -EFAULT;
			return ret;
		}

		nremainingwords = nwords - size;
		ret = copy_from_user(ring->start, (const void __user *)(buf + size), nremainingwords * 4);
		if (ret) {
			ret = -EFAULT;
			return ret;
		}

		ring->cur = ring->start + nwords - size;
	} else {
		ret = copy_from_user(ring->cur, (const void __user *)buf, nwords * 4);
		if (ret) {
			ret = -EFAULT;
			return ret;
		}

		ring->cur += nwords;
	}

	/* Ensure writes to ringbuffer have hit system memory. */
	mb();

	/* Start to execute the commands. */
	gpu_write(gpu, REG_GFX2D_HEAD, ring->cur - ring->start);

	return 0;
}

int gfx2d_flush(struct gfx2d_gpu *gpu)
{
	struct gfx2d_ringbuffer *ring = gpu->rb;

	/* Ensure writes to ringbuffer have hit system memory. */
	mb();

	if (ring->cur != ring->tail)
		gpu_write(gpu, REG_GFX2D_HEAD, ring->cur - ring->start);

	while (gpu_read(gpu, REG_GFX2D_GS) & REG_GFX2D_GS_BUSY);

	return 0;
}

void gfx2d_idle(struct gfx2d_gpu *gpu)
{
	uint32_t wptr = get_wptr(gpu);
	int ret;

	/* wait for CP to drain ringbuffer: */
	ret = spin_until(get_rptr(gpu) == wptr);

	if (ret)
		DRM_ERROR("timeout waiting to drain ringbuffer\n");
}

#ifdef CONFIG_DEBUG_FS
void gfx2d_show(struct gfx2d_gpu *gpu, struct seq_file *m)
{
	int i;

	seq_printf(m, "revision:      %d (%d)\n", gpu->version, gpu->mfn);
	seq_printf(m, "rptr:          %d\n", get_rptr(gpu));
	seq_printf(m, "wptr:          %d\n", get_wptr(gpu));
	//seq_printf(m, "rb freewords:  %d\n", ring_freewords(gpu));

	seq_printf(m, "mmio:\n");
	for (i = 0; i != REG_GFX2D_VERSION; i+=4) {
		uint32_t val = gpu_read(gpu, i);
		seq_printf(m, "IO:R %08x %08x\n", i, val);
	}
}
#endif

static irqreturn_t gfx2d_irq(int irq, void *data)
{
	struct gfx2d_gpu *gpu = data;
	uint32_t status;
	int ret;

	status = gpu_read(gpu, REG_GFX2D_IS);

	if (status & REG_GFX2D_IS_IERR) {
		/*
		 * When an illegal instruction IRQ is rised, the GPU is
		 * automatically disabled.
		 */
		dev_err(&gpu->pdev->dev,"illegal instruction, reinitialize the command buffer\n");
		ret = gfx2d_hw_init(gpu);
		if (ret)
			dev_err(&gpu->pdev->dev,"can't reinitialize the command buffer\n");

		return IRQ_HANDLED;
	}

	if (status & REG_GFX2D_IS_EXEND) {
	}

	if (status & REG_GFX2D_IS_RBEMPTY) {
	}

	return IRQ_HANDLED;
}

static void gfx2d_destroy(struct gfx2d_gpu *gpu)
{
	gfx2d_idle(gpu);
	gpu_write(gpu, REG_GFX2D_GD, REG_GFX2D_GD_DISABLE);
	gfx2d_gpu_cleanup(gpu);
	kfree(gpu);
}

struct gfx2d_gpu *gfx2d_gpu_init(struct device *dev)
{
	struct platform_device *pdev = to_platform_device(dev);
	struct gfx2d_gpu *gpu = NULL;
	struct resource *res;
	int ret;

	if (!pdev) {
		dev_err(dev, "no gfx2d device\n");
		ret = -ENXIO;
		goto fail;
	}

	gpu = kzalloc(sizeof(*gpu), GFP_KERNEL);
	if (!gpu) {
		ret = -ENOMEM;
		goto fail;
	}

	gpu->pdev = pdev;

	gpu->rb = gfx2d_ringbuffer_new(gpu);
	if (IS_ERR(gpu->rb)) {
		ret = PTR_ERR(gpu->rb);
		dev_err(dev, "failed to allocate ringbuffer: %d\n", ret);
		goto fail;
	}

	res = platform_get_resource(pdev, IORESOURCE_MEM, 0);
	gpu->mmio = devm_ioremap_resource(&pdev->dev, res);
	if (IS_ERR(gpu->mmio)) {
		ret = PTR_ERR(gpu->mmio);
		goto fail;
	}

	gpu->irq = platform_get_irq(pdev, 0);
	if (gpu->irq < 0) {
		ret = gpu->irq;
		dev_err(dev, "failed to get irq: %d\n", ret);
		goto fail;
	}

	ret = devm_request_irq(&pdev->dev, gpu->irq, gfx2d_irq,
			       IRQF_TRIGGER_HIGH, "gfx2d", gpu);
	if (ret) {
		dev_err(dev, "failed to request IRQ%u: %d\n", gpu->irq, ret);
		goto fail;
	}

	gpu->periph_clk = devm_clk_get(dev, "periph_clk");
	DBG("periph_clk: %p", gpu->periph_clk);
	if (IS_ERR(gpu->periph_clk))
		gpu->periph_clk = NULL;
	else {
		ret = clk_prepare_enable(gpu->periph_clk);
		if (ret) {
			dev_err(dev, "failed to enable periph_clk\n");
			goto fail;
		}
	}

	disable_irq(gpu->irq);
	ret = gfx2d_hw_init(gpu);
	if (ret)
		goto fail;
	enable_irq(gpu->irq);

	return gpu;

fail:
	if (gpu)
		gfx2d_destroy(gpu);

	return ERR_PTR(ret);
}

void gfx2d_gpu_cleanup(struct gfx2d_gpu *gpu)
{
	gfx2d_ringbuffer_destroy(gpu->rb);
}
