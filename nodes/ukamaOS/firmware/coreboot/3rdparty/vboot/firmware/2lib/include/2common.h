/* Copyright (c) 2014 The Chromium OS Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style license that can be
 * found in the LICENSE file.
 *
 * Common functions between firmware and kernel verified boot.
 */

#ifndef VBOOT_REFERENCE_2COMMON_H_
#define VBOOT_REFERENCE_2COMMON_H_

#include "2api.h"
#include "2gbb.h"
#include "2return_codes.h"
#include "2sha.h"
#include "2struct.h"
#include "2sysincludes.h"

struct vb2_public_key;

/*
 * Return the min/max of A and B.  This is used in macros which calculate the
 * required buffer size, so can't be turned into a static inline function.
 */
#define VB2_MIN(a, b) ({ \
	typeof(a) __vb2_min_a = (a); \
	typeof(b) __vb2_min_b = (b); \
	__vb2_min_a < __vb2_min_b ? __vb2_min_a : __vb2_min_b; \
	})
#define VB2_MAX(a, b) ({ \
	typeof(a) __vb2_max_a = (a); \
	typeof(b) __vb2_max_b = (b); \
	__vb2_max_a > __vb2_max_b ? __vb2_max_a : __vb2_max_b; \
	})

/* Return the number of elements in an array */
#ifndef ARRAY_SIZE
#define ARRAY_SIZE(a) (sizeof(a) / sizeof((a)[0]))
#endif

/* Platform-dependent debug/assert output macros. */
#define VB2_DEBUG(format, args...) \
	vb2ex_printf(__func__, format, ## args)

#define VB2_DEBUG_RAW(format, args...) \
	vb2ex_printf(NULL, format, ## args)

#define VB2_ASSERT(expr) do { \
	if (!(expr)) { \
		VB2_DEBUG("assertion failed: %s at %s:%d\n", \
			  #expr, __FILE__, __LINE__); \
		vb2ex_abort(); \
		for (;;); \
	} \
} while (0)

#define VB2_DIE(format, args...) do { \
	VB2_DEBUG(format, ## args); \
	vb2ex_abort(); \
	for (;;); \
} while (0)

/*
 * Define test_mockable and for mocking functions when compiled for Chrome OS
 * environment (that is, not for firmware).
 */
#ifndef test_mockable
#ifdef CHROMEOS_ENVIRONMENT
#define test_mockable __attribute__((weak))
#else
#define test_mockable
#endif
#endif

// Have a generic fall-through for different versions of C/C++.
// Taken from boringssl.
#if defined(__cplusplus) && __cplusplus >= 201703L
#define VBOOT_FALLTHROUGH [[fallthrough]]
#elif defined(__cplusplus) && __cplusplus >= 201103L && defined(__clang__)
#define VBOOT_FALLTHROUGH [[clang::fallthrough]]
#elif defined(__cplusplus) && __cplusplus >= 201103L && defined(__GNUC__) && \
    __GNUC__ >= 7
#define VBOOT_FALLTHROUGH [[gnu::fallthrough]]
#elif defined(__GNUC__) && __GNUC__ >= 7 // gcc 7
#define VBOOT_FALLTHROUGH __attribute__ ((fallthrough))
#elif defined(__clang__)
#if __has_attribute(fallthrough)
#define VBOOT_FALLTHROUGH __attribute__ ((fallthrough))
#else // clang versions that do not support fallthrough.
#define VBOOT_FALLTHROUGH
#endif
#else // C++11 on gcc 6, and all other cases
#define VBOOT_FALLTHROUGH
#endif

/**
 * Round up a number to a multiple of VB2_WORKBUF_ALIGN
 *
 * @param v		Number to round up
 * @return The number, rounded up.
 */
static inline uint32_t vb2_wb_round_up(uint32_t v)
{
	return (v + VB2_WORKBUF_ALIGN - 1) & ~(VB2_WORKBUF_ALIGN - 1);
}

/* Work buffer */
struct vb2_workbuf {
	uint8_t *buf;
	uint32_t size;
};

/**
 * Initialize a work buffer.
 *
 * @param wb		Work buffer to init
 * @param buf		Pointer to work buffer data
 * @param size		Size of work buffer data in bytes
 */
void vb2_workbuf_init(struct vb2_workbuf *wb, uint8_t *buf, uint32_t size);

/**
 * Allocate space in a work buffer.
 *
 * Note that the returned buffer will always be aligned to VB2_WORKBUF_ALIGN.
 *
 * The work buffer acts like a stack, and detailed tracking of allocs and frees
 * is not done.  The caller must track the size of each allocation and free via
 * vb2_workbuf_free() in the reverse order they were allocated.
 *
 * An acceptable alternate workflow inside a function is to pass in a const
 * work buffer, then make a local copy.  Allocations done to the local copy
 * then don't change the passed-in work buffer, and will effectively be freed
 * when the local copy goes out of scope.
 *
 * @param wb		Work buffer
 * @param size		Requested size in bytes
 * @return A pointer to the allocated space, or NULL if error.
 */
void *vb2_workbuf_alloc(struct vb2_workbuf *wb, uint32_t size);

/**
 * Reallocate space in a work buffer.
 *
 * Note that the returned buffer will always be aligned to VB2_WORKBUF_ALIGN.
 * The work buffer acts like a stack, so this must only be done to the most
 * recently allocated buffer.
 *
 * @param wb		Work buffer
 * @param oldsize	Old allocation size in bytes
 * @param newsize	Requested size in bytes
 * @return A pointer to the allocated space, or NULL if error.
 */
void *vb2_workbuf_realloc(struct vb2_workbuf *wb, uint32_t oldsize,
			  uint32_t newsize);

/**
 * Free the preceding allocation.
 *
 * Note that the work buffer acts like a stack, and detailed tracking of
 * allocs and frees is not done.  The caller must track the size of each
 * allocation and free them in reverse order.
 *
 * @param wb		Work buffer
 * @param size		Size of data to free
 */
void vb2_workbuf_free(struct vb2_workbuf *wb, uint32_t size);

/* Check if a pointer is aligned on an align-byte boundary */
#define vb2_aligned(ptr, align) (!(((uintptr_t)(ptr)) & ((align) - 1)))

/**
 * Safer memcmp() for use in crypto.
 *
 * Compares the buffers to see if they are equal.  Time taken to perform
 * the comparison is dependent only on the size, not the relationship of
 * the match between the buffers.  Note that unlike memcmp(), this only
 * indicates inequality, not which buffer is lesser.
 *
 * @param s1		First buffer
 * @param s2		Second buffer
 * @param size		Number of bytes to compare
 * @return 0 if match or size=0, non-zero if at least one byte mismatched.
 */
vb2_error_t vb2_safe_memcmp(const void *s1, const void *s2, size_t size);

/**
 * Align a buffer and check its size.
 *
 * @param **ptr		Pointer to pointer to align
 * @param *size		Points to size of buffer pointed to by *ptr
 * @param align		Required alignment (must be power of 2)
 * @param want_size	Required size
 * @return VB2_SUCCESS, or non-zero if error.
 */
vb2_error_t vb2_align(uint8_t **ptr, uint32_t *size, uint32_t align,
		      uint32_t want_size);

/**
 * Return offset of ptr from base.
 *
 * @param base		Base pointer
 * @param ptr		Pointer at some offset from base
 * @return The offset of ptr from base.
 */
ptrdiff_t vb2_offset_of(const void *base, const void *ptr);

/**
 * Return member of given object.
 *
 * @param parent	Pointer to parent object
 * @param offset	Offset from base
 * @return Pointer to child object.
 */
void *vb2_member_of(void *parent, ptrdiff_t offset);

/**
 * Return expected signature size for a signature/hash algorithm pair
 *
 * @param sig_alg	Signature algorithm
 * @param hash_alg	Hash algorithm
 * @return The signature size, or zero if error / unsupported algorithm.
 */
uint32_t vb2_sig_size(enum vb2_signature_algorithm sig_alg,
		      enum vb2_hash_algorithm hash_alg);

/**
 * Return a key ID for an unsigned hash algorithm.
 *
 * @param hash_alg	Hash algorithm to return key for
 * @return A pointer to the key ID for that hash algorithm with
 *	   sig_alg=VB2_SIG_NONE, or NULL if error.
 */
const struct vb2_id *vb2_hash_id(enum vb2_hash_algorithm hash_alg);

/* Size of work buffer sufficient for vb2_verify_digest() worst case. */
#define VB2_VERIFY_DIGEST_WORKBUF_BYTES VB2_VERIFY_RSA_DIGEST_WORKBUF_BYTES

/* Size of work buffer sufficient for vb2_verify_data() worst case. */
#define VB2_VERIFY_DATA_WORKBUF_BYTES					\
	(VB2_SHA512_DIGEST_SIZE +					\
	 VB2_MAX(VB2_VERIFY_DIGEST_WORKBUF_BYTES,			\
		 sizeof(struct vb2_digest_context)))

/* Size of work buffer sufficient for vb2_verify_keyblock() worst case. */
#define VB2_KEYBLOCK_VERIFY_WORKBUF_BYTES VB2_VERIFY_DATA_WORKBUF_BYTES

/* Size of work buffer sufficient for vb2_verify_fw_preamble() worst case. */
#define VB2_VERIFY_FIRMWARE_PREAMBLE_WORKBUF_BYTES VB2_VERIFY_DATA_WORKBUF_BYTES

/*
 * Size of work buffer sufficient for vb2_verify_kernel_preamble() worst
 * case.
 */
#define VB2_VERIFY_KERNEL_PREAMBLE_WORKBUF_BYTES VB2_VERIFY_DATA_WORKBUF_BYTES

/**
 * Verify the data pointed to by a subfield is inside the parent data.
 *
 * The subfield has a header pointed to by member, and a separate data
 * field at an offset relative to the header.  That is:
 *
 *   struct parent {
 *     (possibly other parent fields)
 *     struct member {
 *        (member header fields)
 *     };
 *     (possibly other parent fields)
 *   };
 *   (possibly some other parent data)
 *   (member data)
 *   (possibly some other parent data)
 *
 * @param parent		Parent data
 * @param parent_size		Parent size in bytes
 * @param member		Subfield header
 * @param member_size		Size of subfield header in bytes
 * @param member_data_offset	Offset of member data from start of member
 * @param member_data_size	Size of member data in bytes
 * @return VB2_SUCCESS, or non-zero if error.
 */
vb2_error_t vb2_verify_member_inside(const void *parent, size_t parent_size,
				     const void *member, size_t member_size,
				     ptrdiff_t member_data_offset,
				     size_t member_data_size);

/*
 * Helper function to get data pointed to by a public key.
 */
static inline uint8_t *vb2_packed_key_data_mutable(
	struct vb2_packed_key *key)
{
	return (uint8_t *)key + key->key_offset;
}

static inline const uint8_t *vb2_packed_key_data(
	const struct vb2_packed_key *key)
{
	return (const uint8_t *)key + key->key_offset;
}

/**
 * Verify a packed key is fully contained in its parent data
 *
 * @param parent	Parent data
 * @param parent_size	Parent size in bytes
 * @param key		Packed key pointer
 * @return VB2_SUCCESS, or non-zero if error.
 */
static inline vb2_error_t vb2_verify_packed_key_inside(
	const void *parent,
	uint32_t parent_size,
	const struct vb2_packed_key *key)
{
	return vb2_verify_member_inside(parent, parent_size,
					key, sizeof(*key),
					key->key_offset, key->key_size);
}

/*
 * Helper functions to get data pointed to by a public key or signature.
 */
static inline uint8_t *vb2_signature_data_mutable(
	struct vb2_signature *sig)
{
	return (uint8_t *)sig + sig->sig_offset;
}

static inline const uint8_t *vb2_signature_data(
	const struct vb2_signature *sig)
{
	return (const uint8_t *)sig + sig->sig_offset;
}

/**
 * Verify a signature is fully contained in its parent data
 *
 * @param parent	Parent data
 * @param parent_size	Parent size in bytes
 * @param sig		Signature pointer
 * @return VB2_SUCCESS, or non-zero if error.
 */
static inline vb2_error_t vb2_verify_signature_inside(
	const void *parent,
	uint32_t parent_size,
	const struct vb2_signature *sig)
{
	return vb2_verify_member_inside(parent, parent_size,
					sig, sizeof(*sig),
					sig->sig_offset, sig->sig_size);
}

#endif  /* VBOOT_REFERENCE_2COMMON_H_ */
