/*
 * On2/Hantro G1 decoder/pp driver. Single core version.
 *
 * Copyright (C) 2009  Hantro Products Oy.
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 2
 * of the License, or (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, write to the Free Software
 * Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
 */

#include <linux/module.h>
#include <linux/platform_device.h>
#include <linux/of.h>
#include <linux/clk.h>
#include <linux/interrupt.h>
#include <linux/io.h>
#include <linux/err.h>
#include <linux/miscdevice.h>
#include <linux/fs.h>
#include <linux/uaccess.h>
#include <linux/sched.h>

#include "hx170dec.h"
#include "at91_vdec.h"

#define VDEC_MAX_CORES                 1 /* number of cores of the hardware IP */
#define VDEC_NUM_REGS_DEC             60 /* number of registers of the Decoder part */
#define VDEC_NUM_REGS_PP              41 /* number of registers of the Post Processor part */
#define VDEC_DEC_FIRST_REG             0 /* first register (0-based) index */
#define VDEC_DEC_LAST_REG             59 /* last register (0-based) index */
#define VDEC_PP_FIRST_REG             60
#define VDEC_PP_LAST_REG             100

struct vdec_device {
	void __iomem *mmio_base;
	struct clk *clk;
	struct device *dev;
	int irq;
	int num_cores;
	unsigned long iobaseaddr;
	unsigned long iosize;
	wait_queue_head_t dec_wq;
	wait_queue_head_t pp_wq;
	bool dec_irq_done;
	bool pp_irq_done;
	struct semaphore dec_sem;
	struct semaphore pp_sem;
	struct file *dec_owner;
	struct file *pp_owner;
	u32 regs[VDEC_NUM_REGS_DEC + VDEC_NUM_REGS_PP];
};
static struct vdec_device *vdec6731_global;

static inline void vdec_writel(const struct vdec_device *p, unsigned offset, u32 val)
{
	writel(val, p->mmio_base + offset);
}

static inline u32 vdec_readl(const struct vdec_device *p, unsigned offset)
{
	return readl(p->mmio_base + offset);
}

/**
 * Write a range of registers. First register is assumed to be
 * "Interrupt Register" and will be written last.
 */
static int vdec_regs_write(struct vdec_device *p, int begin, int end,
		const struct core_desc *core)
{
	int i;

	if (copy_from_user(&p->regs[begin], core->regs, (end - begin + 1) * 4))
	{
		dev_err(p->dev, "%s: copy_from_user failed\n", __func__);
		return -EFAULT;
	}

	for (i = end; i >= begin; i--)
		vdec_writel(p, 4 * i, p->regs[i]);

	return 0;
}

/**
 * Read a range of registers [begin..end]
 */
static int vdec_regs_read(struct vdec_device *p, int begin, int end,
		const struct core_desc *core)
{
	int i;

	for (i = end; i >= begin; i--)
		p->regs[i] = vdec_readl(p, 4 * i);

	if (copy_to_user(core->regs, &p->regs[begin], (end - begin + 1) * 4))
	{
		dev_err(p->dev, "%s: copy_to_user failed\n", __func__);
		return -EFAULT;
	}

	return 0;
}

/**
 * Misc driver related
 */

static int vdec_misc_open(struct inode *inode, struct file *filp)
{
	struct vdec_device *p = vdec6731_global;
	filp->private_data = p;

	dev_dbg(p->dev, "open\n");
	clk_prepare_enable(p->clk);
	return 0;
}

static int vdec_misc_release(struct inode *inode, struct file *filp)
{
	struct vdec_device *p = filp->private_data;

	if (p->dec_owner == filp) {
		p->dec_irq_done = false;
		init_waitqueue_head(&p->dec_wq);
		sema_init(&p->dec_sem, VDEC_MAX_CORES);
		p->dec_owner = NULL;
	}

	if (p->pp_owner == filp) {
		p->pp_irq_done = false;
		init_waitqueue_head(&p->pp_wq);
		sema_init(&p->pp_sem, 1);
		p->pp_owner = NULL;
	}

	clk_disable_unprepare(p->clk);
	dev_dbg(p->dev, "release\n");
	return 0;
}

static long vdec_misc_ioctl(struct file *filp, unsigned int cmd,
    unsigned long arg)
{
	int ret = 0;
	void __user *argp = (void __user *)arg;
	struct vdec_device *p = vdec6731_global;
	struct core_desc core;
	u32 reg;

	switch (cmd) {
		case HX170DEC_IOX_ASIC_ID:
			reg = vdec_readl(p, VDEC_IDR);
			if (copy_to_user(argp, &reg, sizeof(u32)))
				ret = -EFAULT;
			break;

		case HX170DEC_IOC_MC_OFFSETS:
		case HX170DEC_IOCGHWOFFSET:
			if (copy_to_user(argp, &p->iobaseaddr, sizeof(p->iobaseaddr)))
				ret = -EFAULT;
			break;
		case HX170DEC_IOCGHWIOSIZE: /* in bytes */
			if (copy_to_user(argp, &p->iosize, sizeof(p->iosize)))
				ret = -EFAULT;
			break;
		case HX170DEC_IOC_MC_CORES:
			if (copy_to_user(argp, &p->num_cores, sizeof(p->num_cores)))
				ret = -EFAULT;
			break;

		case HX170DEC_IOCS_DEC_PUSH_REG:
			if (copy_from_user(&core, (void *)arg, sizeof(struct core_desc))) {
				dev_err(p->dev, "copy_from_user (dec push reg) failed\n");
				ret = -EFAULT;
			} else {
				/* Skip VDEC_IDR (ID Register, ro) */
				core.regs++; // core.size -= 4;
				ret = vdec_regs_write(p, VDEC_DEC_FIRST_REG + 1, VDEC_DEC_LAST_REG, &core);
			}
			break;
		case HX170DEC_IOCS_PP_PUSH_REG:
			if (copy_from_user(&core, (void *)arg, sizeof(struct core_desc))) {
				dev_err(p->dev, "copy_from_user (pp push reg) failed\n");
				ret = -EFAULT;
			} else {
				/* Don't consider the 5 lastest registers (ro or unused) */
				ret = vdec_regs_write(p, VDEC_PP_FIRST_REG, VDEC_PP_LAST_REG - 5, &core);
			}
			break;

		case HX170DEC_IOCS_DEC_PULL_REG:
			if (copy_from_user(&core, (void *)arg, sizeof(struct core_desc))) {
				dev_err(p->dev, "copy_from_user (dec pull reg) failed\n");
				ret = -EFAULT;
			} else {
				ret = vdec_regs_read(p, VDEC_DEC_FIRST_REG, VDEC_DEC_LAST_REG, &core);
			}
			break;

		case HX170DEC_IOCS_PP_PULL_REG:
			if (copy_from_user(&core, (void*)arg, sizeof(struct core_desc))) {
				dev_err(p->dev, "copy_from_user (pp pull reg) failed\n");
				ret = -EFAULT;
			} else {
				ret = vdec_regs_read(p, VDEC_PP_FIRST_REG, VDEC_PP_LAST_REG, &core);
			}
			break;

		case HX170DEC_IOCX_DEC_WAIT:
			if (copy_from_user(&core, (void *)arg, sizeof(struct core_desc))) {
				dev_err(p->dev, "copy_from_user (dec wait) failed\n");
				ret = -EFAULT;
			} else {
				ret = wait_event_interruptible(p->dec_wq, p->dec_irq_done);
				p->dec_irq_done = false;
				if (unlikely(ret != 0)) {
					dev_err(p->dev, "wait_event_interruptible dec error %d\n", ret);
				} else {
					/* Update dec registers */
					ret = vdec_regs_read(p, VDEC_DEC_FIRST_REG, VDEC_DEC_LAST_REG, &core);
				}
			}
			break;
		case HX170DEC_IOCX_PP_WAIT:
			if (copy_from_user(&core, (void *)arg, sizeof(struct core_desc))) {
				dev_err(p->dev, "copy_from_user (pp wait) failed\n");
				ret = -EFAULT;
			} else {
				ret = wait_event_interruptible(p->pp_wq, p->pp_irq_done);
				p->pp_irq_done = false;
				if (unlikely(ret != 0)) {
					dev_err(p->dev, "wait_event_interruptible pp error %d\n", ret);
				} else {
					/* Update pp registers */
					ret = vdec_regs_read(p, VDEC_PP_FIRST_REG, VDEC_PP_LAST_REG, &core);
				}
			}
			break;

		case HX170DEC_IOCH_DEC_RESERVE:
			if (likely(down_interruptible(&p->dec_sem) == 0)) {
				p->dec_owner = filp;
				ret = 0; /* core id */
				dev_dbg(p->dev, "down dec_sem (core id %d)\n", ret);
			} else {
				dev_err(p->dev, "down_interruptible dec error\n");
				ret = -ERESTARTSYS;
			}
			break;
		case HX170DEC_IOCT_DEC_RELEASE:
			dev_dbg(p->dev, "up dec_sem\n");
			p->dec_owner = NULL;
			up(&p->dec_sem);
			break;

		case HX170DEC_IOCQ_PP_RESERVE:
			if (likely(down_interruptible(&p->pp_sem) == 0)) {
				p->pp_owner = filp;
				ret = 0; /* core id */
				dev_dbg(p->dev, "down pp_sem (core id %d)\n", ret);
			} else {
				dev_err(p->dev, "down_interruptible pp error\n");
				ret = -ERESTARTSYS;
			}
			break;
		case HX170DEC_IOCT_PP_RELEASE:
			dev_dbg(p->dev, "up pp_sem\n");
			p->pp_owner = NULL;
			up(&p->pp_sem);
			break;

		default:
			dev_warn(p->dev, "unknown ioctl %x\n", cmd);
			ret = -EINVAL;
	}
	return ret;
}

const struct file_operations vdec_misc_fops = {
	.owner          =	THIS_MODULE,
	.llseek         =	no_llseek,
	.open           =	vdec_misc_open,
	.release        =	vdec_misc_release,
	.unlocked_ioctl =	vdec_misc_ioctl,
};

static struct miscdevice vdec_misc_device = {
	MISC_DYNAMIC_MINOR,
	"vdec",
	&vdec_misc_fops
};

/*
 * Platform driver related
 */

/* Should we use spin_lock_irqsave here? */
static irqreturn_t vdec_isr(int irq, void *dev_id)
{
	struct vdec_device *p = dev_id;
	u32 irq_status_dec, irq_status_pp;
	int handled = 0;

	/* interrupt status register read */
	irq_status_dec = vdec_readl(p, VDEC_DIR);
	if (irq_status_dec & VDEC_DIR_ISET) {
		/* Clear IRQ */
		vdec_writel(p, VDEC_DIR, irq_status_dec & ~VDEC_DIR_ISET);

		p->dec_irq_done = true;
		wake_up_interruptible(&p->dec_wq);
		handled++;
	}

	irq_status_pp = vdec_readl(p, VDEC_PPIR);
	if (irq_status_pp & VDEC_PPIR_ISET) {
		/* Clear IRQ */
		vdec_writel(p, VDEC_PPIR, irq_status_pp & ~VDEC_PPIR_ISET);

		p->pp_irq_done = true;
		wake_up_interruptible(&p->pp_wq);
		handled++;
	}

	if (handled == 0) {
		dev_warn(p->dev, "Spurious IRQ (DIR=%08x PPIR=%08x)\n", \
				irq_status_dec, irq_status_pp);
		return IRQ_NONE;
	}

	return IRQ_HANDLED;
}

static int __init vdec_probe(struct platform_device *pdev)
{
	struct vdec_device *p;
	struct resource *res;
	int ret;
	u32 hwid;

	/* Allocate private data */
	p = devm_kzalloc(&pdev->dev, sizeof(struct vdec_device), GFP_KERNEL);
	if (!p) {
		dev_dbg(&pdev->dev, "out of memory\n");
		return -ENOMEM;
	}

	p->dev = &pdev->dev;
	platform_set_drvdata(pdev, p);

	res = platform_get_resource(pdev, IORESOURCE_MEM, 0);
	p->mmio_base = devm_ioremap_resource(&pdev->dev, res);
	if (IS_ERR(p->mmio_base))
		return PTR_ERR(p->mmio_base);

	p->clk = devm_clk_get(&pdev->dev, "vdec_clk");
	if (IS_ERR(p->clk)) {
		dev_err(&pdev->dev, "no vdec_clk clock defined\n");
		return -ENXIO;
	}

	p->irq = platform_get_irq(pdev, 0);
	if (!p->irq) {
		dev_err(&pdev->dev, "could not get irq\n");
		return -ENXIO;
	}

	ret = devm_request_irq(&pdev->dev, p->irq, vdec_isr,
			0, pdev->name, p);
	if (ret) {
		dev_err(&pdev->dev, "unable to request VDEC irq\n");
		return ret;
	}

	/* Register the miscdevice */
	ret = misc_register(&vdec_misc_device);
	if (ret) {
		dev_err(&pdev->dev, "unable to register miscdevice\n");
		return ret;
	}

	p->num_cores = VDEC_MAX_CORES;
	p->iosize = resource_size(res);
	p->iobaseaddr = res->start;
	vdec6731_global = p;

	p->dec_irq_done = false;
	p->pp_irq_done = false;
	p->dec_owner = NULL;
	p->pp_owner = NULL;
	init_waitqueue_head(&p->dec_wq);
	init_waitqueue_head(&p->pp_wq);
	sema_init(&p->dec_sem, VDEC_MAX_CORES);
	sema_init(&p->pp_sem, 1);

	ret = clk_prepare_enable(p->clk);
	if (ret) {
		dev_err(&pdev->dev, "unable to prepare and enable clock\n");
		misc_deregister(&vdec_misc_device);
		return ret;
	}

	dev_info(&pdev->dev, "VDEC controller at 0x%p, irq = %d, misc_minor = %d\n",
			p->mmio_base, p->irq, vdec_misc_device.minor);

	/* Reset Asic (just in case..) */
	vdec_writel(p, VDEC_DIR, VDEC_DIR_ID | VDEC_DIR_ABORT);
	vdec_writel(p, VDEC_PPIR, VDEC_PPIR_ID);

	hwid = vdec_readl(p, VDEC_IDR);
	clk_disable_unprepare(p->clk);

	dev_warn(&pdev->dev, "Product ID: %#x (revision %d.%d.%d)\n", \
			(hwid & VDEC_IDR_PROD_ID) >> 16,
			(hwid & VDEC_IDR_MAJOR_VER) >> 12,
			(hwid & VDEC_IDR_MINOR_VER) >> 4,
			(hwid & VDEC_IDR_BUILD_VER));
	return 0;
}

static int __exit vdec_remove(struct platform_device *pdev)
{
	platform_set_drvdata(pdev, NULL);
	misc_deregister(&vdec_misc_device);
	return 0;
}

static const struct of_device_id vdec_of_match[] = {
	{ .compatible = "on2,sama5d4-g1", .data = NULL },
	{},
};
MODULE_DEVICE_TABLE(of, vdec_of_match);

static struct platform_driver vdec_of_driver = {
	.driver		= {
		.name	= "atmel-vdec",
		.owner	= THIS_MODULE,
		.of_match_table	= vdec_of_match,
	},
	.remove		= vdec_remove,
};

module_platform_driver_probe(vdec_of_driver, vdec_probe);

MODULE_AUTHOR("Hantro Products Oy");
MODULE_DESCRIPTION("G1 decoder/pp driver");
MODULE_LICENSE("GPL");
MODULE_VERSION("0.4");
MODULE_ALIAS("platform:vdec");
