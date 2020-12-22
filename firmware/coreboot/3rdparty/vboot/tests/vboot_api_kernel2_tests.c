/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Tests for vboot_api_kernel, part 2
 */

#include "2common.h"
#include "2misc.h"
#include "2nvstorage.h"
#include "2secdata.h"
#include "host_common.h"
#include "load_kernel_fw.h"
#include "secdata_tpm.h"
#include "test_common.h"
#include "tss_constants.h"
#include "vboot_audio.h"
#include "vboot_common.h"
#include "vboot_display.h"
#include "vboot_kernel.h"
#include "vboot_struct.h"
#include "vboot_test.h"

/* Mock data */
static uint8_t shared_data[VB_SHARED_DATA_MIN_SIZE];
static VbSharedDataHeader *shared = (VbSharedDataHeader *)shared_data;
static LoadKernelParams lkp;
static uint8_t workbuf[VB2_KERNEL_WORKBUF_RECOMMENDED_SIZE];
static struct vb2_context *ctx;
static struct vb2_shared_data *sd;
static struct vb2_gbb_header gbb;

static int audio_looping_calls_left;
static uint32_t vbtlk_retval;
static int vbexlegacy_called;
static enum VbAltFwIndex_t altfw_num;
static uint64_t current_ticks;
static int trust_ec;
static int virtdev_set;
static uint32_t virtdev_retval;
static uint32_t mock_keypress[16];
static uint32_t mock_keyflags[8];
static uint32_t mock_keypress_count;

#define GPIO_SHUTDOWN   1
#define GPIO_PRESENCE   2
#define GPIO_LID_CLOSED 4
typedef struct GpioState {
	uint32_t gpio_flags;
	uint32_t count;
} GpioState;
struct GpioState mock_gpio[8];
static uint32_t mock_gpio_count;

static uint32_t screens_displayed[8];
static uint32_t screens_count = 0;
static uint32_t mock_num_disks[8];
static uint32_t mock_num_disks_count;
static int tpm_set_mode_called;
static enum vb2_tpm_mode tpm_mode;

// Extra character to guarantee null termination.
static char set_vendor_data[VENDOR_DATA_LENGTH + 2];
static int set_vendor_data_called;

/*
 * Mocks the assertion of 1 or more gpios in |gpio_flags| for 100 ticks after
 * an optional |ticks| delay.
 */
static void MockGpioAfter(uint32_t ticks, uint32_t gpio_flags)
{
	uint32_t index = 0;
	if (ticks > 0)
		mock_gpio[index++].count = ticks - 1;

	mock_gpio[index].gpio_flags = gpio_flags;
	mock_gpio[index].count = 100;
}

/* Reset mock data (for use before each test) */
static void ResetMocks(void)
{
	vb2_init_ui();
	memset(VbApiKernelGetFwmp(), 0, sizeof(struct RollbackSpaceFwmp));

	memset(&shared_data, 0, sizeof(shared_data));
	VbSharedDataInit(shared, sizeof(shared_data));

	memset(&lkp, 0, sizeof(lkp));

	TEST_SUCC(vb2api_init(workbuf, sizeof(workbuf), &ctx),
		  "vb2api_init failed");
	vb2_nv_init(ctx);

	sd = vb2_get_sd(ctx);
	sd->vbsd = shared;
	sd->flags |= VB2_SD_FLAG_DISPLAY_AVAILABLE;

	memset(&gbb, 0, sizeof(gbb));

	audio_looping_calls_left = 30;
	vbtlk_retval = 1000;
	vbexlegacy_called = 0;
	altfw_num = -100;
	current_ticks = 0;
	trust_ec = 0;
	virtdev_set = 0;
	virtdev_retval = 0;
	set_vendor_data_called = 0;

	memset(screens_displayed, 0, sizeof(screens_displayed));
	screens_count = 0;

	memset(mock_keypress, 0, sizeof(mock_keypress));
	memset(mock_keyflags, 0, sizeof(mock_keyflags));
	mock_keypress_count = 0;
	memset(mock_gpio, 0, sizeof(mock_gpio));
	mock_gpio_count = 0;
	memset(mock_num_disks, 0, sizeof(mock_num_disks));
	mock_num_disks_count = 0;

	tpm_set_mode_called = 0;
	tpm_mode = VB2_TPM_MODE_ENABLED_TENTATIVE;
}

/* Mock functions */
uint32_t RollbackKernelLock(int recovery_mode)
{
	return TPM_SUCCESS;
}

struct vb2_gbb_header *vb2_get_gbb(struct vb2_context *c)
{
	return &gbb;
}

uint32_t VbExIsShutdownRequested(void)
{
	uint32_t result = 0;
	if (mock_gpio_count >= ARRAY_SIZE(mock_gpio))
		return 0;
	if (mock_gpio[mock_gpio_count].gpio_flags & GPIO_SHUTDOWN)
		result |= VB_SHUTDOWN_REQUEST_POWER_BUTTON;
	if (mock_gpio[mock_gpio_count].gpio_flags & GPIO_LID_CLOSED)
		result |= VB_SHUTDOWN_REQUEST_LID_CLOSED;
	if (mock_gpio[mock_gpio_count].count > 0) {
		--mock_gpio[mock_gpio_count].count;
	} else {
		++mock_gpio_count;
	}
	return result;
}

uint32_t VbExKeyboardRead(void)
{
	return VbExKeyboardReadWithFlags(NULL);
}

uint32_t VbExKeyboardReadWithFlags(uint32_t *key_flags)
{
	if (mock_keypress_count < ARRAY_SIZE(mock_keypress)) {
		if (key_flags != NULL &&
		    mock_keypress_count < ARRAY_SIZE(mock_keyflags))
			*key_flags = mock_keyflags[mock_keypress_count];
		return mock_keypress[mock_keypress_count++];
	} else
		return 0;
}

uint32_t VbExGetSwitches(uint32_t request_mask)
{
	uint32_t result = 0;
	if (mock_gpio_count >= ARRAY_SIZE(mock_gpio))
		return 0;
	if ((request_mask & VB_SWITCH_FLAG_PHYS_PRESENCE_PRESSED) &&
	    (mock_gpio[mock_gpio_count].gpio_flags & GPIO_PRESENCE))
		result |= VB_SWITCH_FLAG_PHYS_PRESENCE_PRESSED;
	if (mock_gpio[mock_gpio_count].count > 0) {
		--mock_gpio[mock_gpio_count].count;
	} else {
		++mock_gpio_count;
	}
	return result;
}

vb2_error_t VbExLegacy(enum VbAltFwIndex_t _altfw_num)
{
	vbexlegacy_called++;
	altfw_num = _altfw_num;

	/* VbExLegacy() can only return failure, or not return at all. */
	return VB2_ERROR_UNKNOWN;
}

void VbExSleepMs(uint32_t msec)
{
	current_ticks += (uint64_t)msec * VB_USEC_PER_MSEC;
}

uint64_t VbExGetTimer(void)
{
	return current_ticks;
}

vb2_error_t VbExDiskGetInfo(VbDiskInfo **infos_ptr, uint32_t *count,
			    uint32_t disk_flags)
{
	if (mock_num_disks_count < ARRAY_SIZE(mock_num_disks)) {
		if (mock_num_disks[mock_num_disks_count] == -1)
			return VB2_ERROR_MOCK;
		else
			*count = mock_num_disks[mock_num_disks_count++];
	} else {
		*count = 0;
	}
	return VB2_SUCCESS;
}

vb2_error_t VbExDiskFreeInfo(VbDiskInfo *infos,
			     VbExDiskHandle_t preserve_handle)
{
	return VB2_SUCCESS;
}

int VbExTrustEC(int devidx)
{
	return trust_ec;
}

int vb2_audio_looping(void)
{
	if (audio_looping_calls_left == 0)
		return 0;
	else if (audio_looping_calls_left > 0)
		audio_looping_calls_left--;

	return 1;
}

vb2_error_t VbTryLoadKernel(struct vb2_context *c, uint32_t get_info_flags)
{
	return vbtlk_retval + get_info_flags;
}

vb2_error_t VbDisplayScreen(struct vb2_context *c, uint32_t screen, int force,
			    const VbScreenData *data)
{
	if (screens_count < ARRAY_SIZE(screens_displayed))
		screens_displayed[screens_count++] = screen;

	return VB2_SUCCESS;
}

vb2_error_t SetVirtualDevMode(int val)
{
	virtdev_set = val;
	return virtdev_retval;
}

vb2_error_t VbExSetVendorData(const char *vendor_data_value)
{
	set_vendor_data_called = 1;
	// set_vendor_data is a global variable, so it is automatically
	// initialized to zero, and so the -1 will ensure the string is null
	// terminated.
	strncpy(set_vendor_data, vendor_data_value, sizeof(set_vendor_data) - 1);

	return VB2_SUCCESS;
}

vb2_error_t vb2ex_tpm_set_mode(enum vb2_tpm_mode mode_val)
{
	tpm_set_mode_called = 1;
	/*
	 * This mock will pretend that any call will fail if the tpm is
	 * already disabled (e.g., as if the code always tries to contact the
	 * tpm to issue a command).  The real version may eventually be changed
	 * to return success if the incoming request is also to disable, but
	 * the point here is to have a way to simulate failure.
	 */
	if (tpm_mode == VB2_TPM_MODE_DISABLED) {
		return VB2_ERROR_UNKNOWN;
	}
	tpm_mode = mode_val;
	return VB2_SUCCESS;
}

/* Tests */

/*
 * Helper function to test VbUserConfirms against a sequence of gpio events.
 * caller specifies a sequence of gpio events and the expected confirm vs.
 * reboot result.
 *
 * Non-asserted gpios are used for 5 events, then 'first' for 5 events,
 * 'second' for 5 events, and 'third' for 5 events.
 * Typically most tests want 5 events of each type (so they'll specify 0 for
 * 'first' and use 'second' through 'fourth'), but a few tests want the
 * shutdown event to be seen once.
 */
static void VbUserConfirmsTestGpio(uint32_t first, uint32_t second,
				   uint32_t third, uint32_t confirm,
				   const char *msg)
{
	ResetMocks();
	mock_gpio[0].gpio_flags = 0;
	mock_gpio[0].count = 4;
	mock_gpio[1].gpio_flags = first;
	mock_gpio[1].count = 4;
	mock_gpio[2].gpio_flags = second;
	mock_gpio[2].count = 4;
	mock_gpio[3].gpio_flags = third;
	mock_gpio[3].count = 4;
	if (confirm) {
		TEST_EQ(VbUserConfirms(ctx,
			VB_CONFIRM_SPACE_MEANS_NO |
			VB_CONFIRM_MUST_TRUST_KEYBOARD),
			1, msg);
	} else {
		TEST_EQ(VbUserConfirms(ctx,
			VB_CONFIRM_SPACE_MEANS_NO |
			VB_CONFIRM_MUST_TRUST_KEYBOARD),
			-1, msg);
	}
}

static void VbUserConfirmsTest(void)
{
	printf("Testing VbUserConfirms()...\n");

	ResetMocks();
	MockGpioAfter(1, GPIO_SHUTDOWN);
	TEST_EQ(VbUserConfirms(ctx, 0), -1, "Shutdown requested");

	ResetMocks();
	mock_keypress[0] = VB_BUTTON_POWER_SHORT_PRESS;
	TEST_EQ(VbUserConfirms(ctx, 0), -1, "Shutdown requested");

	ResetMocks();
	mock_keypress[0] = VB_KEY_ENTER;
	TEST_EQ(VbUserConfirms(ctx, 0), 1, "Enter");

	ResetMocks();
	mock_keypress[0] = VB_KEY_ESC;
	TEST_EQ(VbUserConfirms(ctx, 0), 0, "Esc");

	ResetMocks();
	mock_keypress[0] = ' ';
	MockGpioAfter(1, GPIO_SHUTDOWN);
	TEST_EQ(VbUserConfirms(ctx, VB_CONFIRM_SPACE_MEANS_NO), 0,
		"Space means no");

	ResetMocks();
	mock_keypress[0] = ' ';
	MockGpioAfter(1, GPIO_SHUTDOWN);
	TEST_EQ(VbUserConfirms(ctx, 0), -1, "Space ignored");

	ResetMocks();
	mock_keypress[0] = VB_KEY_ENTER;
	mock_keyflags[0] = VB_KEY_FLAG_TRUSTED_KEYBOARD;
	TEST_EQ(VbUserConfirms(ctx, VB_CONFIRM_MUST_TRUST_KEYBOARD),
		1, "Enter with trusted keyboard");

	ResetMocks();
	mock_keypress[0] = VB_KEY_ENTER;	/* untrusted */
	mock_keypress[1] = ' ';
	TEST_EQ(VbUserConfirms(ctx,
			       VB_CONFIRM_SPACE_MEANS_NO |
			       VB_CONFIRM_MUST_TRUST_KEYBOARD),
		0, "Untrusted keyboard");

	ResetMocks();
	MockGpioAfter(0, GPIO_PRESENCE);
	TEST_EQ(VbUserConfirms(ctx,
			       VB_CONFIRM_SPACE_MEANS_NO |
			       VB_CONFIRM_MUST_TRUST_KEYBOARD),
		1, "Presence button");

	/*
	 * List of possiblities for shutdown and physical presence events that
	 * occur over time.  Time advanced from left to right (where each
	 * represents the gpio[s] that are seen during a given iteration of
	 * the loop).  The meaning of the characters:
	 *
	 *   _ means no gpio
	 *   s means shutdown gpio
	 *   p means presence gpio
	 *   B means both shutdown and presence gpio
	 *
	 *  1: ______ppp______ -> confirm
	 *  2: ______sss______ -> shutdown
	 *  3: ___pppsss______ -> confirm
	 *  4: ___sssppp______ -> shutdown
	 *  5: ___pppBBB______ -> confirm
	 *  6: ___pppBBBppp___ -> shutdown
	 *  7: ___pppBBBsss___ -> confirm
	 *  8: ___sssBBB______ -> confirm
	 *  9: ___sssBBBppp___ -> shutdown
	 * 10: ___sssBBBsss___ -> confirm
	 * 11: ______BBB______ -> confirm
	 * 12: ______BBBsss___ -> confirm
	 * 13: ______BBBppp___ -> shutdown
	 */

	/* 1: presence means confirm */
	VbUserConfirmsTestGpio(GPIO_PRESENCE, 0, 0, 1, "presence");

	/* 2: shutdown means shutdown */
	VbUserConfirmsTestGpio(GPIO_SHUTDOWN, 0, 0, 0, "shutdown");

	/* 3: presence then shutdown means confirm */
	VbUserConfirmsTestGpio(GPIO_PRESENCE, GPIO_SHUTDOWN, 0, 1,
			       "presence then shutdown");

	/* 4: shutdown then presence means shutdown */
	VbUserConfirmsTestGpio(GPIO_SHUTDOWN, GPIO_PRESENCE, 0, 0,
			       "shutdown then presence");

	/* 5: presence then shutdown+presence then none mean confirm */
	VbUserConfirmsTestGpio(GPIO_PRESENCE, GPIO_PRESENCE | GPIO_SHUTDOWN,
			       0, 1, "presence, both, none");

	/* 6: presence then shutdown+presence then presence means shutdown */
	VbUserConfirmsTestGpio(GPIO_PRESENCE, GPIO_PRESENCE | GPIO_SHUTDOWN,
			       GPIO_PRESENCE, 0, "presence, both, presence");

	/* 7: presence then shutdown+presence then shutdown means confirm */
	VbUserConfirmsTestGpio(GPIO_PRESENCE, GPIO_PRESENCE | GPIO_SHUTDOWN,
			       GPIO_SHUTDOWN, 1, "presence, both, shutdown");

	/* 8: shutdown then shutdown+presence then none means confirm */
	VbUserConfirmsTestGpio(GPIO_SHUTDOWN, GPIO_PRESENCE | GPIO_SHUTDOWN,
			       0, 1, "shutdown, both, none");

	/* 9: shutdown then shutdown+presence then presence means shutdown */
	VbUserConfirmsTestGpio(GPIO_SHUTDOWN, GPIO_PRESENCE | GPIO_SHUTDOWN,
			       GPIO_PRESENCE, 0, "shutdown, both, presence");

	/* 10: shutdown then shutdown+presence then shutdown means confirm */
	VbUserConfirmsTestGpio(GPIO_SHUTDOWN, GPIO_PRESENCE | GPIO_SHUTDOWN,
			       GPIO_SHUTDOWN, 1, "shutdown, both, shutdown");

	/* 11: shutdown+presence then none means confirm */
	VbUserConfirmsTestGpio(GPIO_PRESENCE | GPIO_SHUTDOWN, 0, 0, 1,
			       "both");

	/* 12: shutdown+presence then shutdown means confirm */
	VbUserConfirmsTestGpio(GPIO_PRESENCE | GPIO_SHUTDOWN,
			       GPIO_SHUTDOWN, 0, 1, "both, shutdown");

	/* 13: shutdown+presence then presence means shutdown */
	VbUserConfirmsTestGpio(GPIO_PRESENCE | GPIO_SHUTDOWN,
			       GPIO_PRESENCE, 0, 0, "both, presence");

	ResetMocks();
	mock_keypress[0] = VB_KEY_ENTER;
	mock_keypress[1] = 'y';
	mock_keypress[2] = 'z';
	mock_keypress[3] = ' ';
	mock_gpio[0].gpio_flags = GPIO_PRESENCE;
	mock_gpio[0].count = ~0;
	TEST_EQ(VbUserConfirms(ctx,
			       VB_CONFIRM_SPACE_MEANS_NO |
			       VB_CONFIRM_MUST_TRUST_KEYBOARD),
		0, "Recovery button stuck");
	printf("...done.\n");
}

static void VbBootTest(void)
{
	ResetMocks();
	TEST_EQ(VbBootNormal(ctx), 1002, "VbBootNormal()");

	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DISPLAY_REQUEST, 1);
	TEST_EQ(VbBootNormal(ctx), VBERROR_REBOOT_REQUIRED,
		"VbBootNormal() reboot to reset NVRAM display request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DISPLAY_REQUEST), 0,
		"  display request reset");

	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DIAG_REQUEST, 1);
	TEST_EQ(VbBootNormal(ctx), VBERROR_REBOOT_REQUIRED,
		"VbBootNormal() reboot to reset NVRAM diag request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST), 0,
		"  diag request reset");
}

static void VbBootDevTest(void)
{
	int key;

	printf("Testing VbBootDeveloper()...\n");

	/* Proceed after timeout */
	ResetMocks();
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Timeout");
	TEST_EQ(screens_displayed[0], VB_SCREEN_DEVELOPER_WARNING,
		"  warning screen");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST), 0,
		"  recovery reason");
	TEST_EQ(audio_looping_calls_left, 0, "  used up audio");

	/* Proceed to legacy after timeout if GBB flag set */
	ResetMocks();
	gbb.flags |= VB2_GBB_FLAG_DEFAULT_DEV_BOOT_LEGACY |
			VB2_GBB_FLAG_FORCE_DEV_BOOT_LEGACY;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Timeout");
	TEST_EQ(vbexlegacy_called, 1, "  try legacy");
	TEST_EQ(altfw_num, 0, "  check altfw_num");

	/* Proceed to legacy after timeout if GBB flag set */
	ResetMocks();
	gbb.flags |= VB2_GBB_FLAG_DEFAULT_DEV_BOOT_LEGACY |
			VB2_GBB_FLAG_FORCE_DEV_BOOT_LEGACY;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Timeout");
	TEST_EQ(vbexlegacy_called, 1, "  try legacy");
	TEST_EQ(altfw_num, 0, "  check altfw_num");

	/* Proceed to legacy after timeout if boot legacy and default boot
	 * legacy are set */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_DEFAULT_BOOT,
		   VB2_DEV_DEFAULT_BOOT_LEGACY);
	vb2_nv_set(ctx, VB2_NV_DEV_BOOT_LEGACY, 1);
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Timeout");
	TEST_EQ(vbexlegacy_called, 1, "  try legacy");
	TEST_EQ(altfw_num, 0, "  check altfw_num");

	/* Proceed to legacy boot mode only if enabled */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_DEFAULT_BOOT,
		   VB2_DEV_DEFAULT_BOOT_LEGACY);
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Timeout");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");

	/* Proceed to usb after timeout if boot usb and default boot
	 * usb are set */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_DEFAULT_BOOT,
		   VB2_DEV_DEFAULT_BOOT_USB);
	vb2_nv_set(ctx, VB2_NV_DEV_BOOT_USB, 1);
	vbtlk_retval = VB2_SUCCESS - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootDeveloper(ctx), 0, "Ctrl+U USB");

	/* Proceed to usb boot mode only if enabled */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_DEFAULT_BOOT,
		   VB2_DEV_DEFAULT_BOOT_USB);
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Timeout");

	/* If no USB tries fixed disk */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_BOOT_USB, 1);
	vb2_nv_set(ctx, VB2_NV_DEV_DEFAULT_BOOT,
		   VB2_DEV_DEFAULT_BOOT_USB);
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+U enabled");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");

	/* Up arrow is uninteresting / passed to VbCheckDisplayKey() */
	ResetMocks();
	mock_keypress[0] = VB_KEY_UP;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Up arrow");

	/* Shutdown requested in loop */
	ResetMocks();
	mock_gpio[0].gpio_flags = 0;
	mock_gpio[0].count = 2;
	mock_gpio[1].gpio_flags = GPIO_SHUTDOWN;
	TEST_EQ(VbBootDeveloper(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Shutdown requested");
	TEST_NEQ(audio_looping_calls_left, 0, "  aborts audio");

	/* Shutdown requested by keyboard in loop */
	ResetMocks();
	mock_keypress[0] = VB_BUTTON_POWER_SHORT_PRESS;
	TEST_EQ(VbBootDeveloper(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Shutdown requested by keyboard");

	/* Space asks to disable virtual dev switch */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	mock_keypress[0] = ' ';
	mock_keypress[1] = VB_KEY_ENTER;
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Space = tonorm");
	TEST_EQ(screens_displayed[0], VB_SCREEN_DEVELOPER_WARNING,
		"  warning screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_DEVELOPER_TO_NORM,
		"  tonorm screen");
	TEST_EQ(screens_displayed[2], VB_SCREEN_TO_NORM_CONFIRMED,
		"  confirm screen");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DISABLE_DEV_REQUEST), 1,
		"  disable dev request");

	/* Space-space doesn't disable it */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	mock_keypress[0] = ' ';
	mock_keypress[1] = ' ';
	mock_keypress[2] = VB_KEY_ESC;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Space-space");
	TEST_EQ(screens_displayed[0], VB_SCREEN_DEVELOPER_WARNING,
		"  warning screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_DEVELOPER_TO_NORM,
		"  tonorm screen");
	TEST_EQ(screens_displayed[2], VB_SCREEN_DEVELOPER_WARNING,
		"  warning screen");

	/* Enter doesn't by default */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	mock_keypress[0] = VB_KEY_ENTER;
	mock_keypress[1] = VB_KEY_ENTER;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Enter ignored");

	/* Enter does if GBB flag set */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	gbb.flags |= VB2_GBB_FLAG_ENTER_TRIGGERS_TONORM;
	mock_keypress[0] = VB_KEY_ENTER;
	mock_keypress[1] = VB_KEY_ENTER;
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Enter = tonorm");

	/* Tonorm ignored if GBB forces dev switch on */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	gbb.flags |= VB2_GBB_FLAG_FORCE_DEV_SWITCH_ON;
	mock_keypress[0] = ' ';
	mock_keypress[1] = VB_KEY_ENTER;
	TEST_EQ(VbBootDeveloper(ctx), 1002,
		"Can't tonorm gbb-dev");

	/* Shutdown requested at tonorm screen */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	mock_keypress[0] = ' ';
	MockGpioAfter(3, GPIO_SHUTDOWN);
	TEST_EQ(VbBootDeveloper(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Shutdown requested at tonorm");
	TEST_EQ(screens_displayed[0], VB_SCREEN_DEVELOPER_WARNING,
		"  warning screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_DEVELOPER_TO_NORM,
		"  tonorm screen");

	/* Shutdown requested by keyboard at tonorm screen */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	mock_keypress[0] = VB_BUTTON_POWER_SHORT_PRESS;
	TEST_EQ(VbBootDeveloper(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Shutdown requested by keyboard at nonorm");

	/* Ctrl+D dismisses warning */
	ResetMocks();
	mock_keypress[0] = VB_KEY_CTRL('D');
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+D");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST), 0,
		"  recovery reason");
	TEST_NEQ(audio_looping_calls_left, 0, "  aborts audio");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");

	/* Ctrl+D doesn't boot legacy even if GBB flag is set */
	ResetMocks();
	mock_keypress[0] = VB_KEY_CTRL('D');
	gbb.flags |= VB2_GBB_FLAG_DEFAULT_DEV_BOOT_LEGACY;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+D");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");

	/* Ctrl+L tries legacy boot mode only if enabled */
	ResetMocks();
	mock_keypress[0] = VB_KEY_CTRL('L');
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+L normal");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");

	/* Enter altfw menu and time out */
	ResetMocks();
	MockGpioAfter(1000, GPIO_SHUTDOWN);
	gbb.flags |= VB2_GBB_FLAG_FORCE_DEV_BOOT_LEGACY;
	mock_keypress[0] = VB_KEY_CTRL('L');
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_SHUTDOWN_REQUESTED,
		"Ctrl+L force legacy");
	TEST_EQ(vbexlegacy_called, 0, "  try legacy");

	/* Enter altfw menu and select firmware 0 */
	ResetMocks();
	gbb.flags |= VB2_GBB_FLAG_FORCE_DEV_BOOT_LEGACY;
	mock_keypress[0] = VB_KEY_CTRL('L');
	mock_keypress[1] = '0';
	TEST_EQ(VbBootDeveloper(ctx), 1002,
		"Ctrl+L force legacy");
	TEST_EQ(vbexlegacy_called, 1, "  try legacy");
	TEST_EQ(altfw_num, 0, "  check altfw_num");

	/* Enter altfw menu and then exit it */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_BOOT_LEGACY, 1);
	mock_keypress[0] = VB_KEY_CTRL('L');
	mock_keypress[1] = VB_KEY_ESC;
	TEST_EQ(VbBootDeveloper(ctx), 1002,
		"Ctrl+L nv legacy");
	TEST_EQ(vbexlegacy_called, 0, "  try legacy");

	/* Enter altfw menu and select firmware 0 */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_BOOT_LEGACY, 1);
	mock_keypress[0] = VB_KEY_CTRL('L');
	mock_keypress[1] = '0';
	TEST_EQ(VbBootDeveloper(ctx), 1002,
		"Ctrl+L nv legacy");
	TEST_EQ(vbexlegacy_called, 1, "  try legacy");
	TEST_EQ(altfw_num, 0, "  check altfw_num");

	/* Enter altfw menu and select firmware 0 */
	ResetMocks();
	VbApiKernelGetFwmp()->flags |= FWMP_DEV_ENABLE_LEGACY;
	mock_keypress[0] = VB_KEY_CTRL('L');
	mock_keypress[1] = '0';
	TEST_EQ(VbBootDeveloper(ctx), 1002,
		"Ctrl+L fwmp legacy");
	TEST_EQ(vbexlegacy_called, 1, "  fwmp legacy");
	TEST_EQ(altfw_num, 0, "  check altfw_num");

	/* Pressing 1-9 boots alternative firmware only if enabled */
	for (key = '1'; key <= '9'; key++) {
		ResetMocks();
		mock_keypress[0] = key;
		TEST_EQ(VbBootDeveloper(ctx), 1002, "'1' normal");
		TEST_EQ(vbexlegacy_called, 0, "  not legacy");

		ResetMocks();
		gbb.flags |= VB2_GBB_FLAG_FORCE_DEV_BOOT_LEGACY;
		mock_keypress[0] = key;
		TEST_EQ(VbBootDeveloper(ctx), 1002,
			"Ctrl+L force legacy");
		TEST_EQ(vbexlegacy_called, 1, "  try legacy");
		TEST_EQ(altfw_num, key - '0', "  check altfw_num");

		ResetMocks();
		vb2_nv_set(ctx, VB2_NV_DEV_BOOT_LEGACY, 1);
		mock_keypress[0] = key;
		TEST_EQ(VbBootDeveloper(ctx), 1002,
			"Ctrl+L nv legacy");
		TEST_EQ(vbexlegacy_called, 1, "  try legacy");
		TEST_EQ(altfw_num, key - '0', "  check altfw_num");

		ResetMocks();
		VbApiKernelGetFwmp()->flags |= FWMP_DEV_ENABLE_LEGACY;
		mock_keypress[0] = key;
		TEST_EQ(VbBootDeveloper(ctx), 1002,
			"Ctrl+L fwmp legacy");
		TEST_EQ(vbexlegacy_called, 1, "  fwmp legacy");
		TEST_EQ(altfw_num, key - '0', "  check altfw_num");
	}

	/* Ctrl+U boots USB only if enabled */
	ResetMocks();
	mock_keypress[0] = VB_KEY_CTRL('U');
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+U normal");

	/* Ctrl+U enabled, with good USB boot */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_BOOT_USB, 1);
	mock_keypress[0] = VB_KEY_CTRL('U');
	vbtlk_retval = VB2_SUCCESS - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootDeveloper(ctx), 0, "Ctrl+U USB");

	/* Ctrl+U enabled via GBB */
	ResetMocks();
	gbb.flags |= VB2_GBB_FLAG_FORCE_DEV_BOOT_USB;
	mock_keypress[0] = VB_KEY_CTRL('U');
	vbtlk_retval = VB2_SUCCESS - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootDeveloper(ctx), 0, "Ctrl+U force USB");

	/* Ctrl+U enabled via FWMP */
	ResetMocks();
	VbApiKernelGetFwmp()->flags |= FWMP_DEV_ENABLE_USB;
	mock_keypress[0] = VB_KEY_CTRL('U');
	vbtlk_retval = VB2_SUCCESS - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootDeveloper(ctx), 0, "Ctrl+U force USB");

	/* Ctrl+S set vendor data and reboot */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = '4';
	mock_keypress[2] = '3';
	mock_keypress[3] = '2';
	mock_keypress[4] = '1';
	mock_keypress[5] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[6] = VB_KEY_ENTER; // Confirm vendor data
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Ctrl+S set vendor data and reboot");
	TEST_EQ(set_vendor_data_called, 1, "  VbExSetVendorData() called");
	TEST_STR_EQ(set_vendor_data, "4321", "  Vendor data correct");

	/* Ctrl+S extra keys ignored */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = '4';
	mock_keypress[2] = '3';
	mock_keypress[3] = '2';
	mock_keypress[4] = '1';
	mock_keypress[5] = '5';
	mock_keypress[6] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[7] = VB_KEY_ENTER; // Confirm vendor data
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Ctrl+S extra keys ignored");
	TEST_EQ(set_vendor_data_called, 1, "  VbExSetVendorData() called");
	TEST_STR_EQ(set_vendor_data, "4321", "  Vendor data correct");

	/* Ctrl+S converts case */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = 'a';
	mock_keypress[2] = 'B';
	mock_keypress[3] = 'Y';
	mock_keypress[4] = 'z';
	mock_keypress[5] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[6] = VB_KEY_ENTER; // Confirm vendor data
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Ctrl+S converts case");
	TEST_EQ(set_vendor_data_called, 1, "  VbExSetVendorData() called");
	TEST_STR_EQ(set_vendor_data, "ABYZ", "  Vendor data correct");

	/* Ctrl+S backspace works */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = 'A';
	mock_keypress[2] = 'B';
	mock_keypress[3] = 'C';
	mock_keypress[4] = VB_KEY_BACKSPACE;
	mock_keypress[5] = VB_KEY_BACKSPACE;
	mock_keypress[6] = '3';
	mock_keypress[7] = '2';
	mock_keypress[8] = '1';
	mock_keypress[9] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[10] = VB_KEY_ENTER; // Confirm vendor data
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Ctrl+S backspace works");
	TEST_EQ(set_vendor_data_called, 1, "  VbExSetVendorData() called");
	TEST_STR_EQ(set_vendor_data, "A321", "  Vendor data correct");

	/* Ctrl+S invalid chars don't print */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = '4';
	mock_keypress[2] = '-';
	mock_keypress[3] = '^';
	mock_keypress[4] = '&';
	mock_keypress[5] = '$';
	mock_keypress[6] = '.';
	mock_keypress[7] = '3';
	mock_keypress[8] = '2';
	mock_keypress[9] = '1';
	mock_keypress[10] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[11] = VB_KEY_ENTER; // Confirm vendor data
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Ctrl+S invalid chars don't print");
	TEST_EQ(set_vendor_data_called, 1, "  VbExSetVendorData() called");
	TEST_STR_EQ(set_vendor_data, "4321", "  Vendor data correct");

	/* Ctrl+S invalid chars don't print with backspace */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = '4';
	mock_keypress[2] = '-';
	mock_keypress[3] = VB_KEY_BACKSPACE; // Should delete 4
	mock_keypress[4] = '3';
	mock_keypress[5] = '2';
	mock_keypress[6] = '1';
	mock_keypress[7] = '0';
	mock_keypress[8] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[9] = VB_KEY_ENTER; // Confirm vendor data
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Ctrl+S invalid chars don't print with backspace");
	TEST_EQ(set_vendor_data_called, 1, "  VbExSetVendorData() called");
	TEST_STR_EQ(set_vendor_data, "3210", "  Vendor data correct");

	/* Ctrl+S backspace only doesn't underrun */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = 'A';
	mock_keypress[2] = VB_KEY_BACKSPACE;
	mock_keypress[3] = VB_KEY_BACKSPACE;
	mock_keypress[4] = '4';
	mock_keypress[5] = '3';
	mock_keypress[6] = '2';
	mock_keypress[7] = '1';
	mock_keypress[8] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[9] = VB_KEY_ENTER; // Confirm vendor data
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Ctrl+S backspace only doesn't underrun");
	TEST_EQ(set_vendor_data_called, 1, "  VbExSetVendorData() called");
	TEST_STR_EQ(set_vendor_data, "4321", "  Vendor data correct");

	/* Ctrl+S vowels not allowed after first char */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = 'A';
	mock_keypress[2] = 'A';
	mock_keypress[3] = 'B';
	mock_keypress[4] = 'E';
	mock_keypress[5] = 'i';
	mock_keypress[6] = 'C';
	mock_keypress[7] = 'O';
	mock_keypress[8] = 'u';
	mock_keypress[9] = 'D';
	mock_keypress[10] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[11] = VB_KEY_ENTER; // Confirm vendor data
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"Ctrl+S vowels not allowed after first char");
	TEST_EQ(set_vendor_data_called, 1, "  VbExSetVendorData() called");
	TEST_STR_EQ(set_vendor_data, "ABCD", "  Vendor data correct");

	/* Ctrl+S too short */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = '1';
	mock_keypress[2] = '2';
	mock_keypress[3] = '3';
	mock_keypress[4] = VB_KEY_ENTER; // Set vendor data (Nothing happens)
	mock_keypress[5] = VB_KEY_ENTER; // Confirm vendor data (Nothing happens)
	mock_keypress[6] = VB_KEY_ESC;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+S too short");
	TEST_EQ(set_vendor_data_called, 0, "  VbExSetVendorData() not called");

	/* Ctrl+S esc from set screen */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = VB_KEY_ESC;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+S esc from set screen");
	TEST_EQ(set_vendor_data_called, 0, "  VbExSetVendorData() not called");

	/* Ctrl+S esc from set screen with tag */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = '4';
	mock_keypress[2] = '3';
	mock_keypress[3] = '2';
	mock_keypress[4] = '1';
	mock_keypress[5] = VB_KEY_ESC;
	TEST_EQ(VbBootDeveloper(ctx), 1002,
		"Ctrl+S esc from set screen with tag");
	TEST_EQ(set_vendor_data_called, 0, "  VbExSetVendorData() not called");

	/* Ctrl+S esc from confirm screen */
	ResetMocks();
	ctx->flags |= VB2_CONTEXT_VENDOR_DATA_SETTABLE;
	mock_keypress[0] = VB_KEY_CTRL('S');
	mock_keypress[1] = '4';
	mock_keypress[2] = '3';
	mock_keypress[3] = '2';
	mock_keypress[4] = '1';
	mock_keypress[5] = VB_KEY_ENTER; // Set vendor data
	mock_keypress[6] = VB_KEY_ESC;
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+S esc from set screen");
	TEST_EQ(set_vendor_data_called, 0, "  VbExSetVendorData() not called");

	/* If no USB, eventually times out and tries fixed disk */
	ResetMocks();
	vb2_nv_set(ctx, VB2_NV_DEV_BOOT_USB, 1);
	mock_keypress[0] = VB_KEY_CTRL('U');
	TEST_EQ(VbBootDeveloper(ctx), 1002, "Ctrl+U enabled");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST), 0,
		"  recovery reason");
	TEST_EQ(audio_looping_calls_left, 0, "  used up audio");

	/* If dev mode is disabled, goes to TONORM screen repeatedly */
	ResetMocks();
	VbApiKernelGetFwmp()->flags |= FWMP_DEV_DISABLE_BOOT;
	mock_keypress[0] = VB_KEY_ESC;  /* Just causes TONORM again */
	mock_keypress[1] = VB_KEY_ENTER;
	TEST_EQ(VbBootDeveloper(ctx), VBERROR_REBOOT_REQUIRED,
		"FWMP dev disabled");
	TEST_EQ(screens_displayed[0], VB_SCREEN_DEVELOPER_TO_NORM,
		"  tonorm screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_DEVELOPER_TO_NORM,
		"  tonorm screen");
	TEST_EQ(screens_displayed[2], VB_SCREEN_TO_NORM_CONFIRMED,
		"  confirm screen");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DISABLE_DEV_REQUEST), 1,
		"  disable dev request");

	/* Shutdown requested when dev disabled */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	VbApiKernelGetFwmp()->flags |= FWMP_DEV_DISABLE_BOOT;
	MockGpioAfter(1, GPIO_SHUTDOWN);
	TEST_EQ(VbBootDeveloper(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Shutdown requested when dev disabled");
	TEST_EQ(screens_displayed[0], VB_SCREEN_DEVELOPER_TO_NORM,
		"  tonorm screen");

	/* Shutdown requested by keyboard when dev disabled */
	ResetMocks();
	shared->flags = VBSD_BOOT_DEV_SWITCH_ON;
	VbApiKernelGetFwmp()->flags |= FWMP_DEV_DISABLE_BOOT;
	mock_keypress[0] = VB_BUTTON_POWER_SHORT_PRESS;
	TEST_EQ(VbBootDeveloper(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Shutdown requested by keyboard when dev disabled");

	printf("...done.\n");
}

/*
 * Helper function to test VbBootRecovery against a sequence of gpio events.
 * caller specifies a sequence of gpio events and the expected confirm vs.
 * reboot result.
 *
 * Non-asserted gpios are used for 5 events, then 'first' for 5 events,
 * 'second' for 5 events, and 'third' for 5 events.
 */
static void VbBootRecTestGpio(uint32_t first, uint32_t second, uint32_t third,
			      uint32_t confirm, const char *msg)
{
	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	trust_ec = 1;
	mock_keypress[0] = VB_KEY_CTRL('D');
	mock_gpio[0].gpio_flags = 0;
	mock_gpio[0].count = 4;
	mock_gpio[1].gpio_flags = first;
	mock_gpio[1].count = 4;
	mock_gpio[2].gpio_flags = second;
	mock_gpio[2].count = 4;
	mock_gpio[3].gpio_flags = third;
	mock_gpio[3].count = 4;

	if (confirm) {
		TEST_EQ(VbBootRecovery(ctx), VBERROR_EC_REBOOT_TO_RO_REQUIRED,
			msg);
		TEST_EQ(virtdev_set, 1, "  virtual dev mode on");
	} else {
		TEST_EQ(VbBootRecovery(ctx), VBERROR_SHUTDOWN_REQUESTED, msg);
		TEST_EQ(virtdev_set, 0, "  virtual dev mode off");
	}
}

static void VbBootRecTest(void)
{
	printf("Testing VbBootRecovery()...\n");

	/* Shutdown requested in loop */
	ResetMocks();
	MockGpioAfter(10, GPIO_SHUTDOWN);
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Shutdown requested");

	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST), 0,
		"  recovery reason");
	TEST_EQ(screens_displayed[0], VB_SCREEN_OS_BROKEN,
		"  broken screen");

	/* Shutdown requested by keyboard */
	ResetMocks();
	mock_keypress[0] = VB_BUTTON_POWER_SHORT_PRESS;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Shutdown requested by keyboard");

	/* Ignore power button held on boot */
	ResetMocks();
	mock_gpio[0].gpio_flags = GPIO_SHUTDOWN;
	mock_gpio[0].count = 10;
	mock_gpio[1].gpio_flags = 0;
	mock_gpio[1].count = 10;
	mock_gpio[2].gpio_flags = GPIO_SHUTDOWN;
	mock_gpio[2].count = 10;
	mock_gpio[3].gpio_flags = 0;
	mock_gpio[3].count = 100;
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	trust_ec = 1;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Ignore power button held on boot");
	TEST_EQ(screens_displayed[0], VB_SCREEN_RECOVERY_INSERT,
		"  insert screen");
	/* Shutdown should happen while we're sending the 2nd block of events */
	TEST_EQ(mock_gpio_count, 3, "  ignore held button");

	/* Broken screen */
	ResetMocks();
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_num_disks[0] = 1;
	mock_num_disks[1] = 1;
	mock_num_disks[2] = 1;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Broken");
	TEST_EQ(screens_displayed[0], VB_SCREEN_OS_BROKEN,
		"  broken screen");

	/* Broken screen even if dev switch is on */
	ResetMocks();
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_num_disks[0] = 1;
	mock_num_disks[1] = 1;
	shared->flags |= VBSD_BOOT_DEV_SWITCH_ON;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Broken (dev)");
	TEST_EQ(screens_displayed[0], VB_SCREEN_OS_BROKEN,
		"  broken screen");

	/* Force insert screen with GBB flag */
	ResetMocks();
	MockGpioAfter(100, GPIO_SHUTDOWN);
	gbb.flags |= VB2_GBB_FLAG_FORCE_MANUAL_RECOVERY;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Insert (forced by GBB)");
	TEST_EQ(screens_displayed[0], VB_SCREEN_RECOVERY_INSERT,
		"  insert screen");

	/* No removal if recovery button physically pressed */
	ResetMocks();
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_num_disks[0] = 1;
	mock_num_disks[1] = 1;
	shared->flags |= VBSD_BOOT_REC_SWITCH_ON;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"No remove in rec");
	TEST_EQ(screens_displayed[0], VB_SCREEN_OS_BROKEN,
		"  broken screen");

	/* Removal if no disk initially found, but found on second attempt */
	ResetMocks();
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_num_disks[0] = 0;
	mock_num_disks[1] = 1;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Remove");
	TEST_EQ(screens_displayed[0], VB_SCREEN_OS_BROKEN,
		"  broken screen");

	/* Bad disk count doesn't require removal */
	ResetMocks();
	MockGpioAfter(10, GPIO_SHUTDOWN);
	mock_num_disks[0] = -1;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Bad disk count");
	TEST_EQ(screens_displayed[0], VB_SCREEN_OS_BROKEN,
		"  broken screen");

	/* Ctrl+D ignored for many reasons... */
	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_keypress[0] = VB_KEY_CTRL('D');
	trust_ec = 0;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Ctrl+D ignored if EC not trusted");
	TEST_EQ(virtdev_set, 0, "  virtual dev mode off");
	TEST_NEQ(screens_displayed[1], VB_SCREEN_RECOVERY_TO_DEV,
		 "  todev screen");

	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON | VBSD_BOOT_DEV_SWITCH_ON;
	trust_ec = 1;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_keypress[0] = VB_KEY_CTRL('D');
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Ctrl+D ignored if already in dev mode");
	TEST_EQ(virtdev_set, 0, "  virtual dev mode off");
	TEST_NEQ(screens_displayed[1], VB_SCREEN_RECOVERY_TO_DEV,
		 "  todev screen");

	ResetMocks();
	trust_ec = 1;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_keypress[0] = VB_KEY_CTRL('D');
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Ctrl+D ignored if recovery not manually triggered");
	TEST_EQ(virtdev_set, 0, "  virtual dev mode off");
	TEST_NEQ(screens_displayed[1], VB_SCREEN_RECOVERY_TO_DEV,
		 "  todev screen");

	/* Ctrl+D ignored because the physical presence switch is still pressed
	 * and we don't like that.
	 */
	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	trust_ec = 1;
	mock_keypress[0] = VB_KEY_CTRL('D');
	mock_gpio[0].gpio_flags = GPIO_PRESENCE;
	mock_gpio[0].count = 100;
	mock_gpio[1].gpio_flags = GPIO_PRESENCE | GPIO_SHUTDOWN;
	mock_gpio[1].count = 100;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Ctrl+D ignored if phys pres button is still pressed");
	TEST_NEQ(screens_displayed[1], VB_SCREEN_RECOVERY_TO_DEV,
		 "  todev screen");

	/* Ctrl+D then space means don't enable */
	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	trust_ec = 1;
	mock_keypress[0] = VB_KEY_CTRL('D');
	mock_keypress[1] = ' ';
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_SHUTDOWN_REQUESTED,
		"Ctrl+D todev abort");
	TEST_EQ(screens_displayed[0], VB_SCREEN_RECOVERY_INSERT,
		"  insert screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_RECOVERY_TO_DEV,
		"  todev screen");
	TEST_EQ(screens_displayed[2], VB_SCREEN_RECOVERY_INSERT,
		"  insert screen");
	TEST_EQ(virtdev_set, 0, "  virtual dev mode off");

	/* Ctrl+D then enter means enable */
	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	trust_ec = 1;
	mock_keypress[0] = VB_KEY_CTRL('D');
	mock_keypress[1] = VB_KEY_ENTER;
	mock_keyflags[1] = VB_KEY_FLAG_TRUSTED_KEYBOARD;
	TEST_EQ(VbBootRecovery(ctx), VBERROR_EC_REBOOT_TO_RO_REQUIRED,
		"Ctrl+D todev confirm via enter");
	TEST_EQ(virtdev_set, 1, "  virtual dev mode on");

	/*
	 * List of possiblities for shutdown and physical presence events that
	 * occur over time.  Time advanced from left to right (where each
	 * represents the gpio[s] that are seen during a given iteration of
	 * the loop).  The meaning of the characters:
	 *
	 *   _ means no gpio
	 *   s means shutdown gpio
	 *   p means presence gpio
	 *   B means both shutdown and presence gpio
	 *
	 *  1: ______ppp______ -> confirm
	 *  2: ______sss______ -> shutdown
	 *  3: ___pppsss______ -> confirm
	 *  4: ___sssppp______ -> shutdown
	 *  5: ___pppBBB______ -> confirm
	 *  6: ___pppBBBppp___ -> shutdown
	 *  7: ___pppBBBsss___ -> confirm
	 *  8: ___sssBBB______ -> confirm
	 *  9: ___sssBBBppp___ -> shutdown
	 * 10: ___sssBBBsss___ -> confirm
	 * 11: ______BBB______ -> confirm
	 * 12: ______BBBsss___ -> confirm
	 * 13: ______BBBppp___ -> shutdown
	 */

	/* 1: Ctrl+D then presence means enable */
	VbBootRecTestGpio(GPIO_PRESENCE, 0, 0, 1,
			  "Ctrl+D todev confirm via presence");

	/* 2: Ctrl+D then shutdown means shutdown */
	VbBootRecTestGpio(GPIO_SHUTDOWN, 0, 0, 0,
			  "Ctrl+D todev then shutdown");

	/* 3: Ctrl+D then presence then shutdown means confirm */
	VbBootRecTestGpio(GPIO_PRESENCE, GPIO_SHUTDOWN, 0, 1,
			  "Ctrl+D todev confirm via presence then shutdown");

	/* 4: Ctrl+D then 2+ instance shutdown then presence means shutdown */
	VbBootRecTestGpio(GPIO_SHUTDOWN, GPIO_PRESENCE, 0, 0,
			  "Ctrl+D todev then 2+ shutdown then presence");

	/* 5: Ctrl+D then presence then shutdown+presence then none */
	VbBootRecTestGpio(GPIO_PRESENCE, GPIO_PRESENCE | GPIO_SHUTDOWN, 0, 1,
			  "Ctrl+D todev confirm via presence, both, none");

	/* 6: Ctrl+D then presence then shutdown+presence then presence */
	VbBootRecTestGpio(GPIO_PRESENCE, GPIO_PRESENCE | GPIO_SHUTDOWN,
			  GPIO_PRESENCE, 0,
			  "Ctrl+D todev confirm via presence, both, presence");

	/* 7: Ctrl+D then presence then shutdown+presence then shutdown */
	VbBootRecTestGpio(GPIO_PRESENCE, GPIO_PRESENCE | GPIO_SHUTDOWN,
			  GPIO_SHUTDOWN, 1,
			  "Ctrl+D todev confirm via presence, both, shutdown");

	/* 8: Ctrl+D then shutdown then shutdown+presence then none */
	VbBootRecTestGpio(GPIO_SHUTDOWN, GPIO_PRESENCE | GPIO_SHUTDOWN, 0, 1,
			  "Ctrl+D todev then 2+ shutdown, both, none");

	/* 9: Ctrl+D then shutdown then shutdown+presence then presence */
	VbBootRecTestGpio(GPIO_SHUTDOWN, GPIO_PRESENCE | GPIO_SHUTDOWN,
			  GPIO_PRESENCE, 0,
			  "Ctrl+D todev then 2+ shutdown, both, presence");

	/* 10: Ctrl+D then shutdown then shutdown+presence then shutdown */
	VbBootRecTestGpio(GPIO_SHUTDOWN, GPIO_PRESENCE | GPIO_SHUTDOWN,
			  GPIO_SHUTDOWN, 1,
			  "Ctrl+D todev then 2+ shutdown, both, shutdown");

	/* 11: Ctrl+D then shutdown+presence then none */
	VbBootRecTestGpio(GPIO_PRESENCE | GPIO_SHUTDOWN, 0, 0, 1,
			  "Ctrl+D todev confirm via both then none");

	/* 12: Ctrl+D then shutdown+presence then shutdown */
	VbBootRecTestGpio(GPIO_PRESENCE | GPIO_SHUTDOWN, GPIO_SHUTDOWN, 0, 1,
			  "Ctrl+D todev confirm via both then shutdown");

	/* 13: Ctrl+D then shutdown+presence then presence */
	VbBootRecTestGpio(GPIO_PRESENCE | GPIO_SHUTDOWN, GPIO_PRESENCE, 0, 0,
			  "Ctrl+D todev confirm via both then presence");

	/* Handle TPM error in enabling dev mode */
	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	trust_ec = 1;
	mock_keypress[0] = VB_KEY_CTRL('D');
	mock_keypress[1] = VB_KEY_ENTER;
	mock_keyflags[1] = VB_KEY_FLAG_TRUSTED_KEYBOARD;
	virtdev_retval = VB2_ERROR_MOCK;
	TEST_EQ(VbBootRecovery(ctx),
		VBERROR_TPM_SET_BOOT_MODE_STATE,
		"Ctrl+D todev failure");

	/* Test Diagnostic Mode via Ctrl-C - display available */
	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	trust_ec = 1;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_keypress[0] = VB_KEY_CTRL('C');
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST), 0,
		"todiag is zero");
	if (DIAGNOSTIC_UI)
		TEST_EQ(VbBootRecovery(ctx),
			VBERROR_REBOOT_REQUIRED,
			"Ctrl+C todiag - enabled");
	else
		TEST_EQ(VbBootRecovery(ctx),
			VBERROR_SHUTDOWN_REQUESTED,
			"Ctrl+C todiag - disabled");

	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST), DIAGNOSTIC_UI,
		"  todiag is updated for Ctrl-C");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DISPLAY_REQUEST), 0,
		"  todiag doesn't set unneeded DISPLAY_REQUEST");
	TEST_EQ(screens_displayed[0], VB_SCREEN_RECOVERY_INSERT,
		"  insert screen");

	/* Test Diagnostic Mode via F12 - display disabled */
	ResetMocks();
	shared->flags = VBSD_BOOT_REC_SWITCH_ON;
	sd->flags &= ~VB2_SD_FLAG_DISPLAY_AVAILABLE;
	trust_ec = 1;
	vbtlk_retval = VBERROR_NO_DISK_FOUND - VB_DISK_FLAG_REMOVABLE;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_keypress[0] = VB_KEY_F(12);
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST), 0,
		"todiag is zero");
	if (DIAGNOSTIC_UI)
		TEST_EQ(VbBootRecovery(ctx),
			VBERROR_REBOOT_REQUIRED,
			"F12 todiag - enabled");
	else
		TEST_EQ(VbBootRecovery(ctx),
			VBERROR_SHUTDOWN_REQUESTED,
			"F12 todiag - disabled");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST), DIAGNOSTIC_UI,
		"  todiag is updated for F12");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DISPLAY_REQUEST), 0,
		"  todiag doesn't set unneeded DISPLAY_REQUEST");
	TEST_EQ(screens_displayed[0], VB_SCREEN_RECOVERY_INSERT,
		"  insert screen");

	/* Test Diagnostic Mode via Ctrl-C OS broken - display available */
	ResetMocks();
	shared->flags = 0;
	MockGpioAfter(100, GPIO_SHUTDOWN);
	mock_keypress[0] = VB_KEY_CTRL('C');
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST), 0,
		"todiag is zero");
	if (DIAGNOSTIC_UI)
		TEST_EQ(VbBootRecovery(ctx),
			VBERROR_REBOOT_REQUIRED,
			"Ctrl+C todiag os broken - enabled");
	else
		TEST_EQ(VbBootRecovery(ctx),
			VBERROR_SHUTDOWN_REQUESTED,
			"Ctrl+C todiag os broken - disabled");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DIAG_REQUEST), DIAGNOSTIC_UI,
		"  todiag is updated for Ctrl-C");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_DISPLAY_REQUEST), 0,
		"  todiag doesn't set unneeded DISPLAY_REQUEST");
	TEST_EQ(screens_displayed[0], VB_SCREEN_OS_BROKEN,
		"  os broken screen");

	printf("...done.\n");
}

static void VbBootDiagTest(void)
{
	printf("Testing VbBootDiagnostic()...\n");

	/* No key pressed - timeout. */
	ResetMocks();
	TEST_EQ(VbBootDiagnostic(ctx), VBERROR_REBOOT_REQUIRED, "Timeout");
	TEST_EQ(screens_displayed[0], VB_SCREEN_CONFIRM_DIAG,
		"  confirm screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_BLANK,
		"  blank screen");
	TEST_EQ(tpm_set_mode_called, 0, "  no tpm call");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");
	TEST_EQ(current_ticks, 30 * VB_USEC_PER_SEC,
		"  waited for 30 seconds");

	/* Esc key pressed. */
	ResetMocks();
	mock_keypress[0] = VB_KEY_ESC;
	TEST_EQ(VbBootDiagnostic(ctx), VBERROR_REBOOT_REQUIRED, "Esc key");
	TEST_EQ(screens_displayed[0], VB_SCREEN_CONFIRM_DIAG,
		"  confirm screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_BLANK,
		"  blank screen");
	TEST_EQ(tpm_set_mode_called, 0, "  no tpm call");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");
	TEST_EQ(current_ticks, 0, "  didn't wait at all");

	/* Shutdown requested via lid close */
	ResetMocks();
	MockGpioAfter(10, GPIO_LID_CLOSED);
	TEST_EQ(VbBootDiagnostic(ctx), VBERROR_SHUTDOWN_REQUESTED, "Shutdown");
	TEST_EQ(screens_displayed[0], VB_SCREEN_CONFIRM_DIAG,
		"  confirm screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_BLANK,
		"  blank screen");
	TEST_EQ(tpm_set_mode_called, 0, "  no tpm call");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");
	TEST_TRUE(current_ticks < VB_USEC_PER_SEC, "  didn't wait long");

	/* Power button pressed but not released. */
	ResetMocks();
	mock_gpio[0].gpio_flags = GPIO_PRESENCE;
	mock_gpio[0].count = ~0;
	TEST_EQ(VbBootDiagnostic(ctx), VBERROR_REBOOT_REQUIRED, "Power held");
	TEST_EQ(screens_displayed[0], VB_SCREEN_CONFIRM_DIAG,
		"  confirm screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_BLANK,
		"  blank screen");
	TEST_EQ(tpm_set_mode_called, 0, "  no tpm call");
	TEST_EQ(vbexlegacy_called, 0, "  not legacy");

	/* Power button is pressed and released. */
	ResetMocks();
	MockGpioAfter(3, GPIO_PRESENCE);
	TEST_EQ(VbBootDiagnostic(ctx), VBERROR_REBOOT_REQUIRED, "Confirm");
	TEST_EQ(screens_displayed[0], VB_SCREEN_CONFIRM_DIAG,
		"  confirm screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_BLANK,
		"  blank screen");
	TEST_EQ(tpm_set_mode_called, 1, "  tpm call");
	TEST_EQ(tpm_mode, VB2_TPM_MODE_DISABLED, "  tpm disabled");
	TEST_EQ(vbexlegacy_called, 1, "  legacy");
	TEST_EQ(altfw_num, VB_ALTFW_DIAGNOSTIC, "  check altfw_num");
	/*
	 * Ideally we'd that no recovery request was recorded, but
	 * VbExLegacy() can only fail or crash the tests.
	 */
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		VB2_RECOVERY_ALTFW_HASH_FAILED,
		"  recovery request");

        /* Power button confirm, but now with a tpm failure. */
	ResetMocks();
	tpm_mode = VB2_TPM_MODE_DISABLED;
	mock_gpio[0].gpio_flags = 0;
	mock_gpio[0].count = 2;
	mock_gpio[1].gpio_flags = GPIO_PRESENCE;
	mock_gpio[1].count = 2;
	TEST_EQ(VbBootDiagnostic(ctx), VBERROR_REBOOT_REQUIRED,
		"Confirm but tpm fail");
	TEST_EQ(screens_displayed[0], VB_SCREEN_CONFIRM_DIAG,
		"  confirm screen");
	TEST_EQ(screens_displayed[1], VB_SCREEN_BLANK,
		"  blank screen");
	TEST_EQ(tpm_set_mode_called, 1, "  tpm call");
	TEST_EQ(tpm_mode, VB2_TPM_MODE_DISABLED, "  tpm disabled");
	TEST_EQ(vbexlegacy_called, 0, "  legacy not called");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		VB2_RECOVERY_TPM_DISABLE_FAILED,
		"  recovery request");

	printf("...done.\n");
}


int main(void)
{
	VbUserConfirmsTest();
	VbBootTest();
	VbBootDevTest();
	VbBootRecTest();
	if (DIAGNOSTIC_UI)
		VbBootDiagTest();

	return gTestSuccess ? 0 : 255;
}
