/*
 * Physically contiguous memory allocator
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
#include <linux/dma-mapping.h> /* dma_zalloc_coherent, dma_free_coherent */
#include <linux/fs.h>
#include <linux/list.h>
#include <linux/slab.h>
#include <linux/uaccess.h>

#include "memalloc.h"

struct memalloc_drv {
	struct class *class;
	struct device *dev;
	int major;
	struct list_head opened; /* list of opened files (memalloc_file_context) */
	spinlock_t lock;
};

struct memalloc_file_context {
	struct list_head n;
	struct memalloc_drv *parent;
	struct list_head blocks; /* list of allocated blocks (mem_block) */
	spinlock_t lock;
	//unsigned int num_blocks;
};

struct memalloc_block {
	struct list_head n;
	dma_addr_t dma_handle;
	void *virt_addr;
	size_t size; /* page aligned */
	int method;
};

static struct memalloc_drv memalloc_ing = {
	.opened = LIST_HEAD_INIT(memalloc_ing.opened),
	.lock = __SPIN_LOCK_UNLOCKED(memalloc_ing.lock),
};

/**
 * allocate_large_block - Allocate a physically contiguous memory area
 *
 * @param dev driver device node
 * @param pp_block Result block struct address
 * @param p_size Wished size. Output is rounded up to page size (often something like 4k)
 *
 * @return 0 on success, negative value on error
 */
static int allocate_large_block(struct device *dev, struct memalloc_block **pp_block, size_t *p_size)
{
	struct memalloc_block *p;

	p = kmalloc(sizeof(struct memalloc_block), GFP_KERNEL);
	if (p == NULL) {
		dev_info(dev, "unable to alloc block struct.\n");
		return -ENOMEM;
	}

	p->size = PAGE_ALIGN(*p_size);

	/* Multiple of PAGE_SIZE */
	p->virt_addr = dma_alloc_coherent(dev, p->size, &p->dma_handle, GFP_KERNEL);
	if (!p->virt_addr) {
		dev_err(dev, "large alloc failed (%d)\n", p->size);
		kfree(p);
		return -ENOMEM;
	}

	dev_dbg(dev, "large alloc ok: VA=%p PA=0x%llx SZ=%d (requested %d)\n",
			p->virt_addr, (unsigned long long)p->dma_handle, p->size, *p_size);

	*p_size = p->size;
	*pp_block = p;
	return 0;
}

static int free_large_block(struct device *dev, struct memalloc_block *p)
{
	dev_dbg(dev, "large free: VA=%p\n", p->virt_addr);
	dma_free_coherent(dev, p->size, p->virt_addr, p->dma_handle);
	kfree(p);
	return 0;
}

static long memalloc_ioctl(struct file *filp, unsigned int cmd,
	unsigned long arg)
{
	struct memalloc_file_context *fc = filp->private_data;
	struct memalloc_drv *m = fc->parent;

	int ret = -EFAULT;
	MemallocParams mem_params;
	struct memalloc_block *p;
	size_t sz;

	if (!filp || arg == 0)
		return ret;

	switch (cmd) {
		case MEMALLOC_IOCXGETBUFFER:
			spin_lock(&fc->lock);
			if (copy_from_user(&mem_params, (MemallocParams *)arg, sizeof(mem_params)))
				dev_dbg(m->dev, "copy_from_user failed\n");

			sz = mem_params.size;
			ret = allocate_large_block(m->dev, &p, &sz);

			if (!ret) {
				mem_params.busAddress = (unsigned long)p->dma_handle; /* should be 64-bit! */
				mem_params.size = sz;
				if (copy_to_user((MemallocParams *)arg, &mem_params, sizeof(mem_params)))
					dev_dbg(m->dev, "copy_to_user failed\n");

				list_add(&p->n, &fc->blocks);
			}
			spin_unlock(&fc->lock);
			break;

		case MEMALLOC_IOCSFREEBUFFER:
			ret = -EINVAL;
			spin_lock(&fc->lock);
			__get_user(mem_params.busAddress, (unsigned long *)arg);

			/* find memalloc_block */
			list_for_each_entry(p, &fc->blocks, n) {
				if ((unsigned long)p->dma_handle == mem_params.busAddress) {
					list_del(&p->n);
					free_large_block(m->dev, p);
					ret = 0;
					break;
				}
			}
			spin_unlock(&fc->lock);
			break;

		default:
			ret = -ENOIOCTLCMD;
	}
	return ret;
}

static int memalloc_open(struct inode *inode, struct file *filp)
{
	struct memalloc_file_context *fc;
	struct memalloc_drv *m = &memalloc_ing;
	int dev = iminor(inode);

	if (dev != 0) {
		dev_warn(m->dev, "unsupported minor (%d).\n", dev);
		return -EINVAL;
	}

	fc = kmalloc(sizeof(struct memalloc_file_context), GFP_KERNEL);
	if (fc == NULL) {
		dev_err(m->dev, "unable to alloc struct.\n");
		return -ENOMEM;
	}

	INIT_LIST_HEAD(&fc->blocks);
	spin_lock_init(&fc->lock);
	fc->parent = m;

	filp->private_data = fc;

	spin_lock(&m->lock);
	list_add_tail(&fc->n, &m->opened);
	spin_unlock(&m->lock);

	dev_dbg(m->dev, "file open (%p)\n", fc);
	return 0;
}

static int memalloc_release(struct inode *inode, struct file *filp)
{
	struct memalloc_file_context *fc = filp->private_data;
	struct memalloc_drv *m = fc->parent;
	struct memalloc_block *p, *tmp;

	list_for_each_entry_safe(p, tmp, &fc->blocks, n) {
		list_del(&p->n);
		free_large_block(m->dev, p);
	}

	spin_lock(&m->lock);
	list_del(&fc->n);
	spin_unlock(&m->lock);

	kfree(fc);

	dev_dbg(m->dev, "file release (%p)\n", fc);
	return 0;
}

static const struct vm_operations_struct mmap_mem_ops = {
#ifdef CONFIG_HAVE_IOREMAP_PROT
	.access = generic_access_phys
#endif
};

/* This function is based on mmap_mem (drivers/char/mem.c) */
static int memalloc_mmap (struct file *filp, struct vm_area_struct *vma)
{
	struct memalloc_file_context *fc = filp->private_data;
	struct memalloc_block *p;

	int found = 0;
	size_t size = vma->vm_end - vma->vm_start;

	/* Is this a memory chunk provided by our driver ? */
	spin_lock(&fc->lock);
	list_for_each_entry(p, &fc->blocks, n) {
		if (((u64)p->dma_handle == ((u64)vma->vm_pgoff << PAGE_SHIFT)) &&
				(size <= p->size)) {
			found = 1;
			break;
		}
	}
	spin_unlock(&fc->lock);

	if (!found)
		return -EPERM;

	vma->vm_page_prot = phys_mem_access_prot(filp, vma->vm_pgoff,
						 size,
						 vma->vm_page_prot);

	vma->vm_ops = &mmap_mem_ops;

	/* Remap-pfn-range will mark the range VM_IO */
	if (remap_pfn_range(vma,
			    vma->vm_start,
			    vma->vm_pgoff,
			    size,
			    vma->vm_page_prot)) {
		return -EAGAIN;
	}
	return 0;
}

static struct file_operations memalloc_fops = {
	.owner          =	THIS_MODULE,
	.open           =	memalloc_open,
	.release        =	memalloc_release,
	.unlocked_ioctl =	memalloc_ioctl,
	.llseek         =	noop_llseek,
	.mmap           =	memalloc_mmap,
};

static int memalloc_init(void)
{
	struct memalloc_drv *m = &memalloc_ing;
	int ret;

	m->major = register_chrdev(0, "memalloc", &memalloc_fops);
	if (m->major < 0) {
		pr_err("failed to register character device\n");
		return m->major;
	}

	/* create /dev/memalloc */
	m->class = class_create(THIS_MODULE, "memalloc-cls");
	if (IS_ERR(m->class)) {
		ret = PTR_ERR(m->class);
		goto err;
	}
	m->dev = device_create(m->class, NULL, MKDEV(m->major, 0), NULL, "memalloc");
	if (IS_ERR(m->dev)) {
		ret = PTR_ERR(m->dev);
		class_destroy(m->class);
		goto err;
	}

	m->dev->coherent_dma_mask = DMA_BIT_MASK(32);
	dev_dbg(m->dev, "allocator with major = %d\n", m->major);

	return 0;
err:
	unregister_chrdev(m->major, "memalloc");
	return ret;

}
module_init(memalloc_init);

static void memalloc_exit(void)
{
	struct memalloc_drv *m = &memalloc_ing;

	device_destroy(m->class, MKDEV(m->major, 0));
	class_destroy(m->class);
	unregister_chrdev(m->major, "memalloc");
}
module_exit(memalloc_exit);

MODULE_AUTHOR("Hantro Products Oy");
MODULE_DESCRIPTION("Memory allocator for VDEC");
MODULE_LICENSE("GPL");
MODULE_VERSION("0.5");
