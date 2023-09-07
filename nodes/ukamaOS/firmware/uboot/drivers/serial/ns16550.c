/*
 * COM1 NS16550 support
 * originally from linux source (arch/powerpc/boot/ns16550.c)
 * modified to use CONFIG_SYS_ISA_MEM and new defines
 */

#include <common.h>
#include <clk.h>
#include <dm.h>
#include <errno.h>
#include <ns16550.h>
#include <serial.h>
#include <watchdog.h>
#include <linux/types.h>
#include <asm/io.h>

DECLARE_GLOBAL_DATA_PTR;

#define UART_LCRVAL UART_LCR_8N1		/* 8 data, 1 stop, no parity */
#define UART_MCRVAL (UART_MCR_DTR | \
		     UART_MCR_RTS)		/* RTS/DTR */

#ifndef CONFIG_DM_SERIAL
#ifdef CONFIG_SYS_NS16550_PORT_MAPPED
#define serial_out(x, y)	outb(x, (ulong)y)
#define serial_in(y)		inb((ulong)y)
#elif defined(CONFIG_SYS_NS16550_MEM32) && (CONFIG_SYS_NS16550_REG_SIZE > 0)
#define serial_out(x, y)	out_be32(y, x)
#define serial_in(y)		in_be32(y)
#elif defined(CONFIG_SYS_NS16550_MEM32) && (CONFIG_SYS_NS16550_REG_SIZE < 0)
#define serial_out(x, y)	out_le32(y, x)
#define serial_in(y)		in_le32(y)
#else
#define serial_out(x, y)	writeb(x, y)
#define serial_in(y)		readb(y)
#endif
#endif /* !CONFIG_DM_SERIAL */

#if defined(CONFIG_SOC_KEYSTONE)
#define UART_REG_VAL_PWREMU_MGMT_UART_DISABLE   0
#define UART_REG_VAL_PWREMU_MGMT_UART_ENABLE ((1 << 14) | (1 << 13) | (1 << 0))
#undef UART_MCRVAL
#ifdef CONFIG_SERIAL_HW_FLOW_CONTROL
#define UART_MCRVAL             (UART_MCR_RTS | UART_MCR_AFE)
#else
#define UART_MCRVAL             (UART_MCR_RTS)
#endif
#endif

#ifndef CONFIG_SYS_NS16550_IER
#define CONFIG_SYS_NS16550_IER  0x00
#endif /* CONFIG_SYS_NS16550_IER */

static inline void serial_out_shift(void *addr, int shift, int value)
{
#ifdef CONFIG_SYS_NS16550_PORT_MAPPED
	outb(value, (ulong)addr);
#elif defined(CONFIG_SYS_NS16550_MEM32) && defined(CONFIG_SYS_LITTLE_ENDIAN)
	out_le32(addr, value);
#elif defined(CONFIG_SYS_NS16550_MEM32) && defined(CONFIG_SYS_BIG_ENDIAN)
	out_be32(addr, value);
#elif defined(CONFIG_SYS_NS16550_MEM32)
	writel(value, addr);
#elif defined(CONFIG_SYS_BIG_ENDIAN)
	writeb(value, addr + (1 << shift) - 1);
#else
	writeb(value, addr);
#endif
}

static inline int serial_in_shift(void *addr, int shift)
{
#ifdef CONFIG_SYS_NS16550_PORT_MAPPED
	return inb((ulong)addr);
#elif defined(CONFIG_SYS_NS16550_MEM32) && defined(CONFIG_SYS_LITTLE_ENDIAN)
	return in_le32(addr);
#elif defined(CONFIG_SYS_NS16550_MEM32) && defined(CONFIG_SYS_BIG_ENDIAN)
	return in_be32(addr);
#elif defined(CONFIG_SYS_NS16550_MEM32)
	return readl(addr);
#elif defined(CONFIG_SYS_BIG_ENDIAN)
	return readb(addr + (1 << shift) - 1);
#else
	return readb(addr);
#endif
}

#ifdef CONFIG_DM_SERIAL

#ifndef CONFIG_SYS_NS16550_CLK
#define CONFIG_SYS_NS16550_CLK  0
#endif

static void ns16550_writeb(NS16550_t port, int offset, int value)
{
	struct ns16550_platdata *plat = port->plat;
	unsigned char *addr;

	offset *= 1 << plat->reg_shift;
	addr = (unsigned char *)plat->base + offset;

	/*
	 * As far as we know it doesn't make sense to support selection of
	 * these options at run-time, so use the existing CONFIG options.
	 */
	serial_out_shift(addr + plat->reg_offset, plat->reg_shift, value);
}

static int ns16550_readb(NS16550_t port, int offset)
{
	struct ns16550_platdata *plat = port->plat;
	unsigned char *addr;

	offset *= 1 << plat->reg_shift;
	addr = (unsigned char *)plat->base + offset;

	return serial_in_shift(addr + plat->reg_offset, plat->reg_shift);
}

static u32 ns16550_getfcr(NS16550_t port)
{
	struct ns16550_platdata *plat = port->plat;

	return plat->fcr;
}

/* We can clean these up once everything is moved to driver model */
#define serial_out(value, addr)	\
	ns16550_writeb(com_port, \
		(unsigned char *)addr - (unsigned char *)com_port, value)
#define serial_in(addr) \
	ns16550_readb(com_port, \
		(unsigned char *)addr - (unsigned char *)com_port)
#else
static u32 ns16550_getfcr(NS16550_t port)
{
	return UART_FCR_DEFVAL;
}
#endif

int ns16550_calc_divisor(NS16550_t port, int clock, int baudrate)
{
	const unsigned int mode_x_div = 16;

	return DIV_ROUND_CLOSEST(clock, mode_x_div * baudrate);
}

static void NS16550_setbrg(NS16550_t com_port, int baud_divisor)
{
	serial_out(UART_LCR_BKSE | UART_LCRVAL, &com_port->lcr);
	serial_out(baud_divisor & 0xff, &com_port->dll);
	serial_out((baud_divisor >> 8) & 0xff, &com_port->dlm);
	serial_out(UART_LCRVAL, &com_port->lcr);
}

void NS16550_init(NS16550_t com_port, int baud_divisor)
{
#if (defined(CONFIG_SPL_BUILD) && \
		(defined(CONFIG_OMAP34XX) || defined(CONFIG_OMAP44XX)))
	/*
	 * On some OMAP3/OMAP4 devices when UART3 is configured for boot mode
	 * before SPL starts only THRE bit is set. We have to empty the
	 * transmitter before initialization starts.
	 */
	if ((serial_in(&com_port->lsr) & (UART_LSR_TEMT | UART_LSR_THRE))
	     == UART_LSR_THRE) {
		if (baud_divisor != -1)
			NS16550_setbrg(com_port, baud_divisor);
		serial_out(0, &com_port->mdr1);
	}
#endif

	while (!(serial_in(&com_port->lsr) & UART_LSR_TEMT))
		;

	serial_out(CONFIG_SYS_NS16550_IER, &com_port->ier);
#if defined(CONFIG_ARCH_OMAP2PLUS)
	serial_out(0x7, &com_port->mdr1);	/* mode select reset TL16C750*/
#endif
	serial_out(UART_MCRVAL, &com_port->mcr);
	serial_out(ns16550_getfcr(com_port), &com_port->fcr);
	if (baud_divisor != -1)
		NS16550_setbrg(com_port, baud_divisor);
#if defined(CONFIG_ARCH_OMAP2PLUS) || defined(CONFIG_SOC_DA8XX)
	/* /16 is proper to hit 115200 with 48MHz */
	serial_out(0, &com_port->mdr1);
#endif
#if defined(CONFIG_SOC_KEYSTONE)
	serial_out(UART_REG_VAL_PWREMU_MGMT_UART_ENABLE, &com_port->regC);
#endif
}

#ifndef CONFIG_NS16550_MIN_FUNCTIONS
void NS16550_reinit(NS16550_t com_port, int baud_divisor)
{
	serial_out(CONFIG_SYS_NS16550_IER, &com_port->ier);
	NS16550_setbrg(com_port, 0);
	serial_out(UART_MCRVAL, &com_port->mcr);
	serial_out(ns16550_getfcr(com_port), &com_port->fcr);
	NS16550_setbrg(com_port, baud_divisor);
}
#endif /* CONFIG_NS16550_MIN_FUNCTIONS */

void NS16550_putc(NS16550_t com_port, char c)
{
	while ((serial_in(&com_port->lsr) & UART_LSR_THRE) == 0)
		;
	serial_out(c, &com_port->thr);

	/*
	 * Call watchdog_reset() upon newline. This is done here in putc
	 * since the environment code uses a single puts() to print the complete
	 * environment upon "printenv". So we can't put this watchdog call
	 * in puts().
	 */
	if (c == '\n')
		WATCHDOG_RESET();
}

#ifndef CONFIG_NS16550_MIN_FUNCTIONS
char NS16550_getc(NS16550_t com_port)
{
	while ((serial_in(&com_port->lsr) & UART_LSR_DR) == 0) {
#if !defined(CONFIG_SPL_BUILD) && defined(CONFIG_USB_TTY)
		extern void usbtty_poll(void);
		usbtty_poll();
#endif
		WATCHDOG_RESET();
	}
	return serial_in(&com_port->rbr);
}

int NS16550_tstc(NS16550_t com_port)
{
	return (serial_in(&com_port->lsr) & UART_LSR_DR) != 0;
}

#endif /* CONFIG_NS16550_MIN_FUNCTIONS */

#ifdef CONFIG_DEBUG_UART_NS16550

#include <debug_uart.h>

static inline void _debug_uart_init(void)
{
	struct NS16550 *com_port = (struct NS16550 *)CONFIG_DEBUG_UART_BASE;
	int baud_divisor;

	/*
	 * We copy the code from above because it is already horribly messy.
	 * Trying to refactor to nicely remove the duplication doesn't seem
	 * feasible. The better fix is to move all users of this driver to
	 * driver model.
	 */
	baud_divisor = ns16550_calc_divisor(com_port, CONFIG_DEBUG_UART_CLOCK,
					    CONFIG_BAUDRATE);
	serial_dout(&com_port->ier, CONFIG_SYS_NS16550_IER);
	serial_dout(&com_port->mcr, UART_MCRVAL);
	serial_dout(&com_port->fcr, UART_FCR_DEFVAL);

	serial_dout(&com_port->lcr, UART_LCR_BKSE | UART_LCRVAL);
	serial_dout(&com_port->dll, baud_divisor & 0xff);
	serial_dout(&com_port->dlm, (baud_divisor >> 8) & 0xff);
	serial_dout(&com_port->lcr, UART_LCRVAL);
}

static inline void _debug_uart_putc(int ch)
{
	struct NS16550 *com_port = (struct NS16550 *)CONFIG_DEBUG_UART_BASE;

	while (!(serial_din(&com_port->lsr) & UART_LSR_THRE))
		;
	serial_dout(&com_port->thr, ch);
}

DEBUG_UART_FUNCS

#endif

#ifdef CONFIG_DEBUG_UART_OMAP

#include <debug_uart.h>

static inline void _debug_uart_init(void)
{
	struct NS16550 *com_port = (struct NS16550 *)CONFIG_DEBUG_UART_BASE;
	int baud_divisor;

	baud_divisor = ns16550_calc_divisor(com_port, CONFIG_DEBUG_UART_CLOCK,
					    CONFIG_BAUDRATE);
	serial_dout(&com_port->ier, CONFIG_SYS_NS16550_IER);
	serial_dout(&com_port->mdr1, 0x7);
	serial_dout(&com_port->mcr, UART_MCRVAL);
	serial_dout(&com_port->fcr, UART_FCR_DEFVAL);

	serial_dout(&com_port->lcr, UART_LCR_BKSE | UART_LCRVAL);
	serial_dout(&com_port->dll, baud_divisor & 0xff);
	serial_dout(&com_port->dlm, (baud_divisor >> 8) & 0xff);
	serial_dout(&com_port->lcr, UART_LCRVAL);
	serial_dout(&com_port->mdr1, 0x0);
}

static inline void _debug_uart_putc(int ch)
{
	struct NS16550 *com_port = (struct NS16550 *)CONFIG_DEBUG_UART_BASE;

	while (!(serial_din(&com_port->lsr) & UART_LSR_THRE))
		;
	serial_dout(&com_port->thr, ch);
}

DEBUG_UART_FUNCS

#endif

#ifdef CONFIG_DM_SERIAL
static int ns16550_serial_putc(struct udevice *dev, const char ch)
{
	struct NS16550 *const com_port = dev_get_priv(dev);

	if (!(serial_in(&com_port->lsr) & UART_LSR_THRE))
		return -EAGAIN;
	serial_out(ch, &com_port->thr);

	/*
	 * Call watchdog_reset() upon newline. This is done here in putc
	 * since the environment code uses a single puts() to print the complete
	 * environment upon "printenv". So we can't put this watchdog call
	 * in puts().
	 */
	if (ch == '\n')
		WATCHDOG_RESET();

	return 0;
}

static int ns16550_serial_pending(struct udevice *dev, bool input)
{
	struct NS16550 *const com_port = dev_get_priv(dev);

	if (input)
		return (serial_in(&com_port->lsr) & UART_LSR_DR) ? 1 : 0;
	else
		return (serial_in(&com_port->lsr) & UART_LSR_THRE) ? 0 : 1;
}

static int ns16550_serial_getc(struct udevice *dev)
{
	struct NS16550 *const com_port = dev_get_priv(dev);

	if (!(serial_in(&com_port->lsr) & UART_LSR_DR))
		return -EAGAIN;

	return serial_in(&com_port->rbr);
}

static int ns16550_serial_setbrg(struct udevice *dev, int baudrate)
{
	struct NS16550 *const com_port = dev_get_priv(dev);
	struct ns16550_platdata *plat = com_port->plat;
	int clock_divisor;

	clock_divisor = ns16550_calc_divisor(com_port, plat->clock, baudrate);

	NS16550_setbrg(com_port, clock_divisor);

	return 0;
}

int ns16550_serial_probe(struct udevice *dev)
{
	struct NS16550 *const com_port = dev_get_priv(dev);

	com_port->plat = dev_get_platdata(dev);
	NS16550_init(com_port, -1);

	return 0;
}

#if CONFIG_IS_ENABLED(OF_CONTROL)
enum {
	PORT_NS16550 = 0,
	PORT_JZ4780,
};
#endif

#if CONFIG_IS_ENABLED(OF_CONTROL) && !CONFIG_IS_ENABLED(OF_PLATDATA)
int ns16550_serial_ofdata_to_platdata(struct udevice *dev)
{
	struct ns16550_platdata *plat = dev->platdata;
	const u32 port_type = dev_get_driver_data(dev);
	fdt_addr_t addr;
	struct clk clk;
	int err;

	/* try Processor Local Bus device first */
	addr = dev_read_addr(dev);
#if defined(CONFIG_PCI) && defined(CONFIG_DM_PCI)
	if (addr == FDT_ADDR_T_NONE) {
		/* then try pci device */
		struct fdt_pci_addr pci_addr;
		u32 bar;
		int ret;

		/* we prefer to use a memory-mapped register */
		ret = fdtdec_get_pci_addr(gd->fdt_blob, dev_of_offset(dev),
					  FDT_PCI_SPACE_MEM32, "reg",
					  &pci_addr);
		if (ret) {
			/* try if there is any i/o-mapped register */
			ret = fdtdec_get_pci_addr(gd->fdt_blob,
						  dev_of_offset(dev),
						  FDT_PCI_SPACE_IO,
						  "reg", &pci_addr);
			if (ret)
				return ret;
		}

		ret = fdtdec_get_pci_bar32(dev, &pci_addr, &bar);
		if (ret)
			return ret;

		addr = bar;
	}
#endif

	if (addr == FDT_ADDR_T_NONE)
		return -EINVAL;

#ifdef CONFIG_SYS_NS16550_PORT_MAPPED
	plat->base = addr;
#else
	plat->base = (unsigned long)map_physmem(addr, 0, MAP_NOCACHE);
#endif

	plat->reg_offset = dev_read_u32_default(dev, "reg-offset", 0);
	plat->reg_shift = dev_read_u32_default(dev, "reg-shift", 0);

	err = clk_get_by_index(dev, 0, &clk);
	if (!err) {
		err = clk_get_rate(&clk);
		if (!IS_ERR_VALUE(err))
			plat->clock = err;
	} else if (err != -ENOENT && err != -ENODEV && err != -ENOSYS) {
		debug("ns16550 failed to get clock\n");
		return err;
	}

	if (!plat->clock)
		plat->clock = dev_read_u32_default(dev, "clock-frequency",
						   CONFIG_SYS_NS16550_CLK);
	if (!plat->clock) {
		debug("ns16550 clock not defined\n");
		return -EINVAL;
	}

	plat->fcr = UART_FCR_DEFVAL;
	if (port_type == PORT_JZ4780)
		plat->fcr |= UART_FCR_UME;

	return 0;
}
#endif

const struct dm_serial_ops ns16550_serial_ops = {
	.putc = ns16550_serial_putc,
	.pending = ns16550_serial_pending,
	.getc = ns16550_serial_getc,
	.setbrg = ns16550_serial_setbrg,
};

#if CONFIG_IS_ENABLED(OF_CONTROL) && !CONFIG_IS_ENABLED(OF_PLATDATA)
/*
 * Please consider existing compatible strings before adding a new
 * one to keep this table compact. Or you may add a generic "ns16550"
 * compatible string to your dts.
 */
static const struct udevice_id ns16550_serial_ids[] = {
	{ .compatible = "ns16550",		.data = PORT_NS16550 },
	{ .compatible = "ns16550a",		.data = PORT_NS16550 },
	{ .compatible = "ingenic,jz4780-uart",	.data = PORT_JZ4780  },
	{ .compatible = "nvidia,tegra20-uart",	.data = PORT_NS16550 },
	{ .compatible = "snps,dw-apb-uart",	.data = PORT_NS16550 },
	{ .compatible = "ti,omap2-uart",	.data = PORT_NS16550 },
	{ .compatible = "ti,omap3-uart",	.data = PORT_NS16550 },
	{ .compatible = "ti,omap4-uart",	.data = PORT_NS16550 },
	{ .compatible = "ti,am3352-uart",	.data = PORT_NS16550 },
	{ .compatible = "ti,am4372-uart",	.data = PORT_NS16550 },
	{ .compatible = "ti,dra742-uart",	.data = PORT_NS16550 },
	{}
};
#endif /* OF_CONTROL && !OF_PLATDATA */

#if CONFIG_IS_ENABLED(SERIAL_PRESENT)

/* TODO(sjg@chromium.org): Integrate this into a macro like CONFIG_IS_ENABLED */
#if !defined(CONFIG_TPL_BUILD) || defined(CONFIG_TPL_DM_SERIAL)
U_BOOT_DRIVER(ns16550_serial) = {
	.name	= "ns16550_serial",
	.id	= UCLASS_SERIAL,
#if CONFIG_IS_ENABLED(OF_CONTROL) && !CONFIG_IS_ENABLED(OF_PLATDATA)
	.of_match = ns16550_serial_ids,
	.ofdata_to_platdata = ns16550_serial_ofdata_to_platdata,
	.platdata_auto_alloc_size = sizeof(struct ns16550_platdata),
#endif
	.priv_auto_alloc_size = sizeof(struct NS16550),
	.probe = ns16550_serial_probe,
	.ops	= &ns16550_serial_ops,
	.flags	= DM_FLAG_PRE_RELOC,
};
#endif
#endif /* SERIAL_PRESENT */

#endif /* CONFIG_DM_SERIAL */
