/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Tests for api.c
 */

#include <stdio.h>

#include "2api.h"
#include "2common.h"
#include "2misc.h"
#include "2nvstorage.h"
#include "2rsa.h"
#include "2secdata.h"
#include "2sysincludes.h"
#include "test_common.h"
#include "vb2_struct.h"
#include "vb2_common.h"

/* Common context for tests */

static uint8_t workbuf[VB2_FIRMWARE_WORKBUF_RECOMMENDED_SIZE]
	__attribute__ ((aligned (VB2_WORKBUF_ALIGN)));
static struct vb2_context *ctx;
static struct vb2_shared_data *sd;
static struct vb2_gbb_header gbb;

const char mock_body[320] = "Mock body";
const int mock_body_size = sizeof(mock_body);
const int mock_algorithm = VB2_ALG_RSA2048_SHA256;
const int mock_hash_alg = VB2_HASH_SHA256;
static const uint8_t mock_hwid_digest[VB2_GBB_HWID_DIGEST_SIZE] = {
	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
	0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
	0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
};
const int mock_sig_size = 64;
static uint8_t digest_result[VB2_SHA256_DIGEST_SIZE];
static const uint32_t digest_result_size = sizeof(digest_result);

/* Mocked function data */

static enum {
	HWCRYPTO_DISABLED,
	HWCRYPTO_ENABLED,
	HWCRYPTO_FORBIDDEN,
} hwcrypto_state;

static int force_dev_mode;
static vb2_error_t retval_vb2_fw_init_gbb;
static vb2_error_t retval_vb2_check_dev_switch;
static vb2_error_t retval_vb2_check_tpm_clear;
static vb2_error_t retval_vb2_select_fw_slot;
static vb2_error_t retval_vb2_load_fw_keyblock;
static vb2_error_t retval_vb2_load_fw_preamble;
static vb2_error_t retval_vb2_digest_finalize;
static vb2_error_t retval_vb2_verify_digest;

/* Type of test to reset for */

enum reset_type {
	FOR_MISC,
	FOR_EXTEND_HASH,
	FOR_CHECK_HASH,
};

static void reset_common_data(enum reset_type t)
{
	struct vb2_fw_preamble *pre;
	struct vb2_packed_key *k;

	memset(workbuf, 0xaa, sizeof(workbuf));

	TEST_SUCC(vb2api_init(workbuf, sizeof(workbuf), &ctx),
		  "vb2api_init failed");
	sd = vb2_get_sd(ctx);

	vb2_nv_init(ctx);

	vb2api_secdata_firmware_create(ctx);
	vb2_secdata_firmware_init(ctx);

	force_dev_mode = 0;
	retval_vb2_fw_init_gbb = VB2_SUCCESS;
	retval_vb2_check_dev_switch = VB2_SUCCESS;
	retval_vb2_check_tpm_clear = VB2_SUCCESS;
	retval_vb2_select_fw_slot = VB2_SUCCESS;
	retval_vb2_load_fw_keyblock = VB2_SUCCESS;
	retval_vb2_load_fw_preamble = VB2_SUCCESS;
	retval_vb2_digest_finalize = VB2_SUCCESS;
	retval_vb2_verify_digest = VB2_SUCCESS;

	memcpy(&gbb.hwid_digest, mock_hwid_digest,
	       sizeof(gbb.hwid_digest));

	sd->preamble_offset = sd->workbuf_used;
	sd->preamble_size = sizeof(*pre);
	vb2_set_workbuf_used(ctx, sd->preamble_offset + sd->preamble_size);
	pre = vb2_member_of(sd, sd->preamble_offset);
	pre->body_signature.data_size = mock_body_size;
	pre->body_signature.sig_size = mock_sig_size;
	if (hwcrypto_state == HWCRYPTO_FORBIDDEN)
		pre->flags = VB2_FIRMWARE_PREAMBLE_DISALLOW_HWCRYPTO;
	else
		pre->flags = 0;

	sd->data_key_offset = sd->workbuf_used;
	sd->data_key_size = sizeof(*k) + 8;
	vb2_set_workbuf_used(ctx, sd->data_key_offset + sd->data_key_size);
	k = vb2_member_of(sd, sd->data_key_offset);
	k->algorithm = mock_algorithm;

	if (t == FOR_EXTEND_HASH || t == FOR_CHECK_HASH)
		vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY, NULL);

	if (t == FOR_CHECK_HASH)
		vb2api_extend_hash(ctx, mock_body, mock_body_size);

	/* Always clear out the digest result. */
	memset(digest_result, 0, digest_result_size);
};

/* Mocked functions */
struct vb2_gbb_header *vb2_get_gbb(struct vb2_context *c)
{
	return &gbb;
}

vb2_error_t vb2_fw_init_gbb(struct vb2_context *c)
{
	return retval_vb2_fw_init_gbb;
}

vb2_error_t vb2_check_dev_switch(struct vb2_context *c)
{
	if (force_dev_mode)
		sd->flags |= VB2_SD_FLAG_DEV_MODE_ENABLED;
	return retval_vb2_check_dev_switch;
}

vb2_error_t vb2_check_tpm_clear(struct vb2_context *c)
{
	return retval_vb2_check_tpm_clear;
}

vb2_error_t vb2_select_fw_slot(struct vb2_context *c)
{
	return retval_vb2_select_fw_slot;
}

vb2_error_t vb2_load_fw_keyblock(struct vb2_context *c)
{
	return retval_vb2_load_fw_keyblock;
}

vb2_error_t vb2_load_fw_preamble(struct vb2_context *c)
{
	return retval_vb2_load_fw_preamble;
}

vb2_error_t vb2_unpack_key_buffer(struct vb2_public_key *key,
				  const uint8_t *buf, uint32_t size)
{
	struct vb2_packed_key *k = (struct vb2_packed_key *)buf;

	if (size != sizeof(*k) + 8)
		return VB2_ERROR_UNPACK_KEY_SIZE;

	key->sig_alg = vb2_crypto_to_signature(k->algorithm);
	key->hash_alg = vb2_crypto_to_hash(k->algorithm);

	return VB2_SUCCESS;
}

vb2_error_t vb2ex_hwcrypto_digest_init(enum vb2_hash_algorithm hash_alg,
				       uint32_t data_size)
{
	switch (hwcrypto_state) {
	case HWCRYPTO_DISABLED:
		return VB2_ERROR_EX_HWCRYPTO_UNSUPPORTED;
	case HWCRYPTO_ENABLED:
		if (hash_alg != mock_hash_alg)
			return VB2_ERROR_SHA_INIT_ALGORITHM;
		else
			return VB2_SUCCESS;
	case HWCRYPTO_FORBIDDEN:
	default:
		return VB2_ERROR_UNKNOWN;
	}
}

vb2_error_t vb2ex_hwcrypto_digest_extend(const uint8_t *buf,
					 uint32_t size)
{
	if (hwcrypto_state != HWCRYPTO_ENABLED)
		return VB2_ERROR_UNKNOWN;

	return VB2_SUCCESS;
}

static void fill_digest(uint8_t *digest, uint32_t digest_size)
{
	/* Set the result to a known value. */
	memset(digest, 0x0a, digest_size);
}

vb2_error_t vb2ex_hwcrypto_digest_finalize(uint8_t *digest,
				   uint32_t digest_size)
{
	if (hwcrypto_state != HWCRYPTO_ENABLED)
		return VB2_ERROR_UNKNOWN;

	if (retval_vb2_digest_finalize == VB2_SUCCESS)
		fill_digest(digest, digest_size);

	return retval_vb2_digest_finalize;
}

vb2_error_t vb2_digest_init(struct vb2_digest_context *dc,
		    enum vb2_hash_algorithm hash_alg)
{
	if (hwcrypto_state == HWCRYPTO_ENABLED)
		return VB2_ERROR_UNKNOWN;
	if (hash_alg != mock_hash_alg)
		return VB2_ERROR_SHA_INIT_ALGORITHM;

	dc->hash_alg = hash_alg;
	dc->using_hwcrypto = 0;

	return VB2_SUCCESS;
}

vb2_error_t vb2_digest_extend(struct vb2_digest_context *dc, const uint8_t *buf,
			      uint32_t size)
{
	if (hwcrypto_state == HWCRYPTO_ENABLED)
		return VB2_ERROR_UNKNOWN;
	if (dc->hash_alg != mock_hash_alg)
		return VB2_ERROR_SHA_EXTEND_ALGORITHM;

	return VB2_SUCCESS;
}

vb2_error_t vb2_digest_finalize(struct vb2_digest_context *dc, uint8_t *digest,
				uint32_t digest_size)
{
	if (hwcrypto_state == HWCRYPTO_ENABLED)
		return VB2_ERROR_UNKNOWN;

	if (retval_vb2_digest_finalize == VB2_SUCCESS)
		fill_digest(digest, digest_size);

	return retval_vb2_digest_finalize;
}

uint32_t vb2_rsa_sig_size(enum vb2_signature_algorithm sig_alg)
{
	return mock_sig_size;
}

vb2_error_t vb2_rsa_verify_digest(const struct vb2_public_key *key,
				  uint8_t *sig, const uint8_t *digest,
				  const struct vb2_workbuf *wb)
{
	return retval_vb2_verify_digest;
}

/* Tests */

static void misc_tests(void)
{
	/* Test secdata_firmware passthru functions */
	reset_common_data(FOR_MISC);
	/* Corrupt secdata_firmware so initial check will fail */
	ctx->secdata_firmware[0] ^= 0x42;
	TEST_EQ(vb2api_secdata_firmware_check(ctx),
		VB2_ERROR_SECDATA_FIRMWARE_CRC,
		"secdata_firmware check");
	TEST_EQ(vb2api_secdata_firmware_create(ctx), VB2_SECDATA_FIRMWARE_SIZE,
		  "secdata_firmware create");
	TEST_SUCC(vb2api_secdata_firmware_check(ctx),
		  "secdata_firmware check 2");

	/* Test fail passthru */
	reset_common_data(FOR_MISC);
	vb2api_fail(ctx, 12, 34);
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		12, "vb2api_fail request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_SUBCODE),
		34, "vb2api_fail subcode");
}

static void phase1_tests(void)
{
	reset_common_data(FOR_MISC);
	TEST_SUCC(vb2api_fw_phase1(ctx), "phase1 good");
	TEST_EQ(sd->recovery_reason, 0, "  not recovery");
	TEST_EQ(ctx->flags & VB2_CONTEXT_RECOVERY_MODE, 0, "  recovery flag");
	TEST_EQ(ctx->flags & VB2_CONTEXT_CLEAR_RAM, 0, "  clear ram flag");
	TEST_EQ(ctx->flags & VB2_CONTEXT_DISPLAY_INIT,
		0, "  display init context flag");
	TEST_EQ(sd->flags & VB2_SD_FLAG_DISPLAY_AVAILABLE,
		0, "  display available SD flag");

	reset_common_data(FOR_MISC);
	retval_vb2_fw_init_gbb = VB2_ERROR_GBB_MAGIC;
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_RECOVERY,
		"phase1 gbb");
	TEST_EQ(sd->recovery_reason, VB2_RECOVERY_GBB_HEADER,
		"  recovery reason");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_RECOVERY_MODE, 0, "  recovery flag");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_CLEAR_RAM, 0, "  clear ram flag");

	/* Dev switch error in normal mode reboots to recovery */
	reset_common_data(FOR_MISC);
	retval_vb2_check_dev_switch = VB2_ERROR_MOCK;
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_MOCK, "phase1 dev switch");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		VB2_RECOVERY_DEV_SWITCH, "  recovery request");

	/* Dev switch error already in recovery mode just proceeds */
	reset_common_data(FOR_MISC);
	vb2_nv_set(ctx, VB2_NV_RECOVERY_REQUEST, VB2_RECOVERY_RO_UNSPECIFIED);
	retval_vb2_check_dev_switch = VB2_ERROR_MOCK;
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_RECOVERY,
		"phase1 dev switch error in recovery");
	TEST_EQ(sd->recovery_reason, VB2_RECOVERY_RO_UNSPECIFIED,
		"  recovery reason");
	/* Check that DISPLAY_AVAILABLE gets set on recovery mode. */
	TEST_NEQ(ctx->flags & VB2_CONTEXT_DISPLAY_INIT,
		 0, "  display init context flag");
	TEST_NEQ(sd->flags & VB2_SD_FLAG_DISPLAY_AVAILABLE,
		 0, "  display available SD flag");

	reset_common_data(FOR_MISC);
	ctx->secdata_firmware[0] ^= 0x42;
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_RECOVERY,
		"phase1 secdata_firmware");
	TEST_EQ(sd->recovery_reason, VB2_RECOVERY_SECDATA_FIRMWARE_INIT,
		"  recovery reason");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_RECOVERY_MODE, 0, "  recovery flag");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_CLEAR_RAM, 0, "  clear ram flag");

	/* Test secdata_firmware-requested reboot */
	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_SECDATA_WANTS_REBOOT;
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_SECDATA_REBOOT,
		"phase1 secdata_firmware reboot normal");
	TEST_EQ(sd->recovery_reason, 0,	"  recovery reason");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_TPM_REQUESTED_REBOOT),
		1, "  tpm reboot request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		0, "  recovery request");

	reset_common_data(FOR_MISC);
	vb2_nv_set(ctx, VB2_NV_TPM_REQUESTED_REBOOT, 1);
	TEST_SUCC(vb2api_fw_phase1(ctx),
		  "phase1 secdata_firmware reboot back normal");
	TEST_EQ(sd->recovery_reason, 0,	"  recovery reason");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_TPM_REQUESTED_REBOOT),
		0, "  tpm reboot request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		0, "  recovery request");

	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_SECDATA_WANTS_REBOOT;
	memset(ctx->secdata_firmware, 0, sizeof(ctx->secdata_firmware));
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_SECDATA_REBOOT,
		"phase1 secdata_firmware reboot normal, "
		"secdata_firmware blank");
	TEST_EQ(sd->recovery_reason, 0,	"  recovery reason");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_TPM_REQUESTED_REBOOT),
		1, "  tpm reboot request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		0, "  recovery request");

	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_SECDATA_WANTS_REBOOT;
	vb2_nv_set(ctx, VB2_NV_TPM_REQUESTED_REBOOT, 1);
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_RECOVERY,
		"phase1 secdata_firmware reboot normal again");
	TEST_EQ(sd->recovery_reason, VB2_RECOVERY_RO_TPM_REBOOT,
		"  recovery reason");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_TPM_REQUESTED_REBOOT),
		1, "  tpm reboot request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		0, "  recovery request");

	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_SECDATA_WANTS_REBOOT;
	vb2_nv_set(ctx, VB2_NV_RECOVERY_REQUEST, VB2_RECOVERY_RO_UNSPECIFIED);
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_SECDATA_REBOOT,
		"phase1 secdata_firmware reboot recovery");
	/* Recovery reason isn't set this boot because we're rebooting first */
	TEST_EQ(sd->recovery_reason, 0, "  recovery reason not set THIS boot");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_TPM_REQUESTED_REBOOT),
		1, "  tpm reboot request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		VB2_RECOVERY_RO_UNSPECIFIED, "  recovery request not cleared");

	reset_common_data(FOR_MISC);
	vb2_nv_set(ctx, VB2_NV_TPM_REQUESTED_REBOOT, 1);
	vb2_nv_set(ctx, VB2_NV_RECOVERY_REQUEST, VB2_RECOVERY_RO_UNSPECIFIED);
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_RECOVERY,
		"phase1 secdata_firmware reboot back recovery");
	TEST_EQ(sd->recovery_reason, VB2_RECOVERY_RO_UNSPECIFIED,
		"  recovery reason");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_TPM_REQUESTED_REBOOT),
		0, "  tpm reboot request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST), 0,
		"  recovery request cleared");

	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_SECDATA_WANTS_REBOOT;
	vb2_nv_set(ctx, VB2_NV_TPM_REQUESTED_REBOOT, 1);
	vb2_nv_set(ctx, VB2_NV_RECOVERY_REQUEST, VB2_RECOVERY_RO_UNSPECIFIED);
	TEST_EQ(vb2api_fw_phase1(ctx), VB2_ERROR_API_PHASE1_RECOVERY,
		"phase1 secdata_firmware reboot recovery again");
	TEST_EQ(sd->recovery_reason, VB2_RECOVERY_RO_UNSPECIFIED,
		"  recovery reason");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_TPM_REQUESTED_REBOOT),
		1, "  tpm reboot request");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST), 0,
		"  recovery request cleared");

	/* Cases for checking DISPLAY_INIT and DISPLAY_AVAILABLE. */
	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_DISPLAY_INIT;
	TEST_SUCC(vb2api_fw_phase1(ctx), "phase1 with DISPLAY_INIT");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_DISPLAY_INIT,
		 0, "  display init context flag");
	TEST_NEQ(sd->flags & VB2_SD_FLAG_DISPLAY_AVAILABLE,
		 0, "  display available SD flag");

	reset_common_data(FOR_MISC);
	vb2_nv_set(ctx, VB2_NV_DISPLAY_REQUEST, 1);
	TEST_SUCC(vb2api_fw_phase1(ctx), "phase1 with DISPLAY_REQUEST");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_DISPLAY_INIT,
		 0, "  display init context flag");
	TEST_NEQ(sd->flags & VB2_SD_FLAG_DISPLAY_AVAILABLE,
		 0, "  display available SD flag");

	reset_common_data(FOR_MISC);
	force_dev_mode = 1;
	TEST_SUCC(vb2api_fw_phase1(ctx), "phase1 in dev mode");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_DISPLAY_INIT,
		 0, "  display init context flag");
	TEST_NEQ(sd->flags & VB2_SD_FLAG_DISPLAY_AVAILABLE,
		 0, "  display available SD flag");
}

static void phase2_tests(void)
{
	reset_common_data(FOR_MISC);
	TEST_SUCC(vb2api_fw_phase2(ctx), "phase2 good");
	TEST_EQ(ctx->flags & VB2_CONTEXT_CLEAR_RAM, 0, "  clear ram flag");
	TEST_EQ(ctx->flags & VB2_CONTEXT_FW_SLOT_B, 0, "  slot b flag");

	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_DEVELOPER_MODE;
	TEST_SUCC(vb2api_fw_phase2(ctx), "phase2 dev");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_CLEAR_RAM, 0, "  clear ram flag");

	reset_common_data(FOR_MISC);
	retval_vb2_check_tpm_clear = VB2_ERROR_MOCK;
	TEST_EQ(vb2api_fw_phase2(ctx), VB2_ERROR_MOCK, "phase2 tpm clear");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		VB2_RECOVERY_TPM_CLEAR_OWNER, "  recovery reason");

	reset_common_data(FOR_MISC);
	retval_vb2_select_fw_slot = VB2_ERROR_MOCK;
	TEST_EQ(vb2api_fw_phase2(ctx), VB2_ERROR_MOCK, "phase2 slot");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		VB2_RECOVERY_FW_SLOT, "  recovery reason");

	/* S3 resume exits before clearing RAM */
	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_S3_RESUME;
	ctx->flags |= VB2_CONTEXT_DEVELOPER_MODE;
	TEST_SUCC(vb2api_fw_phase2(ctx), "phase2 s3 dev");
	TEST_EQ(ctx->flags & VB2_CONTEXT_CLEAR_RAM, 0, "  clear ram flag");
	TEST_EQ(ctx->flags & VB2_CONTEXT_FW_SLOT_B, 0, "  slot b flag");

	reset_common_data(FOR_MISC);
	ctx->flags |= VB2_CONTEXT_S3_RESUME;
	vb2_nv_set(ctx, VB2_NV_FW_TRIED, 1);
	TEST_SUCC(vb2api_fw_phase2(ctx), "phase2 s3");
	TEST_NEQ(ctx->flags & VB2_CONTEXT_FW_SLOT_B, 0, "  slot b flag");
}

static void get_pcr_digest_tests(void)
{
	uint8_t digest[VB2_PCR_DIGEST_RECOMMENDED_SIZE];
	uint8_t digest_org[VB2_PCR_DIGEST_RECOMMENDED_SIZE];
	uint32_t digest_size;

	reset_common_data(FOR_MISC);
	memset(digest_org, 0, sizeof(digest_org));

	digest_size = sizeof(digest);
	memset(digest, 0, sizeof(digest));
	TEST_SUCC(vb2api_get_pcr_digest(
			ctx, BOOT_MODE_PCR, digest, &digest_size),
		  "BOOT_MODE_PCR");
	TEST_EQ(digest_size, VB2_SHA1_DIGEST_SIZE, "BOOT_MODE_PCR digest size");
	TEST_TRUE(memcmp(digest, digest_org, digest_size),
		  "BOOT_MODE_PCR digest");

	digest_size = sizeof(digest);
	memset(digest, 0, sizeof(digest));
	TEST_SUCC(vb2api_get_pcr_digest(
			ctx, HWID_DIGEST_PCR, digest, &digest_size),
		  "HWID_DIGEST_PCR");
	TEST_EQ(digest_size, VB2_GBB_HWID_DIGEST_SIZE,
		"HWID_DIGEST_PCR digest size");
	TEST_FALSE(memcmp(digest, mock_hwid_digest, digest_size),
		   "HWID_DIGEST_PCR digest");

	digest_size = 1;
	TEST_EQ(vb2api_get_pcr_digest(ctx, BOOT_MODE_PCR, digest, &digest_size),
		VB2_ERROR_API_PCR_DIGEST_BUF,
		"BOOT_MODE_PCR buffer too small");

	TEST_EQ(vb2api_get_pcr_digest(
			ctx, HWID_DIGEST_PCR + 1, digest, &digest_size),
		VB2_ERROR_API_PCR_DIGEST,
		"invalid enum vb2_pcr_digest");
}

static void phase3_tests(void)
{
	reset_common_data(FOR_MISC);
	TEST_SUCC(vb2api_fw_phase3(ctx), "phase3 good");

	reset_common_data(FOR_MISC);
	retval_vb2_load_fw_keyblock = VB2_ERROR_MOCK;
	TEST_EQ(vb2api_fw_phase3(ctx), VB2_ERROR_MOCK, "phase3 keyblock");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		VB2_RECOVERY_RO_INVALID_RW, "  recovery reason");

	reset_common_data(FOR_MISC);
	retval_vb2_load_fw_preamble = VB2_ERROR_MOCK;
	TEST_EQ(vb2api_fw_phase3(ctx), VB2_ERROR_MOCK, "phase3 keyblock");
	TEST_EQ(vb2_nv_get(ctx, VB2_NV_RECOVERY_REQUEST),
		VB2_RECOVERY_RO_INVALID_RW, "  recovery reason");
}

static void init_hash_tests(void)
{
	struct vb2_packed_key *k;
	int wb_used_before;
	uint32_t size;

	/* For now, all we support is body signature hash */
	reset_common_data(FOR_MISC);
	wb_used_before = sd->workbuf_used;
	TEST_SUCC(vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY, &size),
		  "init hash good");
	TEST_EQ(sd->hash_offset, wb_used_before, "hash context offset");
	TEST_EQ(sd->hash_size, sizeof(struct vb2_digest_context),
		"hash context size");
	TEST_EQ(sd->workbuf_used,
		vb2_wb_round_up(sd->hash_offset + sd->hash_size),
		"hash uses workbuf");
	TEST_EQ(sd->hash_tag, VB2_HASH_TAG_FW_BODY, "hash tag");
	TEST_EQ(sd->hash_remaining_size, mock_body_size, "hash remaining");

	wb_used_before = sd->workbuf_used;
	TEST_SUCC(vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY, NULL),
		  "init hash again");
	TEST_EQ(sd->workbuf_used, wb_used_before, "init hash reuses context");

	reset_common_data(FOR_MISC);
	TEST_EQ(vb2api_init_hash(ctx, VB2_HASH_TAG_INVALID, &size),
		VB2_ERROR_API_INIT_HASH_TAG, "init hash invalid tag");

	reset_common_data(FOR_MISC);
	sd->preamble_size = 0;
	TEST_EQ(vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY, &size),
		VB2_ERROR_API_INIT_HASH_PREAMBLE, "init hash preamble");

	reset_common_data(FOR_MISC);
	TEST_EQ(vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY + 1, &size),
		VB2_ERROR_API_INIT_HASH_TAG, "init hash unknown tag");

	reset_common_data(FOR_MISC);
	sd->workbuf_used = sd->workbuf_size + VB2_WORKBUF_ALIGN -
			vb2_wb_round_up(sizeof(struct vb2_digest_context));
	TEST_EQ(vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY, &size),
		VB2_ERROR_API_INIT_HASH_WORKBUF, "init hash workbuf");

	reset_common_data(FOR_MISC);
	sd->data_key_size = 0;
	TEST_EQ(vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY, &size),
		VB2_ERROR_API_INIT_HASH_DATA_KEY, "init hash data key");

	reset_common_data(FOR_MISC);
	sd->data_key_size--;
	TEST_EQ(vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY, &size),
		VB2_ERROR_UNPACK_KEY_SIZE, "init hash data key size");

	reset_common_data(FOR_MISC);
	k = vb2_member_of(sd, sd->data_key_offset);
	k->algorithm--;
	TEST_EQ(vb2api_init_hash(ctx, VB2_HASH_TAG_FW_BODY, &size),
		VB2_ERROR_SHA_INIT_ALGORITHM, "init hash algorithm");
}

static void extend_hash_tests(void)
{
	struct vb2_digest_context *dc;

	reset_common_data(FOR_EXTEND_HASH);
	TEST_SUCC(vb2api_extend_hash(ctx, mock_body, 32),
		"hash extend good");
	TEST_EQ(sd->hash_remaining_size, mock_body_size - 32,
		"hash extend remaining");
	TEST_SUCC(vb2api_extend_hash(ctx, mock_body, mock_body_size - 32),
		"hash extend again");
	TEST_EQ(sd->hash_remaining_size, 0, "hash extend remaining 2");

	reset_common_data(FOR_EXTEND_HASH);
	sd->hash_size = 0;
	TEST_EQ(vb2api_extend_hash(ctx, mock_body, mock_body_size),
		VB2_ERROR_API_EXTEND_HASH_WORKBUF, "hash extend no workbuf");

	reset_common_data(FOR_EXTEND_HASH);
	TEST_EQ(vb2api_extend_hash(ctx, mock_body, mock_body_size + 1),
		VB2_ERROR_API_EXTEND_HASH_SIZE, "hash extend too much");

	reset_common_data(FOR_EXTEND_HASH);
	TEST_EQ(vb2api_extend_hash(ctx, mock_body, 0),
		VB2_ERROR_API_EXTEND_HASH_SIZE, "hash extend empty");

	if (hwcrypto_state != HWCRYPTO_ENABLED) {
		reset_common_data(FOR_EXTEND_HASH);
		dc = (struct vb2_digest_context *)
			vb2_member_of(sd, sd->hash_offset);
		dc->hash_alg = mock_hash_alg + 1;
		TEST_EQ(vb2api_extend_hash(ctx, mock_body, mock_body_size),
			VB2_ERROR_SHA_EXTEND_ALGORITHM, "hash extend fail");
	}
}

static void check_hash_tests(void)
{
	struct vb2_fw_preamble *pre;
	const uint32_t digest_value = 0x0a0a0a0a;

	reset_common_data(FOR_CHECK_HASH);
	TEST_SUCC(vb2api_check_hash(ctx), "check hash good");

	reset_common_data(FOR_CHECK_HASH);
	TEST_SUCC(vb2api_check_hash_get_digest(ctx, digest_result,
			digest_result_size), "check hash good with result");
	/* Check the first 4 bytes to ensure it was copied over. */
	TEST_SUCC(memcmp(digest_result, &digest_value, sizeof(digest_value)),
		"check digest value");

	reset_common_data(FOR_CHECK_HASH);
	TEST_EQ(vb2api_check_hash_get_digest(ctx, digest_result,
			digest_result_size - 1),
		VB2_ERROR_API_CHECK_DIGEST_SIZE, "check digest size");
	TEST_NEQ(memcmp(digest_result, &digest_value, sizeof(digest_value)), 0,
		"check digest wrong size");

	reset_common_data(FOR_CHECK_HASH);
	sd->preamble_size = 0;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_API_CHECK_HASH_PREAMBLE, "check hash preamble");

	reset_common_data(FOR_CHECK_HASH);
	sd->hash_size = 0;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_API_CHECK_HASH_WORKBUF, "check hash no workbuf");

	reset_common_data(FOR_CHECK_HASH);
	sd->hash_remaining_size = 1;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_API_CHECK_HASH_SIZE, "check hash size");

	reset_common_data(FOR_CHECK_HASH);
	sd->workbuf_used = sd->workbuf_size;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_API_CHECK_HASH_WORKBUF_DIGEST, "check hash workbuf");

	reset_common_data(FOR_CHECK_HASH);
	retval_vb2_digest_finalize = VB2_ERROR_MOCK;
	TEST_EQ(vb2api_check_hash(ctx), VB2_ERROR_MOCK, "check hash finalize");

	reset_common_data(FOR_CHECK_HASH);
	sd->hash_tag = VB2_HASH_TAG_INVALID;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_API_CHECK_HASH_TAG, "check hash tag");

	reset_common_data(FOR_CHECK_HASH);
	sd->data_key_size = 0;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_API_CHECK_HASH_DATA_KEY, "check hash data key");

	reset_common_data(FOR_CHECK_HASH);
	sd->data_key_size--;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_UNPACK_KEY_SIZE, "check hash data key size");

	reset_common_data(FOR_CHECK_HASH);
	pre = vb2_member_of(sd, sd->preamble_offset);
	pre->body_signature.sig_size++;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_VDATA_SIG_SIZE, "check hash sig size");

	reset_common_data(FOR_CHECK_HASH);
	retval_vb2_digest_finalize = VB2_ERROR_RSA_VERIFY_DIGEST;
	TEST_EQ(vb2api_check_hash(ctx),
		VB2_ERROR_RSA_VERIFY_DIGEST, "check hash finalize");
}

int main(int argc, char* argv[])
{
	misc_tests();
	phase1_tests();
	phase2_tests();
	phase3_tests();

	fprintf(stderr, "Running hash API tests without hwcrypto support...\n");
	hwcrypto_state = HWCRYPTO_DISABLED;
	init_hash_tests();
	extend_hash_tests();
	check_hash_tests();

	fprintf(stderr, "Running hash API tests with hwcrypto support...\n");
	hwcrypto_state = HWCRYPTO_ENABLED;
	init_hash_tests();
	extend_hash_tests();
	check_hash_tests();

	fprintf(stderr, "Running hash API tests with forbidden hwcrypto...\n");
	hwcrypto_state = HWCRYPTO_FORBIDDEN;
	init_hash_tests();
	extend_hash_tests();
	check_hash_tests();

	get_pcr_digest_tests();

	return gTestSuccess ? 0 : 255;
}
