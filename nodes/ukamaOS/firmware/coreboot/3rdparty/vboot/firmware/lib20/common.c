/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Common functions between firmware and kernel verified boot.
 * (Firmware portion)
 */

#include "2rsa.h"
#include "2sha.h"
#include "2sysincludes.h"
#include "vb2_common.h"

vb2_error_t vb2_verify_digest(const struct vb2_public_key *key,
			      struct vb2_signature *sig, const uint8_t *digest,
			      const struct vb2_workbuf *wb)
{
	/* A signature is destroyed in the process of being verified. */
	uint8_t *sig_data = vb2_signature_data_mutable(sig);

	if (sig->sig_size != vb2_rsa_sig_size(key->sig_alg)) {
		VB2_DEBUG("Wrong data signature size for algorithm, "
			  "sig_size=%d, expected %d for algorithm %d.\n",
			  sig->sig_size, vb2_rsa_sig_size(key->sig_alg),
			  key->sig_alg);
		return VB2_ERROR_VDATA_SIG_SIZE;
	}

	return vb2_rsa_verify_digest(key, sig_data, digest, wb);
}

vb2_error_t vb2_verify_data(const uint8_t *data,
		    uint32_t size,
		    struct vb2_signature *sig,
		    const struct vb2_public_key *key,
		    const struct vb2_workbuf *wb)
{
	struct vb2_workbuf wblocal = *wb;
	struct vb2_digest_context *dc;
	uint8_t *digest;
	uint32_t digest_size;
	vb2_error_t rv;

	if (sig->data_size > size) {
		VB2_DEBUG("Data buffer smaller than length of signed data.\n");
		return VB2_ERROR_VDATA_NOT_ENOUGH_DATA;
	}

	/* Digest goes at start of work buffer */
	digest_size = vb2_digest_size(key->hash_alg);
	if (!digest_size)
		return VB2_ERROR_VDATA_DIGEST_SIZE;

	digest = vb2_workbuf_alloc(&wblocal, digest_size);
	if (!digest)
		return VB2_ERROR_VDATA_WORKBUF_DIGEST;

	/* Hashing requires temp space for the context */
	dc = vb2_workbuf_alloc(&wblocal, sizeof(*dc));
	if (!dc)
		return VB2_ERROR_VDATA_WORKBUF_HASHING;

	rv = vb2_digest_init(dc, key->hash_alg);
	if (rv)
		return rv;

	rv = vb2_digest_extend(dc, data, sig->data_size);
	if (rv)
		return rv;

	rv = vb2_digest_finalize(dc, digest, digest_size);
	if (rv)
		return rv;

	vb2_workbuf_free(&wblocal, sizeof(*dc));

	return vb2_verify_digest(key, sig, digest, &wblocal);
}

vb2_error_t vb2_check_keyblock(const struct vb2_keyblock *block, uint32_t size,
			       const struct vb2_signature *sig)
{
	if(size < sizeof(*block)) {
		VB2_DEBUG("Not enough space for keyblock header.\n");
		return VB2_ERROR_KEYBLOCK_TOO_SMALL_FOR_HEADER;
	}

	if (memcmp(block->magic, VB2_KEYBLOCK_MAGIC, VB2_KEYBLOCK_MAGIC_SIZE)) {
		VB2_DEBUG("Not a valid verified boot keyblock.\n");
		return VB2_ERROR_KEYBLOCK_MAGIC;
	}

	if (block->header_version_major != VB2_KEYBLOCK_VERSION_MAJOR) {
		VB2_DEBUG("Incompatible keyblock header version.\n");
		return VB2_ERROR_KEYBLOCK_HEADER_VERSION;
	}

	if (size < block->keyblock_size) {
		VB2_DEBUG("Not enough data for keyblock.\n");
		return VB2_ERROR_KEYBLOCK_SIZE;
	}

	if (vb2_verify_signature_inside(block, block->keyblock_size, sig)) {
		VB2_DEBUG("Keyblock signature off end of block\n");
		return VB2_ERROR_KEYBLOCK_SIG_OUTSIDE;
	}

	/* Make sure advertised signature data sizes are sane. */
	if (block->keyblock_size < sig->data_size) {
		VB2_DEBUG("Signature calculated past end of block\n");
		return VB2_ERROR_KEYBLOCK_SIGNED_TOO_MUCH;
	}

	/* Verify we signed enough data */
	if (sig->data_size < sizeof(struct vb2_keyblock)) {
		VB2_DEBUG("Didn't sign enough data\n");
		return VB2_ERROR_KEYBLOCK_SIGNED_TOO_LITTLE;
	}

	/* Verify data key is inside the block and inside signed data */
	if (vb2_verify_packed_key_inside(block, block->keyblock_size,
					 &block->data_key)) {
		VB2_DEBUG("Data key off end of keyblock\n");
		return VB2_ERROR_KEYBLOCK_DATA_KEY_OUTSIDE;
	}
	if (vb2_verify_packed_key_inside(block, sig->data_size,
					 &block->data_key)) {
		VB2_DEBUG("Data key off end of signed data\n");
		return VB2_ERROR_KEYBLOCK_DATA_KEY_UNSIGNED;
	}

	return VB2_SUCCESS;
}

vb2_error_t vb2_verify_keyblock(struct vb2_keyblock *block, uint32_t size,
				const struct vb2_public_key *key,
				const struct vb2_workbuf *wb)
{
	struct vb2_signature *sig = &block->keyblock_signature;
	vb2_error_t rv;

	/* Sanity check keyblock before attempting signature check of data */
	rv = vb2_check_keyblock(block, size, sig);
	if (rv)
		return rv;

	VB2_DEBUG("Checking keyblock signature...\n");
	rv = vb2_verify_data((const uint8_t *)block, size, sig, key, wb);
	if (rv) {
		VB2_DEBUG("Invalid keyblock signature.\n");
		return VB2_ERROR_KEYBLOCK_SIG_INVALID;
	}

	/* Success */
	return VB2_SUCCESS;
}

vb2_error_t vb2_verify_fw_preamble(struct vb2_fw_preamble *preamble,
				   uint32_t size,
				   const struct vb2_public_key *key,
				   const struct vb2_workbuf *wb)
{
	struct vb2_signature *sig = &preamble->preamble_signature;

	VB2_DEBUG("Verifying preamble.\n");

	/* Sanity checks before attempting signature of data */
	if(size < sizeof(*preamble)) {
		VB2_DEBUG("Not enough data for preamble header\n");
		return VB2_ERROR_PREAMBLE_TOO_SMALL_FOR_HEADER;
	}
	if (preamble->header_version_major !=
	    FIRMWARE_PREAMBLE_HEADER_VERSION_MAJOR) {
		VB2_DEBUG("Incompatible firmware preamble header version.\n");
		return VB2_ERROR_PREAMBLE_HEADER_VERSION;
	}

	if (preamble->header_version_minor < 1) {
		VB2_DEBUG("Only preamble header 2.1+ supported\n");
		return VB2_ERROR_PREAMBLE_HEADER_OLD;
	}

	if (size < preamble->preamble_size) {
		VB2_DEBUG("Not enough data for preamble.\n");
		return VB2_ERROR_PREAMBLE_SIZE;
	}

	/* Check signature */
	if (vb2_verify_signature_inside(preamble, preamble->preamble_size,
					sig)) {
		VB2_DEBUG("Preamble signature off end of preamble\n");
		return VB2_ERROR_PREAMBLE_SIG_OUTSIDE;
	}

	/* Make sure advertised signature data sizes are sane. */
	if (preamble->preamble_size < sig->data_size) {
		VB2_DEBUG("Signature calculated past end of the block\n");
		return VB2_ERROR_PREAMBLE_SIGNED_TOO_MUCH;
	}

	if (vb2_verify_data((const uint8_t *)preamble, size, sig, key, wb)) {
		VB2_DEBUG("Preamble signature validation failed\n");
		return VB2_ERROR_PREAMBLE_SIG_INVALID;
	}

	/* Verify we signed enough data */
	if (sig->data_size < sizeof(struct vb2_fw_preamble)) {
		VB2_DEBUG("Didn't sign enough data\n");
		return VB2_ERROR_PREAMBLE_SIGNED_TOO_LITTLE;
	}

	/* Verify body signature is inside the signed data */
	if (vb2_verify_signature_inside(preamble, sig->data_size,
					&preamble->body_signature)) {
		VB2_DEBUG("Firmware body signature off end of preamble\n");
		return VB2_ERROR_PREAMBLE_BODY_SIG_OUTSIDE;
	}

	/* Verify kernel subkey is inside the signed data */
	if (vb2_verify_packed_key_inside(preamble, sig->data_size,
					 &preamble->kernel_subkey)) {
		VB2_DEBUG("Kernel subkey off end of preamble\n");
		return VB2_ERROR_PREAMBLE_KERNEL_SUBKEY_OUTSIDE;
	}

	/* Success */
	return VB2_SUCCESS;
}
