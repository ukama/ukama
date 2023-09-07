#ifndef STDDEF_H
#define STDDEF_H

#include <commonlib/helpers.h>

typedef long ptrdiff_t;
#ifndef __SIZE_TYPE__
#define __SIZE_TYPE__ unsigned long
#endif
typedef __SIZE_TYPE__ size_t;
/* There is a GCC macro for a size_t type, but not
 * for a ssize_t type. Below construct tricks GCC
 * into making __SIZE_TYPE__ signed.
 */
#define unsigned signed
typedef __SIZE_TYPE__ ssize_t;
#undef unsigned

typedef int wchar_t;
typedef unsigned int wint_t;

#define NULL ((void *)0)

/* The devicetree data structures are only mutable in ramstage. All other
   stages have a constant devicetree. */
#if !ENV_PAYLOAD_LOADER
#define DEVTREE_EARLY 1
#else
#define DEVTREE_EARLY 0
#endif

#if DEVTREE_EARLY
#define DEVTREE_CONST const
#else
#define DEVTREE_CONST
#endif

#if ENV_STAGE_HAS_DATA_SECTION
#define MAYBE_STATIC_NONZERO static
#else
#define MAYBE_STATIC_NONZERO
#endif

#if ENV_STAGE_HAS_BSS_SECTION
#define MAYBE_STATIC_BSS static
#else
#define MAYBE_STATIC_BSS
#endif

#ifndef __ROMCC__
/* Provide a pointer to address 0 that thwarts any "accessing this is
 * undefined behaviour and do whatever" trickery in compilers.
 * Use when you _really_ need to read32(zeroptr) (ie. read address 0).
 */
extern char zeroptr[];
#endif

#endif /* STDDEF_H */
