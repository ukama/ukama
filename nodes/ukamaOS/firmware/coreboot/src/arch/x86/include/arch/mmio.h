/*
 * This file is part of the coreboot project.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; version 2 of the License.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 */

#ifndef __ARCH_MMIO_H__
#define __ARCH_MMIO_H__

#include <stdint.h>

static __always_inline uint8_t read8(
	const volatile void *addr)
{
	return *((volatile uint8_t *)(addr));
}

static __always_inline uint16_t read16(
	const volatile void *addr)
{
	return *((volatile uint16_t *)(addr));
}

static __always_inline uint32_t read32(
	const volatile void *addr)
{
	return *((volatile uint32_t *)(addr));
}

#ifndef __ROMCC__
static __always_inline uint64_t read64(
	const volatile void *addr)
{
	return *((volatile uint64_t *)(addr));
}
#endif

static __always_inline void write8(volatile void *addr,
	uint8_t value)
{
	*((volatile uint8_t *)(addr)) = value;
}

static __always_inline void write16(volatile void *addr,
	uint16_t value)
{
	*((volatile uint16_t *)(addr)) = value;
}

static __always_inline void write32(volatile void *addr,
	uint32_t value)
{
	*((volatile uint32_t *)(addr)) = value;
}

#ifndef __ROMCC__
static __always_inline void write64(volatile void *addr,
	uint64_t value)
{
	*((volatile uint64_t *)(addr)) = value;
}
#endif

#endif /* __ARCH_MMIO_H__ */
