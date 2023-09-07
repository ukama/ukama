#ifndef __SERIAL_H__
#define __SERIAL_H__

#include <post.h>

struct serial_device {
	/* enough bytes to match alignment of following func pointer */
	char	name[16];

	int	(*start)(void);
	int	(*stop)(void);
	void	(*setbrg)(void);
	int	(*getc)(void);
	int	(*tstc)(void);
	void	(*putc)(const char c);
	void	(*puts)(const char *s);
#if CONFIG_POST & CONFIG_SYS_POST_UART
	void	(*loop)(int);
#endif
	struct serial_device	*next;
};

void default_serial_puts(const char *s);

extern struct serial_device serial_smc_device;
extern struct serial_device serial_scc_device;
extern struct serial_device *default_serial_console(void);

#if	defined(CONFIG_MPC83xx) || defined(CONFIG_MPC85xx) || \
	defined(CONFIG_MPC86xx) || \
	defined(CONFIG_TEGRA) || defined(CONFIG_SYS_COREBOOT) || \
	defined(CONFIG_MICROBLAZE)
extern struct serial_device serial0_device;
extern struct serial_device serial1_device;
#endif

extern struct serial_device eserial1_device;
extern struct serial_device eserial2_device;
extern struct serial_device eserial3_device;
extern struct serial_device eserial4_device;
extern struct serial_device eserial5_device;
extern struct serial_device eserial6_device;

extern void serial_register(struct serial_device *);
extern void serial_initialize(void);
extern void serial_stdio_init(void);
extern int serial_assign(const char *name);
extern void serial_reinit_all(void);

/* For usbtty */
#ifdef CONFIG_USB_TTY

extern int usbtty_getc(void);
extern void usbtty_putc(const char c);
extern void usbtty_puts(const char *str);
extern int usbtty_tstc(void);

#else

/* stubs */
#define usbtty_getc() 0
#define usbtty_putc(a)
#define usbtty_puts(a)
#define usbtty_tstc() 0

#endif /* CONFIG_USB_TTY */

struct udevice;

enum serial_par {
	SERIAL_PAR_NONE,
	SERIAL_PAR_ODD,
	SERIAL_PAR_EVEN
};

/**
 * struct struct dm_serial_ops - Driver model serial operations
 *
 * The uclass interface is implemented by all serial devices which use
 * driver model.
 */
struct dm_serial_ops {
	/**
	 * setbrg() - Set up the baud rate generator
	 *
	 * Adjust baud rate divisors to set up a new baud rate for this
	 * device. Not all devices will support all rates. If the rate
	 * cannot be supported, the driver is free to select the nearest
	 * available rate. or return -EINVAL if this is not possible.
	 *
	 * @dev: Device pointer
	 * @baudrate: New baud rate to use
	 * @return 0 if OK, -ve on error
	 */
	int (*setbrg)(struct udevice *dev, int baudrate);
	/**
	 * getc() - Read a character and return it
	 *
	 * If no character is available, this should return -EAGAIN without
	 * waiting.
	 *
	 * @dev: Device pointer
	 * @return character (0..255), -ve on error
	 */
	int (*getc)(struct udevice *dev);
	/**
	 * putc() - Write a character
	 *
	 * @dev: Device pointer
	 * @ch: character to write
	 * @return 0 if OK, -ve on error
	 */
	int (*putc)(struct udevice *dev, const char ch);
	/**
	 * pending() - Check if input/output characters are waiting
	 *
	 * This can be used to return an indication of the number of waiting
	 * characters if the driver knows this (e.g. by looking at the FIFO
	 * level). It is acceptable to return 1 if an indeterminant number
	 * of characters is waiting.
	 *
	 * This method is optional.
	 *
	 * @dev: Device pointer
	 * @input: true to check input characters, false for output
	 * @return number of waiting characters, 0 for none, -ve on error
	 */
	int (*pending)(struct udevice *dev, bool input);
	/**
	 * clear() - Clear the serial FIFOs/holding registers
	 *
	 * This method is optional.
	 *
	 * This quickly clears any input/output characters from the UART.
	 * If this is not possible, but characters still exist, then it
	 * is acceptable to return -EAGAIN (try again) or -EINVAL (not
	 * supported).
	 *
	 * @dev: Device pointer
	 * @return 0 if OK, -ve on error
	 */
	int (*clear)(struct udevice *dev);
#if CONFIG_POST & CONFIG_SYS_POST_UART
	/**
	 * loop() - Control serial device loopback mode
	 *
	 * @dev: Device pointer
	 * @on: 1 to turn loopback on, 0 to turn if off
	 */
	int (*loop)(struct udevice *dev, int on);
#endif
	/**
	 * setparity() - Set up the parity
	 *
	 * Set up a new parity for this device.
	 *
	 * @dev: Device pointer
	 * @parity: parity to use
	 * @return 0 if OK, -ve on error
	 */
	int (*setparity)(struct udevice *dev, enum serial_par parity);
};

/**
 * struct serial_dev_priv - information about a device used by the uclass
 *
 * @sdev:	stdio device attached to this uart
 *
 * @buf:	Pointer to the RX buffer
 * @rd_ptr:	Read pointer in the RX buffer
 * @wr_ptr:	Write pointer in the RX buffer
 */
struct serial_dev_priv {
	struct stdio_dev *sdev;

	char *buf;
	int rd_ptr;
	int wr_ptr;
};

/* Access the serial operations for a device */
#define serial_get_ops(dev)	((struct dm_serial_ops *)(dev)->driver->ops)

void atmel_serial_initialize(void);
void au1x00_serial_initialize(void);
void mcf_serial_initialize(void);
void mpc85xx_serial_initialize(void);
void mpc8xx_serial_initialize(void);
void mxc_serial_initialize(void);
void ns16550_serial_initialize(void);
void pl01x_serial_initialize(void);
void pxa_serial_initialize(void);
void sh_serial_initialize(void);

#endif
