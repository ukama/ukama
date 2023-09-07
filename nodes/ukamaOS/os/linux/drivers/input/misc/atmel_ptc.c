// SPDX-License-Identifier: GPL-2.0+
/*
 * Atmel PTC subsystem driver for SAMA5D2 devices and compatible.
 *
 * Copyright (C) 2017 Microchip,
 *               2017 Ludovic Desroches <ludovic.desroches@microchip.com>
 *
 */

#include <linux/cdev.h>
#include <linux/clk.h>
#include <linux/delay.h>
#include <linux/device.h>
#include <linux/firmware.h>
#include <linux/gpio.h>
#include <linux/input.h>
#include <linux/interrupt.h>
#include <linux/io.h>
#include <linux/module.h>
#include <linux/moduleparam.h>
#include <linux/of.h>
#include <linux/of_device.h>
#include <linux/platform_device.h>
#include <linux/slab.h>
#include <linux/types.h>
#include <linux/uaccess.h>

#define ATMEL_PTC_MAX_NODES		64
#define ATMEL_PTC_MAX_SCROLLERS		4
#define ATMEL_PTC_MAX_X_LINES		8
#define ATMEL_PTC_MAX_Y_LINES		8
#define ATMEL_PTC_KEYCODE_BASE_OFFSET	0x100

/* ----- PPP ----- */
#define ATMEL_PPP_FW_NAME		"microchip/ptc_fw.bin"
#define ATMEL_PPP_FW_FOOTER_SIZE	16

#define ATMEL_PPP_CONFIG	0x20
#define ATMEL_PPP_CTRL		0x24
#define ATMEL_PPP_CMD		0x28
#define		ATMEL_PPP_CMD_STOP		0x1
#define		ATMEL_PPP_CMD_RESET		0x2
#define		ATMEL_PPP_CMD_RESTART		0x3
#define		ATMEL_PPP_CMD_ABORT		0x4
#define		ATMEL_PPP_CMD_RUN		0x5
#define		ATMEL_PPP_CMD_RUN_LOCKED	0x6
#define		ATMEL_PPP_CMD_RUN_OCD		0x7
#define		ATMEL_PPP_CMD_UNLOCK		0x8
#define		ATMEL_PPP_CMD_NMI		0x9
#define		ATMEL_PPP_CMD_HOST_OCD_RESUME	0xB
#define ATMEL_PPP_ISR		0x33
#define		ATMEL_PPP_IRQ_MASK	GENMASK(7, 4)
#define		ATMEL_PPP_IRQ0		BIT(4)
#define		ATMEL_PPP_IRQ1		BIT(5)
#define		ATMEL_PPP_IRQ2		BIT(6)
#define		ATMEL_PPP_IRQ3		BIT(7)
#define		ATMEL_PPP_NOTIFY_MASK	GENMASK(3, 0)
#define		ATMEL_PPP_NOTIFY0	BIT(0)
#define		ATMEL_PPP_NOTIFY1	BIT(1)
#define		ATMEL_PPP_NOTIFY2	BIT(2)
#define		ATMEL_PPP_NOTIFY3	BIT(3)
#define ATMEL_PPP_IDR		0x34
#define ATMEL_PPP_IER		0x35

#define atmel_ppp_readb(ptc, reg)	readb_relaxed((ptc)->ppp_regs + (reg))
#define atmel_ppp_writeb(ptc, reg, val)	writeb_relaxed(val, (ptc)->ppp_regs + (reg))
#define atmel_ppp_readl(ptc, reg)	readl_relaxed((ptc)->ppp_regs + (reg))
#define atmel_ppp_writel(ptc, reg, val)	writel_relaxed(val, (ptc)->ppp_regs + (reg))

/* ----- QTM ----- */
#define ATMEL_QTM_CONF_NAME		"microchip/ptc_cfg.bin"

#define ATMEL_QTM_MB_OFFSET			0x4000
#define ATMEL_QTM_MB_SIZE			0x1000

#define	ATMEL_QTM_CMD_FIRM_VERSION		8
#define	ATMEL_QTM_CMD_INIT			18
#define	ATMEL_QTM_CMD_RUN			19
#define	ATMEL_QTM_CMD_STOP			21
#define	ATMEL_QTM_CMD_SET_ACQ_MODE_TIMER	24
#define	ATMEL_QTM_SCROLLER_TYPE_SLIDER		0x0
#define	ATMEL_QTM_SCROLLER_TYPE_WHEEL		0x1

static char *firmware_file = ATMEL_PPP_FW_NAME;
static char *configuration_file = ATMEL_QTM_CONF_NAME;
static bool debug_mode;

struct atmel_qtm_conf_header {
	u8	header_version_major;
	u8	header_version_minor;
	u32	header_size;
	char	*fw_version;
	char	*tool_version;
	char	*date;
	char	*description;
};

/* Depends on firmware version */
struct atmel_qtm_mailbox_map {
	unsigned int	cmd_offset;
	unsigned int	cmd_id_offset;
	unsigned int	cmd_addr_offset;
	unsigned int	cmd_data_offset;
	unsigned int	node_group_config_offset;
	unsigned int	node_group_config_count_offset;
	unsigned int	node_config_offset;
	unsigned int	node_config_size;
	unsigned int	node_config_mask_x_offset;
	unsigned int	node_config_mask_y_offset;
	unsigned int	scroller_group_config_offset;
	unsigned int	scroller_group_config_count_offset;
	unsigned int	scroller_config_offset;
	unsigned int	scroller_config_size;
	unsigned int	scroller_config_type_offset;
	unsigned int	scroller_config_key_start_offset;
	unsigned int	scroller_config_key_count_offset;
	unsigned int	scroller_config_resol_deadband_offset;
	unsigned int	scroller_data_offset;
	unsigned int	scroller_data_size;
	unsigned int	scroller_data_status_offset;
	unsigned int	scroller_data_position_offset;
	unsigned int	touch_events_key_event_id;
	unsigned int	touch_events_key_enable_state;
	unsigned int	touch_events_scroller_event_id;
};

static struct atmel_qtm_mailbox_map mailbox_map_v64 = {
	.cmd_offset				= 0x0,
	.cmd_id_offset				= 0,
	.cmd_addr_offset			= 2,
	.cmd_data_offset			= 4,
	.node_group_config_offset		= 0x100,
	.node_group_config_count_offset		= 0,
	.node_config_offset			= 0x106,
	.node_config_size			= 12,
	.node_config_mask_x_offset		= 0,
	.node_config_mask_y_offset		= 2,
	.scroller_group_config_offset		= 0x816,
	.scroller_group_config_count_offset	= 2,
	.scroller_config_offset			= 0x81a,
	.scroller_config_size			= 10,
	.scroller_config_type_offset		= 0,
	.scroller_config_key_start_offset	= 2,
	.scroller_config_key_count_offset	= 4,
	.scroller_config_resol_deadband_offset	= 5,
	.scroller_data_offset			= 0x842,
	.scroller_data_size			= 10,
	.scroller_data_status_offset		= 0,
	.scroller_data_position_offset		= 6,
	.touch_events_key_event_id		= 0x880,
	.touch_events_key_enable_state		= 0x888,
	.touch_events_scroller_event_id		= 0x890,
};

struct atmel_qtm_cmd {
	u16	id;
	u16	addr;
	u32	data;
};

struct atmel_ptc_pin {
	const char *name;
	unsigned int id;
};

struct atmel_ptc_pins {
	struct atmel_ptc_pin	x_lines[ATMEL_PTC_MAX_X_LINES];
	struct atmel_ptc_pin	y_lines[ATMEL_PTC_MAX_Y_LINES];
};

struct atmel_ptc {
	void __iomem			*ppp_regs;
	void __iomem			*firmware;
	void __iomem			*qtm_mb;
	struct clk			*clk_per;
	struct clk			*clk_int_osc;
	struct clk			*clk_slow;
	struct device			*dev;
	struct input_dev		*buttons_input;
	struct input_dev		*scroller_input[ATMEL_PTC_MAX_SCROLLERS];
	struct atmel_qtm_conf_header	conf;
	struct atmel_qtm_mailbox_map	*mb_map;
	const struct atmel_ptc_pins	*pins;
	bool				*x_lines_requested;
	bool				*y_lines_requested;
	int				irq;
	u8				imr;
	struct completion		ppp_ack;
	struct tasklet_struct		tasklet;
	unsigned int			button_keycode[ATMEL_PTC_MAX_NODES];
	bool				buttons_registered;
	bool				scroller_registered[ATMEL_PTC_MAX_SCROLLERS];
	u32				button_event[ATMEL_PTC_MAX_NODES / 32];
	u32				button_state[ATMEL_PTC_MAX_NODES / 32];
	u32				scroller_event;
	bool				scroller_tracking;
	char				fw_version[ATMEL_PPP_FW_FOOTER_SIZE];
};

static void atmel_ppp_irq_enable(struct atmel_ptc *ptc, u8 mask)
{
	ptc->imr |= mask;
	atmel_ppp_writeb(ptc, ATMEL_PPP_IER, mask & ATMEL_PPP_IRQ_MASK);
}

static void atmel_ppp_irq_disable(struct atmel_ptc *ptc, u8 mask)
{
	ptc->imr &= ~mask;
	atmel_ppp_writeb(ptc, ATMEL_PPP_IDR, mask & ATMEL_PPP_IRQ_MASK);
}

static void atmel_ppp_notify(struct atmel_ptc *ptc, u8 mask)
{
	if (mask & ATMEL_PPP_NOTIFY_MASK) {
		u8 notify = atmel_ppp_readb(ptc, ATMEL_PPP_ISR)
			| (mask & ATMEL_PPP_NOTIFY_MASK);

		atmel_ppp_writeb(ptc, ATMEL_PPP_ISR, notify);
	}
}

static void atmel_ppp_irq_pending_clr(struct atmel_ptc *ptc, u8 mask)
{
	if (mask & ATMEL_PPP_IRQ_MASK) {
		u8 irq = atmel_ppp_readb(ptc, ATMEL_PPP_ISR) & ~mask;

		atmel_ppp_writeb(ptc, ATMEL_PPP_ISR, irq);
	}
}

static void atmel_ppp_cmd_send(struct atmel_ptc *ptc, u32 cmd)
{
	atmel_ppp_writel(ptc, ATMEL_PPP_CMD, cmd);
}

static void atmel_qtm_set_cmd_id(struct atmel_ptc *ptc, u16 val)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->cmd_offset
			+ ptc->mb_map->cmd_id_offset;

	writew_relaxed(val, addr);
}

static void atmel_qtm_set_cmd_addr(struct atmel_ptc *ptc, u16 val)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->cmd_offset
			+ ptc->mb_map->cmd_addr_offset;

	writew_relaxed(val, addr);
}

static u32 atmel_qtm_get_cmd_data(struct atmel_ptc *ptc)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->cmd_offset
			+ ptc->mb_map->cmd_data_offset;

	return readl_relaxed(addr);
}

static void atmel_qtm_set_cmd_data(struct atmel_ptc *ptc, u32 val)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->cmd_offset
			+ ptc->mb_map->cmd_data_offset;

	writel_relaxed(val, addr);
}

static u16 atmel_qtm_get_node_mask_x(struct atmel_ptc *ptc,
				     unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->node_config_offset
			+ ptc->mb_map->node_config_size * index
			+ ptc->mb_map->node_config_mask_x_offset;

	return readw_relaxed(addr);
}

static u32 atmel_qtm_get_node_mask_y(struct atmel_ptc *ptc,
					      unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->node_config_offset
			+ ptc->mb_map->node_config_size * index
			+ ptc->mb_map->node_config_mask_y_offset;

	return readl_relaxed(addr);
}

static u16 atmel_qtm_get_key_count(struct atmel_ptc *ptc)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->node_group_config_offset
			+ ptc->mb_map->node_group_config_count_offset;

	return readw_relaxed(addr);
}

static u8 atmel_qtm_get_scroller_group_config_count(struct atmel_ptc *ptc)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->scroller_group_config_offset
			+ ptc->mb_map->scroller_group_config_count_offset;

	return readb_relaxed(addr);
}

static u8 atmel_qtm_get_scroller_type(struct atmel_ptc *ptc,
				      unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->scroller_config_offset
			+ ptc->mb_map->scroller_config_size * index
			+ ptc->mb_map->scroller_config_type_offset;

	return readb_relaxed(addr);
}

static u16 atmel_qtm_get_scroller_key_start(struct atmel_ptc *ptc,
					    unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->scroller_config_offset
			+ ptc->mb_map->scroller_config_size * index
			+ ptc->mb_map->scroller_config_key_start_offset;

	return readw_relaxed(addr);
}

static u8 atmel_qtm_get_scroller_key_count(struct atmel_ptc *ptc,
					   unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->scroller_config_offset
			+ ptc->mb_map->scroller_config_size * index
			+ ptc->mb_map->scroller_config_key_count_offset;

	return readb_relaxed(addr);
}

static u8 atmel_qtm_get_scroller_resolution(struct atmel_ptc *ptc,
					    unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->scroller_config_offset
			+ ptc->mb_map->scroller_config_size * index
			+ ptc->mb_map->scroller_config_resol_deadband_offset;

	return (1 << (readb_relaxed(addr) >> 4));
}

static u8 atmel_qtm_get_scroller_status(struct atmel_ptc *ptc,
					unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->scroller_data_offset
			+ ptc->mb_map->scroller_data_size * index
			+ ptc->mb_map->scroller_data_status_offset;

	return readb_relaxed(addr);
}

static u16 atmel_qtm_get_scroller_position(struct atmel_ptc *ptc,
					   unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->scroller_data_offset
			+ ptc->mb_map->scroller_data_size * index
			+ ptc->mb_map->scroller_data_position_offset;

	return readw_relaxed(addr);
}

static u32 atmel_qtm_get_touch_events_key_event_id(struct atmel_ptc *ptc,
						   unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->touch_events_key_event_id
			+ sizeof(u32) * index;

	return readw_relaxed(addr);
}

static u32 atmel_qtm_get_touch_events_key_enable_state(struct atmel_ptc *ptc,
						       unsigned int index)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->touch_events_key_enable_state
			+ sizeof(u32) * index;

	return readw_relaxed(addr);
}

static u32 atmel_qtm_get_touch_events_scroller_event_id(struct atmel_ptc *ptc)
{
	void __iomem *addr = ptc->qtm_mb
			+ ptc->mb_map->touch_events_scroller_event_id;

	return readw_relaxed(addr);
}

static void atmel_ptc_irq_scroller_event(struct atmel_ptc *ptc)
{
	unsigned int status, i;

	if (!ptc->scroller_event || ptc->scroller_tracking)
		return;

	/*
	 * Report the touch event and let the tasklet tracking the position
	 * until the scrollers are no longer touched.
	 */
	for (i = 0 ; i < ATMEL_PTC_MAX_SCROLLERS; i++) {
		if (!ptc->scroller_input[i])
			break;

		status = atmel_qtm_get_scroller_status(ptc, i);

		input_report_key(ptc->scroller_input[i], BTN_TOUCH,
				 status & 0x1);
		input_sync(ptc->scroller_input[i]);

		if (status & 0x1) {
			ptc->scroller_tracking = true;
			tasklet_schedule(&ptc->tasklet);
		}
	}
}

static void atmel_ptc_irq_button_event(struct atmel_ptc *ptc)
{
	unsigned long i;

	for_each_set_bit(i, (unsigned long *)&ptc->button_event, ATMEL_PTC_MAX_NODES) {
		u32 state = ptc->button_state[i / 32] & BIT(i % 32);

		input_report_key(ptc->buttons_input,
				 ptc->button_keycode[i], !!state);
		input_sync(ptc->buttons_input);
	}
}

static void atmel_ptc_irq_touch_event(struct atmel_ptc *ptc)
{
	if (ptc->scroller_input[0])
		atmel_ptc_irq_scroller_event(ptc);

	if (ptc->buttons_input)
		atmel_ptc_irq_button_event(ptc);
}

static irqreturn_t atmel_ptc_irq_handler(int irq, void *data)
{
	struct atmel_ptc *ptc = data;
	u32 isr = atmel_ppp_readb(ptc, ATMEL_PPP_ISR) & ptc->imr;

	/* QTM CMD acknowledgment */
	if (isr & ATMEL_PPP_IRQ0) {
		atmel_ppp_irq_disable(ptc, ATMEL_PPP_IRQ0);
		atmel_ppp_irq_pending_clr(ptc, ATMEL_PPP_IRQ0);
		complete(&ptc->ppp_ack);
	}
	/* QTM touch event */
	if (isr & ATMEL_PPP_IRQ1) {
		int i;

		for (i = 0; i < ATMEL_PTC_MAX_NODES / 32; i++) {
			ptc->button_event[i] = atmel_qtm_get_touch_events_key_event_id(ptc, i);
			ptc->button_state[i] = atmel_qtm_get_touch_events_key_enable_state(ptc, i);
		}
		ptc->scroller_event = atmel_qtm_get_touch_events_scroller_event_id(ptc);

		atmel_ppp_irq_pending_clr(ptc, ATMEL_PPP_IRQ1);

		atmel_ptc_irq_touch_event(ptc);
	}
	/* Debug event */
	if (isr & ATMEL_PPP_IRQ2)
		atmel_ppp_irq_pending_clr(ptc, ATMEL_PPP_IRQ2);

	return IRQ_HANDLED;
}

static void atmel_ptc_free_pins(struct atmel_ptc *ptc)
{
	int i;

	for (i = 0; i < ATMEL_PTC_MAX_X_LINES; i++) {
		if (ptc->x_lines_requested[i]) {
			gpio_free(ptc->pins->x_lines[i].id);
			ptc->x_lines_requested[i] = false;
		}
	}

	for (i = 0; i < ATMEL_PTC_MAX_Y_LINES; i++) {
		if (ptc->y_lines_requested[i]) {
			gpio_free(ptc->pins->y_lines[i].id);
			ptc->y_lines_requested[i] = false;
		}
	}
}

static int atmel_ptc_request_pins(struct atmel_ptc *ptc)
{
	u16 key_count = atmel_qtm_get_key_count(ptc);
	int i, j, ret = 0;

	/*
	 * One line can be used by several nodes. For that reason, we have
	 * to take care about not requesting a GPIO more than once.
	 */
	for (i = 0; i < key_count; i++) {
		u16 mask_x = atmel_qtm_get_node_mask_x(ptc, i);
		u32 mask_y = atmel_qtm_get_node_mask_y(ptc, i);

		for_each_set_bit(j, (unsigned long *)&mask_x, 16) {
			if (ptc->x_lines_requested[j])
				continue;

			if (gpio_request_one(ptc->pins->x_lines[j].id, GPIOF_DIR_IN, ptc->pins->x_lines[j].name)) {
				dev_err(ptc->dev, "Can't get %s\n", ptc->pins->x_lines[j].name);
				return -ENXIO;
			}
			ptc->x_lines_requested[j] = true;
		}

		for_each_set_bit(j, (unsigned long *)&mask_y, 32) {
			if (ptc->y_lines_requested[j])
				continue;

			if (gpio_request_one(ptc->pins->y_lines[j].id, GPIOF_DIR_IN, ptc->pins->y_lines[j].name)) {
				dev_err(ptc->dev, "Can't get %s\n", ptc->pins->y_lines[j].name);
				return -ENXIO;
			}
			ptc->y_lines_requested[j] = true;
		}
	}

	return ret;
}

static int atmel_ptc_conf_input_devs(struct atmel_ptc *ptc)
{
	struct input_dev *input_buttons = ptc->buttons_input;
	u16 key_count = atmel_qtm_get_key_count(ptc);
	bool *buttons;
	int i, j, ret = 0;

	buttons = kmalloc_array(key_count, sizeof(*buttons), GFP_KERNEL);
	if (!buttons)
		return -ENOMEM;
	memset(buttons, true, key_count);

	for (i = 0; i < atmel_qtm_get_scroller_group_config_count(ptc); i++) {
		struct input_dev *scroller = ptc->scroller_input[i];
		u16 key_start = atmel_qtm_get_scroller_key_start(ptc, i);
		u8 key_count = atmel_qtm_get_scroller_key_count(ptc, i);

		if (key_count) {
			/* If a key is part of a scroller, it's not a button. */
			for (j = key_start; j < key_start + key_count; j++)
				buttons[j] = false;

			/*
			 * Prevent several allocations for the same scroller if
			 * atmel_ptc_conf_input_devs() is call more than
			 * once.
			 */
			if (!scroller) {
				scroller = devm_input_allocate_device(ptc->dev);
				if (!scroller)
					return -ENOMEM;
				ptc->scroller_input[i] = scroller;
				scroller->dev.parent = ptc->dev;
			}

			switch (atmel_qtm_get_scroller_type(ptc, i)) {
			case ATMEL_QTM_SCROLLER_TYPE_SLIDER:
				input_set_capability(scroller, EV_ABS, ABS_X);
				input_set_capability(scroller, EV_KEY, BTN_TOUCH);
				input_set_abs_params(scroller, ABS_X, 0, atmel_qtm_get_scroller_resolution(ptc, i), 0, 0);
				break;
			case ATMEL_QTM_SCROLLER_TYPE_WHEEL:
				input_set_capability(scroller, EV_ABS, ABS_WHEEL);
				input_set_capability(scroller, EV_KEY, BTN_TOUCH);
				input_set_abs_params(scroller, ABS_WHEEL, 0, atmel_qtm_get_scroller_resolution(ptc, i), 0, 0);
				break;
			default:
				ret = -EINVAL;
				goto out;
			}
		}
	}

	for (i = 0; i < key_count; i++) {
		if (buttons[i]) {
			int keycode = ATMEL_PTC_KEYCODE_BASE_OFFSET + i;

			if (!input_buttons) {
				input_buttons = devm_input_allocate_device(ptc->dev);
				if (!input_buttons)
					return -ENOMEM;

				input_buttons->dev.parent = ptc->dev;
				input_buttons->keycode = ptc->button_keycode;
				input_buttons->keycodesize = sizeof(ptc->button_keycode[0]);
				input_buttons->keycodemax = ATMEL_PTC_MAX_NODES;
				ptc->buttons_input = input_buttons;
			}
			ptc->button_keycode[i] = keycode;
			input_set_capability(input_buttons, EV_KEY, keycode);
		}
	}

out:
	kfree(buttons);
	return ret;
}

static void atmel_ptc_unregister_input_devices(struct atmel_ptc *ptc)
{
	int i;

	if (ptc->buttons_registered) {
		input_unregister_device(ptc->buttons_input);
		ptc->buttons_registered = false;
	}

	for (i = 0; i < ATMEL_PTC_MAX_SCROLLERS; i++) {
		struct input_dev *scroller = ptc->scroller_input[i];

		if (!scroller || !ptc->scroller_registered[i])
			continue;

		input_unregister_device(scroller);
		ptc->scroller_registered[i] = false;
	}
}

static int atmel_ptc_register_input_devices(struct atmel_ptc *ptc)
{
	int i, ret = 0, id = 0;

	if (ptc->buttons_input) {
		struct input_dev *buttons = ptc->buttons_input;
		buttons->name = devm_kasprintf(&buttons->dev, GFP_KERNEL,
					       "atmel_ptc%d", id++);
		if (!buttons->name)
			return -ENOMEM;
		ret = input_register_device(buttons);
		if (ret) {
			dev_err(ptc->dev, "can't register input button device.\n");
			atmel_ptc_unregister_input_devices(ptc);
			return ret;
		}

		ptc->buttons_registered = true;
	}

	for (i = 0; i < ATMEL_PTC_MAX_SCROLLERS; i++) {
		struct input_dev *scroller = ptc->scroller_input[i];

		if (!scroller)
			continue;

		scroller->name = devm_kasprintf(&scroller->dev, GFP_KERNEL,
						"atmel_ptc%d", id++);
		if (!scroller->name) {
			atmel_ptc_unregister_input_devices(ptc);
			return -ENOMEM;
		}
		ret = input_register_device(scroller);
		if (ret) {
			dev_err(ptc->dev, "can't register input scroller device.\n");
			atmel_ptc_unregister_input_devices(ptc);
			return ret;
		}

		ptc->scroller_registered[i] = true;
	}

	return ret;
}

static int atmel_ptc_cmd_send(struct atmel_ptc *ptc, struct atmel_qtm_cmd *cmd)
{
	int ret;

	dev_dbg(ptc->dev, "%s: cmd=0x%x, addr=0x%x, data=0x%x\n",
		__func__, cmd->id, cmd->addr, cmd->data);

	/*
	 * Configure and register input devices only when QTM is started
	 * since some information from the QTM configuration is needed.
	 * It could be done when loading the configuration file but if
	 * the debug mode is enabled, the configuration can be changed
	 * before sending the run command.
	 */
	if (cmd->id == ATMEL_QTM_CMD_RUN) {
		ret = atmel_ptc_request_pins(ptc);
		if (ret)
			return ret;

		ret = atmel_ptc_conf_input_devs(ptc);
		if (ret)
			return ret;

		ret = atmel_ptc_register_input_devices(ptc);
		if (ret)
			return ret;
	}

	atmel_qtm_set_cmd_id(ptc, cmd->id);
	atmel_qtm_set_cmd_addr(ptc, cmd->addr);
	atmel_qtm_set_cmd_data(ptc, cmd->data);

	/* Once command performed, we'll get an IRQ. */
	atmel_ppp_irq_enable(ptc, ATMEL_PPP_IRQ0);
	/* Notify PPP that we have sent a command. */
	atmel_ppp_notify(ptc, ATMEL_PPP_NOTIFY0);
	/* Wait for IRQ from PPP. */
	wait_for_completion(&ptc->ppp_ack);

	if (cmd->id == ATMEL_QTM_CMD_STOP) {
		atmel_ptc_unregister_input_devices(ptc);
		atmel_ptc_free_pins(ptc);
	}

	cmd->data = atmel_qtm_get_cmd_data(ptc);
	return 0;
}

static void atmel_ptc_tasklet(unsigned long priv)
{
	struct atmel_ptc *ptc = (struct atmel_ptc *)priv;
	unsigned int scroller_type, position, status, i;

	for (i = 0 ; i < ATMEL_PTC_MAX_SCROLLERS; i++) {
		if (!ptc->scroller_input[i])
			break;

		scroller_type = atmel_qtm_get_scroller_type(ptc, i);
		position = atmel_qtm_get_scroller_position(ptc, i);
		status = atmel_qtm_get_scroller_status(ptc, i);

		if (status & 0x1) {
			if (scroller_type == ATMEL_QTM_SCROLLER_TYPE_WHEEL)
				input_report_abs(ptc->scroller_input[i],
						 ABS_WHEEL, position);
			else
				input_report_abs(ptc->scroller_input[i],
						 ABS_X, position);
			input_sync(ptc->scroller_input[i]);
			tasklet_schedule(&ptc->tasklet);
		} else {
			input_report_key(ptc->scroller_input[i], BTN_TOUCH, 0);
			input_sync(ptc->scroller_input[i]);

			ptc->scroller_tracking = false;
		}
	}
}

static int atmel_ptc_conf_load(struct atmel_ptc *ptc)
{
	const struct firmware *conf;
	struct atmel_qtm_cmd cmd;
	int ret;
	u16 key_count;
	char *dst;

	dev_info(ptc->dev, "Loading configuration: %s\n", configuration_file);
	ret = request_firmware(&conf, configuration_file, ptc->dev);
	if (ret) {
		dev_err(ptc->dev, "Can't load configuration %s\n", configuration_file);
		return ret;
	}

	ptc->conf.header_version_major = conf->data[0];
	ptc->conf.header_version_minor = conf->data[1];
	switch (ptc->conf.header_version_major) {
	case (1):
		ptc->conf.header_size = 96;
		ptc->conf.fw_version = (char *) conf->data + 16;
		ptc->conf.tool_version = (char *) conf->data + 32;
		ptc->conf.date = (char *) conf->data + 48;
		ptc->conf.description = (char *) conf->data + 64;
		break;
	default:
		release_firmware(conf);
		dev_err(ptc->dev, "Unsupported header version: %u.%u\n",
			ptc->conf.header_version_major,
			ptc->conf.header_version_minor);
		return -EINVAL;
	};

	dev_info(ptc->dev, "firmware version: %s, tool version: %s\n",
		 ptc->conf.fw_version, ptc->conf.tool_version);
	dev_info(ptc->dev, "date: %s, description: %s\n",
		 ptc->conf.date, ptc->conf.description);

	/*
	 * TODO: check the version of the firmware loaded vs the version of the
	 * firmware needed by the configuration file.
	 */
	if (strcmp(ptc->fw_version, ptc->conf.fw_version))
		dev_warn(ptc->dev, "be careful the configuration requires firmware %s, current firmware is %s\n",
			 ptc->conf.fw_version, ptc->fw_version);

	atmel_ppp_irq_enable(ptc, ATMEL_PPP_IRQ1);
	atmel_ppp_irq_disable(ptc, ATMEL_PPP_IRQ2 | ATMEL_PPP_IRQ3);

	cmd.id = ATMEL_QTM_CMD_STOP;
	atmel_ptc_cmd_send(ptc, &cmd);

	dst = (char *)ptc->qtm_mb;
	/* Need to use _memcpy_toio, otherwise configuration is not well loaded. */
	_memcpy_toio(dst, conf->data + ptc->conf.header_size,
		     conf->size - ptc->conf.header_size);

	key_count = atmel_qtm_get_key_count(ptc);

	/* Start QTM. */
	cmd.id = ATMEL_QTM_CMD_INIT;
	cmd.data = key_count;
	atmel_ptc_cmd_send(ptc, &cmd);

	cmd.id = ATMEL_QTM_CMD_SET_ACQ_MODE_TIMER;
	cmd.data = 1;
	atmel_ptc_cmd_send(ptc, &cmd);

	cmd.id = ATMEL_QTM_CMD_RUN;
	cmd.data = key_count;
	atmel_ptc_cmd_send(ptc, &cmd);

	cmd.id = ATMEL_QTM_CMD_FIRM_VERSION;
	cmd.data = 0;
	atmel_ptc_cmd_send(ptc, &cmd);
	dev_dbg(ptc->dev, "firmware version: v%u.%u\n",
		(cmd.data >> 16) & 0xffff, cmd.data & 0xffff);

	release_firmware(conf);

	return ret;
}

static int atmel_ptc_fw_load(struct atmel_ptc *ptc)
{
	const struct firmware *fw;
	int ret;

	dev_dbg(ptc->dev, "loading firmware: %s\n", firmware_file);
	ret = request_firmware(&fw, firmware_file, ptc->dev);
	if (ret) {
		dev_err(ptc->dev, "error while requesting the firmware\n");
		return ret;
	}

	strncpy(ptc->fw_version, fw->data + fw->size - ATMEL_PPP_FW_FOOTER_SIZE,
		ATMEL_PPP_FW_FOOTER_SIZE);
	dev_dbg(ptc->dev, "version: %s\n", ptc->fw_version);

	if (!strcmp(ptc->fw_version, "PPP_VER_006.004")) {
		ptc->mb_map = &mailbox_map_v64;
	} else {
		dev_err(ptc->dev, "unsupported firmware version\n");
		ret = -EINVAL;
		goto out;
	}

	/* Memset needed to avoid firmware unexpected behavior. */
	memset(ptc->firmware, 0, ATMEL_QTM_MB_OFFSET + sizeof(*ptc->qtm_mb));
	ptc->imr = 0;
	/* Command sequence to start from a clean state. */
	atmel_ppp_cmd_send(ptc, ATMEL_PPP_CMD_ABORT);
	atmel_ppp_irq_pending_clr(ptc, ATMEL_PPP_IRQ_MASK);
	atmel_ppp_cmd_send(ptc, ATMEL_PPP_CMD_RESET);

	dev_dbg(ptc->dev, "downloading %zu bytes\n", fw->size);
	memcpy(ptc->firmware, fw->data, fw->size);

	atmel_ppp_cmd_send(ptc, ATMEL_PPP_CMD_RUN);

out:
	release_firmware(fw);
	return ret;
}

static void atmel_ptc_clk_disable_unprepare(void *data)
{
	struct atmel_ptc *ptc = data;

	clk_disable_unprepare(ptc->clk_slow);
	clk_disable_unprepare(ptc->clk_per);
	clk_disable_unprepare(ptc->clk_int_osc);
}

static inline struct atmel_ptc *kobj_to_atmel_ptc(struct kobject *kobj)
{
	struct device *dev = kobj_to_dev(kobj);

	return dev->driver_data;
}

static ssize_t atmel_qtm_mb_read(struct file *filp, struct kobject *kobj,
				 struct bin_attribute *attr,
				 char *buf, loff_t off, size_t count)
{
	struct atmel_ptc *ptc = kobj_to_atmel_ptc(kobj);
	char *qtm_mb = (char *)ptc->qtm_mb;

	dev_dbg(ptc->dev, "%s: off=0x%llx, count=%zu\n", __func__, off, count);

	memcpy_fromio(buf, qtm_mb + off, count);

	return count;
}

static ssize_t atmel_qtm_mb_write(struct file *filp, struct kobject *kobj,
				  struct bin_attribute *attr,
				  char *buf, loff_t off, size_t count)
{
	struct atmel_ptc *ptc = kobj_to_atmel_ptc(kobj);
	char *qtm_mb = (char *)ptc->qtm_mb;

	dev_dbg(ptc->dev, "%s: off=0x%llx, count=%zu\n", __func__, off, count);

	if (off == 0 && count == sizeof(struct atmel_qtm_cmd))
		atmel_ptc_cmd_send(ptc, (struct atmel_qtm_cmd *)buf);
	else
		memcpy_toio(qtm_mb + off, buf, count);

	return count;
}

static BIN_ATTR_RW(atmel_qtm_mb, ATMEL_QTM_MB_SIZE);

static struct bin_attribute *atmel_ptc_qtm_mb_attrs[] = {
	&bin_attr_atmel_qtm_mb,
	NULL,
};

static const struct attribute_group atmel_ptc_qtm_mb_attr_group = {
	.bin_attrs = atmel_ptc_qtm_mb_attrs,
};

static int atmel_ptc_probe(struct platform_device *pdev)
{
	struct atmel_ptc *ptc;
	struct resource	*res;
	void *shared_memory;
	int ret;

	ptc = devm_kzalloc(&pdev->dev, sizeof(*ptc), GFP_KERNEL);
	if (!ptc)
		return -ENOMEM;

	platform_set_drvdata(pdev, ptc);
	ptc->dev = &pdev->dev;
	ptc->dev->driver_data = ptc;

	res = platform_get_resource(pdev, IORESOURCE_MEM, 0);
	if (!res)
		return -ENODEV;

	shared_memory = devm_ioremap_resource(&pdev->dev, res);
	if (IS_ERR(shared_memory))
		return PTR_ERR(shared_memory);

	ptc->firmware = shared_memory;
	ptc->qtm_mb = shared_memory + ATMEL_QTM_MB_OFFSET;

	res = platform_get_resource(pdev, IORESOURCE_MEM, 1);
	if (!res)
		return -EINVAL;

	ptc->ppp_regs = devm_ioremap_resource(&pdev->dev, res);
	if (IS_ERR(ptc->ppp_regs))
		return PTR_ERR(ptc->ppp_regs);

	ptc->irq = platform_get_irq(pdev, 0);
	if (ptc->irq <= 0) {
		if (!ptc->irq)
			ptc->irq = -ENXIO;

		return ptc->irq;
	}

	ptc->pins = of_device_get_match_data(&pdev->dev);
	if (!ptc->pins) {
		dev_err(ptc->dev, "can't retrieve pins\n");
		return -EINVAL;
	}

	ptc->x_lines_requested = devm_kzalloc(ptc->dev,
		ATMEL_PTC_MAX_X_LINES * sizeof(*ptc->x_lines_requested),
		GFP_KERNEL);
	if (!ptc->x_lines_requested)
		return -ENOMEM;

	ptc->y_lines_requested = devm_kzalloc(ptc->dev,
		ATMEL_PTC_MAX_Y_LINES * sizeof(*ptc->y_lines_requested),
		GFP_KERNEL);
	if (!ptc->y_lines_requested)
		return -ENOMEM;

	ptc->clk_per = devm_clk_get(&pdev->dev, "ptc_clk");
	if (IS_ERR(ptc->clk_per))
		return PTR_ERR(ptc->clk_per);

	ptc->clk_int_osc = devm_clk_get(&pdev->dev, "ptc_int_osc");
	if (IS_ERR(ptc->clk_int_osc))
		return PTR_ERR(ptc->clk_int_osc);

	ptc->clk_slow = devm_clk_get(&pdev->dev, "slow_clk");
	if (IS_ERR(ptc->clk_slow))
		return PTR_ERR(ptc->clk_slow);

	ret = clk_prepare_enable(ptc->clk_int_osc);
	if (ret)
		return ret;

	ret = clk_prepare_enable(ptc->clk_per);
	if (ret) {
		clk_disable_unprepare(ptc->clk_int_osc);
		return ret;
	}

	ret = clk_prepare_enable(ptc->clk_slow);
	if (ret) {
		clk_disable_unprepare(ptc->clk_per);
		clk_disable_unprepare(ptc->clk_int_osc);
		return ret;
	}

	ret = devm_add_action_or_reset(&pdev->dev,
				       atmel_ptc_clk_disable_unprepare,
				       ptc);
	if (ret)
		return ret;

	ret = devm_request_irq(&pdev->dev, ptc->irq, atmel_ptc_irq_handler, 0,
			       pdev->dev.driver->name, ptc);
	if (ret)
		return ret;

	init_completion(&ptc->ppp_ack);

	ret = atmel_ptc_fw_load(ptc);
	if (ret)
		return ret;

	ret = atmel_ptc_conf_load(ptc);
	if (ret)
		return ret;

	/*
	 * Expose a file to give an access to the QTM mailbox to a user space
	 * application in order to configure it or to send commands.
	 */
	if (debug_mode)
		ret = sysfs_create_group(&ptc->dev->kobj, &atmel_ptc_qtm_mb_attr_group);

	tasklet_init(&ptc->tasklet, atmel_ptc_tasklet, (unsigned long)ptc);

	return ret;
}

static int atmel_ptc_remove(struct platform_device *pdev)
{
	struct atmel_ptc *ptc = platform_get_drvdata(pdev);

	tasklet_kill(&ptc->tasklet);
	atmel_ptc_unregister_input_devices(ptc);
	atmel_ptc_free_pins(ptc);

	if (debug_mode)
		sysfs_remove_group(&ptc->dev->kobj, &atmel_ptc_qtm_mb_attr_group);

	return 0;
}

/*
 * x_lines and y_lines have to be described in the ascending order i.e.
 * PTC_X0 must be at index 0, PTC_X1 must be at index 1, etc.
 */
static struct atmel_ptc_pins atmel_ptc_pins_sama5d2 = {
	.x_lines = {
		{ .name = "PTC_X0", .id = 99  /* PD3 */  },
		{ .name = "PTC_X1", .id = 100 /* PD4 */  },
		{ .name = "PTC_X2", .id = 101 /* PD5 */  },
		{ .name = "PTC_X3", .id = 102 /* PD6 */  },
		{ .name = "PTC_X4", .id = 103 /* PD7 */  },
		{ .name = "PTC_X5", .id = 104 /* PD8 */  },
		{ .name = "PTC_X6", .id = 105 /* PD9 */  },
		{ .name = "PTC_X7", .id = 106 /* PD10 */ },
	},
	.y_lines = {
		{ .name = "PTC_Y0", .id = 107 /* PD11 */ },
		{ .name = "PTC_Y1", .id = 108 /* PD12 */ },
		{ .name = "PTC_Y2", .id = 109 /* PD13 */ },
		{ .name = "PTC_Y3", .id = 110 /* PD14 */ },
		{ .name = "PTC_Y4", .id = 111 /* PD15 */ },
		{ .name = "PTC_Y5", .id = 112 /* PD16 */ },
		{ .name = "PTC_Y6", .id = 113 /* PD17 */ },
		{ .name = "PTC_Y7", .id = 114 /* PD18 */ },
	},
};

static const struct of_device_id atmel_ptc_dt_match[] = {
	{
		.compatible = "atmel,sama5d2-ptc",
		.data = &atmel_ptc_pins_sama5d2,
	}, {
		/* sentinel */
	}
};
MODULE_DEVICE_TABLE(of, atmel_ptc_dt_match);

static struct platform_driver atmel_ptc_driver = {
	.probe = atmel_ptc_probe,
	.remove = atmel_ptc_remove,
	.driver = {
		.name = "atmel_ptc",
		.of_match_table = of_match_ptr(atmel_ptc_dt_match),
	},
};
module_platform_driver(atmel_ptc_driver)

module_param(firmware_file, charp, 0444);
MODULE_PARM_DESC(firmware_file, "Name of the firmware file");
module_param(configuration_file, charp, 0444);
MODULE_PARM_DESC(configuration_file, "Name of the configuration file");
module_param(debug_mode, bool, 0444);
MODULE_PARM_DESC(debug_mode, "The debug mode provides an interface to the mailbox through sysfs");
MODULE_AUTHOR("Ludovic Desroches <ludovic.desroches@microchip.com>");
MODULE_DESCRIPTION("Atmel PTC subsystem");
MODULE_LICENSE("GPL v2");
MODULE_FIRMWARE(ATMEL_PPP_FW_NAME);
MODULE_FIRMWARE(ATMEL_QTM_CONF_NAME);
