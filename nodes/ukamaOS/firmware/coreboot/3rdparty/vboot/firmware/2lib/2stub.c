/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Stub API implementations which should be implemented by the caller.
 */

#include <stdarg.h>
#include <stdio.h>

#include "2api.h"
#include "2sysincludes.h"

__attribute__((weak))
void vb2ex_printf(const char *func, const char *fmt, ...)
{
#ifdef VBOOT_DEBUG
	va_list ap;
	va_start(ap, fmt);
	if (func)
		fprintf(stderr, "%s: ", func);
	vfprintf(stderr, fmt, ap);
	va_end(ap);
#endif
}

__attribute__((weak))
vb2_error_t vb2ex_tpm_clear_owner(struct vb2_context *ctx)
{
	fprintf(stderr, "%s: function not implemented\n", __func__);
	return VB2_ERROR_EX_UNIMPLEMENTED;
}

__attribute__((weak))
vb2_error_t vb2ex_read_resource(struct vb2_context *ctx,
				enum vb2_resource_index index, uint32_t offset,
				void *buf, uint32_t size)
{
	fprintf(stderr, "%s: function not implemented\n", __func__);
	return VB2_ERROR_EX_UNIMPLEMENTED;
}

__attribute__((weak))
vb2_error_t vb2ex_hwcrypto_digest_init(enum vb2_hash_algorithm hash_alg,
				       uint32_t data_size)
{
	return VB2_ERROR_EX_HWCRYPTO_UNSUPPORTED;
}

__attribute__((weak))
vb2_error_t vb2ex_hwcrypto_digest_extend(const uint8_t *buf, uint32_t size)
{
	return VB2_ERROR_SHA_EXTEND_ALGORITHM;  /* Should not be called. */
}

__attribute__((weak))
vb2_error_t vb2ex_hwcrypto_digest_finalize(uint8_t *digest,
					   uint32_t digest_size)
{
	return VB2_ERROR_SHA_FINALIZE_ALGORITHM;  /* Should not be called. */
}

__attribute__((weak))
vb2_error_t vb2ex_tpm_set_mode(enum vb2_tpm_mode mode_val)
{
	fprintf(stderr, "%s: function not implemented\n", __func__);
	return VB2_ERROR_EX_UNIMPLEMENTED;
}

__attribute__((weak))
void vb2ex_abort(void)
{
	/* Stub simply exits. */
	exit(1);
}
