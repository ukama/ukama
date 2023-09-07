/* Copyright (c) 2013 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * APIs between calling firmware and vboot_reference
 *
 * General notes:
 *
 * TODO: split this file into a vboot_entry_points.h file which contains the
 * entry points for the firmware to call vboot_reference, and a
 * vboot_firmware_exports.h which contains the APIs to be implemented by the
 * calling firmware and exported to vboot_reference.
 *
 * Notes:
 *    * Assumes this code is never called in the S3 resume path.  TPM resume
 *      must be done elsewhere, and VB2_NV_DEBUG_RESET_MODE is ignored.
 */

#ifndef VBOOT_REFERENCE_2API_H_
#define VBOOT_REFERENCE_2API_H_

#include "2constants.h"
#include "2crypto.h"
#include "2fw_hash_tags.h"
#include "2gbb_flags.h"
#include "2id.h"
#include "2recovery_reasons.h"
#include "2return_codes.h"

/* TODO(chromium:972956): Remove once coreboot is using updated names */
#define secdata secdata_firmware
#define secdatak secdata_kernel
#define vb2api_secdata_check vb2api_secdata_firmware_check
#define vb2api_secdata_create vb2api_secdata_firmware_create
#define vb2api_secdatak_check vb2api_secdata_kernel_check
#define vb2api_secdatak_create vb2api_secdata_kernel_create

/* Modes for vb2ex_tpm_set_mode. */
enum vb2_tpm_mode {
	/*
	 * TPM is enabled tentatively, and may be set to either
	 * ENABLED or DISABLED mode.
	 */
	VB2_TPM_MODE_ENABLED_TENTATIVE = 0,

	/* TPM is enabled, and mode may not be changed. */
	VB2_TPM_MODE_ENABLED = 1,

	/* TPM is disabled, and mode may not be changed. */
	VB2_TPM_MODE_DISABLED = 2,
};

/* Flags for vb2_context.
 *
 * Unless otherwise noted, flags are set by verified boot and may be read (but
 * not set or cleared) by the caller.
 */
enum vb2_context_flags {

	/*
	 * Verified boot has changed nvdata[].  Caller must save nvdata[] back
	 * to its underlying storage, then may clear this flag.
	 */
	VB2_CONTEXT_NVDATA_CHANGED = (1 << 0),

	/*
	 * Verified boot has changed secdata_firmware[].  Caller must save
	 * secdata_firmware[] back to its underlying storage, then may clear
	 * this flag.
	 */
	VB2_CONTEXT_SECDATA_FIRMWARE_CHANGED = (1 << 1),
	/* TODO: Remove once coreboot has switched over */
	VB2_CONTEXT_SECDATA_CHANGED = (1 << 1),

	/* Recovery mode is requested this boot */
	VB2_CONTEXT_RECOVERY_MODE = (1 << 2),

	/* Developer mode is requested this boot */
	VB2_CONTEXT_DEVELOPER_MODE = (1 << 3),

	/*
	 * Force recovery mode due to physical user request.  Caller may set
	 * this flag when initializing the context.
	 */
	VB2_CONTEXT_FORCE_RECOVERY_MODE = (1 << 4),

	/*
	 * Force developer mode enabled.  Caller may set this flag when
	 * initializing the context.  Previously used for forcing developer
	 * mode with physical dev switch.
	 *
	 * Deprecated as part of chromium:942901.
	 */
	VB2_CONTEXT_DEPRECATED_FORCE_DEVELOPER_MODE = (1 << 5),

	/* Using firmware slot B.  If this flag is clear, using slot A. */
	VB2_CONTEXT_FW_SLOT_B = (1 << 6),

	/* RAM should be cleared by caller this boot */
	VB2_CONTEXT_CLEAR_RAM = (1 << 7),

	/* Wipeout by the app should be requested. */
	VB2_CONTEXT_FORCE_WIPEOUT_MODE = (1 << 8),

	/* Erase TPM developer mode state if it is enabled. */
	VB2_CONTEXT_DISABLE_DEVELOPER_MODE = (1 << 9),

	/*
	 * Verified boot has changed secdata_kernel[].  Caller must save
	 * secdata_kernel[] back to its underlying storage, then may clear
	 * this flag.
	 */
	VB2_CONTEXT_SECDATA_KERNEL_CHANGED = (1 << 10),

	/*
	 * Allow kernel verification to roll forward the version in
	 * secdata_kernel[].  Caller may set this flag before calling
	 * vb2api_kernel_phase3().
	 */
	VB2_CONTEXT_ALLOW_KERNEL_ROLL_FORWARD = (1 << 11),

	/*
	 * Boot optimistically: don't touch failure counters.  Caller may set
	 * this flag when initializing the context.
	 */
	VB2_CONTEXT_NOFAIL_BOOT = (1 << 12),

	/*
	 * secdata is not ready this boot, but should be ready next boot.  It
	 * would like to reboot.  The decision whether to reboot or not must be
	 * deferred until vboot, because rebooting all the time before then
	 * could cause a device with malfunctioning secdata to get stuck in an
	 * unrecoverable crash loop.
	 */
	VB2_CONTEXT_SECDATA_WANTS_REBOOT = (1 << 13),

	/*
	 * Boot is S3->S0 resume, not S5->S0 normal boot.  Caller may set this
	 * flag when initializing the context.
	 */
	VB2_CONTEXT_S3_RESUME = (1 << 14),

	/*
	 * System supports EC software sync.  Caller may set this flag at any
	 * time before calling VbSelectAndLoadKernel().
	 */
	VB2_CONTEXT_EC_SYNC_SUPPORTED = (1 << 15),

	/*
	 * EC software sync is slow to update; warning screen should be
	 * displayed.  Caller may set this flag at any time before calling
	 * VbSelectAndLoadKernel().
	 */
	VB2_CONTEXT_EC_SYNC_SLOW = (1 << 16),

	/*
	 * EC firmware supports early firmware selection; two EC images exist,
	 * and EC may have already verified and jumped to EC-RW prior to EC
	 * software sync.
	 */
	VB2_CONTEXT_EC_EFS = (1 << 17),

	/*
	 * NV storage uses data format V2.  Data is size VB2_NVDATA_SIZE_V2,
	 * not VB2_NVDATA_SIZE.
	 *
	 * Caller must set this flag when initializing the context to use V2.
	 * (Vboot cannot infer the data size from the data itself, because the
	 * data provided by the caller could be uninitialized.)
	 */
	VB2_CONTEXT_NVDATA_V2 = (1 << 18),

	/* Allow vendor data to be set via the vendor data ui. */
	VB2_CONTEXT_VENDOR_DATA_SETTABLE = (1 << 19),

	/*
	 * Caller may set this before running vb2api_fw_phase1.  In this case,
	 * it means: "Display is available on this boot.  Please advertise
	 * as such to downstream vboot code and users."
	 *
	 * vboot may also set this before returning from vb2api_fw_phase1.
	 * In this case, it means: "Please initialize display so that it is
	 * available to downstream vboot code and users."  This is used when
	 * vboot encounters some internally-generated request for display
	 * support.
	 */
	VB2_CONTEXT_DISPLAY_INIT = (1 << 20),

	/*
	 * Caller may set this before running vb2api_kernel_phase1.  It means
	 * that there is no FWMP on this system, and thus default values should
	 * be used instead.
	 *
	 * Caller should *not* set this when FWMP is available but invalid.
	 */
	VB2_CONTEXT_NO_SECDATA_FWMP = (1 << 21),
};

/* Helper for aligning fields in vb2_context. */
#define VB2_PAD_STRUCT3(size, align, count) \
	uint8_t _pad##count[align - (((size - 1) % align) + 1)]
#define VB2_PAD_STRUCT2(size, align, count) VB2_PAD_STRUCT3(size, align, count)
#define VB2_PAD_STRUCT(size, align) VB2_PAD_STRUCT2(size, align, __COUNTER__)

/*
 * Context for firmware verification.  Pass this to all vboot APIs.
 *
 * Context is stored as part of vb2_shared_data, initialized with vb2api_init().
 * Subsequent retrieval of the context object should be done by calling
 * vb2api_reinit(), e.g. if switching firmware applications.
 *
 * The context struct can be seen as the "publicly accessible" portion of
 * vb2_shared_data, and thus does not require its own magic and version fields.
 */
struct vb2_context {

	/**********************************************************************
	 * Fields caller must initialize before calling any API functions.
	 */

	/*
	 * Flags; see vb2_context_flags.  Some flags may only be set by caller
	 * prior to calling vboot functions.
	 */
	uint64_t flags;

	/*
	 * Non-volatile data.  Caller must fill this from some non-volatile
	 * location before calling vb2api_fw_phase1.  If the
	 * VB2_CONTEXT_NVDATA_CHANGED flag is set when a vb2api function
	 * returns, caller must save the data back to the non-volatile location
	 * and then clear the flag.
	 */
	uint8_t nvdata[VB2_NVDATA_SIZE_V2];
	VB2_PAD_STRUCT(VB2_NVDATA_SIZE_V2, 8);

	/*
	 * Secure data for firmware verification stage.  Caller must fill this
	 * from some secure non-volatile location before calling
	 * vb2api_fw_phase1.  If the VB2_CONTEXT_SECDATA_CHANGED flag is set
	 * when a function returns, caller must save the data back to the
	 * secure non-volatile location and then clear the flag.
	 */
	uint8_t secdata_firmware[VB2_SECDATA_FIRMWARE_SIZE];
	VB2_PAD_STRUCT(VB2_SECDATA_FIRMWARE_SIZE, 8);

	/**********************************************************************
	 * Fields caller must initialize before calling vb2api_kernel_phase1().
	 */

	/*
	 * Secure data for kernel verification stage.  Caller must fill this
	 * from some secure non-volatile location before calling
	 * vb2api_kernel_phase1.  If the VB2_CONTEXT_SECDATA_KERNEL_CHANGED
	 * flag is set when a function returns, caller must save the data back
	 * to the secure non-volatile location and then clear the flag.
	 */
	uint8_t secdata_kernel[VB2_SECDATA_KERNEL_SIZE];
	VB2_PAD_STRUCT(VB2_SECDATA_KERNEL_SIZE, 8);

	/*
	 * Firmware management parameters (FWMP) secure data.  Caller must fill
	 * this from some secure non-volatile location before calling
	 * vb2api_kernel_phase1.  Since FWMP is a variable-size space, caller
	 * should initially fill in VB2_SECDATA_FWMP_MIN_SIZE bytes, and call
	 * vb2_secdata_fwmp_check() to see whether more should be read.  If the
	 * VB2_CONTEXT_SECDATA_FWMP_CHANGED flag is set when a function
	 * returns, caller must save the data back to the secure non-volatile
	 * location and then clear the flag.
	 */
	uint8_t secdata_fwmp[VB2_SECDATA_FWMP_MAX_SIZE];
	VB2_PAD_STRUCT(VB2_SECDATA_FWMP_MAX_SIZE, 8);
};

/* Resource index for vb2ex_read_resource() */
enum vb2_resource_index {

	/* Google binary block */
	VB2_RES_GBB,

	/*
	 * Firmware verified boot block (keyblock+preamble).  Use
	 * VB2_CONTEXT_FW_SLOT_B to determine whether this refers to slot A or
	 * slot B; vboot will set that flag to the proper state before reading
	 * the vblock.
	 */
	VB2_RES_FW_VBLOCK,

	/*
	 * Kernel verified boot block (keyblock+preamble) for the current
	 * kernel partition.  Used only by vb2api_kernel_load_vblock().
	 * Contents are allowed to change between calls to that function (to
	 * allow multiple kernels to be examined).
	 */
	VB2_RES_KERNEL_VBLOCK,
};

/* Digest ID for vbapi_get_pcr_digest() */
enum vb2_pcr_digest {
	/* Digest based on current developer and recovery mode flags */
	BOOT_MODE_PCR,

	/* SHA-256 hash digest of HWID, from GBB */
	HWID_DIGEST_PCR,
};

/******************************************************************************
 * APIs provided by verified boot.
 *
 * At a high level, call functions in the order described below.  After each
 * call, examine vb2_context.flags to determine whether nvdata or secdata
 * needs to be written.
 *
 * If you need to cause the boot process to fail at any point, call
 * vb2api_fail().  Then check vb2_context.flags to see what data needs to be
 * written.  Then reboot.
 *
 *	Load nvdata from wherever you keep it.
 *
 *	Load secdata_firmware from wherever you keep it.
 *
 *      	If it wasn't there at all (for example, this is the first boot
 *		of a new system in the factory), call
 *		vb2api_secdata_firmware_create() to initialize the data.
 *
 *		If access to your storage is unreliable (reads/writes may
 *		contain corrupt data), you may call
 *		vb2api_secdata_firmware_check() to determine if the data was
 *		valid, and retry reading if it wasn't.  (In that case, you
 *		should also read back and check the data after any time you
 *		write it, to make sure it was written correctly.)
 *
 *	Call vb2api_fw_phase1().  At present, this nominally decides whether
 *	recovery mode is needed this boot.
 *
 *	Call vb2api_fw_phase2().  At present, this nominally decides which
 *	firmware slot will be attempted (A or B).
 *
 *	Call vb2api_fw_phase3().  At present, this nominally verifies the
 *	firmware keyblock and preamble.
 *
 *	Lock down wherever you keep secdata_firmware.  It should no longer be
 *	writable this boot.
 *
 *	Verify the hash of each section of code/data you need to boot the RW
 *	firmware.  For each section:
 *
 *		Call vb2_init_hash() to see if the hash exists.
 *
 *		Load the data for the section.  Call vb2_extend_hash() on the
 *		data as you load it.  You can load it all at once and make one
 *		call, or load and hash-extend a block at a time.
 *
 *		Call vb2_check_hash() to see if the hash is valid.
 *
 *			If it is valid, you may use the data and/or execute
 *			code from that section.
 *
 *			If the hash was invalid, you must reboot.
 *
 * At this point, firmware verification is done, and vb2_context contains the
 * kernel key needed to verify the kernel.  That context should be preserved
 * and passed on to kernel selection.  The kernel selection process may be
 * done by the same firmware image, or may be done by the RW firmware.  The
 * recommended order is:
 *
 *	Load secdata_kernel from wherever you keep it.
 *
 *      	If it wasn't there at all (for example, this is the first boot
 *		of a new system in the factory), call
 *		vb2api_secdata_kernel_create() to initialize the data.
 *
 *		If access to your storage is unreliable (reads/writes may
 *		contain corrupt data), you may call
 *		vb2api_secdata_kernel_check() to determine if the data was
 *		valid, and retry reading if it wasn't.  (In that case, you
 *		should also read back and check the data after any time you
*		write it, to make sure it was written correctly.)
 *
 *	Call vb2api_kernel_phase1().  At present, this decides which key to
 *	use to verify kernel data - the recovery key from the GBB, or the
 *	kernel subkey from the firmware verification stage.
 *
 *	Kernel phase 2 is finding loading, and verifying the kernel partition.
 *
 *	Find a boot device (you're on your own here).
 *
 *	Call vb2api_load_kernel_vblock() for each kernel partition on the
 *	boot device, until one succeeds.
 *
 *	When that succeeds, call vb2api_get_kernel_size() to determine where
 *	the kernel is located in the stream and how big it is.  Load or map
 *	the kernel.  (Again, you're on your own.  This is the responsibility of
 *	the caller so that the caller can choose whether to allocate a buffer,
 *	load the kernel data into a predefined area of RAM, or directly map a
 *	kernel file into the address space.  Note that technically it doesn't
 *	matter whether the kernel data is even in the same file or stream as
 *	the vblock, as long as the caller loads the right data.
 *
 *	Call vb2api_verify_kernel_data() on the kernel data.
 *
 *	If you ran out of kernels before finding a good one, call vb2api_fail()
 *	with an appropriate recovery reason.
 *
 *	Set the VB2_CONTEXT_ALLOW_KERNEL_ROLL_FORWARD flag if the current
 *	kernel partition has the successful flag (that is, it's already known
 *	or assumed to be a functional kernel partition).
 *
 *	Call vb2api_kernel_phase3().  This cleans up from kernel verification
 *	and updates the secure data if needed.
 *
 *	Lock down wherever you keep secdata_kernel.  It should no longer be
 *	writable this boot.
 */

/**
 * Initialize verified boot data structures.
 *
 * Needs to be called once per boot, before using any API functions that
 * accept a vb2_context object.  Sets up the vboot work buffer, as well as
 * vb2_shared_data and vb2_context.  A pointer to the context object is
 * written to ctxptr.  After transitioning between different firmware
 * applications, or any time the context pointer is lost, vb2api_reinit()
 * should be used to restore access to the context and data on the workbuf.
 *
 * If the workbuf needs to be relocated, call vb2api_relocate() instead
 * of copying memory manually.
 *
 * @param workbuf	Workbuf memory location to initialize
 * @param size		Size of workbuf being initialized
 * @param ctxptr	Pointer to a context pointer to be filled in
 * @return VB2_SUCCESS, or non-zero error code.
 */
vb2_error_t vb2api_init(void *workbuf, uint32_t size,
			struct vb2_context **ctxptr);

/**
 * Reinitialize vboot data structures.
 *
 * After transitioning between different firmware applications, or any time the
 * context pointer is lost, this function should be called to restore access to
 * the workbuf.  A pointer to the context object is written to ctxptr.  Returns
 * an error if the vboot work buffer is inconsistent.
 *
 * If the workbuf needs to be relocated, call vb2api_relocate() instead
 * of copying memory manually.
 *
 * @param workbuf	Workbuf memory location to check
 * @param ctxptr	Pointer to a context pointer to be filled in
 * @return VB2_SUCCESS, or non-zero error code.
 */
vb2_error_t vb2api_reinit(void *workbuf, struct vb2_context **ctxptr);

/**
 * Relocate vboot data structures.
 *
 * Move the vboot work buffer from one memory location to another, and expand
 * or contract the workbuf to fit.  The target memory location may be the same
 * as the original (used for a "resize" operation), and it is safe to call this
 * function with overlapping memory regions.
 *
 * A pointer to the context object is written to ctxptr.  Returns an error if
 * the vboot work buffer is inconsistent, or if the new memory space is too
 * small to contain the work buffer.
 *
 * @param new_workbuf	Target workbuf memory location
 * @param cur_workbuf	Original workbuf memory location to relocate
 * @param size		Target size of relocated workbuf
 * @param ctxptr	Pointer to a context pointer to be filled in
 * @return VB2_SUCCESS, or non-zero error code.
 */
vb2_error_t vb2api_relocate(void *new_workbuf, void *cur_workbuf, uint32_t size,
			    struct vb2_context **ctxptr);

/**
 * Check the validity of firmware secure storage context.
 *
 * Checks version and CRC.
 *
 * @param ctx		Context pointer
 * @return VB2_SUCCESS, or non-zero error code if error.
 */
vb2_error_t vb2api_secdata_firmware_check(struct vb2_context *ctx);

/**
 * Create fresh data in firmware secure storage context.
 *
 * Use this only when initializing the secure storage context on a new machine
 * the first time it boots.  Do NOT simply use this if
 * vb2api_secdata_firmware_check() (or any other API in this library) fails;
 * that could allow the secure data to be rolled back to an insecure state.
 *
 * @param ctx		Context pointer
 * @return size of created firmware secure storage data in bytes
 */
uint32_t vb2api_secdata_firmware_create(struct vb2_context *ctx);

/**
 * Check the validity of kernel secure storage context.
 *
 * Checks version, UID, and CRC.
 *
 * @param ctx		Context pointer
 * @return VB2_SUCCESS, or non-zero error code if error.
 */
vb2_error_t vb2api_secdata_kernel_check(struct vb2_context *ctx);

/**
 * Create fresh data in kernel secure storage context.
 *
 * Use this only when initializing the secure storage context on a new machine
 * the first time it boots.  Do NOT simply use this if
 * vb2api_secdata_kernel_check() (or any other API in this library) fails; that
 * could allow the secure data to be rolled back to an insecure state.
 *
 * @param ctx		Context pointer
 * @return size of created kernel secure storage data in bytes
 */
uint32_t vb2api_secdata_kernel_create(struct vb2_context *ctx);

/**
 * Check the validity of firmware management parameters (FWMP) space.
 *
 * Checks size, version, and CRC.  If the struct size is larger than the size
 * passed in, the size pointer is set to the expected full size of the struct,
 * and VB2_ERROR_SECDATA_FWMP_INCOMPLETE is returned.  The caller should
 * re-read the returned number of bytes, and call this function again.
 *
 * @param ctx		Context pointer
 * @param size		Amount of struct which has been read
 * @return VB2_SUCCESS, or non-zero error code if error.
 */
vb2_error_t vb2api_secdata_fwmp_check(struct vb2_context *ctx, uint8_t *size);

/**
 * Report firmware failure to vboot.
 *
 * If the failure occurred after choosing a firmware slot, and the other
 * firmware slot is not known-bad, try the other firmware slot after reboot.
 *
 * If the failure occurred before choosing a firmware slot, or both slots have
 * failed in successive boots, request recovery.
 *
 * This may be called before vb2api_phase1() to indicate errors in the boot
 * process prior to the start of vboot.  On return, the calling firmware should
 * check for updates to secdata and/or nvdata, then reboot.
 *
 * @param reason	Recovery reason
 * @param subcode	Recovery subcode
 */
void vb2api_fail(struct vb2_context *ctx, uint8_t reason, uint8_t subcode);

/**
 * Firmware selection, phase 1.
 *
 * If the returned error is VB2_ERROR_API_PHASE1_RECOVERY, the calling firmware
 * should jump directly to recovery-mode firmware without rebooting.
 *
 * For other errors, the calling firmware should check for updates to secdata
 * and/or nvdata, then reboot.
 *
 * @param ctx		Vboot context
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_fw_phase1(struct vb2_context *ctx);

/**
 * Firmware selection, phase 2.
 *
 * On error, the calling firmware should check for updates to secdata and/or
 * nvdata, then reboot.
 *
 * @param ctx		Vboot context
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_fw_phase2(struct vb2_context *ctx);

/**
 * Firmware selection, phase 3.
 *
 * On error, the calling firmware should check for updates to secdata and/or
 * nvdata, then reboot.
 *
 * On success, the calling firmware should lock down secdata before continuing
 * with the boot process.
 *
 * @param ctx		Vboot context
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_fw_phase3(struct vb2_context *ctx);

/**
 * Same, but for new-style structs.
 */
vb2_error_t vb21api_fw_phase3(struct vb2_context *ctx);

/**
 * Initialize hashing data for the specified tag.
 *
 * @param ctx		Vboot context
 * @param tag		Tag to start hashing (enum vb2_hash_tag)
 * @param size		If non-null, expected size of data for tag will be
 *			stored here on output.
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_init_hash(struct vb2_context *ctx, uint32_t tag,
			     uint32_t *size);

/**
 * Same, but for new-style structs.
 */
vb2_error_t vb21api_init_hash(struct vb2_context *ctx, const struct vb2_id *id,
			      uint32_t *size);

/**
 * Extend the hash started by vb2api_init_hash() with additional data.
 *
 * (This is the same for both old and new style structs.)
 *
 * @param ctx		Vboot context
 * @param buf		Data to hash
 * @param size		Size of data in bytes
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_extend_hash(struct vb2_context *ctx, const void *buf,
			       uint32_t size);

/**
 * Check the hash value started by vb2api_init_hash().
 *
 * @param ctx		Vboot context
 * @return VB2_SUCCESS, or error code on error.
 */
int vb2api_check_hash(struct vb2_context *ctx);

/**
 * Same, but for new-style structs.
 */
vb2_error_t vb21api_check_hash(struct vb2_context *ctx);

/**
 * Check the hash value started by vb2api_init_hash() while retrieving
 * calculated digest.
 *
 * @param ctx			Vboot context
 * @param digest_out		optional pointer to buffer to store digest
 * @param digest_out_size	optional size of buffer to store digest
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_check_hash_get_digest(struct vb2_context *ctx,
					 void *digest_out,
					 uint32_t digest_out_size);

/**
 * Get a PCR digest
 *
 * @param ctx		Vboot context
 * @param which_digest	PCR index of the digest
 * @param dest		Destination where the digest is copied.
 * 			Recommended size is VB2_PCR_DIGEST_RECOMMENDED_SIZE.
 * @param dest_size	IN: size of the buffer pointed by dest
 * 			OUT: size of the copied digest
 * @return VB2_SUCCESS, or error code on error
 */
vb2_error_t vb2api_get_pcr_digest(struct vb2_context *ctx,
				  enum vb2_pcr_digest which_digest,
				  uint8_t *dest, uint32_t *dest_size);

/**
 * Prepare for kernel verification stage.
 *
 * Must be called before other vb2api kernel functions.
 *
 * @param ctx		Vboot context
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_kernel_phase1(struct vb2_context *ctx);

/**
 * Load the verified boot block (vblock) for a kernel.
 *
 * This function may be called multiple times, to load and verify the
 * vblocks from multiple kernel partitions.
 *
 * @param ctx		Vboot context
 * @param stream	Kernel stream
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_load_kernel_vblock(struct vb2_context *ctx);

/**
 * Get the size and offset of the kernel data for the most recent vblock.
 *
 * Valid after a successful call to vb2api_load_kernel_vblock().
 *
 * @param ctx		Vboot context
 * @param offset_ptr	Destination for offset in bytes of kernel data as
 *			reported by vblock.
 * @param size_ptr      Destination for size of kernel data in bytes.
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_get_kernel_size(struct vb2_context *ctx,
				   uint32_t *offset_ptr, uint32_t *size_ptr);

/**
 * Verify kernel data using the previously loaded kernel vblock.
 *
 * Valid after a successful call to vb2api_load_kernel_vblock().  This allows
 * the caller to load or map the kernel data, as appropriate, and pass the
 * pointer to the kernel data into vboot.
 *
 * @param ctx		Vboot context
 * @param buf		Pointer to kernel data
 * @param size		Size of kernel data in bytes
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_verify_kernel_data(struct vb2_context *ctx, const void *buf,
				      uint32_t size);

/**
 * Clean up after kernel verification.
 *
 * Call this after successfully loading a vblock and verifying kernel data,
 * or if you've run out of boot devices and/or kernel partitions.
 *
 * This cleans up intermediate data structures in the vboot context, and
 * updates the version in the secure data if necessary.
 */
vb2_error_t vb2api_kernel_phase3(struct vb2_context *ctx);

/**
 * Read the hardware ID from the GBB, and store it onto the given buffer.
 *
 * @param ctx		Vboot context.
 * @param hwid		Buffer to store HWID, which will be null-terminated.
 * @param size		Maximum size of HWID including null terminator.  HWID
 * 			length may not exceed 256 (VB2_GBB_HWID_MAX_SIZE), so
 * 			this value is suggested.  If size is too small, then
 * 			VB2_ERROR_INVALID_PARAMETER is returned.  Actual size
 * 			of the output HWID string is returned in this pointer,
 * 			also including null terminator.
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2api_gbb_read_hwid(struct vb2_context *ctx, char *hwid,
				 uint32_t *size);

/**
 * Retrieve current GBB flags.
 *
 * See enum vb2_gbb_flag in 2gbb_flags.h for a list of all GBB flags.
 *
 * @param ctx		Vboot context.
 *
 * @return vb2_gbb_flags_t representing current GBB flags.
 */
vb2_gbb_flags_t vb2api_gbb_get_flags(struct vb2_context *ctx);

/*****************************************************************************/
/* APIs provided by the caller to verified boot */

/**
 * Clear the TPM owner.
 *
 * @param ctx		Vboot context
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2ex_tpm_clear_owner(struct vb2_context *ctx);

/**
 * Read a verified boot resource.
 *
 * @param ctx		Vboot context
 * @param index		Resource index to read
 * @param offset	Byte offset within resource to start at
 * @param buf		Destination for data
 * @param size		Amount of data to read
 * @return VB2_SUCCESS, or error code on error.
 */
vb2_error_t vb2ex_read_resource(struct vb2_context *ctx,
				enum vb2_resource_index index, uint32_t offset,
				void *buf, uint32_t size);

/**
 * Print debug output
 *
 * This should work like printf().  If func!=NULL, it will be a string with
 * the current function name; that can be used to generate prettier debug
 * output.  If func==NULL, don't print any extra header/trailer so that this
 * can be used to composite a bigger output string from several calls - for
 * example, when doing a hex dump.
 *
 * @param func		Function name generating output, or NULL.
 * @param fmt		Printf format string
 */
void vb2ex_printf(const char *func, const char *fmt, ...);

/**
 * Initialize the hardware crypto engine to calculate a block-style digest.
 *
 * @param hash_alg	Hash algorithm to use
 * @param data_size	Expected total size of data to hash
 * @return VB2_SUCCESS, or non-zero error code (HWCRYPTO_UNSUPPORTED not fatal).
 */
vb2_error_t vb2ex_hwcrypto_digest_init(enum vb2_hash_algorithm hash_alg,
				       uint32_t data_size);

/**
 * Extend the hash in the hardware crypto engine with another block of data.
 *
 * @param buf		Next data block to hash
 * @param size		Length of data block in bytes
 * @return VB2_SUCCESS, or non-zero error code.
 */
vb2_error_t vb2ex_hwcrypto_digest_extend(const uint8_t *buf, uint32_t size);

/**
 * Finalize the digest in the hardware crypto engine and extract the result.
 *
 * @param digest	Destination buffer for resulting digest
 * @param digest_size	Length of digest buffer in bytes
 * @return VB2_SUCCESS, or non-zero error code.
 */
vb2_error_t vb2ex_hwcrypto_digest_finalize(uint8_t *digest,
					   uint32_t digest_size);

/*
 * Set the current TPM mode value, and validate that it was changed.  If one
 * of the following occurs, the function call fails:
 *   - TPM does not understand the instruction (old version)
 *   - TPM has already left the TpmModeEnabledTentative mode
 *   - TPM responds with a mode other than the requested mode
 *   - Some other communication error occurs
 *  Otherwise, the function call succeeds.
 *
 * @param mode_val       Desired TPM mode to set.  May be one of ENABLED
 *                       or DISABLED from vb2_tpm_mode enum.
 * @returns VB2_SUCCESS, or non-zero error code.
 */
vb2_error_t vb2ex_tpm_set_mode(enum vb2_tpm_mode mode_val);

/*
 * Abort vboot flow due to a failed assertion or broken assumption.
 *
 * Likely due to caller misusing vboot (e.g. calling API functions
 * out-of-order, filling in vb2_context fields inappropriately).
 * Implementation should reboot or halt the machine, or fall back to some
 * alternative boot flow.  Retrying vboot is unlikely to succeed.
 */
void vb2ex_abort(void);

#endif  /* VBOOT_REFERENCE_2API_H_ */
